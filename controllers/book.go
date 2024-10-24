package controllers

import (
	"database/sql"
	"encoding/json"
	"github.com/Melanjnk/REST_API_GO/m/models"
	bookRepository "github.com/Melanjnk/REST_API_GO/m/repositories/book"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Controller struct{}
type Book struct {
	models.Book
}

var books []models.Book

func (c Controller) GetBooks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		books = make([]models.Book, 0)
		book := &models.Book{}
		books := bookRepository.BookRepository{}.GetBooks(db, *book, books)
		json.NewEncoder(w).Encode(books)
	}
}

func (c Controller) GetBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		book := &models.Book{}
		params := mux.Vars(r)

		id, _ := strconv.Atoi(params["id"])

		book = bookRepository.BookRepository{}.GetBookByID(db, book, id)

		json.NewEncoder(w).Encode(&book)
	}
}

func (c Controller) AddBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		book := &models.Book{}
		_ = json.NewDecoder(r.Body).Decode(&book)
		log.Println(book)

		book = bookRepository.BookRepository{}.AddBook(db, book)

		json.NewEncoder(w).Encode(book)
	}
}

func (c Controller) UpdateBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		book := &models.Book{}
		_ = json.NewDecoder(r.Body).Decode(&book)
		log.Println(book)
		rowsUpdated := bookRepository.BookRepository{}.UpdateBook(db, book, book.ID)
		json.NewEncoder(w).Encode(rowsUpdated)
	}
}

func (c Controller) RemoveBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		id, _ := strconv.Atoi(params["id"])
		rowsDeleted := bookRepository.BookRepository{}.RemoveBook(db, id)
		json.NewEncoder(w).Encode(rowsDeleted)
	}
}
