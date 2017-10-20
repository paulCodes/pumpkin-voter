package category

import (
	"github.com/paulCodes/pumpkin-voter/domain"
	"github.com/paulCodes/pumpkin-voter/domain/mysql"
	"strconv"
	"github.com/paulCodes/pumpkin-voter/pvform"
)

type CategoryLister struct {
	domain.Category
	mysql.Registry `json:"-"`
}

func (e CategoryLister) GorpTitle(fieldName string) string {
	return domain.GorpTitleDeterminer(e, fieldName)
}

func (e CategoryLister) AdminListFields() [][]string {
	return [][]string{
		{"Title", "text", "text"},
		{"Active", "text", "text"},
	}
}

func (e CategoryLister) AdminListFields3Col() []pvform.ThreeCol {
	return []pvform.ThreeCol{
		{
			GroupName:  "Category",
			LabelWidth: 3,
			InputWidth: 9,
			Fields: []pvform.FormField{
				{Title: "Title", Type: "text", ClarifyingText: "CT_PrepAdmin_ExamEdit_Name", IsRequired: true},
				{Title: "Active", Type: "select", ClarifyingText: "CT_PrepAdmin_ExamEdit_Active"},
			},
		},
	}
}

func (e CategoryLister) ByField(s string) interface{} {
	if s == "Id" {
		return e.Id
	}
	if s == "Title" {
		return e.Title
	}
	if s == "Active" {
		return e.Active
	}
	panic("ByField: field not found: " + s)
}

func (e CategoryLister) ByFieldForList(s string) interface{} {
	return e.ByField(s)
}

func (e CategoryLister) ByFieldAsString(s string) string {
	if s == "Active" {
		return strconv.FormatBool(e.ByField(s).(bool))
	}
	return e.ByField(s).(string)
}

func (e CategoryLister) ByFieldChoice(s string) string {
	choices := e.FieldChoices(s)
	for _, v := range choices {
		if e.ByFieldAsString(s) == v[0] {
			return v[1]
		}
	}
	return "CHOICE_NOT_FOUND"
}

func (e CategoryLister) FieldChoices(s string) [][]string {
	if s == "Active" {
		return [][]string{
			{"false", "No"},
			{"true", "Yes"},
		}
	}
	panic("FieldChoices: field not found: " + s)
}
