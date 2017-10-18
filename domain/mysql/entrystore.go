package mysql

import (
	"github.com/paulCodes/pumpkin-voter/domain"
	"strconv"
	"strings"
)

type MysqlEntryStore MysqlStore

func (s MysqlStore) EntryStore() MysqlEntryStore {
	return MysqlEntryStore(s)
}

func (s MysqlEntryStore) getIdColumn() string {
	return "id"
}

func (s MysqlEntryStore) getNonIdColumnArray(includeCreated bool, includeUpdated bool) []string {
	cols := []string{"title"}
	return cols
}

func (s MysqlEntryStore) getColumns(withId bool) string {
	c := strings.Join(s.getNonIdColumnArray(true, true), ", ")
	if withId {
		return strings.Join([]string{s.getIdColumn(), c}, ", ")
	}
	return c
}

func (s MysqlEntryStore) fillObjRow(entry *domain.Entry, row RowScanner) error {
	return row.Scan(&entry.Id, &entry.Title)
}

func (s MysqlEntryStore) Add(entry *domain.Entry) (err error) {
	stmt, err := s.Prepare("INSERT INTO entry (" + s.getColumns(false) + ") VALUES (?, ?)")
	if err != nil {
		return
	}

	res, err := stmt.Exec(entry.Id, entry.Title)
	if err != nil {
		return
	}

	newId, err := res.LastInsertId()
	if err != nil {
		return
	}

	entry.Id = strconv.FormatInt(newId, 10)
	return
}

func (s MysqlEntryStore) Replace(entry *domain.Entry) (err error) {
	setParams := []string{}
	for _, column := range s.getNonIdColumnArray(false, true) {
		setParams = append(setParams, column+" = ?")
	}
	stmt, err := s.Prepare("UPDATE entry SET " + strings.Join(setParams, ", ") + " WHERE id = ?")
	if err != nil {
		return
	}

	_, err = stmt.Exec(entry.Id, entry.Title)
	if err != nil {
		return
	}
	return
}

func (s MysqlEntryStore) Delete(id string) (err error) {
	stmt, err := s.Prepare("DELETE FROM entry WHERE id = ?")

	if err != nil {
		return
	}
	_, err = stmt.Exec(id)
	return
}

func (s MysqlEntryStore) FindById(id string) (entry domain.Entry, err error) {
	stmtOut, err := s.Prepare("SELECT " + s.getColumns(true) + " FROM entry WHERE id = ?")
	if err != nil {
		return
	}
	defer func() {
		if cerr := stmtOut.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()
	err = s.fillObjRow(&entry, stmtOut.QueryRow(id))
	return
}
