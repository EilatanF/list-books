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

	//Parse the flags, providing both shorthand and long flags
	flag.StringVar(&terms, "s", "", "Keyword to search in Goodreads (shorthand)")
	flag.StringVar(&terms, "search", "", "Keyword to search in Goodreads")
	flag.StringVar(&field, "sort", "title", "Sort the results by the specified field")
	flag.StringVar(&hostname, "h", "127.0.0.1", "the hostname or ip address where the server(shorthand)")
	flag.StringVar(&hostname, "host", "127.0.0.1", "the hostname or ip address where the server")
	flag.StringVar(&page, "p", "", "the page number of results to show. Default to show all results")

	flag.Parse()
	client := &http.Client{}

	//generate the request to the server (I would use https instead of http in real life)
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:8080/list_book", hostname), nil)
	if err != nil {
		log.Print(err)
	}

	//check for required field
	if terms == "" {
		log.Print(fmt.Errorf("keyword must not be blank, please try again"))
		return
	}

	query := req.URL.Query()
	query.Add("q", terms)
	if page != "" {
		query.Add("page", page)
	}
	query.Add("field", field)
	req.URL.RawQuery = query.Encode()

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	//parse the body and return it
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}

	//ask to retry if status code is not OK
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		log.Println(fmt.Errorf("error making the request: %s, please try again", body))
		return
	}

	//I would store the data somewhere like S3 or local file instead of printing it out in a normal case
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, body, "", "    ")
	fmt.Println(prettyJSON.String())
}
