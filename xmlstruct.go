package main

import (
	"encoding/xml"
	"math"
)

const base = 1.3
const textRev = 100.0

var bloomfilter Bloomfilter

// GoodreadsResponse ... Contains extracted fields
type GoodreadsResponse struct {
	XMLName xml.Name `xml:"GoodreadsResponse"`
	Book    book     `xml:"book"`
}

type book struct {
	Title       string   `xml:"title"`
	ID          string   `xml:"id"`
	Description string   `xml:"description"`
	Work        work     `xml:"work"`
	Authors     []string `xml:"authors>author>name"`
	NumPages    uint32   `xml:"num_pages"`
	Tags        []tag    `xml:"popular_shelves>shelf"`
	Similar     []string `xml:"similar_books>book>id"`
}

type work struct {
	RatingsSum   uint32 `xml:"ratings_sum"`
	RatingsCount uint32 `xml:"ratings_count"`
	TextRevCount uint32 `xml:"text_reviews_count"`
}

type tag struct {
	Name  string  `xml:"name,attr"`
	Count float64 `xml:"count,attr"`
}

// Retrieve ... Extracts important xml fields
func (g *GoodreadsResponse) Retrieve(xmlData []byte) {
	xml.Unmarshal(xmlData, g)
}

func filterOut(avg float64, rCount, tRev uint32) bool {
	if avg < 4 || rCount < 50000 || tRev < 2000 {
		return true
	}

	return false
}

// CalcScore ... Calculates Aggregate Scores
func (g *GoodreadsResponse) calcScore() float64 {
	rCount := g.Book.Work.RatingsCount
	rSum := g.Book.Work.RatingsSum
	tRev := g.Book.Work.TextRevCount

	var avg = float64(rSum) / float64(rCount)

	if filterOut(avg, rCount, tRev) {
		return 0
	}

	score := math.Pow(base, avg) * math.Log2(float64(rCount*tRev))
	score = score / math.Sqrt(float64(g.Book.NumPages))

	return score
}

// ConstructBook ... Returns a book object
func (g *GoodreadsResponse) ConstructBook() (bool, *Book) {
	book := new(Book)

	book.Title = g.Book.Title

	book.Authors = g.Book.Authors

	book.Description = g.Book.Description
	book.Score = g.calcScore()

	if book.Score == 0 {
		return false, nil
	}

	book.Tag = make(map[string]float64)

	for _, tag := range g.Book.Tags {
		book.Tag[tag.Name] = tag.Count
	}

	return true, book
}

// NextIds ... Returns next Ids
func (g *GoodreadsResponse) NextIds() []string {

	idList := make([]string, 0, 15)

	for _, id := range g.Book.Similar {
		if bloomfilter.IsPresent(id) == false {
			idList = append(idList, id)
			bloomfilter.Add(id)
		}
	}

	return idList
}
