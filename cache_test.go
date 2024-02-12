package cache

import (
	"bytes"
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	var myDB Cache
	myDB = CreateDatabase()
	var numOfKeys = 100
	// add key0 -> value0, key1 -> value1,... to myDB
	for n := 0; n < numOfKeys; n++ {
		key := fmt.Sprintf("key%d", n)
		value := fmt.Sprintf("value%d", n)
		myDB.Set([]byte(key), []byte(value), 10)
	}

	// check existing key value
	value, _ := myDB.Get([]byte("key76"))
	expectedValue := "value76"
	if !bytes.Equal(value, []byte(expectedValue)) {
		t.Errorf("Incorrect value, got: %s, expects: %s.", value, expectedValue)
	}

	// check for key that does not exist in DB
	value, _ = myDB.Get([]byte("key206"))
	if value != nil {
		t.Errorf("value should not exist")
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

	myDB = CreateDatabase()
	for n := 0; n < b.N; n++ {
		myDB.Set([]byte("name"), []byte("john"), 10)
	}

	if myDB.GetSize() != 1 {
		b.Errorf("cache should produce only 1 allocation if changes are made on same key in cache")
	}
}
