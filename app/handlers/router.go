package handlers

import (
	"net/http"

	"github.com/aarjan/blog/app/api"
	"github.com/gorilla/mux"
)

type cat struct {
	name string
}

func Handlers(app *api.AppService) *mux.Router {
	m := mux.NewRouter()
	m.StrictSlash(true)

	m.PathPrefix("/api/v1/website/").Handler(http.StripPrefix("/api/v1/website/", http.FileServer(http.Dir("../dist"))))

	m.Handle("/api/v1/register", AppHandler{app, Register}).Methods("POST")
	m.Handle("/api/v1/login", AppHandler{app, Login}).Methods("POST")
	m.Handle("/api/v1/logout", AppHandler{app, Logout})
	m.Handle("/api/v1/", AppHandler{app, Index}).Methods("GET")
	m.Handle("/api/v1/register", AppHandler{app, RegisterPage}).Methods("GET")

	// GET all items
	m.Handle("/api/v1/posts", AuthMiddleware(AppHandler{app, GetPosts})).Methods("GET")
	m.Handle("/api/v1/tags", AuthMiddleware(AppHandler{app, GetTags})).Methods("GET")
	m.Handle("/api/v1/categories", AuthMiddleware(AppHandler{app, GetCategories})).Methods("GET")

	// GET particular item
	m.Handle("/api/v1/post/{id:[0-9]+}", AuthMiddleware(AppHandler{app, GetPost})).Methods("GET")
	m.Handle("/api/v1/tag/{id:[0-9]+}", AuthMiddleware(AppHandler{app, GetTag})).Methods("GET")
	m.Handle("/api/v1/category/{id:[0-9]+}", AuthMiddleware(AppHandler{app, GetCategory})).Methods("GET")

	// POST
	m.Handle("/api/v1/post", AuthMiddleware(AppHandler{app, InsertPost})).Methods("POST")
	m.Handle("/api/v1/tag", AuthMiddleware(AppHandler{app, InsertTag})).Methods("POST")
	m.Handle("/api/v1/category", AuthMiddleware(AppHandler{app, InsertCategory})).Methods("POST")

	// DELETE
	m.Handle("/api/v1/post/{id:[0-9]+}", AuthMiddleware(AppHandler{app, DeletePost})).Methods("DELETE")
	m.Handle("/api/v1/tag/{id:[0-9]+}", AuthMiddleware(AppHandler{app, DeleteTag})).Methods("DELETE")
	m.Handle("/api/v1/category/{id:[0-9]+}", AuthMiddleware(AppHandler{app, DeleteCategory})).Methods("DELETE")

	// PUT
	m.Handle("/api/v1/tag/{id:[0-9]+}", AuthMiddleware(AppHandler{app, UpdateTag})).Methods("PUT")
	m.Handle("/api/v1/category/{id:[0-9]+}", AuthMiddleware(AppHandler{app, UpdateCategory})).Methods("PUT")

	// REFRESH
	m.Handle("/api/v1/refresh", http.HandlerFunc(Refresh)).Methods("GET")
	return m
}
