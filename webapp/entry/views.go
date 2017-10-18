package entry

import (
	"github.com/flosch/pongo2"
	"github.com/paulCodes/pumpkin-voter/domain"
	"net/http"

	"github.com/paulCodes/pumpkin-voter/pvhelpers"
	"github.com/paulCodes/pumpkin-voter/webtypes"
)

type EntryApp struct {
	Env webtypes.Env
}

func (app EntryApp) Entries(w http.ResponseWriter, r *http.Request) {
	entries, err := app.Env.Registry.Entry.Select(&domain.Entry{}, "select * from entry")
	if err != nil {
		panic("An error has occured accessing your entries.")
	}

	ctx := pongo2.Context{
		"entries": entries,
	}

	pvhelpers.RenderTemplate(w, r, "entry/list.html", ctx, "vote")
}

//
//func (app EntryApp) EntryCreate(w http.ResponseWriter, r *http.Request) {
//	var form *gforms.FormInstance
//	var err error
//	var action string
//	vars := mux.Vars(r)
//	action = vars["action"]
//
//	entry := domain.Entry{}
//
//	if entryId, ok := vars["entryId"]; ok {
//		entry, err = app.StoreRegistry.EntryStore.FindById(entryId)
//		if err != nil {
//			app.WebApp.AddErrorAndLog(w, r,"Entry was not found or is not available.", err)
//			http.Redirect(w, r, app.Reverse("entry/list.html"), 302)
//		}
//	}
//
//	form = EntryForm(app.WebApp)(r)
//	httphelpers.PopulateFormInstanceFromModel(form, &entry)
//
//	if r.Method == "POST" {
//		if form.IsValid() {
//			form.MapTo(&entry)
//			if action == "create" {
//				err = app.StoreRegistry.EntryStore.Add(&entry)
//			} else if action == "edit" {
//				err = app.StoreRegistry.EntryStore.Replace(&entry)
//			}
//			if err != nil {
//				app.WebApp.AddErrorAndLog(w,r,"There was an error saving.", err)
//			} else {
//				app.WebApp.AddSuccess(w,r,"the entry was saved successfully")
//				http.Redirect(w, r, app.Reverse("entry/list.html"), 302)
//				return
//			}
//
//		}
//	}
//
//	ctx := pongo2.Context{
//		"entries": entry,
//		"form": form,
//		"action": vars["action"],
//	}
//	httphelpers.RenderTemplate(app.WebApp, w, r, "entry/create.html", ctx)
//}
