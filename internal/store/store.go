// Package store abstracts interaction for implementing a data repository.
package store

// Store describes a data repository.
type Store interface {
	Create(string, []string, ...Row) error
}

// DBConfig holds the configuration used for instantiating a new database instance.
type DBConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

// Row is a store's entry.
type Row []interface{}
