package webapp

import (
	"github.com/gorilla/mux"
	"github.com/paulCodes/pumpkin-voter/webtypes"
)

type PVApp struct {
	Env webtypes.Env
}

func AddRoutes(router *mux.Router, env webtypes.Env) {
	prep := router.StrictSlash(true).PathPrefix("/voter").Subrouter()

	app := PVApp{Env: env}
	prep.HandleFunc("", app.Index).Methods("GET", "POST")
	prep.HandleFunc("/contest", app.Index).Methods("GET")
	prep.HandleFunc("/category", app.Index).Methods("GET")
	prep.HandleFunc("/entry", app.Index).Methods("POST")
	prep.HandleFunc("/vote", app.Index).Methods("POST")
}
