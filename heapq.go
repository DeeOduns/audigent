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

// return true is a < b
func comparator(a, b []byte) bool {
	result := bytes.Compare(a, b)
	return result < 0
}

func (db *Database) Swap(a int, b int) {
	db.records[a], db.records[b] = db.records[b], db.records[a]
}

func (db *Database) Percolate(pos int) {
	if leftIdx(pos) > db.size {
		return
	}

	var cur = &db.records[pos]
	var left = &db.records[leftIdx(pos)]
	var right *Record
	if rightIdx(pos) <= db.size {
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
	if db.size >= len(db.records) {
		return fmt.Errorf("no space, cache is full")
	}
	db.size++
	cur := db.size

	db.records[cur] = record
	for comparator((&db.records[cur]).key, (&db.records[parentIdx(cur)]).key) {
		db.Swap(cur, parentIdx(cur))
		cur = parentIdx(cur)
	}

	return nil
}

// find record using binary search
func (db *Database) Find(key []byte) (record Record, index int, err error) {
	high := db.size
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

func (db *Database) Pop(key []byte) (record Record, err error) {
	if db.size < 1 {
		return record, fmt.Errorf("cache is empty")
	}

	record, index, err := db.Find(key)
	if err != nil {
		return record, err
	}

	db.records[index] = db.records[db.size]
	db.size--
	db.Percolate(index)
	return record, nil
}
