package model

import (
	"database/sql"
	"errors"
	"time"
)

/*
	Define a Snippet type to hold the data for an individual snippet.
*/

type Blog struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a blogModel type which wraps a sql.DB connection pool.
type BlogModel struct {
	/*
		Go provides three different methods for executing database queries:
		DB.Query() is used for SELECT queries which return multiple rows.
		DB.QueryRow() is used for SELECT queries which return a single row.
		DB.Exec() is used for statements which don’t return rows (like INSERT and DELETE).
	*/
	DB *sql.DB
}

// This will insert a new blog into the database.
func (m *BlogModel) Insert(title string, content string, expires int) (int, error) {
	/*
		Write the SQL statement we want to execute. I've split it over two lines
		for readability (which is why it's surrounded with backquotes instead
		of normal double quotes).
	*/
	stmt := `INSERT INTO blogs (title, content, created, expires) 
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	/*
		Use the Exec() method on the embedded connection pool to execute the
		statement. The first parameter is the SQL statement, followed by the
		title, content and expiry values for the placeholder parameters. This
		method returns a sql.Result type, which contains some basic
		information about what happened when the statement was executed.
	*/
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	/*
		Use the LastInsertId() method on the result to get the ID of our
		newly inserted record in the snippets table.
	*/
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// This will return a specific blog based on its id
func (m *BlogModel) Get(id int) (*Blog, error) {

	stmt := `SELECT id, title, content, created, expires FROM blogs
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	/*
		Use the QueryRow() method on the connection pool to execute our
		SQL statement, passing in the untrusted id variable as the value for the
		placeholder parameter. This returns a pointer to a sql.Row object which
		holds the result from the database.
	*/
	row := m.DB.QueryRow(stmt, id)

	// Initialize a pointer to a new zeroed Blogs struct.
	s := &Blog{}

	/*
		Use row.Scan() to copy the values from each field in sql.Row to the
		corresponding field in the blog struct. Notice that the arguments
		to row.Scan are *pointers* to the place you want to copy the data into,
		and the number of arguments must be exactly the same as the number of
		columns returned by your statement.
	*/
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		/*
			If the query returns no rows, then row.Scan() will return a
			sql.ErrNoRows error. We use the errors.Is() function check for that
			error specifically, and return our own ErrNoRecord error
			instead (we'll create this in a moment).
		*/
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord // from errors.go file
		} else {
			return nil, err
		}
	}

	return s, nil
}

// This will return the 10 most recently created blogs.
func (m *BlogModel) Latest() ([]*Blog, error) {
	return nil, nil
}
