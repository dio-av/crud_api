package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func readTestData(t *testing.T, name string) []byte {
	t.Helper()
	content, err := os.ReadFile("./testdata/" + name)
	if err != nil {
		t.Errorf("Could not read %v", name)
	}

	return content
}

func TestRecipesHandlerCRUD_Integration(t *testing.T) {

	store := NewMemStore()
	recipesHandler := NewBooksHandler(store)

	postumas := readTestData(t, "memorias_postumas.json")
	postumasReader := bytes.NewReader(postumas)

	_1984 := readTestData(t, "1984.json")
	_1984Reader := bytes.NewReader(_1984)

	// CREATE - add a new book
	req := httptest.NewRequest(http.MethodPost, "/books", postumasReader)
	w := httptest.NewRecorder()
	recipesHandler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, 200, res.StatusCode)

	saved, _ := store.List()
	assert.Len(t, saved, 1)

	// GET - find the record we just added
	req = httptest.NewRequest(http.MethodGet, "/recipes/memorias-postumas-de-bras-cubas", nil)
	w = httptest.NewRecorder()
	recipesHandler.ServeHTTP(w, req)

	res = w.Result()
	defer res.Body.Close()
	assert.Equal(t, 200, res.StatusCode)

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	assert.JSONEq(t, string(postumas), string(data))

	// UPDATE - bras cubas -> 1948
	req = httptest.NewRequest(http.MethodPut, "/books/memorias-postumas-de-bras-cubas", _1984Reader)
	w = httptest.NewRecorder()
	recipesHandler.ServeHTTP(w, req)

	res = w.Result()
	defer res.Body.Close()
	assert.Equal(t, 200, res.StatusCode)

	updatedBrasCubas, err := store.Get("memorias-postumas-de-bras-cubas")
	assert.NoError(t, err)

	assert.Contains(t, updatedBrasCubas.Name, "1984")

	//DELETE - remove the ham and cheese recipe
	req = httptest.NewRequest(http.MethodDelete, "/books/memorias-postumas-de-bras-cubas", nil)
	w = httptest.NewRecorder()
	recipesHandler.ServeHTTP(w, req)

	res = w.Result()
	defer res.Body.Close()
	assert.Equal(t, 200, res.StatusCode)

	saved, _ = store.List()
	assert.Len(t, saved, 0)

}
