package main

import (
	"bufio"
	"os"
)

const capacity = 10000

//Queue ... queue to store ids
type Queue struct {
	list              [capacity]string
	count, head, tail uint16
}

//Init ... initialize queue
func (q *Queue) Init() {
	q.count = 0
	q.tail = 0
	q.head = 0
}

//IsEmpty ... is empty?
func (q *Queue) IsEmpty() bool {
	if q.count == 0 {
		return true
	}

	return false
}

//IsFull ... is full?
func (q *Queue) IsFull() bool {
	if q.count == capacity {
		return true
	}

	return false
}

//Push ... Push a string
func (q *Queue) Push(id string) {
	if !q.IsFull() {
		q.list[q.head] = id
		q.head = (q.head + 1) % capacity
		q.count++
	}
}

//Add ... Add a list of strings
func (q *Queue) Add(ids []string) {
	for _, id := range ids {
		q.Push(id)
	}
}

//Pop ... Pop an id
func (q *Queue) Pop() string {
	var temp string
	if !q.IsEmpty() {
		temp = q.list[q.tail]
		q.tail = (q.tail + 1) % capacity
		q.count--
	}

	return temp
}

// Read ... Read from a file
func (q *Queue) Read(filename, def string) {
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		q.Push(def)
	} else {
		scanner := bufio.NewScanner(file)
		count := 0
		for scanner.Scan() {
			count++
			q.Push(scanner.Text())
		}

		if count == 0 {
			q.Push(def)
		}
	}
}

// Write ... Write to a file
func (q *Queue) Write(filename string) {

	if q.IsEmpty() {
		return
	}

	file, _ := os.Create(filename)
	defer file.Close()

	for !q.IsEmpty() {
		file.WriteString(q.Pop() + "\n")
	}

}
