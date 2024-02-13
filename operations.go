package cache

import (
	"bytes"
	"fmt"
	"slices"
	"sort"
)

// return true if key1 < key2
func comparator(key1, key2 []byte) bool {
	result := bytes.Compare(key1, key2)
	return result < 0
}

func (db *Database) Sort() {
	sort.Slice(db.records, func(i, j int) bool {
		return comparator(db.records[i].key, db.records[j].key)
	})
}

func (db *Database) Add(record Record) {
	db.records = append(db.records, record)
	db.Sort()
}

// find record using binary search
func (db *Database) Find(key []byte) (record *Record, idx int, err error) {
	idx, found := slices.BinarySearchFunc(
		db.records, Record{key: key}, func(a, b Record) int {
			return bytes.Compare(a.key, b.key)
		},
	)

	if found {
		return &db.records[idx], idx, nil
	} else {
		return record, idx, fmt.Errorf("record is not found in cache")
	}
}

// remove record at specified idx
func (db *Database) RemoveAtIndex(idx int) error {
	dbSize := db.GetSize() - 1
	if idx > dbSize {
		return fmt.Errorf("idx is out of range")
	}

	db.records[idx], db.records[dbSize] = db.records[dbSize], db.records[idx]
	db.records = db.records[:dbSize]
	db.Sort()
	return nil
}

// remove record with specified key name, if exists
func (db *Database) RemoveAtKey(key []byte) error {
	_, idx, err := db.Find(key)
	if err == nil {
		return db.RemoveAtIndex(idx)
	} else {
		return err
	}
}

// get number of records in cache
func (db *Database) GetSize() int {
	return len(db.records)
}
