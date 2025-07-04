package model

import (
	"database/sql"
	"time"
)

/*
	Define a Snippet type to hold the data for an individual snippet.
*/

type Blogs struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a blogModel type which wraps a sql.DB connection pool.
type BlogModel struct {
	DB *sql.DB
}

// This will insert a new blog into the database.
func (m *BlogModel) Insert(title string, content string, expires int) (int, error) {
	return 0, nil
}

// This will return a specific snippet based on its id
func (m *BlogModel) Get(id int) (*Blogs, error) {
	return nil, nil
}

func (m *BlogModel) Latest() ([]*Blogs, error) {
	return nil, nil
}
