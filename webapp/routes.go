package webapp

import (
	"github.com/gorilla/mux"
	"github.com/paulCodes/pumpkin-voter/webtypes"
	"github.com/paulCodes/pumpkin-voter/webapp/contest"
	"github.com/paulCodes/pumpkin-voter/webapp/category"
	"github.com/paulCodes/pumpkin-voter/webapp/entry"
)

type PVApp struct {
	Env webtypes.Env
}

func AddRoutes(router *mux.Router, env webtypes.Env) {
	prep := router.StrictSlash(true).PathPrefix("/voter").Subrouter()

	app := PVApp{Env: env}
	prep.HandleFunc("", app.Index).Methods("GET", "POST")

	contestApp := contest.ContestApp{ Env: env}
	prep.HandleFunc("/contest", contestApp.Contests).Methods("GET", "POST")
	prep.HandleFunc("/contest/create", contestApp.Create).Methods("GET", "POST")
	prep.HandleFunc("/contest/edit/{contestId}", contestApp.Edit).Methods("GET", "POST")
	prep.HandleFunc("/contest/delete/{contestId}", contestApp.Delete).Methods("GET", "POST")
	prep.HandleFunc("/contest/results/{contestId}", contestApp.Results).Methods("GET", "POST")
	prep.HandleFunc("/contest/vote/{contestId}", contestApp.Vote).Methods("GET", "POST")

	categoryApp := category.CategoryApp{ Env: env}
	prep.HandleFunc("/category", categoryApp.Categories).Methods("GET", "POST")
	prep.HandleFunc("/category/create", categoryApp.Create).Methods("GET", "POST")
	prep.HandleFunc("/category/edit/{categoryId}", categoryApp.Edit).Methods("GET", "POST")
	prep.HandleFunc("/category/delete/{categoryId}", categoryApp.Delete).Methods("GET", "POST")


	entryApp := entry.EntryApp{ Env: env}
	prep.HandleFunc("/entry", entryApp.Entries).Methods("GET", "POST")
	prep.HandleFunc("/entry/create", entryApp.Create).Methods("GET", "POST")
	prep.HandleFunc("/entry/edit/{entryId}", entryApp.Edit).Methods("GET", "POST")
	prep.HandleFunc("/entry/delete/{entryId}", entryApp.Delete).Methods("GET", "POST")
}
