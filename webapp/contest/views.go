package contest

import (
	"github.com/flosch/pongo2"
	"github.com/paulCodes/pumpkin-voter/domain"
	"github.com/paulCodes/pumpkin-voter/pvhelpers"
	"github.com/paulCodes/pumpkin-voter/webtypes"
	"net/http"
)

type ContestApp struct {
	Env webtypes.Env
}

func (app ContestApp) Contests(w http.ResponseWriter, r *http.Request) {
	contests, err := app.Env.Registry.Entry.Select(&domain.Contest{}, "select * from contest")
	if err != nil {
		panic("An error has occured accessing your entries.")
	}

	ctx := pongo2.Context{
		"entries": contests,
	}

	pvhelpers.RenderTemplate(w, r, "contest/list.html", ctx, "vote")
}
