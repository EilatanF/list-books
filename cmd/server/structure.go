package main

type Results struct {
	Works []Book `xml:"search>results>work"`
}
type Book struct {
	Title  string `xml:"best_book>title"`
	Author string `xml:"best_book>author>name"`
	Image  string `xml:"best_book>image_url"`
}
