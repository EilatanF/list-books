package main

import (
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

//Get id of user who authorized OAuth
//Get an xml response with the Goodreads user_id for the user who authorized access using OAuth. You'll need to register your app (required).
//URL: https://www.goodreads.com/api/auth_user
//HTTP method: GET

//Find books by title, author, or ISBN
//Get an xml response with the most popular books for the given query. This will search all books in the title/author/ISBN fields and show matches, sorted by popularity on Goodreads. There will be cases where a result is shown on the Goodreads site, but not through the API. This happens when the result is an Amazon-only edition and we have to honor Amazon's terms of service.
//URL: https://www.goodreads.com/search/index.xml    (sample url)
//HTTP method: GET
//Parameters:
//q: The query text to match against book title, author, and ISBN fields. Supports boolean operators and phrase searching.
//page: Which page to return (default 1, optional)
//key: Developer key (required).
//search[field]: Field to search, one of 'title', 'author', or 'all' (default is 'all')

func main() {
	client := &http.Client{}
	r := gin.Default()
	r.GET("/list_book", func(c *gin.Context) {
		req, err := http.NewRequest("GET", "https://www.goodreads.com/search/index.xml", nil)
		if err != nil {
			log.Print(err)
		}
		query := c.Request.URL.Query()
		query.Add("key", "RDfV4oPehM6jNhxfNQzzQ")

		req.URL.RawQuery = query.Encode()

		res, err := client.Do(req)
		if err != nil {
			log.Println(err)
			c.String(http.StatusInternalServerError, fmt.Sprintf("error making the query: %s, please retry", err))
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			c.String(http.StatusInternalServerError, fmt.Sprintf("error parsing the data: %s, please retry", err))
		}
		var parsed Results
		err = xml.Unmarshal(body, &parsed)
		if err != nil {
			log.Println(err)
			c.String(http.StatusInternalServerError, fmt.Sprintf("error parsing the data: %s, please retry", err))

		}
		c.JSON(http.StatusOK, parsed)
	})
	r.Run(":8080")
}
