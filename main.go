package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

var books []Book

func main() {
	router := mux.NewRouter()

	books = append(books, Book{ID: 1, Title: "Golang pointers", Author: "Mr. Golang", Year: 2010},
		Book{ID: 2, Title: "Goroutines", Author: "Mr. Goroutines", Year: 2011},
		Book{ID: 3, Title: "Golang routers", Author: "Mr. Routers", Year: 2012},
		Book{ID: 4, Title: "Golang concurrency", Author: "Mr. Currency", Year: 2013},
		Book{ID: 5, Title: "Golang good parts", Author: "Mr. Good", Year: 2014})

	//router.MethodNotAllowedHandler =
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
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	i, _ := strconv.Atoi(params["id"])

	for _, book := range books {
		if book.ID == i {
			json.NewEncoder(w).Encode(&book)
		}
	}
}

func (b *Book) String() string {
	return fmt.Sprintf("ID: %d, Title: %s, Author: %s", b.ID, b.Title, b.Author)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	book := &Book{}

	_ = json.NewDecoder(r.Body).Decode(book)
	books := append(books, *book)

	json.NewEncoder(w).Encode(books)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	book := &Book{}

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Println(err)
	}

	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])

	for index, item := range books {
		if item.ID == id {
			books[index] = *book
		}
	}

	json.NewEncoder(w).Encode(books)
	log.Println("Update one book")
}

func removeBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])

	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
		}
	}
	json.NewEncoder(w).Encode(books)
}
