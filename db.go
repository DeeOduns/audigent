package cache

import (
	"bytes"
	"time"
)

// Record represents the db record that is added or retrieved
type Record struct {
	key         []byte
	value       []byte
	ttl         time.Duration
	expiry_time time.Time
}

// Database represents the cache
type Database struct {
	size    int
	records []Record
}

func CreateDatabase(capacity int) *Database {
	return &Database{
		size:    0,
		records: make([]Record, capacity+1, capacity+1),
	}
}

// Sets a record into the cache
func (db *Database) Set(key, value []byte, ttl time.Duration) {
	// add absolute expiry time for the record, for later clean up
	expiry_time := time.Now().Add(ttl)

	// Check if the key already exists, if yes, update the value
	for i := range db.records {
		if bytes.Equal(db.records[i].key, key) {
			db.records[i].value = value
			db.records[i].ttl = ttl
			db.records[i].expiry_time = expiry_time
			return
		}
	}

	// If the key doesn't exist, add the new record
	new_record := Record{
		key:         key,
		value:       value,
		ttl:         ttl,
		expiry_time: expiry_time,
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

// Remove deletes the key-value pair associated with the given key.
func (m *Database) Remove(key []byte) {
	for i, record := range m.records {
		if bytes.Equal(record.key, key) {
			// Delete the record by replacing it with the last record and truncating the slice
			m.records[i] = m.records[len(m.records)-1]
			m.records = m.records[:len(m.records)-1]
			return
		}
	}
}
