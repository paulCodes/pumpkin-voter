package mysql

import (
	"fmt"
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
	var t int
	err = s.db.SelectOne(&t, `select count(*) from contest`)
	_, err = s.db.Select(&contests, `select * from contest`)
	println(fmt.Sprintf("t %v", t))
	return
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
