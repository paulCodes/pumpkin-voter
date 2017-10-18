package contest

import (
	"github.com/flosch/pongo2"
	"github.com/paulCodes/pumpkin-voter/domain"
	"github.com/paulCodes/pumpkin-voter/pvhelpers"
	"github.com/paulCodes/pumpkin-voter/webtypes"
	"net/http"
	"strings"
	"fmt"
	"github.com/gorilla/mux"
)

type ContestApp struct {
	Env webtypes.Env
}

func (app ContestApp) Contests(w http.ResponseWriter, r *http.Request) {
	contests, err := app.Env.Registry.Contest.All()
	if err != nil {
		panic("An error has occured accessing your contests." + err.Error())
	}
	println(fmt.Sprintf("contests %v", contests))
	models := []ContestLister{}
	for _, contest := range contests {
		models = append(models, ContestLister{
			Contest:  contest,
			Registry: app.Env.Registry,
		})
	}
	println(fmt.Sprintf("models %v", models))
	pvhelpers.RenderTemplate(w, r, "templates/contest/list.html",
		pongo2.Context{
			"point_to": "contest",
			"models":   models,
			"stub":     &ContestLister{},
		},
		"vote")
}

func (app ContestApp) Create(w http.ResponseWriter, r *http.Request) {
	session, _ := pvhelpers.Store.Get(r, "vote")
	contest := domain.Contest{}
	contest.Id = "-1"
	if r.Method == "POST" {
		contest.Title = strings.TrimSpace(r.FormValue("Title"))
		contest.CategoryIds = strings.TrimSpace(r.FormValue("CategoryIds"))
		contest.Active = strings.TrimSpace(r.FormValue("Active")) == "true"

		//TODO add validation
		contest.Id = pvhelpers.GenerateUUID()
		err := app.Env.Registry.Contest.Add(contest)

		if err != nil {
			pvhelpers.AddFlash(w,r,pvhelpers.FlashMessage{MsgType: "danger", Msg: "Error saving contest"}, "vote")
		} else {
		pvhelpers.AddFlash(w, r, pvhelpers.FlashMessage{MsgType: "success", Msg: "Contest created successfully"}, "vote")
		session.Save(r, w)
		http.Redirect(w, r, "/voter", http.StatusFound)
		return
		}
	}

	session.Save(r, w)
	ctx := pongo2.Context{
		"point_to":         "contest",
		"model":            ContestLister{Contest: contest, Registry: app.Env.Registry},
		"id":               contest.Id,
	}
	pvhelpers.RenderTemplate(w, r, "templates/contest/create.html", ctx, "voter")
}

func (app ContestApp) Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contestId := vars["contestId"]
	session, _ := pvhelpers.Store.Get(r, "vote")
	contest, err := app.Env.Registry.Contest.GetID(contestId)
	if err != nil {
		pvhelpers.LogErrorObject(err)
		pvhelpers.AddFlash(w, r, pvhelpers.FlashMessage{MsgType: "danger", Msg: "Could not find contest"}, "vote")
		session.Save(r, w)
		http.Redirect(w, r, "/voter/contest", http.StatusFound)
		return
	}
	if r.Method == "POST" {
		contest.Title = strings.TrimSpace(r.FormValue("Title"))
		contest.CategoryIds = strings.TrimSpace(r.FormValue("CategoryIds"))
		contest.Active = strings.TrimSpace(r.FormValue("Active")) == "true"

		//TODO add validation
		//contest.Id = pvhelpers.GenerateUUID()
		err := app.Env.Registry.Contest.Replace(contest)

		if err != nil {
			pvhelpers.AddFlash(w,r,pvhelpers.FlashMessage{MsgType: "danger", Msg: "Error saving contest"}, "vote")
		} else {
			pvhelpers.AddFlash(w, r, pvhelpers.FlashMessage{MsgType: "success", Msg: "Contest created successfully"}, "vote")
			session.Save(r, w)
			http.Redirect(w, r, "/voter/contest", http.StatusFound)
			return
		}
	}

	session.Save(r, w)
	ctx := pongo2.Context{
		"point_to":         "contest",
		"model":            ContestLister{Contest: contest, Registry: app.Env.Registry},
		"id":               contest.Id,
	}
	pvhelpers.RenderTemplate(w, r, "templates/contest/create.html", ctx, "voter")
}

func (app ContestApp) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contestId := vars["contestId"]
	session, _ := pvhelpers.Store.Get(r, "vote")
	contest, err := app.Env.Registry.Contest.GetID(contestId)
	if err != nil {
		pvhelpers.LogErrorObject(err)
		pvhelpers.AddFlash(w, r, pvhelpers.FlashMessage{MsgType: "danger", Msg: "Could not find contest"}, "vote")
		session.Save(r, w)
		http.Redirect(w, r, "/voter/contest", http.StatusFound)
		return
	}

	err = app.Env.Registry.Contest.Delete(contest)
	if err != nil {
		pvhelpers.LogErrorObject(err)
		pvhelpers.AddFlash(w, r, pvhelpers.FlashMessage{MsgType: "danger", Msg: "Error deleting contest"}, "vote")
		session.Save(r, w)
		http.Redirect(w, r, "/voter/contest", http.StatusFound)
		return
	}

	session.Save(r, w)
	http.Redirect(w, r, "/voter/contest", http.StatusFound)
}