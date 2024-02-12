package main

import (
	"bytes"
	"fmt"
	"time"
)

// KeyValue represents a key-value pair.
type KeyValue struct {
	key   []byte
	value []byte
	ttl   time.Duration
}

// DB represents a custom map data structure.
type DB struct {
	entries []KeyValue
}

// Put adds a key-value pair to the map.
func (m *DB) Put(key, value []byte, ttl time.Duration) {
	// Check if the key already exists, if yes, update the value
	for i := range m.entries {
		if bytes.Equal(m.entries[i].key, key) {
			m.entries[i].value = value
			return
		}
	}
	// If the key doesn't exist, append a new entry
	m.entries = append(m.entries, KeyValue{key: key, value: value, ttl: ttl})
}

// Get retrieves the value associated with the given key.
func (m *DB) Get(key []byte) ([]byte, bool) {
	for _, entry := range m.entries {
		if bytes.Equal(entry.key, key) {
			return entry.value, true
		}
	}
	return nil, false
}

// Remove deletes the key-value pair associated with the given key.
func (m *DB) Remove(key []byte) {
	for i, entry := range m.entries {
		if bytes.Equal(entry.key, key) {
			// Delete the entry by replacing it with the last entry and truncating the slice
			m.entries[i] = m.entries[len(m.entries)-1]
			m.entries = m.entries[:len(m.entries)-1]
			return
		}
	}
}

func main() {
	// Create a new DB
	myMap := DB{}

	// Put some key-value pairs
	myMap.Put([]byte("name"), []byte("John"), 1000000000)
	myMap.Put([]byte("age"), []byte("30"), 1000000000)
	myMap.Put([]byte("city"), []byte("New York"), 1000000000)

	// Retrieve values
	if value, ok := myMap.Get([]byte("name")); ok {
		fmt.Println("Name:", string(value))
	}
	if value, ok := myMap.Get([]byte("age")); ok {
		fmt.Println("Age:", string(value))
	}

	// Remove a key-value pair
	myMap.Remove([]byte("city"))

	// Check if a key exists
	if _, ok := myMap.Get([]byte("city")); !ok {
		fmt.Println("City not found")
	}
}
