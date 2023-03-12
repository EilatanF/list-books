package main

import (
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strconv"
)

func main() {
	client := &http.Client{}
	r := gin.Default()
	r.GET("/list_book", func(c *gin.Context) {
		req, err := http.NewRequest("GET", "https://www.goodreads.com/search/index.xml", nil)
		if err != nil {
			log.Print(err)
			c.String(http.StatusInternalServerError, fmt.Sprintf("error making query: %s ", err))
		}

		// add auth key (I would ask the client to send over the auth key in real life,
		//but to keep things simple, I let the server handle the key)
		query := c.Request.URL.Query()
		query.Add("key", "*********")

		req.URL.RawQuery = query.Encode()
		var books []Book

		//get all pages if not using pagination
		data := getPaginatedResult(c, client, req)
		books = append(books, data.Works...)
		for data.ResultsEnd < data.TotalResults {
			page, _ := strconv.Atoi(query.Get("page"))
			query.Set("page", strconv.Itoa(page+1))
			req.URL.RawQuery = query.Encode()

			data = getPaginatedResult(c, client, req)
			books = append(books, data.Works...)
		}
		c.JSON(http.StatusOK, books)
	})
	r.Run(":8080")
}

// get book data per page
func getPaginatedResult(c *gin.Context, client *http.Client, req *http.Request) *Results {
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("error calling query: %s ", err))
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("error parsing data: %s ", err))
	}
	var parsed Results
	err = xml.Unmarshal(body, &parsed)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("error unmarshaling xml: %s ", err))

	}
	return &parsed
}
