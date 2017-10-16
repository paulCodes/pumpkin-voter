package mysql

import (
	"github.com/paulCodes/pumpkin-voter/domain"
	"strconv"
	"strings"
)

type MysqlCategoryStore MysqlStore

func (s MysqlStore) CategoryStore() domain.CategoryStore  {
return MysqlCategoryStore(s)
}

func (s MysqlCategoryStore) getIdColumn() string {
return "id"
}

func (s MysqlCategoryStore) getNonIdColumnArray(includeCreated bool, includeUpdated bool) []string {
cols := []string{"title"}
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
return row.Scan(&category.ID, &category.Title)
}

func (s MysqlCategoryStore) Add(category *domain.Category) (err error) {
stmt, err := s.Prepare("INSERT INTO category (" + s.getColumns(false) + ") VALUES (?, ?)")
if err != nil {
return
}

res, err := stmt.Exec(category.ID, category.Title)
if err != nil {
return
}

newId, err := res.LastInsertId()
if err != nil {
return
}

category.ID= strconv.FormatInt(newId, 10)
return
}

func (s MysqlCategoryStore) Replace(category *domain.Category) (err error) {
setParams := []string{}
for _, column := range s.getNonIdColumnArray(false, true) {
setParams = append(setParams, column+" = ?")
}
stmt, err := s.Prepare("UPDATE category SET " + strings.Join(setParams, ", ") + " WHERE id = ?")
if err != nil {
return
}

_, err = stmt.Exec(category.ID, category.Title)
if err != nil {
return
}
return
}

func (s MysqlCategoryStore) Delete(id string) (err error) {
stmt, err := s.Prepare("DELETE FROM category WHERE id = ?")

if err != nil {
return
}
_, err = stmt.Exec(id)
return
}

func (s MysqlCategoryStore) FindById(id string) (category domain.Category, err error) {
stmtOut, err := s.Prepare("SELECT " + s.getColumns(true) + " FROM category WHERE id = ?")
if err != nil {
return
}
defer func() {
if cerr := stmtOut.Close(); cerr != nil && err == nil {
err = cerr
}
}()
err = s.fillObjRow(&category, stmtOut.QueryRow(id))
return
}