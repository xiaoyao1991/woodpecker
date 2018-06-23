package main

import (
	"log"
	"os"
	"fmt"
	"bufio"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/ethereum/go-ethereum/ethdb"
	"time"
	"github.com/syndtr/goleveldb/leveldb"
	"sync"
	"github.com/xiaoyao1991/woodpecker/snapshot"
	"github.com/xiaoyao1991/go-ethereum/rlp"
	"github.com/golang/protobuf/proto"
	"gopkg.in/cheggaaa/pb.v1"
)

const numRecords = 37973685
const batchSize = 1000
const snapshotFile = "snapshot.dat"
const edgeLDBPath = "./edge"
const adsLDBPath = "./data"
const numThreads = 63

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func frontlineWithVerification() {
	db, err := leveldb.OpenFile(edgeLDBPath, nil)
	failOnError(err, "Fail to open ldb")
	defer db.Close()

	file, err := os.Open(snapshotFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	queue := make(chan string, numThreads * 2)
	var wg sync.WaitGroup

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()
			batchCount := 0
			var batch *leveldb.Batch

			for s := range queue {
				if batch == nil {
					batch = new(leveldb.Batch)
					batchCount = 0
				}

				// verify
				account := &snapshot.Account{}
				err := proto.Unmarshal(common.Hex2Bytes(s), account)
				failOnError(err, "Fail to unmarshal")

				// verify storage root
				if account.Detail.Storage != nil && len(account.Detail.Storage.Items) > 0 {
					diskdb, _ := ethdb.NewMemDatabase()
					db := trie.NewDatabase(diskdb)
					tree, _ := trie.New(common.Hash{}, db)

					for _, item := range account.Detail.Storage.Items {
						err := tree.TryUpdate(item.Key, item.Val)
						failOnError(err, "Fail during constructing storage trie")
					}

					if tree.Hash() != common.BytesToHash(account.Detail.StorageRoot) {
						log.Fatal("StorageRoot not match!", common.Bytes2Hex(tree.Root()), "\t", common.Bytes2Hex(account.Detail.StorageRoot))
					}
				}

				var accountDetails [][]byte
				accountDetails = append(accountDetails, account.Detail.Nonce)
				accountDetails = append(accountDetails, account.Detail.Balance)
				accountDetails = append(accountDetails, account.Detail.StorageRoot)
				accountDetails = append(accountDetails, account.Detail.CodeHash)
				encoded, err := rlp.EncodeToBytes(accountDetails)
				failOnError(err, "Fail to RLP encode")
				if common.BytesToHash(encoded) != common.BytesToHash(account.Detail.Rlp) {
					log.Fatal("RLP not match!")
				}

				batch.Put(account.Addr, common.Hex2Bytes(s))
				batchCount ++

				if batchCount >= batchSize {
					if err := db.Write(batch, nil); err != nil {
						log.Fatal(err)
					}

					batchCount = 0
					batch = nil
				}
			}

			if err := db.Write(batch, nil); err != nil {
				log.Fatal(err)
			}

			fmt.Println("[META] Closing gid:", gid)
		}(i)
	}

	bar := pb.StartNew(numRecords)
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024*1024)
	for scanner.Scan() {
		bar.Increment()
		queue <- scanner.Text()
	}
	bar.FinishPrint("End")
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	close(queue)

	wg.Wait()
}

func fullTableScan() {
	db, err := leveldb.OpenFile(edgeLDBPath, nil)
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
	ldb, err := leveldb.OpenFile(edgeLDBPath, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ldb.Close()

	diskdb, err := ethdb.NewLDBDatabase(adsLDBPath, 1024, 16)
	failOnError(err, "Error creating DB")

	db := trie.NewDatabase(diskdb)
	tree, err := trie.New(common.Hash{}, db)
	failOnError(err, "Error creating trie")

	bar := pb.StartNew(numRecords)
	iter := ldb.NewIterator(nil, nil)
	for iter.Next() {
		bar.Increment()

		key := iter.Key()
		val := iter.Value()
		k := make([]byte, len(key))
		v := make([]byte, len(val))
		copy(k, key)
		copy(v, val)

		account := &snapshot.Account{}
		err := proto.Unmarshal(v, account)
		failOnError(err, "Fail to unmarshal")

		tree.TryUpdate(k, account.Detail.Rlp)

		if account.Detail.Storage != nil && len(account.Detail.Storage.Items) > 0 {
			batch := diskdb.NewBatch()
			for _, item := range account.Detail.Storage.Items {
				batch.Put(item.Key, item.Val)
			}
			err := batch.Write()
			failOnError(err, "Fail to write storage batch")
		}
	}

	bar.FinishPrint("End")
	iter.Release()
	failOnError(iter.Error(), "Iter error")

	fmt.Println("Root is now: ", common.Bytes2Hex(tree.Root()))
}

func main() {
	startTs := time.Now()
	frontlineWithVerification()
	elapsedTime := time.Since(startTs)
	fmt.Println("[BUILDING FRONTLINE + VERIFICATION] Time taken: ", elapsedTime.Seconds())

	startTs = time.Now()
	fullTableScan()
	elapsedTime = time.Since(startTs)
	fmt.Println("[FULL TABLE SCAN] Time taken: ", elapsedTime.Seconds())

	startTs = time.Now()
	rebuildFromFrontline()
	elapsedTime = time.Since(startTs)
	fmt.Println("[REBUILDING ADS] Time taken: ", elapsedTime.Seconds())
}