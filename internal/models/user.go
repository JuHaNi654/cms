package models

import (
	"database/sql"
	"fmt"

	"github.com/JuHaNi654/cms/internal/database"
)

type User struct {
	Id        int
	Firstname string
	Lastname  string
	Email     string
}

func GetUser(db *database.SqlClient, id int) (*User, error) {
	user := &User{}
	q := `SELECT id, firstname, lastname, email FROM users WHERE id = $1;`
	row := db.Db.QueryRow(q, id)
	err := row.Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Email)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found (id=%d)", id)
	} else if err != nil {
		return nil, err
	}

	return user, nil
}
