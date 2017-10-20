package mysql

import (
	"github.com/paulCodes/pumpkin-voter/domain"
	"strings"
)

type MysqlVoteStore struct {
	db MysqlStore
}

func (s MysqlStore) VoteStore() MysqlVoteStore {
	return MysqlVoteStore{db: s}
}

func (s MysqlVoteStore) All() (votes []domain.Vote, err error) {
	_, err = s.db.Select(&votes, `select * from vote`)
	return
}

func (s MysqlVoteStore) Add(vote domain.Vote) error {
	return s.db.Insert(&vote)
}

func (e MysqlVoteStore) GetID(id string) (vote domain.Vote, err error) {
	err = e.db.SelectOne(&vote, `select * from vote where id = ?`, id)
	return
}

func (e MysqlVoteStore) Replace(vote domain.Vote) error {
	_, err := e.db.Update(&vote)
	return err
}

func (e MysqlVoteStore) Delete(vote domain.Vote) error {
	_, err := e.db.Delete(&vote)
	return err
}

func (s MysqlVoteStore) getIdColumn() string {
	return "id"
}

func (s MysqlVoteStore) getNonIdColumnArray(includeCreated bool, includeUpdated bool) []string {
	cols := []string{"title", "active"}
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
	return row.Scan(&vote.Id, &vote.ContestId, &vote.CategoryId, &vote.EntryId)
}
