package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/aarjan/blog/app/api"

	"strconv"

	"github.com/aarjan/blog/app/models"
	"github.com/aarjan/blog/app/renderer"
	"github.com/aarjan/blog/app/shared"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// RegisterPage GET Requests
func RegisterPage(w http.ResponseWriter, r *http.Request, app *api.AppService) error {
	cookie := shared.GetCookie(r)
	return renderer.RenderTemplate(w, "register.tmpl", cookie)
}

// Register POST
func Register(w http.ResponseWriter, r *http.Request, app *api.AppService) error {
	var user models.User
	var data shared.Meta

	r.ParseForm()
	name := r.FormValue("username")
	pass := r.FormValue("password")
	passAgain := r.FormValue("password_again")

	if name == "" || pass == "" {
		data.Code = 401
		data.Message = "Please enter all the fields"
		shared.SetCookie(w, data, 1000)
		http.Redirect(w, r, "/api/v1/register", http.StatusFound)
		return errors.New(data.Message)
	}
	if pass != passAgain {
		data.Code = 401
		data.Message = "Password do not match"
		shared.SetCookie(w, data, 1000)
		http.Redirect(w, r, "/api/v1/register", http.StatusFound)
		return errors.New(data.Message)
	}
	// verify username
	user.UserByName(app.DB, name)
	if user.ID != 0 {
		data.Code = 401
		data.Message = "Username already exists"
		shared.SetCookie(w, data, 1000)
		http.Redirect(w, r, "/api/v1/register", http.StatusFound)
		return errors.New(data.Message)
	}

	accessToken := shared.GenerateAccessToken(name)
	password, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	user = models.User{
		Username:    name,
		Password:    password,
		AccessToken: accessToken,
	}

	// save record to database
	err := user.CreateUser(app.DB)
	if err != nil {
		data.Code = 505
		data.Message = err.Error()
		shared.SetCookie(w, data, 1000)
		http.Redirect(w, r, "/api/v1/register", http.StatusTemporaryRedirect)
		return err
	}

	// Set cookie
	data = shared.Meta{true, 302, "Registration Succesful", user.AccessToken}
	shared.SetCookie(w, data, 60000)

	http.Redirect(w, r, "/api/v1/posts", http.StatusFound)
	return err
}

// Index GET
func Index(w http.ResponseWriter, r *http.Request, app *api.AppService) error {
	cookie := shared.GetCookie(r)
	return renderer.RenderTemplate(w, "index.tmpl", cookie)
}

// Login POST
func Login(w http.ResponseWriter, r *http.Request, app *api.AppService) error {
	var data shared.Meta

	r.ParseForm()
	name := r.FormValue("username")
	pass := r.FormValue("password")

	if name == "" || pass == "" {
		data.Code = 401
		data.Message = "Please enter all the fields"
		shared.SetCookie(w, data, 1000)
		http.Redirect(w, r, "/api/v1", http.StatusFound)
		return errors.New(data.Message)
	}

	user := models.User{}

	// user verification
	err := user.VerifyUser(name, pass, app.DB)

	// return error message and redirect
	if err != nil {
		data.Code = 401
		data.Message = err.Error()
		shared.SetCookie(w, data, 1000)
		http.Redirect(w, r, "/api/v1", http.StatusFound)
		return err
	}

	cookie := shared.GetCookie(r)

	// Set cookie and redirect
	if cookie.AccessToken == nil {
		data.AccessToken = user.AccessToken
		data.Code = 302
		data.Status = true
		data.Message = "Login Successful"
		shared.SetCookie(w, data, 60000)
	}
	http.Redirect(w, r, "/api/v1/posts", http.StatusFound)
	return nil
}

// Logout GET
func Logout(w http.ResponseWriter, r *http.Request, app *api.AppService) error {
	// Should have use shared.DeleteCookie(r); not working
	shared.SetCookie(w, shared.Meta{}, 0)

	http.Redirect(w, r, "/api/v1", http.StatusFound)
	return nil
}

// GET
func GetPosts(w http.ResponseWriter, r *http.Request, app *api.AppService) error {

	posts, err := models.GetPosts(app.DB)
	if err != nil {
		return err
	}
	tags, err := models.GetTags(app.DB)
	if err != nil {
		return err
	}

	categories, err := models.GetCategories(app.DB)
	if err != nil {
		return err
	}

	data := struct {
		Posts      []models.Post
		Tags       []models.Tag
		Categories []models.Category
	}{posts, tags, categories}

	cookie := shared.GetCookie(r)
	wrapper := shared.Wrapper{data, cookie}
	return renderer.RenderTemplate(w, "home.tmpl", wrapper)
}

// GetPost retrieves a particluar post
func GetPost(w http.ResponseWriter, r *http.Request, app *api.AppService) error {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	post := &models.Post{ID: id}

	err := post.GetPost(app.DB)
	if err != nil {
		return err
	}
	cookie := shared.GetCookie(r)
	data := shared.Wrapper{*post, cookie}
	return renderer.RenderTemplate(w, "post.tmpl", data)
}

