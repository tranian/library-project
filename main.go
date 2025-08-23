package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CheckRequest struct {
	Username string `json:"username"`
}

func main() {
	InitDB()
	defer db.Close()

	router := mux.NewRouter()

	// Routes
	router.HandleFunc("/books", GetBooksHandler).Methods("GET")
	router.HandleFunc("/books", CreateBookHandler).Methods("POST")
	router.HandleFunc("/books/{id}", GetBookHandler).Methods("GET")
	router.HandleFunc("/books/{id}/update", UpdateBookHandler).Methods("POST")
	router.HandleFunc("/books/{id}/delete", DeleteBookHandler).Methods("POST")
	router.HandleFunc("/books/{id}/checkout", CheckoutBookHandler).Methods("POST")
	router.HandleFunc("/books/{id}/return", ReturnBookHandler).Methods("POST")

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))

}

func GetBooksHandler(w http.ResponseWriter, r *http.Request) {
	books, err := GetAllBooks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(books)
}

func CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	var b Book
	json.NewDecoder(r.Body).Decode(&b)
	id, err := CreateBook(b.Title, b.Author, b.Year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b.ID = int(id)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(b)
}

func GetBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	b, err := GetBook(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(b)
}

func UpdateBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	var b Book
	json.NewDecoder(r.Body).Decode(&b)
	err := UpdateBook(id, b.Title, b.Author, b.Year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(b)
}


func DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	err := DeleteBook(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func CheckoutBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var c CheckRequest
	json.NewDecoder(r.Body).Decode(&c)
	username := c.Username
	id, _ := strconv.Atoi(vars["id"])
	err := CheckoutBook(id, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}
 
func ReturnBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var c CheckRequest
	json.NewDecoder(r.Body).Decode(&c)
	username := c.Username
	id, _ := strconv.Atoi(vars["id"])
	err := ReturnBook(id, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}
