package main

// Simple program for playing with boltdb.
//
// To dump the raw database file:
//
//  hexdump -C bolt.db

import (
	"log"

	"github.com/boltdb/bolt"
)

var world = []byte("world")

func main() {
	db, err := bolt.Open("bolt.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//key1 := []byte("key1")
	//value1 := []byte("value333")
	key2 := []byte("anotherkey3")
	value2 := []byte("value333")

	// store some data
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(world)
		if err != nil {
			return err
		}

		//_ = bucket.Put(key1, value1)
		_ = bucket.Put(key2, value2)

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
