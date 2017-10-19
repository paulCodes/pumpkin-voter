package mysql

import (
	"github.com/paulCodes/pumpkin-voter/domain"
	"strings"
)

type MysqlCategoryStore struct {
	db MysqlStore
}

func (s MysqlStore) CategoryStore() MysqlCategoryStore {
	return MysqlCategoryStore{db: s}
}

func (s MysqlCategoryStore) All() (categories []domain.Category, err error) {
	_, err = s.db.Select(&categories, `select * from category`)
	return
}

func (s MysqlCategoryStore) Add(category domain.Category) error {
	return s.db.Insert(&category)
}

func (e MysqlCategoryStore) GetID(id string) (category domain.Category, err error) {
	err = e.db.SelectOne(&category, `select * from category where id = ?`, id)
	return
}

func (e MysqlCategoryStore) GetAllAsSelect() (categories [][]string) {
	_, _ = e.db.Select(&categories, `select id, title from category where active = '1'`)
	return
}

func (e MysqlCategoryStore) Replace(category domain.Category) error {
	_, err := e.db.Update(&category)
	return err
}

func (e MysqlCategoryStore) Delete(category domain.Category) error {
	_, err := e.db.Delete(&category)
	return err
}

func (s MysqlCategoryStore) getIdColumn() string {
	return "id"
}

func (s MysqlCategoryStore) getNonIdColumnArray(includeCreated bool, includeUpdated bool) []string {
	cols := []string{"title", "active"}
	return cols
}

func (s MysqlCategoryStore) getColumns(withId bool) string {
	c := strings.Join(s.getNonIdColumnArray(true, true), ", ")
	if withId {
		return strings.Join([]string{s.getIdColumn(), c}, ", ")
	}
	return c
}

func (s MysqlCategoryStore) fillObjRow(category *domain.Category, row RowScanner) error {
	return row.Scan(&category.Id, &category.Title, &category.Active)
}
