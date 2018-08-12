// Package store abstracts interaction for implementing a data repository.
package store

// Store describes a data repository.
type Store interface {
	Create(string, []string, ...Row) error
}

// Row is a store's entry.
type Row []interface{}
