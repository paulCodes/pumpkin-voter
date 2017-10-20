package mysql

import (
	"github.com/paulCodes/pumpkin-voter/domain"
	"strings"
)

type MysqlEntryStore struct {
	db MysqlStore
}

func (s MysqlStore) EntryStore() MysqlEntryStore {
	return MysqlEntryStore{db: s}
}

func (s MysqlEntryStore) All() (entries []domain.Entry, err error) {
	_, err = s.db.Select(&entries, `select * from entry`)
	return
}

func (s MysqlEntryStore) FindAllForCategoryId(categoryId string) (entries []domain.Entry, err error) {
	_, err = s.db.Select(&entries, `select * from entry where FIND_IN_SET(?, category_ids) `, categoryId)
	return
}

func (s MysqlEntryStore) FindAllCategoryIdFromContest(contestId string) (categoryIds []string, err error) {
	_, err = s.db.Select(&categoryIds, `select category_ids from entry where ? = contest_id`, contestId)
	return
}

func (s MysqlEntryStore) Add(entry domain.Entry) error {
	return s.db.Insert(&entry)
}

func (e MysqlEntryStore) GetID(id string) (entry domain.Entry, err error) {
	err = e.db.SelectOne(&entry, `select * from entry where id = ?`, id)
	return
}

func (e MysqlEntryStore) Replace(entry domain.Entry) error {
	_, err := e.db.Update(&entry)
	return err
}

func (e MysqlEntryStore) Delete(entry domain.Entry) error {
	_, err := e.db.Delete(&entry)
	return err
}

func (s MysqlEntryStore) getIdColumn() string {
	return "id"
}

func (s MysqlEntryStore) getNonIdColumnArray(includeCreated bool, includeUpdated bool) []string {
	cols := []string{"title", "category_ids", "contest_id"}
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
	return row.Scan(&entry.Id, &entry.Title, &entry.CategoryIds, &entry.ContestId)
}
