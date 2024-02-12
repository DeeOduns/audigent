package cache

import (
	"bytes"
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	var myDB Cache
	myDB = CreateDatabase(1000)
	var numOfKeys = 100
	// add key0 -> value0, key1 -> value1,... to myDB
	for i := 0; i < numOfKeys; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		myDB.Set([]byte(key), []byte(value), 10)
	}

	// check existing key value
	value, _ := myDB.Get([]byte("key76"))
	expected_value := "value76"
	if !bytes.Equal(value, []byte(expected_value)) {
		t.Errorf("Incorrect value, got: %s, expects: %s.", value, expected_value)
	}

	// check key that does not exist in DB
	value, _ = myDB.Get([]byte("key206"))
	if value != nil {
		t.Errorf("value should not exist")
	}
}

func BenchmarkSet(b *testing.B) {

}
