package mysql

import (
	"database/sql"
	"gopkg.in/gorp.v1"
)

type MysqlOps interface {
	Begin() (*gorp.Transaction, error)
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Select(i interface{}, query string, args ...interface{}) ([]interface{}, error)
	SelectOne(holder interface{}, query string, args ...interface{}) error
	Insert(list ...interface{}) error
	Update(list ...interface{}) (int64, error)
	Delete(list ...interface{}) (int64, error)
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
