package main

import (
	"log"
	"os"
	"strings"
	"fmt"
	"bufio"
	"github.com/ethereum/go-ethereum/common"
	"time"
	"sync"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/xiaoyao1991/woodpecker/snapshot"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/golang/protobuf/proto"
)

const batchSize = 1000
const numRecords = 37973685
const numStorageLines = 80279576
const numThreads = 63

// emptyRoot is the known root hash of an empty trie.
var emptyRoot = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

// emptyState is the known hash of an empty state trie entry.
var emptyState = crypto.Keccak256Hash(nil)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func getStorage(dbpath string) {
	f, _ := os.OpenFile("storage.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	file, err := os.Open("data.txt")
	failOnError(err, "Fail to open data.txt")
	defer file.Close()

	diskdb, err := ethdb.NewLDBDatabase(dbpath, 1024, 16)
	failOnError(err, "Error creating DB")
	defer diskdb.Close()

	queue := make(chan string, numThreads * 2)
	var wg sync.WaitGroup

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()
			var buf strings.Builder

			for s := range queue {
				buf.Reset()

				kvs := strings.Split(s, ":")
				val := common.Hex2Bytes(kvs[1])

				var decoded [][]byte
				err = rlp.DecodeBytes(val, &decoded)
				failOnError(err, "fail to RLP decode")

				if len(decoded) < 4 {
					log.Fatal("The value for the address must be an array of 4 elements")
				}

				storageRoot := decoded[2]
				storageRootHash := common.BytesToHash(storageRoot)
				if storageRootHash == emptyRoot || storageRootHash == emptyState {
					continue
				}

				db := trie.NewDatabase(diskdb)
				tree, err := trie.NewSecure(storageRootHash, db, 1024)
				failOnError(err, "Error creating Trie")

				it := trie.NewIterator(tree.NodeIterator(nil))
				for it.Next() {
					buf.WriteString(fmt.Sprintf("%s:%s:%s\n", common.Bytes2Hex(storageRoot), common.Bytes2Hex(it.Key), common.Bytes2Hex(it.Value)))
				}

				if _, err := f.Write([]byte(buf.String())); err != nil {
					log.Fatal(err)
				}
			}

			fmt.Println("[META] Closing gid:", gid)
		}(i)
	}

	scanner := bufio.NewScanner(file)
	count := 0
	progress := 0.00
	for scanner.Scan() {
		count ++
		if count >= numRecords / 10000 {
			count = 0
			progress += 0.01
			fmt.Println(progress)
		}
		queue <- scanner.Text()
	}
	failOnError(scanner.Err(), "Error scanning")

	close(queue)

	wg.Wait()
}

func processStorage() {
	db, err := leveldb.OpenFile("./storage", nil)
	failOnError(err, "Fail to open db")
	defer db.Close()

	file, err := os.Open("storage.txt")
	failOnError(err, "Fail to open storage.txt")
	defer file.Close()

	currStorageRoot := ""
	storageAgg := make(map[string]string)
	serialized := false
	countDistinctStorageTries := 0

	count := 0
	progress := 0.00

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		count ++
		if count >= numStorageLines / 10000 {
			count = 0
			progress += 0.01
			fmt.Println(progress)
		}

		kvs := strings.Split(scanner.Text(), ":")

		storageRoot := kvs[0]
		storageKey := kvs[1]
		storageVal := kvs[2]

		if currStorageRoot == storageRoot {
			storageAgg[storageKey] = storageVal
		} else {
			if !serialized {
				var items []*snapshot.StorageItem
				for k, v := range storageAgg {
					storageItem := &snapshot.StorageItem{Key: common.Hex2Bytes(k), Val: common.Hex2Bytes(v)}
					items = append(items, storageItem)
				}

				storageItems := &snapshot.StorageItems{Items: items}
				out, err := proto.Marshal(storageItems)
				failOnError(err, "Fail to pb encode")

				err = db.Put(common.Hex2Bytes(currStorageRoot), out, nil)
				failOnError(err, "Fail to write to db")
				countDistinctStorageTries ++
				serialized = true
			}

			// start next agg
			hasKey, _ := db.Has(common.Hex2Bytes(storageRoot), nil)
			if hasKey {
				continue
			}
			currStorageRoot = storageRoot
			storageAgg = make(map[string]string)
			serialized = false
		}
	}

	// last one
	if !serialized {
		var items []*snapshot.StorageItem
		for k, v := range storageAgg {
			storageItem := &snapshot.StorageItem{Key: common.Hex2Bytes(k), Val: common.Hex2Bytes(v)}
			items = append(items, storageItem)
		}

		storageItems := &snapshot.StorageItems{Items: items}
		out, err := proto.Marshal(storageItems)
		failOnError(err, "Fail to pb encode")

		err = db.Put(common.Hex2Bytes(currStorageRoot), out, nil)
		failOnError(err, "Fail to write to db")
		countDistinctStorageTries ++
		serialized = true
	}
	failOnError(scanner.Err(), "Error scanning")

	fmt.Printf("[Process] Processed %d distinct storage tries\n", countDistinctStorageTries)
}

func main() {
	startTs := time.Now()
	//getStorage()
	processStorage()
	elapsedTime := time.Since(startTs)
	fmt.Println("[BUILDING FRONTLINE] Time taken: ", elapsedTime.Seconds())
}