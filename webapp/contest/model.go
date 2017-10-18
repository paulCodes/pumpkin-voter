package contest

import (
	"github.com/paulCodes/pumpkin-voter/domain"
	"github.com/paulCodes/pumpkin-voter/domain/mysql"
	"strconv"
	"github.com/paulCodes/pumpkin-voter/pvform"
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

func (e ContestLister) AdminListFields3Col() []pvform.ThreeCol {
	return []pvform.ThreeCol{
		{
			GroupName:  "Contest",
			LabelWidth: 3,
			InputWidth: 9,
			Fields: []pvform.FormField{
				{Title: "Title", Type: "text", ClarifyingText: "CT_PrepAdmin_ExamEdit_Name", IsRequired: true},
				{Title: "CategoryIds", Type: "text", ClarifyingText: "CT_PrepAdmin_ExamEdit_Tag",},
				{Title: "Active", Type: "select", ClarifyingText: "CT_PrepAdmin_ExamEdit_Active"},
			},
		},
	}
}

func (e ContestLister) ByField(s string) interface{} {
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

func (e ContestLister) ByFieldChoice(s string) string {
	choices := e.FieldChoices(s)
	for _, v := range choices {
		if e.ByFieldAsString(s) == v[0] {
			return v[1]
		}
	}
	return "CHOICE_NOT_FOUND"
}

func (e ContestLister) FieldChoices(s string) [][]string {
	if s == "Active" {
		return [][]string{
			{"false", "No"},
			{"true", "Yes"},
		}
	}
	panic("FieldChoices: field not found: " + s)
}
