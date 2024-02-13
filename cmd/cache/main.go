package main

import (
	"cache"
	"fmt"
	"time"
)

func main() {
	// Create a new DB/cache
	var myDB cache.Cache
	myDB = cache.CreateDatabase()

	// Add some values to cache
	ttl := time.Duration(10) * time.Minute
	myDB.Set([]byte("name"), []byte("John"), ttl)
	myDB.Set([]byte("age"), []byte("30"), ttl)
	myDB.Set([]byte("city"), []byte("New York"), ttl)

	// Retrieve values
	if value, _ := myDB.Get([]byte("name")); value != nil {
		fmt.Println("name:", string(value))
	}
	if value, _ := myDB.Get([]byte("age")); value != nil {
		fmt.Println("age:", string(value))
	}
}
