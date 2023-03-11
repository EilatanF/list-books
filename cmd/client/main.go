package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	var terms string
	var field string
	var page string
	var hostname string

	flag.StringVar(&terms, "s", "", "Keyword to search in Goodreads (shorthand)")
	flag.StringVar(&terms, "search", "", "Keyword to search in Goodreads")
	flag.StringVar(&field, "sort", "title", "Sort the results by the specified field")
	flag.StringVar(&hostname, "h", "127.0.0.1", "the hostname or ip address where the server(shorthand)")
	flag.StringVar(&hostname, "host", "127.0.0.1", "the hostname or ip address where the server")
	flag.StringVar(&page, "p", "1", "the page number of results to show")

	flag.Parse()
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:8080/list_book", hostname), nil)
	if err != nil {
		log.Print(err)
	}

	if terms == "" {
		log.Print(fmt.Errorf("keyword should not be blank"))
	}

	query := req.URL.Query()
	query.Add("q", terms)
	query.Add("page", page)
	query.Add("field", field)
	req.URL.RawQuery = query.Encode()

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, body, "", "    ")
	fmt.Println(prettyJSON.String())
}
