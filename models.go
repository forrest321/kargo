package main

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	ISBN        string `json:"isbn" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Author      string `json:"author" binding:"required"`
	Status      string `json:"status"`
	Cover       string `json:"cover"`
	Pages       int    `json:"pages"`
	CurrentPage int    `json:"current_page"`
}
