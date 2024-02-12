package cache

import (
	"bytes"
	"fmt"
)

func parentIdx(pos int) int {
	return pos / 2
}

func leftIdx(pos int) int {
	return pos * 2
}

func rightIdx(pos int) int {
	return pos*2 + 1
}

// return true if key1 < key2
func comparator(key1, key2 []byte) bool {
	result := bytes.Compare(key1, key2)
	return result < 0
}

func (db *Database) Swap(a int, b int) {
	db.records[a], db.records[b] = db.records[b], db.records[a]
}

func (db *Database) Percolate(pos int) {
	var cur = &db.records[pos]
	var left = &db.records[leftIdx(pos)]
	var right *Record
	if rightIdx(pos) <= db.GetSize() {
		right = &db.records[rightIdx(pos)]
	}

	if comparator(left.key, cur.key) || comparator(right.key, cur.key) {
		if comparator(left.key, right.key) {
			db.Swap(pos, leftIdx(pos))
			db.Percolate(leftIdx(pos))
		} else {
			db.Swap(pos, rightIdx(pos))
			db.Percolate(rightIdx(pos))
		}
	}
}

func (db *Database) Push(record Record) error {
	db.records = append(db.records, record)
	cur := db.GetSize()
	for comparator((&db.records[cur]).key, (&db.records[parentIdx(cur)]).key) {
		db.Swap(cur, parentIdx(cur))
		cur = parentIdx(cur)
	}

	return nil
}

// find record using binary search
func (db *Database) Find(key []byte) (record Record, index int, err error) {
	high := db.GetSize()
	low := 1
	var mid int

	for low <= high {
		mid = (high + low) / 2
		if bytes.Equal(db.records[mid].key, key) {
			return db.records[mid], mid, nil
		} else if comparator(key, db.records[mid].key) {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return record, -1, fmt.Errorf("record is not found in cache")
}

// remove record at specified index
func (db *Database) PopAtIndex(index int) error {
	dbSize := db.GetSize()
	if index > dbSize {
		return fmt.Errorf("index is out of range")
	}

	db.records[index], db.records[dbSize] = db.records[dbSize], db.records[index]
	db.records = db.records[:dbSize]
	db.Percolate(index)
	return nil
}

// remove record with specified key name, if exists
func (db *Database) PopAtKey(key []byte) error {
	_, index, err := db.Find(key)
	if err != nil {
		return err
	}

	dbSize := db.GetSize()
	db.records[index], db.records[dbSize] = db.records[dbSize], db.records[index]
	db.records = db.records[:dbSize]
	db.Percolate(index)
	return nil
}

// get number of records in cache
func (db *Database) GetSize() int {
	return len(db.records) - 1
}
