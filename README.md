# kargo
Tech Assessment for Kargo

This is a REST API that manages a reading list.  The API uses Gin as the server and can be accessed 
with a client like Postman.

The base path is localhost:8080
The routes are as follows:

Get Books: GET /books
Optional query params of sort and direction
Sort can be one of: title, isbn, author
Direction can be one of: asc, desc
The default is title and asc

Get Book: GET /books/id
ID is the numerical ID in the SQLite DB of the requested book

Create Book: POST /books
Required values are: Title, ISBN, and Author

Update Book: PUT /books
Required values are: ID, Title, ISBN, and Author

Delete Book: DELETE /books
Required value is ID

Export to YAML: GET /books/export

Export to Pantry: GET /books/export/PantryID/BasketName

Upon create and update, an external call to the OpenLibrary API is made to retrieve the book cover and number of pages.