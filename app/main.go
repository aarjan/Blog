package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/aarjan/blog/app/api"
	h "github.com/aarjan/blog/app/handlers"
	"github.com/aarjan/blog/app/shared"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "user=postgres password=1234 dbname=blog")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	shared.InitCookieDBStore(db)

	app := &api.AppService{DB: db}
	log.Fatal(http.ListenAndServe(":3000", h.Handlers(app)))
}
