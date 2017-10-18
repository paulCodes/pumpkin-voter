package main

import (
	"encoding/gob"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/paulCodes/pumpkin-voter/domain/mysql"
	"github.com/paulCodes/pumpkin-voter/pvdb"
	"github.com/paulCodes/pumpkin-voter/pvhelpers"
	"github.com/paulCodes/pumpkin-voter/pvtemplate"
	"github.com/paulCodes/pumpkin-voter/webapp"
	"github.com/paulCodes/pumpkin-voter/webtypes"
	"gopkg.in/gcfg.v1"
	"net/http"
	"os"
)

func main() {
	argsWithoutProg := os.Args[1:]
	configFilename := "config.gcfg"
	if len(argsWithoutProg) > 0 {
		configFilename = argsWithoutProg[0]
	}
	pvhelpers.EncryptionKey = []byte("a CRCH 1133 BLRG very sEkrE7 KeY") // 32 bytes

	err := gcfg.ReadFileInto(&pvhelpers.Cfg, configFilename)
	if err != nil {
		panic(err)
	}

	pvdb.InitDbModels()
	pvtemplate.RegisterFilters()
	pvhelpers.PrintServerConnectionDetails()
	gob.Register(&pvhelpers.FlashMessage{})

	env := webtypes.Env{
		Registry: mysql.NewRegistry(pvhelpers.GorpDb),
	}

	router := mux.NewRouter()
	// frontend.
	webapp.AddRoutes(router, env)

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/version.html")
	})

	n := negroni.New(pvhelpers.NewNoStaticLogger(), negroni.NewRecovery())
	n.UseHandler(router)
	n.Run(pvhelpers.Cfg.Deploy.Listen)
}
