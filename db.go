package cache

import (
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
	return &Database{records: make([]Record, 1, 1)}
}

// Sets a record into the cache
func (db *Database) Set(key, value []byte, ttl time.Duration) {
	// add absolute expiry time for the record, for later clean up
	expiryTime := time.Now().Add(ttl)

	// Check if the key already exists, if yes, update the value
	_, index, err := db.Find(key)
	if err == nil {
		db.records[index].value = value
		db.records[index].ttl = ttl
		db.records[index].expiryTime = expiryTime
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
func (db *Database) RemoveExpiredRecords() {
	checkNum := 20
	repeatCleanUp := true
	repeatCleanUpThreshold := int(checkNum / 4)

	dbSize := db.GetSize()
	rand.Seed(time.Now().UnixNano())
	for repeatCleanUp && dbSize >= checkNum {
		selectedValues := make([]Record, checkNum)

		for i := 0; i < checkNum; i++ {
			randomIndex := rand.Intn(dbSize-1) + 1 // Generate random index
			selectedValues[i] = db.records[randomIndex]
		}

		numRemoved := 0
		for i := 0; i < len(selectedValues); i++ {
			if selectedValues[i].expiryTime.After(time.Now()) {
				db.PopAtIndex(i) // Remove record from database
				numRemoved++
			}
		}
		repeatCleanUp = (numRemoved > repeatCleanUpThreshold)
	}
}
