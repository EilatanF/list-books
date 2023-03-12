package main

// Results struct for parsing XML
type Results struct {
	Works        []Book `xml:"search>results>work"`
	ResultsEnd   int    `xml:"search>results-end"`
	TotalResults int    `xml:"search>total-results"`
}

// Book struct for book data
type Book struct {
	Title  string `xml:"best_book>title"`
	Author string `xml:"best_book>author>name"`
	Image  string `xml:"best_book>image_url"`
}
