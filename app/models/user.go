package models

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                 int    `json:"id"`
	Username           string `json:"username"`
	Password           []byte `json:"password"`
	AccessToken        []byte `json:"access_token"`
	VerificationToken  string `json:"verification_token,omitempty"`
	VerificationStatus bool   `json:"verification_status,omitempty"`
}

func (u *User) CreateUser(db *sql.DB) error {
	const sqlQuery = `INSERT into users( ` +
		`username,password,access_token,verification_token,verification_status ` +
		`)VALUES( ` +
		`$1,$2,$3,$4,$5) RETURNING id;`
	return db.QueryRow(sqlQuery, u.Username, u.Password, u.AccessToken, u.VerificationToken, u.VerificationStatus).Scan(&u.ID)
}

func GetUsers(db *sql.DB) ([]User, error) {
	const sqlQuery = `SELECT * FROM users`
	query, err := db.Query(sqlQuery)
	defer query.Close()
	if err != nil {
		return nil, err
	}
	res := []User{}
	u := User{}
	for query.Next() {
		query.Scan(&u.ID, &u.Username, &u.Password, &u.AccessToken, &u.VerificationToken, &u.VerificationStatus)
		res = append(res, u)
	}
	return res, nil
}

func (u *User) UserByName(db *sql.DB, name string) error {
	const sqlQuery = `SELECT ` +
		`id,username,password,access_token,verification_token,verification_status ` +
		`FROM public.users ` +
		`WHERE username = $1`
	return db.QueryRow(sqlQuery, name).Scan(&u.ID, &u.Username, &u.Password, &u.AccessToken, &u.VerificationToken, &u.VerificationStatus)
}

func (user *User) VerifyUser(name, pass string, db *sql.DB) error {
	err := user.UserByName(db, name)
	if err != nil {
		return errors.New(" Invalid username")
	}

	if err = bcrypt.CompareHashAndPassword(user.Password, []byte(pass)); err != nil {
		return errors.New(" Invalid Password")
	}
	return nil
}
