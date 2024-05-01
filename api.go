package main

import (
	"log"
	"net/http"
	"time"
)

func routes() *http.ServeMux {
	mux := http.NewServeMux()

	store := NewMemStore()
	booksHandler := NewBooksHandler(store)

	mux.Handle("/", &homeHandler{})
	mux.Handle("api/v1/books", booksHandler)
	mux.Handle("api/v1/books/", booksHandler)

	return mux
}

func runServer() {
	server := &http.Server{
		Handler:      routes(),
		Addr:         ":8080",
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}

func (b *BooksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && BookRegx.MatchString(r.URL.Path):
		b.CreateBook(w, r)
		return
	case r.Method == http.MethodGet && BookRegx.MatchString(r.URL.Path):
		b.ListBooks(w, r)
		return
	case r.Method == http.MethodGet && BookRegxWithID.MatchString(r.URL.Path):
		b.GetBook(w, r)
		return
	case r.Method == http.MethodPut && BookRegxWithID.MatchString(r.URL.Path):
		b.UpdateBook(w, r)
		return
	case r.Method == http.MethodDelete && BookRegxWithID.MatchString(r.URL.Path):
		b.DeleteBook(w, r)
		return
	default:
		NotFoundHandler(w, r)
		return
	}
}

func NewBooksHandler(s bookStorage) *BooksHandler {
	return &BooksHandler{
		store: s,
	}
}
