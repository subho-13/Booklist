package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const queueFile string = "/home/subho/queueIds.txt"
const bookFile string = "/home/subho/Ultimate_Reading_List.txt"

func fetchBooks(queue *Queue, bookChan chan<- *Book, sigChan <-chan os.Signal) {
	reqString1 := "https://www.goodreads.com/book/show/"
	reqString2 := ".xml?key=0SZLTTz8zQXu2a5qGnJAA"
	var stop bool = false

	for range time.Tick(time.Second) {
		if !queue.IsEmpty() && !stop {
			id := queue.Pop()
			url := reqString1 + id + reqString2

			resp, err := http.Get(url)

			if err != nil {
				panic("Couldn't fetch url")
			}

			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)
			var goodreads GoodreadsResponse
			goodreads.Retrieve(body)

			okay, book := goodreads.ConstructBook()

			if okay {
				bookChan <- book
				fmt.Printf("c(o-o) =>")
			} else {
				fmt.Printf("\b=>")
			}

			queue.Add(goodreads.NextIds())
		} else {
			close(bookChan)
			return
		}

		select {
		case <-sigChan:
			stop = true
		default:
			continue
		}
	}
}

func storeBooks(bookChan <-chan *Book, bookList *List, done chan<- bool) {
	for books := range bookChan {
		bookList.Append(books)
	}

	done <- true
}

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	bookChan := make(chan *Book)
	queue := new(Queue)

	var id string
	fmt.Print("Enter your Id here :: ")
	fmt.Scanln(&id)

	queue.Init()
	queue.Push(id)

	go fetchBooks(queue, bookChan, sigChan)

	bookList := new(List)

	done := make(chan bool)

	go storeBooks(bookChan, bookList, done)

	<-done

	fmt.Println("\n\n -- Recieved Signal. Writing Files. -- ")
	bookList.Write(bookFile)
	bookList.Print()
}
