// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package mysql

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/paulCodes/pumpkin-voter/domain"
)

var createTableStatements = []string{
	`CREATE DATABASE IF NOT EXISTS pumpkin DEFAULT CHARACTER SET = 'utf8' DEFAULT COLLATE 'utf8_general_ci';`,
	`USE pumpkin;`,
	`CREATE TABLE IF NOT EXISTS entry (
		id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		title VARCHAR(255) NULL,
		PRIMARY KEY (id)
	)`,
	`CREATE TABLE IF NOT EXISTS vote (
		id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		category_id INT UNSIGNED NOT NULL,
		entry_id INT UNSIGNED NOT NULL,
		contest_id INT UNSIGNED NOT NULL,
		PRIMARY KEY (id)
	)`,
	`CREATE TABLE IF NOT EXISTS contest (
		id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		title VARCHAR(255) NULL,
		active VARCHAR(255) NULL,,
		PRIMARY KEY (id)
	)`,
	`CREATE TABLE IF NOT EXISTS entry_category_contest (
		id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		category_id INT UNSIGNED NOT NULL,
		entry_id INT UNSIGNED NOT NULL,
		contest_id INT UNSIGNED NOT NULL,
		PRIMARY KEY (id)
	)`,
}

// mysqlDB persists entries to a MySQL instance.
type mysqlDB struct {
	conn *sql.DB

	list   *sql.Stmt
	listBy *sql.Stmt
	insert *sql.Stmt
	get    *sql.Stmt
	update *sql.Stmt
	delete *sql.Stmt
}

// Ensure mysqlDB conforms to the BookDatabase interface.
var _ EntryStore = &mysqlDB{}

type MySQLConfig struct {
	// Optional.
	Username, Password string

	// Host of the MySQL instance.
	//
	// If set, UnixSocket should be unset.
	Host string

	// Port of the MySQL instance.
	//
	// If set, UnixSocket should be unset.
	Port int

	// UnixSocket is the filepath to a unix socket.
	//
	// If set, Host and Port should be unset.
	UnixSocket string
}

// dataStoreName returns a connection string suitable for sql.Open.
func (c MySQLConfig) dataStoreName(databaseName string) string {
	var cred string
	// [username[:password]@]
	if c.Username != "" {
		cred = c.Username
		if c.Password != "" {
			cred = cred + ":" + c.Password
		}
		cred = cred + "@"
	}

	if c.UnixSocket != "" {
		return fmt.Sprintf("%sunix(%s)/%s", cred, c.UnixSocket, databaseName)
	}
	return fmt.Sprintf("%stcp([%s]:%d)/%s", cred, c.Host, c.Port, databaseName)
}

// newMySQLDB creates a new entryStore backed by a given MySQL server.
func newMySQLDB(config MySQLConfig) (EntryStore, error) {
	// Check database and table exists. If not, create it.
	if err := config.ensureTableExists(); err != nil {
		return nil, err
	}

	conn, err := sql.Open("mysql", config.dataStoreName("library"))
	if err != nil {
		return nil, fmt.Errorf("mysql: could not get a connection: %v", err)
	}
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("mysql: could not establish a good connection: %v", err)
	}

	db := &mysqlDB{
		conn: conn,
	}

	// Prepared statements. The actual SQL queries are in the code near the
	// relevant method (e.g. addBook).
	if db.list, err = conn.Prepare(listStatement); err != nil {
		return nil, fmt.Errorf("mysql: prepare list: %v", err)
	}
	if db.get, err = conn.Prepare(getStatement); err != nil {
		return nil, fmt.Errorf("mysql: prepare get: %v", err)
	}
	if db.insert, err = conn.Prepare(insertStatement); err != nil {
		return nil, fmt.Errorf("mysql: prepare insert: %v", err)
	}
	if db.update, err = conn.Prepare(updateStatement); err != nil {
		return nil, fmt.Errorf("mysql: prepare update: %v", err)
	}
	if db.delete, err = conn.Prepare(deleteStatement); err != nil {
		return nil, fmt.Errorf("mysql: prepare delete: %v", err)
	}

	return db, nil
}

// Close closes the database, freeing up any resources.
func (db *mysqlDB) Close() {
	db.conn.Close()
}