// InsertPost inserts a record in the posts table
func InsertPost(w http.ResponseWriter, r *http.Request, app *api.AppService) error {
	r.ParseForm()
	name := r.FormValue("name")
	content := r.FormValue("content")
	// categoryID, _ := strconv.Atoi(r.FormValue("category_id"))

	post := &models.Post{
		Name:       name,
		Content:    content,
		CategoryID: 1,
	}

	// Check if post exists
	post.GetPostByName(app.DB)
	if post.ID != 0 {
		// Exists
		return errors.New("Record already exists")
	}

	// Create new post
	err := post.CreatePost(app.DB)
	return err
}

func DeletePost(w http.ResponseWriter, r *http.Request, app *api.AppService) error {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	post := models.Post{ID: id}

	// Check if record exists
	post.GetPostByID(app.DB)
	if post.ID == 0 {
		return errors.New("Record not found")
	}

	// Delete record
	return post.Delete(app.DB)
}

func GetTags(w http.ResponseWriter, r *http.Request, app *api.AppService) error {

	tags, err := models.GetTags(app.DB)
	if err != nil {
		return err
	}
	cookie := shared.GetCookie(r)
	data := shared.Wrapper{tags, cookie}
	return renderer.RenderTemplate(w, "tags.tmpl", data)
}

func GetTag(w http.ResponseWriter, r *http.Request, app *api.AppService) error {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	tag := &models.Tag{ID: id}
	err := tag.GetTag(app.DB)
	if err != nil {
		return err
	}
	cookie := shared.GetCookie(r)
	data := shared.Wrapper{tag, cookie}
	return renderer.RenderTemplate(w, "tag.tmpl", data)
}

// InsertTag inserts a record in the posts table
func InsertTag(w http.ResponseWriter, r *http.Request, app *api.AppService) error {
	r.ParseForm()
	name := r.FormValue("name")

	tag := &models.Tag{
		Name: name,
	}

	// Check if post exists
	tag.GetTagByName(app.DB)
	if tag.ID != 0 {
		// Exists
		return errors.New("Record already exists")
	}

	// Create new tag
	err := tag.CreateTag(app.DB)
	return err
}

func DeleteTag(w http.ResponseWriter, r *http.Request, app *api.AppService) error {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	tag := models.Tag{ID: id}

	// Check if record exists
	tag.GetTagByID(app.DB)
	if tag.ID == 0 {
		return errors.New("Record not found")
	}

	// Delete record
	return tag.Delete(app.DB)
}

func UpdateTag(w http.ResponseWriter, r *http.Request, app *api.AppService) error {

	r.ParseForm()
	name := r.FormValue("name")

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	tag := models.Tag{
		ID:   id,
		Name: name,
	}
	return tag.UpdateTag(app.DB)
}

func GetCategories(w http.ResponseWriter, r *http.Request, app *api.AppService) error {

	categories, err := models.GetCategories(app.DB)
	if err != nil {
		return err
	}
	cookie := shared.GetCookie(r)
	data := shared.Wrapper{categories, cookie}
	return renderer.RenderTemplate(w, "categories.tmpl", data)
}

func GetCategory(w http.ResponseWriter, r *http.Request, app *api.AppService) error {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	category := &models.Category{ID: id}

	err := category.GetCategory(app.DB)
	if err != nil {
		return err
	}
	cookie := shared.GetCookie(r)
	data := shared.Wrapper{category, cookie}
	return renderer.RenderTemplate(w, "category.tmpl", data)
}

func InsertCategory(w http.ResponseWriter, r *http.Request, app *api.AppService) error {
	r.ParseForm()
	name := r.FormValue("name")

	category := &models.Category{
		Name: name,
	}

	// Check if post exists
	category.GetCategoryByName(app.DB)
	if category.ID != 0 {
		// Exists
		return errors.New("Record already exists")
	}

	// Create new category
	err := category.CreateCategory(app.DB)

	return err
}

func DeleteCategory(w http.ResponseWriter, r *http.Request, app *api.AppService) error {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	category := models.Category{ID: id}

	// Check if record exists
	category.GetCategoryByID(app.DB)
	if category.ID == 0 {
		return errors.New("Record not found")
	}

	// Delete record
	err := category.Delete(app.DB)
	return err
}

func UpdateCategory(w http.ResponseWriter, r *http.Request, app *api.AppService) error {

	r.ParseForm()
	name := r.FormValue("name")

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	category := models.Category{
		ID:   id,
		Name: name,
	}
	err := category.UpdateCategory(app.DB)
	return err

}

func Refresh(w http.ResponseWriter, r *http.Request) {
	renderer.LoadTemplates()
	http.Redirect(w, r, "/api/v1/", http.StatusFound)
	fmt.Println("TEMPLATES REFRESHED !")
}
