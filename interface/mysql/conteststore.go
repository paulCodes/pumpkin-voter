package mysql

import (
	"github.com/paulCodes/pumpkin-voter/domain"
	"strconv"
	"strings"
)

type MysqlContestStore MysqlStore

func (s MysqlStore) ContestStore() domain.ContestStore  {
	return MysqlContestStore(s)
}

func (s MysqlContestStore) getIdColumn() string {
	return "id"
}

func (s MysqlContestStore) getNonIdColumnArray(includeCreated bool, includeUpdated bool) []string {
	cols := []string{"title","active"}
	return cols
}


func (s MysqlContestStore) getColumns(withId bool) string {
	c := strings.Join(s.getNonIdColumnArray(true, true), ", ")
	if withId {
		return strings.Join([]string{s.getIdColumn(), c}, ", ")
	}
	return c
}

func (s MysqlContestStore) fillObjRow(contest *domain.Contest, row RowScanner) error {
	return row.Scan(&contest.ID, &contest.Title, &contest.Active)
}

func (s MysqlContestStore) Add(contest *domain.Contest) (err error) {
	stmt, err := s.Prepare("INSERT INTO contest (" + s.getColumns(false) + ") VALUES (?, ?)")
	if err != nil {
		return
	}

	res, err := stmt.Exec(contest.ID, contest.Title)
	if err != nil {
		return
	}

	newId, err := res.LastInsertId()
	if err != nil {
		return
	}

	contest.ID= strconv.FormatInt(newId, 10)
	return
}

func (s MysqlContestStore) Replace(contest *domain.Contest) (err error) {
	setParams := []string{}
	for _, column := range s.getNonIdColumnArray(false, true) {
		setParams = append(setParams, column+" = ?")
	}
	stmt, err := s.Prepare("UPDATE contest SET " + strings.Join(setParams, ", ") + " WHERE id = ?")
	if err != nil {
		return
	}

	_, err = stmt.Exec(contest.ID, contest.Title)
	if err != nil {
		return
	}
	return
}

func (s MysqlContestStore) Delete(id string) (err error) {
	stmt, err := s.Prepare("DELETE FROM contest WHERE id = ?")

	if err != nil {
		return
	}
	_, err = stmt.Exec(id)
	return
}

func (s MysqlContestStore) FindById(id string) (contest domain.Contest, err error) {
	stmtOut, err := s.Prepare("SELECT " + s.getColumns(true) + " FROM contest WHERE id = ?")
	if err != nil {
		return
	}
	defer func() {
		if cerr := stmtOut.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()
	err = s.fillObjRow(&contest, stmtOut.QueryRow(id))
	return
}