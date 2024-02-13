package cache

import (
	"fmt"
	"math/rand"
	"time"
)

// Record represents the db record that is added or retrieved
type Record struct {
	key        []byte
	value      []byte
	ttl        time.Duration
	expiryTime time.Time
}

// Database represents the cache
type Database struct {
	records []Record
}

func CreateDatabase() *Database {
	return &Database{records: []Record{}}
}

// Sets a record into the cache
func (db *Database) Set(key, value []byte, ttl time.Duration) {
	// add absolute expiry time for the record, for later clean up
	expiryTime := time.Now().Add(ttl)

	// Check if the key already exists, if yes, update the value
	_, idx, err := db.Find(key)
	if err == nil {
		db.records[idx].value = value
		db.records[idx].ttl = ttl
		db.records[idx].expiryTime = expiryTime
		return
	}

	// If the key doesn't exist, add the new record
	new_record := Record{
		key:        key,
		value:      value,
		ttl:        ttl,
		expiryTime: expiryTime,
	}
	db.Push(new_record)
}

// Get retrieves the value associated with the given key.
func (db *Database) Get(key []byte) ([]byte, time.Duration) {
	record, _, err := db.Find(key)
	if err != nil {
		return nil, -1
	}
	return record.value, record.ttl
}

// Implements active removal of expired records in cache
func (db *Database) RemoveExpiredRecords(checkSize int) {
	repeatCleanUp := true
	repeatCleanUpThreshold := int(checkSize / 4)

	dbSize := db.GetSize()
	fmt.Printf("Active removal process begins, current size: %d...\n", dbSize)
	for repeatCleanUp && dbSize >= checkSize {
		numRemoved := 0
		for i, dbSize := 0, db.GetSize(); i < checkSize && dbSize > 1; i, dbSize = i+1, db.GetSize() {
			randomCheckIndex := rand.Intn(dbSize-1) + 1 // Generate random idx
			if db.records[randomCheckIndex].expiryTime.Before(time.Now()) {
				db.RemoveAtIndex(randomCheckIndex)
				numRemoved++
			}
		}
		repeatCleanUp = (numRemoved > repeatCleanUpThreshold)
	}
	fmt.Printf("Active removal complete! new size: %d\n", db.GetSize())
}
