package controllers

import (
	"database/sql"
	"encoding/json"
	"github.com/Melanjnk/REST_API_GO/m/drivers"
	"github.com/Melanjnk/REST_API_GO/m/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Controller struct{}
type Book struct {
	models.Book
}

var books []Book

func (c Controller) GetBooks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		books = make([]Book, 0)
		book := &Book{}
		rows, err := db.Query("select * from books")
		drivers.LogFatal(err)

		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
			drivers.LogFatal(err)
			books = append(books, *book)
		}
		json.NewEncoder(w).Encode(books)
	}
}

func (c Controller) GetBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		book := &Book{}
		params := mux.Vars(r)

		id, _ := strconv.Atoi(params["id"])

		rows := db.QueryRow("select * from books where id = $1", id)
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		drivers.LogFatal(err)

		json.NewEncoder(w).Encode(&book)
	}
}

func (c Controller) AddBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		book := &Book{}
		_ = json.NewDecoder(r.Body).Decode(&book)
		log.Println(book)

		row := db.QueryRow("insert into books (title, author, year) values ($1, $2, $3) RETURNING id;", book.Title, book.Author, book.Year)
		err := row.Scan(&book.ID)
		drivers.LogFatal(err)
		json.NewEncoder(w).Encode(book)
	}
}

func (c Controller) UpdateBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		book := &Book{}
		_ = json.NewDecoder(r.Body).Decode(&book)
		log.Println(book)

		result, err := db.Exec("update books set title=$1, author=$2, year=$3 where id = $4 RETURNING id;",
			book.Title, book.Author, book.Year, book.ID)
		drivers.LogFatal(err)
		rowsUpdated, err := result.RowsAffected()
		drivers.LogFatal(err)
		json.NewEncoder(w).Encode(rowsUpdated)
	}
}

func (c Controller) RemoveBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		id, _ := strconv.Atoi(params["id"])

		result, err := db.Exec("delete from books where id = $1",
			id)
		drivers.LogFatal(err)

		rowsDeleted, err := result.RowsAffected()
		drivers.LogFatal(err)

		json.NewEncoder(w).Encode(rowsDeleted)
	}
}
