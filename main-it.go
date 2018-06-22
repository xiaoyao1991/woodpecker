package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"time"
	"log"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/ethereum/go-ethereum/ethdb"
	"os"
	"strings"
	"bufio"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func stream(root string, dbpath string) {
	diskdb, err := ethdb.NewLDBDatabase(dbpath, 1024, 16)
	failOnError(err, "Error creating DB")
	defer diskdb.Close()

	// data logger
	f, _ := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	db := trie.NewDatabase(diskdb)
	tree, err := trie.NewSecure(common.HexToHash(root), db, 1024)
	failOnError(err, "Error creating Trie")

	valueCount := 0

	it := trie.NewIterator(tree.NodeIterator(nil))
	for it.Next() {
		valueCount++
		fmt.Println(common.Bytes2Hex(it.Key), common.Bytes2Hex(it.Value))
		if _, err := f.Write([]byte(fmt.Sprintf("%s:%s\n", common.Bytes2Hex(it.Key), common.Bytes2Hex(it.Value)))); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("[ITERATOR] Value Count: %d\n", valueCount)
}

func rebuild(dataFile string) {
	file, err := os.Open(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//diskdb, err := ethdb.NewLDBDatabase("/data3/xiaoyao/mydata", 1024, 16)
	diskdb, err := ethdb.NewMemDatabase()
	failOnError(err, "Error creating DB")

	db := trie.NewDatabase(diskdb)
	trie, err := trie.New(common.Hash{}, db)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		kvs := strings.Split(scanner.Text(), ":")
		trie.TryUpdate(common.Hex2Bytes(kvs[0]), common.Hex2Bytes(kvs[1]))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Root is now: ", common.Bytes2Hex(trie.Root()))
}

func main() {
	startTs := time.Now()
	stream("1177207cfca75f3eb517d97be70e071c66db3850691b2a69386a704a419bf9bb", "/Users/xiaoyaoqian/projects/cs598am/ethereum-bootstrap/data/geth/chaindata")
	//stream("587ee10ea8af4d603d243b0146275a906d552ec061e1a424a0fa5d4dcea042c0", "/data3/xiaoyao/data/geth/chaindata/")
	elapsedTime := time.Since(startTs)
	fmt.Println("[STREAMING] Time taken: ", elapsedTime.Seconds())

	startTs = time.Now()
	rebuild("data.txt")
	elapsedTime = time.Since(startTs)
	fmt.Println("[REBUILDING] Time taken: ", elapsedTime.Seconds())
}