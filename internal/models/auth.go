package models

import (
	"database/sql"
	"fmt"

	"github.com/JuHaNi654/cms/internal/database"
	"github.com/JuHaNi654/cms/internal/password"
)

type Password interface {
	GetPassword() string
	GetMatchingPassword() string
}

type Login struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required"`
}

func (l Login) Authenticate(db *database.SqlClient) (match bool, id int, err error) {
	var userPassword string
	q := `SELECT id, password FROM users WHERE email = $1;`

	row := db.Db.QueryRow(q, l.Email)
	err = row.Scan(&id, &userPassword)
	if err == sql.ErrNoRows {
		return false, id, nil
	} else if err != nil {
		return false, id, err
	}

	match, err = password.Compare(l.Password, userPassword)
	if err != nil {
		return false, id, err
	}

	if match {
		return true, id, nil
	}

	return false, 0, nil
}

type Install struct {
	Firstname  string `form:"firstname" validate:"required"`
	Lastname   string `form:"firstname" validate:"required"`
	Email      string `form:"email" validate:"required,email"`
	Password   string `form:"password" validate:"required,match-passwords,gte=8,lte=32"`
	RePassword string `form:"repassword"`
}

func (i Install) GetPassword() string {
	return i.Password
}

func (i Install) GetMatchingPassword() string {
	return i.RePassword
}

func (i Install) SaveUser(db *database.SqlClient) error {
	var (
		q   string
		err error
	)

	q = `INSERT INTO users(firstname, lastname, email, password) values (?, ?, ?, ?);`

	hashPassword, err := password.Hash(i.Password)
	if err != nil {
		return err
	}

	_, err = db.Db.Exec(q, i.Firstname, i.Lastname, i.Email, hashPassword)
	if err != nil {
		return err
	}

	// TODO maybe add where clause where select specific metadata row
	// but there should be only 1 row
	q = `UPDATE metadata SET ready = true;`
	_, err = db.Db.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

func GetErrorMessage(field, tag, param string) string {
	switch tag {
	case "match-passwords":
		return "Passwords must match"
	case "email":
		return "Invalid email"
	case "required":
		return fmt.Sprintf("%s cannot be empty", field)
	case "gte":
		return fmt.Sprintf("%s length must be at least %s", field, param)
	case "lte":
		return fmt.Sprintf("%s length cannot be larger than %s", field, param)
	}

	return ""
}
