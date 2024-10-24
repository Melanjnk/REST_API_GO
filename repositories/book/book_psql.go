package bookRepository

import (
	"database/sql"
	"github.com/Melanjnk/REST_API_GO/m/drivers"
	"github.com/Melanjnk/REST_API_GO/m/models"
	"log"
)

type BookRepository struct{}

//func NewBookRepository() *BookRepository {
//	return
//}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (br BookRepository) GetBooks(db *sql.DB, book models.Book, books []models.Book) []models.Book {
	rows, err := db.Query("select * from books")
	drivers.LogFatal(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		logFatal(err)
		drivers.LogFatal(err)
		books = append(books, book)
	}
	return books
}

func (br BookRepository) GetBookByID(db *sql.DB, book *models.Book, id int) *models.Book {
	rows := db.QueryRow("select * from books where id = $1", id)
	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	logFatal(err)
	return book
}

func (br BookRepository) AddBook(db *sql.DB, book *models.Book) *models.Book {
	row := db.QueryRow("insert into books (title, author, year) values ($1, $2, $3) RETURNING id;", book.Title, book.Author, book.Year)
	err := row.Scan(&book.ID)
	drivers.LogFatal(err)
	return book
}

func (br BookRepository) UpdateBook(db *sql.DB, book *models.Book, id int) int64 {
	result, err := db.Exec("update books set title=$1, author=$2, year=$3 where id = $4 RETURNING id;",
		book.Title, book.Author, book.Year, book.ID)
	drivers.LogFatal(err)
	rowsUpdated, err := result.RowsAffected()
	drivers.LogFatal(err)
	return rowsUpdated
}

func (br BookRepository) RemoveBook(db *sql.DB, id int) int64 {
	result, err := db.Exec("delete from books where id = $1",
		id)
	drivers.LogFatal(err)

	rowsDeleted, err := result.RowsAffected()
	drivers.LogFatal(err)
	return rowsDeleted
}
