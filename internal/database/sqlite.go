package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type SqlClient struct {
	Db *sql.DB
}

func NewSQLClient(path string) (*SqlClient, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return &SqlClient{
		Db: db,
	}, nil
}

func (r *SqlClient) Migrate(file string) error {
	query, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	if _, err := r.Db.Exec(string(query)); err != nil {
		return err
	}

	return nil
}
