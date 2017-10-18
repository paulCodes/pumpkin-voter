package webapp

import (
	"fmt"
	"github.com/flosch/pongo2"
	"github.com/paulCodes/pumpkin-voter/pvhelpers"
	contestApp "github.com/paulCodes/pumpkin-voter/webapp/contest"
	"net/http"
)

func (app PVApp) Index(w http.ResponseWriter, r *http.Request) {
	contests, err := app.Env.Registry.Contest.All()
	if err != nil {
		panic("An error has occured accessing your contests." + err.Error())
	}
	println(fmt.Sprintf("contests %v", contests))
	models := []contestApp.ContestLister{}
	for _, contest := range contests {
		models = append(models, contestApp.ContestLister{
			Contest:  contest,
			Registry: app.Env.Registry,
		})
	}
	println(fmt.Sprintf("models %v", models))
	pvhelpers.RenderTemplate(w, r, "templates/index.html",
		pongo2.Context{
			"point_to": "dashboard",
			"models":   models,
			"stub":     &contestApp.ContestLister{},
		},
		"vote")
}
