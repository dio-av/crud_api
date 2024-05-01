package main

import (
	"encoding/json"
	"net/http"
	"regexp"
)

var (
	BookRegx       = regexp.MustCompile(`^/recipes/*$`)
	BookRegxWithID = regexp.MustCompile(`^/recipes/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
)

func (b *BooksHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book Book

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	resourceID := book.Name

	if err := b.store.Create(resourceID, book); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (b *BooksHandler) ListBooks(w http.ResponseWriter, r *http.Request) {
	booksList, err := b.store.List()
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	response, err := json.Marshal(booksList)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (b *BooksHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	matches := BookRegxWithID.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	recipe, err := b.store.Get(matches[1])
	if err != nil {
		if err == ErrNotFound {
			NotFoundHandler(w, r)
			return
		}

		InternalServerErrorHandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(recipe)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (b *BooksHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	matches := BookRegxWithID.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	if err := b.store.Update(matches[1], book); err != nil {
		if err == ErrNotFound {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (b *BooksHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	matches := BookRegxWithID.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	if err := b.store.Delete(matches[1]); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}
