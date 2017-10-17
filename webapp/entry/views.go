package entry

import (
	"github.com/paulCodes/pumpkin-voter/webapp/webtype"
	"net/http"
	"github.com/flosch/pongo2"
	"github.com/paulCodes/pumpkin-voter/webapp/httphelpers"
	"github.com/bluele/gforms"
	"github.com/gorilla/mux"
	"github.com/paulCodes/pumpkin-voter/domain"
)

type EntryApp struct {
	webtype.WebApp
}

func (app EntryApp) Entries(w http.ResponseWriter, r *http.Request) {
	entries, err := app.StoreRegistry.EntryStore.ListAll()
	if err !=nil {
		app.WebApp.AddErrorAndLog(w,r,"An error has occured accessing your entries.", err)
		http.Redirect(w, r, app.Reverse("index"), 302)
	}

	ctx := pongo2.Context{
		"entries" : entries,
	}

	httphelpers.RenderTemplate(app.WebApp, w, r, "entry/list.html", ctx)
}

func (app EntryApp) EntryCreate(w http.ResponseWriter, r *http.Request) {
	var form *gforms.FormInstance
	var err error
	var action string
	vars := mux.Vars(r)
	action = vars["action"]

	entry := domain.Entry{}

	if entryId, ok := vars["entryId"]; ok {
		entry, err = app.StoreRegistry.EntryStore.FindById(entryId)
		if err != nil {
			app.WebApp.AddErrorAndLog(w, r,"Entry was not found or is not available.", err)
			http.Redirect(w, r, app.Reverse("entry/list.html"), 302)
		}
	}

	form = EntryForm(app.WebApp)(r)
	httphelpers.PopulateFormInstanceFromModel(form, &entry)

	if r.Method == "POST" {
		if form.IsValid() {
			form.MapTo(&entry)
			if action == "create" {
				err = app.StoreRegistry.EntryStore.Add(&entry)
			} else if action == "edit" {
				err = app.StoreRegistry.EntryStore.Replace(&entry)
			}
			if err != nil {
				app.WebApp.AddErrorAndLog(w,r,"There was an error saving.", err)
			} else {
				app.WebApp.AddSuccess(w,r,"the entry was saved successfully")
				http.Redirect(w, r, app.Reverse("entry/list.html"), 302)
				return
			}

		}
	}

	ctx := pongo2.Context{
		"entries": entry,
		"form": form,
		"action": vars["action"],
	}
	httphelpers.RenderTemplate(app.WebApp, w, r, "entry/create.html", ctx)
}