// rowScanner is implemented by sql.Row and sql.Rows
type rowScanner interface {
	Scan(dest ...interface{}) error
}

// scanBook reads a book from a sql.Row or sql.Rows
func scanEntry(s rowScanner) (*domain.Entry, error) {
	var (
		id            int64
		title         sql.NullString
	)
	if err := s.Scan(&id, &title); err != nil {
		return nil, err
	}

	entry := &domain.Entry{
		ID:            id,
		Title:         title.String,
	}
	return entry, nil
}

const listStatement = `SELECT * FROM entry ORDER BY title`

// ListBooks returns a list of books, ordered by title.
func (db *mysqlDB) ListEntries() ([]*domain.Entry, error) {
	rows, err := db.list.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entrys []*domain.Entry
	for rows.Next() {
		entry, err := scanEntry(rows)
		if err != nil {
			return nil, fmt.Errorf("mysql: could not read row: %v", err)
		}

		entrys = append(entrys, entry)
	}

	return entrys, nil
}

const getStatement = "SELECT * FROM books WHERE id = ?"

// GetBook retrieves a book by its ID.
func (db *mysqlDB) GetEntry(id int64) (*domain.Entry, error) {
	book, err := scanEntry(db.get.QueryRow(id))
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("mysql: could not find book with id %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("mysql: could not get book: %v", err)
	}
	return book, nil
}

const insertStatement = `
  INSERT INTO entry (title) VALUES (?)`

// AddEntry saves a given entry, assigning it a new ID.
func (db *mysqlDB) AddEntry(b *domain.Entry) (id int64, err error) {
	r, err := execAffectingOneRow(db.insert, b.Title)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := r.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("mysql: could not get last insert ID: %v", err)
	}
	return lastInsertID, nil
}

const deleteStatement = `DELETE FROM entry WHERE id = ?`

// DeleteBook removes a given book by its ID.
func (db *mysqlDB) DeleteBook(id int64) error {
	if id == 0 {
		return errors.New("mysql: entry with unassigned ID passed into deleteEntry")
	}
	_, err := execAffectingOneRow(db.delete, id)
	return err
}

const updateStatement = `
  UPDATE entry
  SET title=?
  WHERE id = ?`

// UpdateEntry updates the entry for a given entry.
func (db *mysqlDB) UpdateEntry(b *domain.Entry) error {
	if b.ID == 0 {
		return errors.New("mysql: entry with unassigned ID passed into updateEntry")
	}

	_, err := execAffectingOneRow(db.update, b.Title, b.ID)
	return err
}

// ensureTableExists checks the table exists. If not, it creates it.
func (config MySQLConfig) ensureTableExists() error {
	conn, err := sql.Open("mysql", config.dataStoreName(""))
	if err != nil {
		return fmt.Errorf("mysql: could not get a connection: %v", err)
	}
	defer conn.Close()

	// Check the connection.
	if conn.Ping() == driver.ErrBadConn {
		return fmt.Errorf("mysql: could not connect to the database. " +
			"could be bad address, or this address is not whitelisted for access.")
	}

	if _, err := conn.Exec("USE pumpkin"); err != nil {
		// MySQL error 1049 is "database does not exist"
		if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1049 {
			return createTable(conn)
		}
	}

	if _, err := conn.Exec("DESCRIBE entry"); err != nil {
		// MySQL error 1146 is "table does not exist"
		if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1146 {
			return createTable(conn)
		}
		// Unknown error.
		return fmt.Errorf("mysql: could not connect to the database: %v", err)
	}
	return nil
}

// createTable creates the table, and if necessary, the database.
func createTable(conn *sql.DB) error {
	for _, stmt := range createTableStatements {
		_, err := conn.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

// execAffectingOneRow executes a given statement, expecting one row to be affected.
func execAffectingOneRow(stmt *sql.Stmt, args ...interface{}) (sql.Result, error) {
	r, err := stmt.Exec(args...)
	if err != nil {
		return r, fmt.Errorf("mysql: could not execute statement: %v", err)
	}
	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return r, fmt.Errorf("mysql: could not get rows affected: %v", err)
	} else if rowsAffected != 1 {
		return r, fmt.Errorf("mysql: expected 1 row affected, got %d", rowsAffected)
	}
	return r, nil
}
