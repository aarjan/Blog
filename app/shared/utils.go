package shared

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Wrapper struct {
	Data interface{}
	Meta
}

type Meta struct {
	Status      bool
	Code        int
	Message     string
	AccessToken []byte
}

func GenerateAccessToken(s string) []byte {
	str := s + fmt.Sprint(time.Now().UnixNano())
	token, _ := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	return token
}

type CookieKey struct {
	ID      int
	HashKey []byte
}

func (c *CookieKey) GetCookieKey(db *sql.DB) error {
	const sqlQuery = `SELECT id,hash_key FROM cookie_key WHERE id = $1`
	return db.QueryRow(sqlQuery, DEFAULT_COOKIE_KEY).Scan(&c.ID, &c.HashKey)
}

func (c *CookieKey) Insert(db *sql.DB) error {
	var err error
	const sqlQuery = `INSERT INTO public.cookie_key (` +
		`id,hash_key` +
		`)VALUES (` +
		`$1,$2)`

	_, err = db.Exec(sqlQuery, c.ID, c.HashKey)
	return err
}
