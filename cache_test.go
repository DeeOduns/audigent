package cache

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkGet(b *testing.B) {
	myDB := CreateDatabase()
	// add key0:value0, key1:value1,... to myDB
	for n := 0; n <= b.N; n++ {
		key := fmt.Sprintf("key%d", n)
		value := fmt.Sprintf("value%d", n)
		myDB.Set([]byte(key), []byte(value), 10)
	}

	// check randomly for an existing key-value in cache
	var checkIndex = 0
	if b.N == 1 {
		checkIndex = 1
	} else {
		checkIndex = rand.Intn(myDB.GetSize()-1) + 1
	}
	randomKey := fmt.Sprintf("key%d", checkIndex)
	expectedValue := fmt.Sprintf("value%d", checkIndex)
	cacheValue, _ := myDB.Get([]byte(randomKey))
	if !bytes.Equal(cacheValue, []byte(expectedValue)) {
		b.Errorf("Incorrect value, got: %s, expects: %s.", cacheValue, expectedValue)
	}

	// check for a key that does not exist in cache
	nonExistingKey := fmt.Sprintf("key%d", b.N+10)
	cacheValue, _ = myDB.Get([]byte(nonExistingKey))
	if cacheValue != nil {
		b.Errorf("Cached value should not exist, got %s", cacheValue)
	}
}

func BenchmarkSetAllocation(b *testing.B) {
	myDB := CreateDatabase()
	for n := 0; n < b.N; n++ {
		key := fmt.Sprintf("key%d", n)
		value := fmt.Sprintf("value%d", n)
		myDB.Set([]byte(key), []byte(value), 10)
	}
	if myDB.GetSize() != b.N {
		b.Errorf("cache should produce one 1 allocation per operation")
	}

	// set record with same key to cache
	myDB = CreateDatabase()
	for n := 0; n < b.N; n++ {
		myDB.Set([]byte("name"), []byte("john"), 10)
	}

	if myDB.GetSize() != 1 {
		b.Errorf("cache should produce only 1 allocation if changes are made on same key in cache")
	}
}

func BenchmarkRemoveExpiredRecords(b *testing.B) {
	myDB := CreateDatabase()

	// create keys with expired ttl
	expiredTtl := time.Duration(-10)
	for n := 0; n < b.N; n++ {
		key := fmt.Sprintf("key%d", n)
		value := fmt.Sprintf("value%d", n)
		myDB.Set([]byte(key), []byte(value), expiredTtl)
	}

	var defaultCheckSize = 25
	if b.N > defaultCheckSize {
		myDB.RemoveExpiredRecords(defaultCheckSize)
		if myDB.GetSize() >= defaultCheckSize {
			b.Errorf("Expired records not deleted. Added %d records, %d records left", b.N, myDB.GetSize())
		}
	} else {
		fmt.Printf("Cache size is to small for cleanup. Only added %d records.\n", myDB.GetSize())
	}
}

func BenchmarkRemoveExpiredRecordsNoOp(b *testing.B) {
	myDB := CreateDatabase()

	// create keys that has not expired
	futureTtl := time.Duration(10) * time.Hour
	for n := 0; n < b.N; n++ {
		key := fmt.Sprintf("key%d", n)
		value := fmt.Sprintf("value%d", n)
		myDB.Set([]byte(key), []byte(value), futureTtl)
	}

	var defaultCheckSize = 25
	if b.N > defaultCheckSize {
		myDB.RemoveExpiredRecords(defaultCheckSize)
		if myDB.GetSize() != b.N {
			b.Errorf("Unexpectedly removed records. Added %d records, %d records left", b.N, myDB.GetSize())
		}
	} else {
		fmt.Printf("Cache size is to small for cleanup. Only added %d records.\n", myDB.GetSize())
	}
}
