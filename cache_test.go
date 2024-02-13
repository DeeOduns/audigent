package cache

import (
	"bytes"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	var cache Cache
	cache = CreateDatabase()
	cache.Set([]byte("name"), []byte("John Doe"), time.Duration(10))

	value, _ := cache.Get([]byte("name"))
	if !bytes.Equal(value, []byte("John Doe")) {
		t.Errorf("Unexpected cache value")
	}
}
