package entry

import (
	"github.com/flosch/pongo2"
	"github.com/paulCodes/pumpkin-voter/domain"
	"github.com/paulCodes/pumpkin-voter/pvhelpers"
	"github.com/paulCodes/pumpkin-voter/webtypes"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"log"
)

type EntryApp struct {
	Env webtypes.Env
}

func (app EntryApp) Entries(w http.ResponseWriter, r *http.Request) {
	entries, err := app.Env.Registry.Entry.All()
	if err != nil {
		panic("An error has occured accessing your entrys." + err.Error())
	}

	models := []EntryLister{}
	for _, entry := range entries {
		models = append(models, EntryLister{
			Entry:  entry,
			Registry: app.Env.Registry,
		})
	}

	pvhelpers.RenderTemplate(w, r, "templates/entry/list.html",
		pongo2.Context{
			"point_to": "entry",
			"models":   models,
			"stub":     &EntryLister{},
		},
		"vote")
}

func (app EntryApp) Create(w http.ResponseWriter, r *http.Request) {
	session, _ := pvhelpers.Store.Get(r, "vote")
	entry := domain.Entry{}
	entry.Id = "-1"
	r.ParseForm()
	if r.Method == "POST" {
		entry.Title = strings.TrimSpace(r.FormValue("Title"))
		entry.CategoryIds = r.FormValue("CategoryIds")
		entry.ContestId = strings.TrimSpace(r.FormValue("ContestId"))

		//TODO add validation
		entry.Id = pvhelpers.GenerateUUID()
		err := app.Env.Registry.Entry.Add(entry)

		if err != nil {
			pvhelpers.AddFlash(w,r,pvhelpers.FlashMessage{MsgType: "danger", Msg: "Error saving entry"}, "vote")
		} else {
		pvhelpers.AddFlash(w, r, pvhelpers.FlashMessage{MsgType: "success", Msg: "Entry created successfully"}, "vote")
		session.Save(r, w)
		http.Redirect(w, r, "/voter/entry", http.StatusFound)
		return
		}
	}

	session.Save(r, w)
	ctx := pongo2.Context{
		"point_to": "entry",
		"model":    EntryLister{Entry: entry, Registry: app.Env.Registry},
		"id":       entry.Id,
	}
	pvhelpers.RenderTemplate(w, r, "templates/entry/create.html", ctx, "voter")
}

func (app EntryApp) Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entryId := vars["entryId"]
	session, _ := pvhelpers.Store.Get(r, "vote")
	r.ParseForm()
	entry, err := app.Env.Registry.Entry.GetID(entryId)
	if err != nil {
		pvhelpers.LogErrorObject(err)
		pvhelpers.AddFlash(w, r, pvhelpers.FlashMessage{MsgType: "danger", Msg: "Could not find entry"}, "vote")
		session.Save(r, w)
		http.Redirect(w, r, "/voter/entry", http.StatusFound)
		return
	}
	if r.Method == "POST" {
		entry.Title = strings.TrimSpace(r.FormValue("Title"))
		entry.CategoryIds = strings.Join(r.Form["CategoryIds"], ",")
		entry.ContestId = strings.TrimSpace(r.FormValue("ContestId"))
log.Printf("categoryIds %v", entry.CategoryIds)
		//TODO add validation
		//entry.Id = pvhelpers.GenerateUUID()
		err := app.Env.Registry.Entry.Replace(entry)

		if err != nil {
			pvhelpers.AddFlash(w,r,pvhelpers.FlashMessage{MsgType: "danger", Msg: "Error saving entry"}, "vote")
		} else {
			pvhelpers.AddFlash(w, r, pvhelpers.FlashMessage{MsgType: "success", Msg: "entry created successfully"}, "vote")
			session.Save(r, w)
			http.Redirect(w, r, "/voter/entry", http.StatusFound)
			return
		}
	}

	session.Save(r, w)
	ctx := pongo2.Context{
		"point_to": "entry",
		"model":    EntryLister{Entry: entry, Registry: app.Env.Registry},
		"id":       entry.Id,
	}
	pvhelpers.RenderTemplate(w, r, "templates/entry/create.html", ctx, "voter")
}

func (app EntryApp) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entryId := vars["entryId"]
	session, _ := pvhelpers.Store.Get(r, "vote")
	entry, err := app.Env.Registry.Entry.GetID(entryId)
	if err != nil {
		pvhelpers.LogErrorObject(err)
		pvhelpers.AddFlash(w, r, pvhelpers.FlashMessage{MsgType: "danger", Msg: "Could not find entry"}, "vote")
		session.Save(r, w)
		http.Redirect(w, r, "/voter/entry", http.StatusFound)
		return
	}

	err = app.Env.Registry.Entry.Delete(entry)
	if err != nil {
		pvhelpers.LogErrorObject(err)
		pvhelpers.AddFlash(w, r, pvhelpers.FlashMessage{MsgType: "danger", Msg: "Error deleting entry"}, "vote")
		session.Save(r, w)
		http.Redirect(w, r, "/voter/entry", http.StatusFound)
		return
	}

	session.Save(r, w)
	http.Redirect(w, r, "/voter/entry", http.StatusFound)
}