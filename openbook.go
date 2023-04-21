package main

import (
	"encoding/json"

	"fmt"
	"io"

	"net/http"
)

func getOpenBook(book *Book) {
	if book.ISBN == "" {
		return
	}

	url := fmt.Sprintf("https://openlibrary.org/isbn/%s.json", book.ISBN)
	fmt.Println(url)
	var openBook OpenBook
	//get the data from the url
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	//read the data from the url
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	//unmarshall the data
	err = json.Unmarshal(body, &openBook)
	if err != nil {
		fmt.Println(err)
		return
	}
	//set relevant values in book
	book.Pages = openBook.NumberOfPages
	book.Cover = formatCoverLink(book.ISBN, "M")

}

func formatCoverLink(isbn, size string) string {
	if isbn == "" {
		return ""
	}
	if size == "" {
		size = "M"
	}
	return fmt.Sprintf("https://covers.openlibrary.org/b/isbn/%s-%s.jpg", isbn, size)
}

type OpenBook struct {
	Title   string `json:"title"`
	Authors []struct {
		Key string `json:"key"`
	} `json:"authors"`
	Isbn13        []string `json:"isbn_13"`
	FirstSentence string   `json:"first_sentence,omitempty"`
	NumberOfPages int      `json:"number_of_pages"`
}
