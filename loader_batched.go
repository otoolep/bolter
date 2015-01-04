package main

// Simple program for loading BoltDB with data.
//
// To dump the raw database file:
//
//  hexdump -C bolt.db

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/boltdb/bolt"
)

var series = []byte("series")

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var op = flag.String("op", "write", "perform read or write. Default is write")
var rkey = flag.String("rkey", "", "when reading, get this key")
var rn = flag.Int("rn", 1, "when reading, number of subsequent keys to fetch. Default is 1")
var wn = flag.Int("wn", 1000000, "when writing, number of keys to write. Default is 1,000,000")
var wb = flag.Int("wb", 1000, "when writing, writes per transaction. Default is 1,000")
var verbose = flag.Bool("v", false, "display progress")

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// Open the database.
	db, err := bolt.Open("bolt.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if *op == "write" {
		// Start writing data.
		for i := 0; i < *wn; i += 1000 {
			err = db.Update(func(tx *bolt.Tx) error {
				bucket, err := tx.CreateBucketIfNotExists(series)
				if err != nil {
					return err
				}

				for j := 0; j < 10000; j++ {
					n := fmt.Sprintf("%08d", i+j)
					key := []byte(n + "cpu.load")
					_ = bucket.Put(key, []byte(n))
				}
				if *verbose {
					fmt.Println(i, "keys inserted")
				}
				return nil
			})
		}
	} else {
		if *rkey == "" {
			fmt.Println("No read key specified")
			os.Exit(1)
		}

		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(series)
			if b == nil {
				fmt.Println("Bucket not found")
				os.Exit(1)
			}

			if *rn == 1 {
				v := b.Get([]byte(*rkey))
				if *verbose {
					fmt.Println("Key:", *rkey, "value:", string(v))
				}
			} else {
				c := b.Cursor()
				for k, v := c.Seek([]byte(*rkey)); k != nil && *rn > 0; k, v = c.Next() {
					if *verbose {
						fmt.Println("Key:", string(k), "value:", string(v))
					}
					*rn--
				}
			}
			return nil
		})

	}

}
