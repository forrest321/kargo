package main

import (
	"fmt"
	"net/http"

	"github.com/forrest321/kargo/pantry"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	router.PUT("/books", updateBook)
	router.DELETE("/books", deleteBook)
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
	if err := c.ShouldBindQuery(&book); err != nil {
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
	var newBookDetails Book
	bookId := c.Query("id")
	if err := db.First(&book, bookId).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err := c.ShouldBindQuery(&newBookDetails); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if book.Author != newBookDetails.Author && newBookDetails.Author != "" {
		book.Author = newBookDetails.Author
	}

	if book.Title != newBookDetails.Title && newBookDetails.Title != "" {
		book.Title = newBookDetails.Title
	}

	if book.ISBN != newBookDetails.ISBN && newBookDetails.ISBN != "" {
		book.ISBN = newBookDetails.ISBN
	}

	book.CurrentPage = newBookDetails.CurrentPage
	book.Status = newBookDetails.Status

	getOpenBook(&book)
	db.Save(&book)
	c.JSON(http.StatusOK, book)
}

// Delete a book by ID
func deleteBook(c *gin.Context) {
	var book Book
	bookId := c.Query("id")
	if bookId == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := db.First(&book, bookId).Error; err != nil {
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
		return
	}
	var books []Book
	db.Find(&books)

	booksJson := formatBooksForPantry(books)
	resp, err := pantry.Export(pantryID, basketName, booksJson)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// Workaround: Pantry will not allow an array of items to be posted,
// so all books must be a part of the same object
func formatBooksForPantry(books []Book) string {
	formattedBooks := "{"
	for i, book := range books {
		c := i + 1
		formattedBooks = formattedBooks + fmt.Sprintf(`"ISBN-%v": "%s",`, c, book.ISBN)
		formattedBooks = formattedBooks + fmt.Sprintf(`"Title-%v": "%s",`, c, book.Title)
		formattedBooks = formattedBooks + fmt.Sprintf(`"Status-%v": "%s",`, c, book.Status)
		formattedBooks = formattedBooks + fmt.Sprintf(`"CurrentPage-%v": %v,`, c, book.CurrentPage)
		if book.Pages != 0 {
			formattedBooks = formattedBooks + fmt.Sprintf(`"Pages-%v": %v,`, c, book.Pages)
			if book.CurrentPage > 0 {
				pc := (float32(book.CurrentPage) / float32(book.Pages)) * 100
				percentCompleted := int(pc)
				formattedBooks = formattedBooks + fmt.Sprintf(`"Completed-%v": "%v%s",`, c, percentCompleted, "%")
			}
		}

		if book.Cover != "" {
			formattedBooks = formattedBooks + fmt.Sprintf(`"Cover-%v": "%s",`, c, book.Cover)
		}

		if i+1 == len(books) {
			formattedBooks = formattedBooks + fmt.Sprintf(`"Author-%v": "%s"`, c, book.Author)
		} else {
			formattedBooks = formattedBooks + fmt.Sprintf(`"Author-%v": "%s",`, c, book.Author)
		}

	}
	formattedBooks = formattedBooks + "}"
	return formattedBooks
}
