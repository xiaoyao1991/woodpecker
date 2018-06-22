package main

import (
	"github.com/boltdb/bolt"
	"log"
	"os"
	"strings"
	"fmt"
	"bufio"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/ethereum/go-ethereum/ethdb"
	"time"
	"sync"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func frontline() {
	db, err := bolt.Open("frontline.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	file, err := os.Open("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("statetrie"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	numThreads := 63
	queue := make(chan string, numThreads * 2)
	var wg sync.WaitGroup

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()

			for s := range queue {
				kvs := strings.Split(s, ":")
				db.Batch(func(tx *bolt.Tx) error {
					b := tx.Bucket([]byte("statetrie"))
					err := b.Put(common.Hex2Bytes(kvs[0]), common.Hex2Bytes(kvs[1]))
					return err
				})
			}

			fmt.Println("[META] Closing gid:", gid)
		}(i)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		queue <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	close(queue)

	wg.Wait()

}

func fullTableScan() {
	db, err := bolt.Open("frontline.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("statetrie"))

		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
		}

		return nil
	})
}

func rebuildFromFrontline() {
	bdb, err := bolt.Open("frontline.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer bdb.Close()

	//diskdb, err := ethdb.NewLDBDatabase("/data3/xiaoyao/mydata", 1024, 16)
	diskdb, err := ethdb.NewMemDatabase()
	failOnError(err, "Error creating DB")

	db := trie.NewDatabase(diskdb)
	trie, err := trie.New(common.Hash{}, db)

	bdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("statetrie"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			trie.TryUpdate(k, v)
		}

		return nil
	})

	fmt.Println("Root is now: ", common.Bytes2Hex(trie.Root()))
}

func main() {
	startTs := time.Now()
	frontline()
	elapsedTime := time.Since(startTs)
	fmt.Println("[BUILDING FRONTLINE] Time taken: ", elapsedTime.Seconds())

	startTs = time.Now()
	fullTableScan()
	elapsedTime = time.Since(startTs)
	fmt.Println("[FULL TABLE SCAN] Time taken: ", elapsedTime.Seconds())

	startTs = time.Now()
	rebuildFromFrontline()
	elapsedTime = time.Since(startTs)
	fmt.Println("[REBUILDING] Time taken: ", elapsedTime.Seconds())
}