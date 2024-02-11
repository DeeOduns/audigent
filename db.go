package main

import "fmt"

// KeyValue represents a key-value pair.
type KeyValue struct {
	Key   string
	Value interface{}
}

// CustomMap represents a custom map data structure.
type CustomMap struct {
	entries []KeyValue
}

// Put adds a key-value pair to the map.
func (m *CustomMap) Put(key string, value interface{}) {
	// Check if the key already exists, if yes, update the value
	for i := range m.entries {
		if m.entries[i].Key == key {
			m.entries[i].Value = value
			return
		}
	}
	// If the key doesn't exist, append a new entry
	m.entries = append(m.entries, KeyValue{Key: key, Value: value})
}

// Get retrieves the value associated with the given key.
func (m *CustomMap) Get(key string) (interface{}, bool) {
	for _, entry := range m.entries {
		if entry.Key == key {
			return entry.Value, true
		}
	}
	return nil, false
}

// Remove deletes the key-value pair associated with the given key.
func (m *CustomMap) Remove(key string) {
	for i, entry := range m.entries {
		if entry.Key == key {
			// Delete the entry by replacing it with the last entry and truncating the slice
			m.entries[i] = m.entries[len(m.entries)-1]
			m.entries = m.entries[:len(m.entries)-1]
			return
		}
	}
}

func main() {
	// Create a new CustomMap
	myMap := CustomMap{}

	// Put some key-value pairs
	myMap.Put("name", "John")
	myMap.Put("age", 30)
	myMap.Put("city", "New York")

	// Retrieve values
	if value, ok := myMap.Get("name"); ok {
		fmt.Println("Name:", value)
	}
	if value, ok := myMap.Get("age"); ok {
		fmt.Println("Age:", value)
	}

	// Remove a key-value pair
	myMap.Remove("city")

	// Check if a key exists
	if _, ok := myMap.Get("city"); !ok {
		fmt.Println("City not found")
	}
}
