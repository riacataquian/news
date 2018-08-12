package store

// This file contains implementations of a Postgresql store or repository.

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// Repo is an interface to a Postgresql database.
// Repo satisfies the Store interface.
type Repo struct {
	*sqlx.DB
}

// New returns a Store which is a handler to a database pool.
//
// It opens and establishes a connection to a Postgresql database pool if none is available.
func New() Store {
	constr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"))
	db, err := sqlx.Open("postgres", constr) // Validates the connection string.
	if err != nil {
		panic(err)
	}

	// Open a connection to the database.
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return &Repo{DB: db}
}

// Create performs Postgresql's `copy` to insert the supplied rows given a list of `cols` columns.
func (repo *Repo) Create(table string, cols []string, rows ...Row) error {
	tx, err := repo.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(pq.CopyIn(table, cols...))
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, row := range rows {
		_, err = stmt.Exec(row...)
		if err != nil {
			return fmt.Errorf("persisting rows: %v", err)
		}
	}

	// Clear any buffered data; plan pq execution.
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("clearing unbuffered data: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("committing changes to the database: %v", err)
	}

	return nil
}
