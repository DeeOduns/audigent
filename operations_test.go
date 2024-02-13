package cache

import "testing"

func TestSort(t *testing.T) {
	var cache = CreateDatabase()

	var unorderedKeys = []string{"name", "age", "rack_number", "compute_time"}
	for _, element := range unorderedKeys {
		cache.Set([]byte(element), []byte(""), 0)
	}

	var expectedOrderedKeys = []string{"age", "compute_time", "name", "rack_number"}
	for i := 0; i < len(unorderedKeys); i++ {
		if expectedOrderedKeys[i] != string(cache.records[i].key) {
			t.Errorf("Keys not sorted correctly in cache")
			break
		}
	}
}
