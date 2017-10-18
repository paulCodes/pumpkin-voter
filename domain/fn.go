package domain

import (
	"github.com/paulCodes/pumpkin-voter/pvhelpers"
	"gopkg.in/gorp.v1"
)

func InitializeDatabase(dbmap *gorp.DbMap, createTables bool) {
	dbmap.AddTableWithName(Category{}, "category").SetKeys(false, "Id")
	dbmap.AddTableWithName(Contest{}, "contest").SetKeys(false, "Id")
	dbmap.AddTableWithName(Entry{}, "entry").SetKeys(false, "Id")
	dbmap.AddTableWithName(Vote{}, "vote").SetKeys(false, "Id")

	if createTables {
		err := dbmap.CreateTablesIfNotExists()
		if err != nil {
			pvhelpers.LogError(err.Error())
			panic("Create tables failed")
		}
	}

	pvhelpers.LogInfo("InitDB complete")
}
