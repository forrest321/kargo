package main

import (
	"encoding/json"
	"github.com/forrest321/kargo/pantry"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

var db *gorm.DB

func main() {
	setupDb()

	// Set up Gin router
	router := gin.Default()

	// Define routes
	router.GET("/books", getBooks)
	router.GET("/books/:id", getBook)
	router.POST("/books", createBook)
	router.PUT("/books/:id", updateBook)
	router.DELETE("/books/:id", deleteBook)
	router.GET("/books/export", exportToYAML)
	router.GET("/books/export/:pantryID/:basketName", exportToPantry)

	// Start server
	router.Run(":8080")
}

// Set up SQLite database
func setupDb() {
	var err error
	db, err = gorm.Open(sqlite.Open("readinglist.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Book{})
}

// Get all books
func getBooks(c *gin.Context) {
	var books []Book
	db.Find(&books)
	c.JSON(http.StatusOK, books)
}

// Get a single book by ID
func getBook(c *gin.Context) {
	var book Book
	if err := db.First(&book, c.Param("id")).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, book)
}

// Create a new book
func createBook(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// get openBook data
	getOpenBook(&book)

	db.Create(&book)
	c.JSON(http.StatusCreated, book)
}

// Update an existing book by ID
func updateBook(c *gin.Context) {
	var book Book
	if err := db.First(&book, c.Param("id")).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err := c.ShouldBindJSON(&book); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	db.Save(&book)
	c.JSON(http.StatusOK, book)
}

// Delete a book by ID
func deleteBook(c *gin.Context) {
	var book Book
	if err := db.First(&book, c.Param("id")).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	db.Delete(&book)
	c.Status(http.StatusNoContent)
}

func exportToYAML(c *gin.Context) {
	var books []Book
	db.Find(&books)
	c.YAML(http.StatusOK, books)
}

func exportToPantry(c *gin.Context) {
	pantryID := c.Param("pantryID")
	basketName := c.Param("basketName")
	if pantryID == "" || basketName == "" {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	var books []Book
	db.Find(&books)
	var booksJson string
	bookBytes, err := json.Marshal(books)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	booksJson = string(bookBytes)
	resp, err := pantry.Export(pantryID, basketName, booksJson)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, resp)
}
