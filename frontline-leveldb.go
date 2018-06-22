package main

import (
	"log"
	"os"
	"strings"
	"fmt"
	"bufio"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/ethereum/go-ethereum/ethdb"
	"time"
	"github.com/syndtr/goleveldb/leveldb"
	"sync"
)

const batchSize = 1000

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func frontline() {
	db, err := leveldb.OpenFile("./data", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	file, err := os.Open("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	numThreads := 63
	queue := make(chan string, numThreads * 2)
	var wg sync.WaitGroup

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()
			count := 0
			var batch *leveldb.Batch

			for s := range queue {
				if batch == nil {
					batch = new(leveldb.Batch)
					count = 0
				}

				kvs := strings.Split(s, ":")
				batch.Put(common.Hex2Bytes(kvs[0]), common.Hex2Bytes(kvs[1]))
				count ++

				if count >= batchSize {
					if err := db.Write(batch, nil); err != nil {
						log.Fatal(err)
					}

					count = 0
					batch = nil
				}
			}

			if err := db.Write(batch, nil); err != nil {
				log.Fatal(err)
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
	db, err := leveldb.OpenFile("./data", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
	}
	iter.Release()
}

func rebuildFromFrontline() {
	ldb, err := leveldb.OpenFile("./data", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ldb.Close()

	//diskdb, err := ethdb.NewLDBDatabase("./mydata", 1024, 16)
	diskdb, err := ethdb.NewMemDatabase()
	failOnError(err, "Error creating DB")

	db := trie.NewDatabase(diskdb)
	trie, err := trie.New(common.Hash{}, db)

	iter := ldb.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		val := iter.Value()

		k := make([]byte, len(key))
		v := make([]byte, len(val))

		copy(k, key)
		copy(v, val)

		trie.TryUpdate(k, v)
	}

	iter.Release()
	err = iter.Error()
	if err != nil{
		log.Fatal(err)
	}

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