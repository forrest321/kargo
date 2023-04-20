package main

import (
	"encoding/json"

	"fmt"
	"io/ioutil"

	"net/http"
)

func getOpenBook(book *Book) {
	if book.ISBN == "" {
		return
	}

	url := fmt.Sprintf("https://openlibrary.org/isbn/%s.json", book.ISBN)
	var openBook OpenBook
	//get the data from the url
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	//read the data from the url
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	//unmarshall the data
	err = json.Unmarshal(body, &openBook)
	if err != nil {
		return
	}
	//set relevant values in book
	book.Pages = openBook.NumberOfPages

}

func formatCoverLink(coverId, size string) string {
	if coverId == "" {
		return ""
	}
	if size == "" {
		size = "M"
	}
	return fmt.Sprintf("https://covers.openlibrary.org/b/id/%s-%s.jpg", coverId, size)
}

type OpenBook struct {
	Identifiers struct {
		Goodreads    []string `json:"goodreads"`
		Librarything []string `json:"librarything"`
		Amazon       []string `json:"amazon"`
	} `json:"identifiers"`
	Title   string `json:"title"`
	Authors []struct {
		Key string `json:"key"`
	} `json:"authors"`
	PublishDate   string   `json:"publish_date"`
	Publishers    []string `json:"publishers"`
	Isbn10        []string `json:"isbn_10"`
	Covers        []int    `json:"covers"`
	Ocaid         string   `json:"ocaid"`
	Contributions []string `json:"contributions"`
	Languages     []struct {
		Key string `json:"key"`
	} `json:"languages"`
	SourceRecords []string `json:"source_records"`
	Isbn13        []string `json:"isbn_13"`
	LocalID       []string `json:"local_id"`
	Type          struct {
		Key string `json:"key"`
	} `json:"type"`
	FirstSentence struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"first_sentence"`
	Key           string `json:"key"`
	NumberOfPages int    `json:"number_of_pages"`
	Works         []struct {
		Key string `json:"key"`
	} `json:"works"`
	LatestRevision int `json:"latest_revision"`
	Revision       int `json:"revision"`
	Created        struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"created"`
	LastModified struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"last_modified"`
}
