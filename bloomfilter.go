package main

import (
	"fmt"
	"hash/fnv"
)

const bloomsize = 7235225

// Bloomfilter ... filter out weeds
type Bloomfilter struct {
	bitset [bloomsize]bool
}

// Add ... Add an id
func (bloom *Bloomfilter) Add(id string) {
	hash := fnv.New32a()
	var loc uint32 = 0

	for i := 0; i < 5; i++ {
		hash.Write([]byte(id))
		loc = hash.Sum32() % bloomsize
		bloom.bitset[loc] = true
		id = id + fmt.Sprintf("%d", loc)
	}
}

// IsPresent ... Check if present
func (bloom *Bloomfilter) IsPresent(id string) bool {
	hash := fnv.New32a()
	var loc uint32 = 0
	for i := 0; i < 5; i++ {
		hash.Write([]byte(id))
		loc = hash.Sum32() % bloomsize
		if bloom.bitset[loc] == false {
			return false
		}
		id = id + fmt.Sprintf("%d", loc)
	}

	return true
}
