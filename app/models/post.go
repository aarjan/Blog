package models

import (
	"database/sql"
)

type Post struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Date       string `json:"date"`
	Content    string `json:"content"`
	CategoryID int    `json:"category_id,omitempty"`

	// Extra
	Tags []Tag     `json:"tags,omitempty"`
	Cat  *Category `json:"category,omitempty"`
}

func (p *Post) CreatePost(db *sql.DB) error {
	var err error

	const sqlQuery = `INSERT INTO public.posts (` +
		`name,date,content,category_id` +
		`)VALUES(` +
		`$1,$2,$3 ` +
		`) RETURNING id`

	err = db.QueryRow(sqlQuery, p.Name, p.Date, p.Content, p.CategoryID).Scan(&p.ID)
	if err != nil {
		return err
	}
	return nil
}

func (p *Post) GetPostByID(db *sql.DB) error {
	const sqlQuery = `SELECT id FROM public.posts WHERE id=$1;`
	// TODO: Make it less hacky
	p1 := Post{}
	err := db.QueryRow(sqlQuery, p.ID).Scan(&p1.ID)
	p.ID = p1.ID
	return err
}

func (p *Post) GetPostByName(db *sql.DB) error {
	const sqlQuery = `SELECT id FROM public.posts WHERE name=$1;`
	return db.QueryRow(sqlQuery, p.Name).Scan(&p.ID)
}

func (p *Post) Delete(db *sql.DB) error {
	const sqlQuery = `DELETE FROM public.posts WHERE id=$1`
	_, err := db.Exec(sqlQuery, p.ID)
	return err
}
func (p *Post) GetPost(db *sql.DB) error {
	const sqlQuery = `SELECT id,date,name,content,category_id FROM public.posts WHERE id=$1;`
	err := db.QueryRow(sqlQuery, p.ID).Scan(&p.ID, &p.Date, &p.Name, &p.Content, &p.CategoryID)
	if err != nil {
		return err
	}

	// Get the tags associated with the post
	const tagsQuery = `SELECT tags.* ` +
		`FROM posts ` +
		`INNER JOIN tags_posts on tags_posts.post_id = posts.id ` +
		`INNER JOIN tags on tags.id = tags_posts.tag_id ` +
		`WHERE posts.id=$1`

	query, err := db.Query(tagsQuery, p.ID)
	defer query.Close()
	if err != nil {
		return err
	}
	t := Tag{}
	for query.Next() {
		query.Scan(&t.ID, &t.Name)
		p.Tags = append(p.Tags, t)
	}

	c := &Category{}
	// Get the category from category_id
	const catQuery = `SELECT id,name FROM categories WHERE id=$1`
	err = db.QueryRow(catQuery, p.CategoryID).Scan(&c.ID, &c.Name)

	p.Cat = c
	return err
}

func GetPosts(db *sql.DB) ([]Post, error) {
	const sqlQuery = `SELECT id,name,date,content FROM posts;`
	query, err := db.Query(sqlQuery)
	defer query.Close()
	if err != nil {
		return nil, err
	}
	var posts []Post
	p := Post{}
	for query.Next() {
		query.Scan(&p.ID, &p.Name, &p.Date, &p.Content)
		posts = append(posts, p)
	}
	return posts, nil
}

// func GetCategoryID(name string, db *sql.DB) (int, error) {
// 	const sqlQuery = `SELECT id FROM categories where name=$1`
// 	var id int
// 	err := db.QueryRow(sqlQuery, name).Scan(&id)
// 	return id, err
// }
