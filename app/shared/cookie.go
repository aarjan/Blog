package shared

import (
	"database/sql"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/gorilla/securecookie"
)

var cookieHandler *securecookie.SecureCookie

const DEFAULT_COOKIE_KEY = 1

func InitCookieDBStore(db *sql.DB) {
	key := &CookieKey{}
	if cookieHandler == nil {
		log.Info("Initializing Cookie Handler")

		if err := key.GetCookieKey(db); err != nil {
			log.Warn("cookie handler is corrupt", err)
		}

		key.ID = DEFAULT_COOKIE_KEY
		key.HashKey = securecookie.GenerateRandomKey(32)
		key.Insert(db)
	}
	cookieHandler = securecookie.New(key.HashKey, nil)

}
func SetCookie(w http.ResponseWriter, data Meta, expire int64) {
	encoded, _ := cookieHandler.Encode("session", data)

	cookie := &http.Cookie{
		Name:    "session",
		Value:   encoded,
		Expires: time.Now().Add(time.Duration(expire) * time.Millisecond),
	}
	http.SetCookie(w, cookie)
}

func GetCookie(r *http.Request) (decodedVal Meta) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return Meta{}
	}
	if err := cookieHandler.Decode("session", cookie.Value, &decodedVal); err != nil {
		log.Println("Problem getting value out of cookie")
	}
	return decodedVal
}

func DeleteCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "session",
		Value:    "nil",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	http.SetCookie(w, cookie)
}
