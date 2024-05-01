package main

import "errors"

type bookStorage interface {
	Create(name string, book Book) error
	Get(name string) (Book, error)
	List() (map[string]Book, error)
	Update(name string, book Book) error
	Delete(name string) error
}

type BooksHandler struct {
	store bookStorage
}

var (
	ErrNotFound = errors.New("not found")
)

type MemStore struct {
	list map[string]Book
}

func NewMemStore() *MemStore {
	list := make(map[string]Book)
	return &MemStore{
		list,
	}
}

func (m MemStore) Create(name string, book Book) error {
	m.list[name] = book
	return nil
}

func (m MemStore) Get(name string) (Book, error) {
	if val, ok := m.list[name]; ok {
		return val, nil
	}

	return Book{}, ErrNotFound
}

func (m MemStore) List() (map[string]Book, error) {
	return m.list, nil
}

func (m MemStore) Update(name string, book Book) error {

	if _, ok := m.list[name]; ok {
		m.list[name] = book
		return nil
	}

	return ErrNotFound
}

func (m MemStore) Delete(name string) error {
	delete(m.list, name)
	return nil
}
