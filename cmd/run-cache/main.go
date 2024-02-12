package main

import (
	"cache"
	"fmt"
	"time"
)

type Cache interface {
	// Set will store the key value pair with a given TTL.
	Set(key, value []byte, ttl time.Duration)

	// Get returns the value stored using `key`.
	//
	// If the key is not present value will be set to nil.
	Get(key []byte) (value []byte, ttl time.Duration)
}

func main() {
	// Create a new DB
	var myMap *cache.Database
	//var myMap Cache
	myMap = cache.CreateDatabase(1000)

	// Put some key-value pairs
	myMap.Set([]byte("name"), []byte("John"), 1000000000)
	myMap.Set([]byte("age"), []byte("30"), 1000000000)
	myMap.Set([]byte("city"), []byte("New York"), 1000000000)
	myMap.Set([]byte("city2"), []byte("New York2"), 1000000000)

	// Retrieve values
	if value, _ := myMap.Get([]byte("name")); value != nil {
		fmt.Println("Name:", string(value))
	}
	if value, _ := myMap.Get([]byte("city2")); value != nil {
		fmt.Println("Age:", string(value))
	}

	// Remove a key-value pair
	myMap.Remove([]byte("city2"))

	// Check if a key exists
	if value, _ := myMap.Get([]byte("city")); value == nil {
		fmt.Println("City not found")
	}
}
