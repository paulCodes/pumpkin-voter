package mysql

import (
	"github.com/paulCodes/pumpkin-voter/domain"
	"strconv"
	"strings"
)

type MysqlVoteStore MysqlStore

func (s MysqlStore) VoteStore() domain.VoteStore  {
	return MysqlVoteStore(s)
}

func (s MysqlVoteStore) getIdColumn() string {
	return "id"
}

func (s MysqlVoteStore) getNonIdColumnArray(includeCreated bool, includeUpdated bool) []string {
	cols := []string{"title","active"}
	return cols
}


func (s MysqlVoteStore) getColumns(withId bool) string {
	c := strings.Join(s.getNonIdColumnArray(true, true), ", ")
	if withId {
		return strings.Join([]string{s.getIdColumn(), c}, ", ")
	}
	return c
}

func (s MysqlVoteStore) fillObjRow(vote *domain.Vote, row RowScanner) error {
	return row.Scan(&vote.ID, &vote.Contest, &vote.Entry, &vote.Category)
}

func (s MysqlVoteStore) Add(vote *domain.Vote) (err error) {
	stmt, err := s.Prepare("INSERT INTO vote (" + s.getColumns(false) + ") VALUES (?, ?)")
	if err != nil {
		return
	}

	res, err := stmt.Exec(vote.ID, vote.Contest, vote.Entry, vote.Category)
	if err != nil {
		return
	}

	newId, err := res.LastInsertId()
	if err != nil {
		return
	}

	vote.ID= strconv.FormatInt(newId, 10)
	return
}

func (s MysqlVoteStore) Replace(vote *domain.Vote) (err error) {
	setParams := []string{}
	for _, column := range s.getNonIdColumnArray(false, true) {
		setParams = append(setParams, column+" = ?")
	}
	stmt, err := s.Prepare("UPDATE vote SET " + strings.Join(setParams, ", ") + " WHERE id = ?")
	if err != nil {
		return
	}

	_, err = stmt.Exec(vote.ID, vote.Contest, vote.Entry, vote.Category)
	if err != nil {
		return
	}
	return
}

func (s MysqlVoteStore) Delete(id string) (err error) {
	stmt, err := s.Prepare("DELETE FROM vote WHERE id = ?")

	if err != nil {
		return
	}
	_, err = stmt.Exec(id)
	return
}

func (s MysqlVoteStore) FindById(id string) (vote domain.Vote, err error) {
	stmtOut, err := s.Prepare("SELECT " + s.getColumns(true) + " FROM vote WHERE id = ?")
	if err != nil {
		return
	}
	defer func() {
		if cerr := stmtOut.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()
	err = s.fillObjRow(&vote, stmtOut.QueryRow(id))
	return
}