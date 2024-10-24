package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Melanjnk/REST_API_GO/m/drivers"
	"github.com/Melanjnk/REST_API_GO/m/models"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
	_ "github.com/subosito/gotenv"
	"log"
	"net/http"
	"strconv"
)

type Book struct {
	models.Book
}

var db *sql.DB
var books []Book

func init() {
	gotenv.Load()
}

func main() {
	db = drivers.ConnectDB()
	router := mux.NewRouter()

	// Handle all preflight request
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fmt.Printf("OPTIONS")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		w.WriteHeader(http.StatusNoContent)
		return
	})

	router.StrictSlash(true)

	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", removeBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))

	fmt.Println("Hello World")
}

func getBooks(w http.ResponseWriter, r *http.Request) {
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

func getBook(w http.ResponseWriter, r *http.Request) {
	book := &Book{}
	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])

	rows := db.QueryRow("select * from books where id = $1", id)
	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	drivers.LogFatal(err)

	json.NewEncoder(w).Encode(&book)
}

func (b *Book) String() string {
	return fmt.Sprintf("ID: %d, Title: %s, Author: %s", b.ID, b.Title, b.Author)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	book := &Book{}
	_ = json.NewDecoder(r.Body).Decode(&book)
	log.Println(book)

	row := db.QueryRow("insert into books (title, author, year) values ($1, $2, $3) RETURNING id;", book.Title, book.Author, book.Year)
	err := row.Scan(&book.ID)
	drivers.LogFatal(err)
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
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

func removeBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])

	result, err := db.Exec("delete from books where id = $1",
		id)
	drivers.LogFatal(err)

	rowsDeleted, err := result.RowsAffected()
	drivers.LogFatal(err)

	json.NewEncoder(w).Encode(rowsDeleted)
}
