package main

import (
	"database/sql"
	"fmt"
	"github.com/Melanjnk/REST_API_GO/m/controllers"
	"github.com/Melanjnk/REST_API_GO/m/drivers"
	"github.com/Melanjnk/REST_API_GO/m/models"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
	_ "github.com/subosito/gotenv"
	"log"
	"net/http"
)

type Book struct {
	models.Book
}

var db *sql.DB

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

	conrollers := controllers.Controller{}

	router.HandleFunc("/books", conrollers.GetBooks(db)).Methods("GET")
	router.HandleFunc("/books/{id}", conrollers.GetBook(db)).Methods("GET")
	router.HandleFunc("/books", conrollers.AddBook(db)).Methods("POST")
	router.HandleFunc("/books/{id}", conrollers.UpdateBook(db)).Methods("PUT")
	router.HandleFunc("/books/{id}", conrollers.RemoveBook(db)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))

	fmt.Println("Hello World")
}
