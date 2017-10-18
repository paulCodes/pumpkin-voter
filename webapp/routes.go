package webapp

import (
	"github.com/gorilla/mux"
	"github.com/paulCodes/pumpkin-voter/webtypes"
	"github.com/paulCodes/pumpkin-voter/webapp/contest"
)

type PVApp struct {
	Env webtypes.Env
}

func AddRoutes(router *mux.Router, env webtypes.Env) {
	prep := router.StrictSlash(true).PathPrefix("/voter").Subrouter()

	app := PVApp{Env: env}
	prep.HandleFunc("", app.Index).Methods("GET", "POST")

	contestA := contest.ContestApp{ Env: env}
	prep.HandleFunc("/contest", contestA.Contests).Methods("GET", "POST")
	prep.HandleFunc("/contest/create", contestA.Create).Methods("GET", "POST")
	prep.HandleFunc("/contest/edit/{contestId}", contestA.Edit).Methods("GET", "POST")
	prep.HandleFunc("/contest/delete/{contestId}", contestA.Delete).Methods("GET", "POST")
	prep.HandleFunc("/category", app.Index).Methods("GET")
	prep.HandleFunc("/entry", app.Index).Methods("POST")
	prep.HandleFunc("/vote", app.Index).Methods("POST")
}
