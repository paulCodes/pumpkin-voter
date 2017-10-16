package mysql

import "database/sql"

type MysqlOps interface {
	Begin() (*sql.Tx, error)
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type mysqlStore struct {
	MysqlOps
}

type MysqlStore mysqlStore

func NewMysqlStore(client MysqlOps) MysqlStore {
	return MysqlStore{
		MysqlOps: client,
	}
}

type RowScanner interface {
	Scan(dest ...interface{}) error
}
