package models

import (
	"database/sql"
)

type Tag struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Posts []Post `json:"posts,omitempty"`
}

func (t *Tag) GetTag(db *sql.DB) error {
	const sqlQuery = `SELECT id,name FROM tags WHERE id=$1;`
	err := db.QueryRow(sqlQuery, t.ID).Scan(&t.ID, &t.Name)
	if err != nil {
		return err
	}

	// Get all the posts
	const postQuery = `select posts.id,posts.name,posts.date,posts.content ` +
		`FROM tags ` +
		`inner join tags_posts on tags.id=tags_posts.tag_id ` +
		`inner join posts  on tags_posts.post_id=posts.id ` +
		`where tags.id=$1;`
	query, err := db.Query(postQuery, t.ID)
	defer query.Close()
	if err != nil {
		return err
	}
	p := Post{}
	for query.Next() {
		query.Scan(&p.ID, &p.Name, &p.Date, &p.Content)
		t.Posts = append(t.Posts, p)
	}
	return nil
}

func (t *Tag) CreateTag(db *sql.DB) error {
	const sqlQuery = `INSERT INTO tags (name` +
		`)VALUES(` +
		`$1) RETURNING id;`
	return db.QueryRow(sqlQuery, t.Name).Scan(&t.ID)
}

func GetTags(db *sql.DB) ([]Tag, error) {
	const sqlQuery = `SELECT * FROM tags;`
	query, err := db.Query(sqlQuery)
	defer query.Close()
	if err != nil {
		return nil, err
	}
	var tags []Tag
	t := Tag{}
	for query.Next() {
		query.Scan(&t.ID, &t.Name)
		tags = append(tags, t)
	}
	return tags, nil
}

func (c *Tag) Delete(db *sql.DB) error {
	const sqlQuery = `DELETE FROM public.tags WHERE id=$1`
	_, err := db.Exec(sqlQuery, c.ID)
	return err
}

func (t *Tag) GetTagByName(db *sql.DB) error {
	const sqlQuery = `SELECT id FROM public.tags WHERE name=$1;`
	return db.QueryRow(sqlQuery, t.Name).Scan(&t.ID)
}

func (t *Tag) GetTagByID(db *sql.DB) error {
	const sqlQuery = `SELECT id FROM public.tags WHERE id=$1;`
	// TODO: Make it less hacky
	t1 := Tag{}
	err := db.QueryRow(sqlQuery, t.ID).Scan(&t1.ID)
	t.ID = t1.ID
	return err
}

func (t *Tag) UpdateTag(db *sql.DB) error {
	const sqlQuery = `UPDATE tags SET name=$1 WHERE id=$2 RETURNING id,name`
	return db.QueryRow(sqlQuery, t.Name, t.ID).Scan(&t.ID, &t.Name)
}
