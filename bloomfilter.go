package main

import (
	"hash/fnv"
)

const bloomsize = 6235225

// Bloomfilter ... filter out weeds
type Bloomfilter struct {
	bitset [bloomsize]bool
}

func (bloom *Bloomfilter) add(id string) {
	hash := fnv.New32a()
	var loc uint32 = 0

	for i := 0; i < 5; i++ {
		hash.Write([]byte(id))
		loc = hash.Sum32() % bloomsize
		bloom.bitset[loc] = true
		id = id + string(loc)
	}
}

func (bloom *Bloomfilter) isPresent(id string) bool {
	hash := fnv.New32a()
	var loc uint32 = 0
	for i := 0; i < 5; i++ {
		hash.Write([]byte(id))
		loc = hash.Sum32() % bloomsize
		if bloom.bitset[loc] == false {
			return false
		}
		id = id + string(loc)
	}

	return true
}
