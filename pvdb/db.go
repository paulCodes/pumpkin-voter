package pvdb

import (
	"database/sql"
	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/paulCodes/pumpkin-voter/domain"
	"github.com/paulCodes/pumpkin-voter/pvhelpers"
	"gopkg.in/gorp.v1"
	"log"
	"os"
)

func InitDbModels() {
	doSync := true

	dbInfo := pvhelpers.Cfg.Database

	dbUrl := dbInfo.Username + ":" + dbInfo.Password + "@" + dbInfo.Host + "/" + dbInfo.DB + "?charset=utf8&parseTime=true"

	// GORP
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		pvhelpers.LogErrorObject(err)
		panic("sql.Open failed")
	}

	// construct a gorp DbMap
	pvhelpers.GorpDb = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	if pvhelpers.Cfg.Database.Debug {
		pvhelpers.GorpDb.TraceOn("[gorp]", log.New(os.Stdout, "myapp:", log.Lmicroseconds))
	} else {
		pvhelpers.GorpDb.TraceOff()
	}
	domain.InitializeDatabase(pvhelpers.GorpDb, doSync)
	//if doSync {
	//	regadmin.StubTestData(ccshelpers.GorpDb)
	//}
}
