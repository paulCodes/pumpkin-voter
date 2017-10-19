package mysql

import (
	"github.com/paulCodes/pumpkin-voter/domain"
	"strings"
)

type MysqlContestStore struct {
	db MysqlStore
}

func (s MysqlStore) ContestStore() MysqlContestStore {
	return MysqlContestStore{db: s}
}

func (s MysqlContestStore) All() (contests []domain.Contest, err error) {
	_, err = s.db.Select(&contests, `select * from contest`)
	return
}

func (s MysqlContestStore) GetAllAsSelect() (contests [][]string) {
	_, _ = s.db.Select(&contests, `select id, title from contest where active = '1'`)
	return
}

func (s MysqlContestStore) Add(contest domain.Contest) error {
	return s.db.Insert(&contest)
}

func (e MysqlContestStore) GetID(id string) (contest domain.Contest, err error) {
	err = e.db.SelectOne(&contest, `select * from contest where id = ?`, id)
	return
}

func (e MysqlContestStore) Replace(contest domain.Contest) error {
	_, err := e.db.Update(&contest)
	return err
}

func (e MysqlContestStore) Delete(contest domain.Contest) error {
	_, err := e.db.Delete(&contest)
	return err
}

func (s MysqlContestStore) getIdColumn() string {
	return "id"
}

func (s MysqlContestStore) getNonIdColumnArray(includeCreated bool, includeUpdated bool) []string {
	cols := []string{"title", "active"}
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
	return row.Scan(&contest.Id, &contest.Title, &contest.Active)
}
