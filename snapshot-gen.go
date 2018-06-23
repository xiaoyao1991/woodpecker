package main

import (
	"log"
	"os"
	"strings"
	"fmt"
	"bufio"
	"github.com/ethereum/go-ethereum/common"
	"time"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/xiaoyao1991/woodpecker/snapshot"
	"github.com/golang/protobuf/proto"
	"github.com/ethereum/go-ethereum/crypto"
	"gopkg.in/cheggaaa/pb.v1"
)

const accountFile = "data.txt"
const storageLDBPath = "./storage"
const outputSnapshotFile = "snapshot.dat"
const numRecords = 37973685

// emptyRoot is the known root hash of an empty trie.
var emptyRoot = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

// emptyState is the known hash of an empty state trie entry.
var emptyState = crypto.Keccak256Hash(nil)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func buildSnapshot() {
	db, err := leveldb.OpenFile(storageLDBPath, nil)
	failOnError(err, "Fail to open db")
	defer db.Close()

	accountFile, err := os.Open(accountFile)
	failOnError(err, "Fail to open account file")
	defer accountFile.Close()

	out, _ := os.OpenFile(outputSnapshotFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer out.Close()

	bar := pb.StartNew(numRecords)
	scanner := bufio.NewScanner(accountFile)
	for scanner.Scan() {
		bar.Increment()

		kvs := strings.Split(scanner.Text(), ":")
		val := common.Hex2Bytes(kvs[1])

		var decoded [][]byte
		err = rlp.DecodeBytes(val, &decoded)
		failOnError(err, "Fail to RLP decode")

		if len(decoded) < 4 {
			log.Fatal("The value for the address must be an array of 4 elements")
		}

		accountDetail := &snapshot.AccountDetail{
			Nonce: decoded[0],
			Balance: decoded[1],
			StorageRoot: decoded[2],
			CodeHash: decoded[3],
			Rlp: val,
			Storage: nil,
		}

		storageRoot := decoded[2]
		storageRootHash := common.BytesToHash(storageRoot)
		if storageRootHash != emptyRoot && storageRootHash != emptyState {
			storageVal, err := db.Get(storageRoot, nil)
			failOnError(err, fmt.Sprintf("Fail to retrieve storageRoot %s", common.Bytes2Hex(storageRoot)))

			storageItems := &snapshot.StorageItems{}
			err = proto.Unmarshal(storageVal, storageItems)
			failOnError(err, "Fail to unmarshall")

			accountDetail.Storage = storageItems
		}

		account := &snapshot.Account{
			Addr: common.Hex2Bytes(kvs[0]),
			Detail: accountDetail,
		}

		b, err := proto.Marshal(account)
		failOnError(err, "Fail to marshall")
		out.WriteString(fmt.Sprintf("%s\n", common.Bytes2Hex(b)))
	}
	bar.FinishPrint("End")

	failOnError(scanner.Err(), "Fail to scan account file")
}

func main() {
	startTs := time.Now()
	buildSnapshot()
	elapsedTime := time.Since(startTs)
	fmt.Println("[BUILD SNAPSHOT] Time taken: ", elapsedTime.Seconds())
}