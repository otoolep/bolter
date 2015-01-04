package main

// Simple program for loading BoltDB with data.
//
// To dump the raw database file:
//
//  hexdump -C bolt.db

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"strconv"

	"github.com/boltdb/bolt"
)

const (
	secondsPerDay = 86400
	numDays       = 1
	pointsPerTx   = 100
)

var series = []byte("series")

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	// Profiling requested?
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

	// Start writing data.
	value := []byte("12345678")
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(series)
		if err != nil {
			return err
		}

		for i := 0; i < 1000000; i++ {
			key := []byte(strconv.Itoa(i) + "cpu.load")
			_ = bucket.Put(key, value)
		}
		return nil
	})
}
