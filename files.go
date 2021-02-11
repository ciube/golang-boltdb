package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

type Config struct {
	Birthday time.Time
	Height   float64
}

type Weight struct {
	Date   time.Time
	Weight float64
}

type Entry struct {
	Date     time.Time
	Food     string
	Calories int
}

// setupDB
// set the test.db and add 3 buckets (DB [root], WEIGHT, ENTRIES)
func setupDB() (*bolt.DB, error) {
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("could not open db, %v", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("DB"))
		if err != nil {
			return fmt.Errorf("could not create root bucket: %v", err)
		}
		_, err = root.CreateBucketIfNotExists([]byte("WEIGHT"))
		if err != nil {
			return fmt.Errorf("could not create weight bucket: %v", err)
		}
		_, err = root.CreateBucketIfNotExists([]byte("ENTRIES"))
		if err != nil {
			return fmt.Errorf("could not create entry bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not set up buckets, %v", err)
	}
	fmt.Println("DB Setup Done")
	return db, nil
}

func setConfig(db *bolt.DB, config Config) error {
	confBytes, err := json.Marshal((config))
	if err != nil {
		return fmt.Errorf("could not marshal config json: %v", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		err = tx.Bucket([]byte("DB")).Put([]byte("CONFIG"), confBytes)
		if err != nil {
			return fmt.Errorf("could not set config: %v", err)
		}
		return nil
	})
	fmt.Println("Set Config")
	return err
}

func main() {
	db, _ := setupDB()
	err := setConfig(db, Config{time.Now(), 84.15})
	if err != nil {
		fmt.Println("something wrong, %v", err)
	}
}
