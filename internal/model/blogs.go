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
	stmt := `SELECT id, title, content, created, expires FROM blogs
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	/*
		We defer rows.Close() to ensure the sql.Rows resultset is
		always properly closed before the Latest() method returns. This defer
		statement should come *after* you check for an error from the Query()
		method. Otherwise, if Query() returns an error, you'll get a panic
		trying to close a nil resultset.

		Important: Closing a resultset with defer rows.Close() is critical in the code above. As
		long as a resultset is open it will keep the underlying database connection open… so if
		something goes wrong in this method and the resultset isn’t closed, it can rapidly lead
		to all the connections in your pool being used up.

	*/
	defer rows.Close()

	// Initialize an empty slice to hold the Snippet structs.
	blogs := []*Blog{} // pointer save memory compare to struct slice

	for rows.Next() {
		s := &Blog{}

		/*
			Use rows.Scan() to copy the values from each field in the row to the
			new Snippet object that we created. Again, the arguments to row.Scan()
			must be pointers to the place you want to copy the data into, and the
			number of arguments must be exactly the same as the number of
			columns returned by your statement.
		*/
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

		if err != nil {
			return nil, err
		}

		blogs = append(blogs, s)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return blogs, nil
}
