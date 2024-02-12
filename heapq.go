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

// return true is key1 < key2
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

func (db *Database) Pop(key []byte) (record Record, err error) {
	record, index, err := db.Find(key)
	if err != nil {
		return record, err
	}

	db.records[index] = db.records[db.GetSize()]
	db.records = db.records[:db.GetSize()]
	db.Percolate(index)
	return record, nil
}

func (db *Database) GetSize() int {
	return len(db.records) - 1
}
