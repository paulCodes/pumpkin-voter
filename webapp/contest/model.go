package contest

import (
	"fmt"
	"github.com/paulCodes/pumpkin-voter/domain"
	"github.com/paulCodes/pumpkin-voter/domain/mysql"
	"strconv"
)

type ContestLister struct {
	domain.Contest
	mysql.Registry `json:"-"`
}

func (e ContestLister) GorpTitle(fieldName string) string {
	return domain.GorpTitleDeterminer(e, fieldName)
}

func (e ContestLister) AdminListFields() [][]string {
	return [][]string{
		{"Id", "text", "text"},
		{"Title", "text", "text"},
		{"CategoryIds", "text", "text"},
		{"Active", "text", "text"},
	}
}

func (e ContestLister) ByField(s string) interface{} {
	println(fmt.Sprintf("s %v", s))
	if s == "Id" {
		return e.Id
	}
	if s == "Title" {
		return e.Title
	}
	if s == "CategoryIds" {
		return e.CategoryIds
	}
	if s == "Active" {
		return e.Active
	}
	panic("ByField: field not found: " + s)
}

func (e ContestLister) ByFieldForList(s string) interface{} {
	return e.ByField(s)
}

func (e ContestLister) ByFieldAsString(s string) string {
	if s == "Active" {
		return strconv.FormatBool(e.ByField(s).(bool))
	}
	return e.ByField(s).(string)
}
