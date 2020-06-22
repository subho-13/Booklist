package main

import (
	"fmt"
	"math"
	"os"
)

const size uint16 = 40000

// Book ... Stores necessary attributes of a book
type Book struct {
	Title       string
	Authors     []string
	Description string
	Score       float64
	Tag         map[string]float64
}

func (book *Book) print() {
	fmt.Println()
	fmt.Println("Title = ", book.Title)
	fmt.Println("Authors = ", book.Authors)
	fmt.Println("Description = ", book.Description)
	fmt.Println("Score = ", book.Score)
	/*
		fmt.Printf("Tags = [ ")

		for key := range book.Tag {
			fmt.Printf("%s ", key)
		}

		fmt.Printf("] \n")
	*/
	fmt.Println()
}

func normalize(book *Book) {
	sum := 0.0

	for _, element := range book.Tag {
		sum += float64(element)
	}

	for key := range book.Tag {
		book.Tag[key] /= sum
	}
}

func calcCost(book1, book2 *Book) float64 {

	similar := 0.0
	dissimilar := 0.0
	count := 0.0

	for key, element1 := range book1.Tag {
		element2, ok := book2.Tag[key]

		if !ok {
			element2 = 0.0
		}

		dissimilar += math.Pow(element1-element2, 2)
		similar += 1.0 - dissimilar

		count = count + 1
	}

	for key, element1 := range book2.Tag {
		element2, ok := book1.Tag[key]

		if ok {
			continue
		}

		element2 = 0.0

		dissimilar += math.Pow(element1-element2, 2)
		similar += 1.0 - dissimilar

		count = count + 1
	}

	similar = similar * 1000 / count
	dissimilar = dissimilar * 1000 / count

	return similar * dissimilar / math.Pow(1.3, math.Abs(book1.Score-book2.Score))
}

type bookCell struct {
	book *Book
	next int32
	cost float64
}

// List ... Stores the booklist
type List struct {
	cell  [size]bookCell
	begin int32
	len   int32
}

// Init ... Initialize struct
func (list *List) Init() {
	list.len = 0
	list.begin = 0
}

func (list *List) insertAfter(pos int32, book *Book) {
	list.cell[list.len].book = book

	list.cell[list.len].next = list.cell[pos].next
	list.cell[pos].next = list.len

	list.cell[pos].cost = calcCost(list.cell[pos].book, book)

	next := list.cell[list.len].next

	if next != -1 {
		list.cell[list.len].cost = calcCost(list.cell[next].book, book)
	}

	list.len++
}

func (list *List) insertHead(book *Book) {
	curr := list.len

	list.cell[curr].book = book
	list.cell[curr].next = -1

	if list.len != 0 {
		list.cell[curr].cost = calcCost(list.cell[list.begin].book, book)
		list.cell[curr].next = list.begin
		list.begin = curr
	}

	list.len++
}

func (list *List) calcIncrease(pos int32, book *Book) float64 {
	if pos == -1 {
		return calcCost(list.cell[list.begin].book, book)
	} else if list.cell[pos].next == -1 {
		return calcCost(list.cell[pos].book, book)
	} else {
		cost1n := calcCost(list.cell[pos].book, book)
		next := list.cell[pos].next
		cost2n := calcCost(list.cell[next].book, book)

		return cost1n + cost2n - list.cell[pos].cost
	}
}

// Append ... Add a new book
func (list *List) Append(book *Book) {
	normalize(book)

	if list.len == 0 {
		list.insertHead(book)
	} else if list.len == 1 {
		list.insertAfter(0, book)
	} else {
		var pos int32
		var maxIncrease float64 = 0

		var i int32

		for i = 0; i < list.len; i++ {
			if maxIncrease < list.calcIncrease(i, book) {
				maxIncrease = list.calcIncrease(i, book)
				pos = i
			}
		}

		list.insertAfter(pos, book)
	}
}

// Print ... Print the list
func (list *List) Print() {
	i := list.begin

	fmt.Println()
	fmt.Println("=====\t=====\t=====\t=====\t=====\t=====\t=====\t")
	fmt.Println("!\t\tList " + fmt.Sprintf("%d", list.len) + "- Bitches\t\t!")

	if list.len == 0 {
		return
	}

	for i != -1 {
		list.cell[i].book.print()
		fmt.Println("Cost = ", list.cell[i].cost)
		i = list.cell[i].next
	}

	fmt.Println("!\t\tEnjoy - " + fmt.Sprintf("%d", list.len) + " - Suckers\t\t!")
	fmt.Println("=====\t=====\t=====\t=====\t=====\t=====\t=====\t")
	fmt.Println()
}

// Write ... Store the list
func (list *List) Write(filename string) {
	f, _ := os.Create(filename)
	defer f.Close()

	s := "List of " + fmt.Sprintf("%d", list.len) + " Books\n"
	f.WriteString(s)

	if list.len == 0 {
		return
	}

	i := list.begin
	count := 1

	for i != -1 {
		book := list.cell[i].book
		s = ""
		s += "\n\n" + fmt.Sprintf("%d", count) + ". Title = " + book.Title
		s += "\nAuthor = [ "
		for _, author := range book.Authors {
			s += "\"" + author + "\"" + " "
		}

		s += "]\nScore = " + fmt.Sprintf("%f", book.Score)
		s += "\nDescription :: \n" + book.Description

		/*
			s += "\nTags = [ "

			for key := range book.Tag {
				s += key + " "
			}
		*/

		s += "\n\nCost = " + fmt.Sprintf("%f", list.cell[i].cost)

		i = list.cell[i].next
		count = count + 1

		f.WriteString(s)
	}
}
