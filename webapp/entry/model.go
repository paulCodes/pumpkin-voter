package entry

import (
	"github.com/paulCodes/pumpkin-voter/domain"
	"github.com/paulCodes/pumpkin-voter/domain/mysql"
	"github.com/paulCodes/pumpkin-voter/pvform"
	"github.com/paulCodes/pumpkin-voter/pvhelpers"
	"log"
	"strings"
)

type EntryLister struct {
	domain.Entry
	mysql.Registry `json:"-"`
}

func (e EntryLister) GorpTitle(fieldName string) string {
	return domain.GorpTitleDeterminer(e, fieldName)
}

func (e EntryLister) AdminListFields() [][]string {
	return [][]string{
		{"Title", "text", "text"},
		{"CategoryIds", "text", "text"},
		{"ContestId", "text", "text"},
	}
}

func (e EntryLister) AdminListFields3Col() []pvform.ThreeCol {
	return []pvform.ThreeCol{
		{
			GroupName:  "Entry",
			LabelWidth: 3,
			InputWidth: 9,
			Fields: []pvform.FormField{
				{Title: "Title", Type: "text", ClarifyingText: "CT_PrepAdmin_ExamEdit_Name", IsRequired: true},
				{Title: "CategoryIds", Type: "multiselect", ClarifyingText: "CT_PrepAdmin_ExamEdit_Tag",},
				{Title: "ContestId", Type: "select", ClarifyingText: "CT_PrepAdmin_ExamEdit_Tag",},
			},
		},
	}
}

func (e EntryLister) ByField(s string) interface{} {
	if s == "Id" {
		return e.Id
	}
	if s == "Title" {
		return e.Title
	}
	if s == "CategoryIds" {
		return e.CategoryIds
	}
	if s == "ContestId" {
		return e.ContestId
	}
	panic("ByField: field not found: " + s)
}

func (e EntryLister) ByFieldForList(s string) interface{} {
	return e.ByField(s)
}

func (e EntryLister) ByFieldAsString(s string) string {
	return e.ByField(s).(string)
}

func (e EntryLister) ByFieldAsSelect(s string) []string {
	return strings.Split(e.ByField(s).(string),",")

}

func (e EntryLister) ByFieldChoice(s string) string {
	choices := e.FieldChoices(s)
	for _, v := range choices {
		if e.ByFieldAsString(s) == v[0] {
			return v[1]
		}
	}
	return "CHOICE_NOT_FOUND"
}

func (e EntryLister) FieldChoices(s string) [][]string {
	if s == "CategoryIds" {
		return e.getCategoriesAsSelect()
	}
	if s == "ContestId" {
		return e.getContestAsSelect()
	}
	panic("FieldChoices: field not found: " + s)
}

func (e EntryLister) getCategoriesAsSelect() (categories [][]string) {
	results, err := pvhelpers.GorpDb.Db.Query(`select id, title from category where active = '1'`)
log.Printf("error %v", err)
	for results.Next() {
		var (
			id   string
			title string
		)
		if err := results.Scan(&id, &title); err != nil {
			log.Fatal(err)
		}
		temp := [][]string{{id,title}}
		categories = append(categories,temp... )
	}
	return
}

func (e EntryLister) getContestAsSelect() (categories [][]string) {
	results, _ := pvhelpers.GorpDb.Db.Query(`select id, title from contest where active = '1'`)
	for results.Next() {
		var (
			id   string
			title string
		)
		if err := results.Scan(&id, &title); err != nil {
			log.Fatal(err)
		}
		temp := [][]string{{id,title}}
		categories = append(categories,temp... )
	}
	return
}