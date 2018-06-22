package main

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/ethereum/go-ethereum/rlp"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"sync/atomic"
	"time"
	"sync"
	"github.com/streadway/amqp"
	"log"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/ethereum/go-ethereum/ethdb"
	"os"
	"bufio"
	"strings"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func putToQueue(ch *amqp.Channel, q amqp.Queue, b []byte) bool {
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        b,
		})
	failOnError(err, "Failed to publish a message")

	return err == nil
}

func dig(root string, dbpath string) {
	// leveldb
	db, _ := leveldb.OpenFile(dbpath, nil)
	defer db.Close()

	// rabbitmq
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"statetrie", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	// meta
	numThreads := 64

	// data logger
	f, _ := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	putToQueue(ch, q, common.Hex2Bytes(root))

	var count uint64 = 1
	var valueCount uint64 = 0
	var wg sync.WaitGroup

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()

			msgs, err := ch.Consume(
				q.Name, // queue
				"",     // consumer
				true,   // auto-ack
				false,  // exclusive
				false,  // no-local
				false,  // no-wait
				nil,    // args
			)
			failOnError(err, "Failed to register a consumer")

			for {
				select {
				case msg := <- msgs:
					//fmt.Printf("[MAIN LOOP] Getting %s...\n", common.ToHex(msg.Body))
					data, err := db.Get(msg.Body, nil)
					var localCount uint64

					failOnError(err, "leveldb not found")

					var decoded [][]byte
					err = rlp.DecodeBytes(data, &decoded)
					failOnError(err, "fail to RLP decode")

					if len(decoded) == 17 {
						for j := 0; j < len(decoded) - 1; j++ {
							if len(decoded[j]) > 0 {
								//fmt.Printf("\t[BRANCH] Found at %d:	%s\n", j, common.ToHex(decoded[j]))
								localCount++
								putToQueue(ch, q, decoded[j])
							}
						}

						//if len(decoded[16]) > 0 {
						//	fmt.Printf("\t[BRANCH] Found branch value: %s\n", common.ToHex(decoded[16]))
						//}
					} else {
						key := common.ToHex(decoded[0])
						val := common.ToHex(decoded[1])
						fmt.Printf("\t[NON-BRANCH] key: %s, val: %s\n", key, val)
						localCount++
						atomic.AddUint64(&valueCount, 1)

						if _, err := f.Write([]byte(fmt.Sprintf("%s:%s\n", key, val))); err != nil {
							log.Fatal(err)
						}
					}

					atomic.AddUint64(&count, localCount)

				case <-time.After(10 * time.Second):
					fmt.Println("[META] Closing gid:", gid)
					return
				}
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("[MAIN LOOP] Node Count: %d, Value Count: %d\n", atomic.LoadUint64(&count), atomic.LoadUint64(&valueCount))
}

func rebuild(dataFile string) {
	file, err := os.Open(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	diskdb, err := ethdb.NewLDBDatabase("/data3/xiaoyao/mydata", 1024, 16)
	//diskdb, err := ethdb.NewMemDatabase()
	failOnError(err, "Error creating DB")

	db := trie.NewDatabase(diskdb)
	trie, err := trie.New(common.Hash{}, db)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		kvs := strings.Split(scanner.Text(), ":")
		trie.TryUpdate(common.Hex2Bytes(kvs[0][2:]), common.Hex2Bytes(kvs[1][2:]))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Root is now: ", common.Bytes2Hex(trie.Root()))
}

func main() {
	startTs := time.Now()
	//dig("587ee10ea8af4d603d243b0146275a906d552ec061e1a424a0fa5d4dcea042c0", "/data3/xiaoyao/data/geth/chaindata/")
	dig("1177207cfca75f3eb517d97be70e071c66db3850691b2a69386a704a419bf9bb", "/Users/xiaoyaoqian/projects/cs598am/ethereum-bootstrap/data/geth/chaindata")
	//rebuild("data.txt")

	elapsedTime := time.Since(startTs)
	fmt.Println("[TIME] Time taken: ", elapsedTime.Seconds())
}