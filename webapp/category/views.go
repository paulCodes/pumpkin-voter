package category

import (
	"github.com/flosch/pongo2"
	"github.com/paulCodes/pumpkin-voter/domain"
	"github.com/paulCodes/pumpkin-voter/pvhelpers"
	"github.com/paulCodes/pumpkin-voter/webtypes"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type CategoryApp struct {
	Env webtypes.Env
}

func (app CategoryApp) Categories(w http.ResponseWriter, r *http.Request) {
	categories, err := app.Env.Registry.Category.All()
	if err != nil {
		panic("An error has occured accessing your contests." + err.Error())
	}

	models := []CategoryLister{}
	for _, category := range categories {
		models = append(models, CategoryLister{
			Category:  category,
			Registry: app.Env.Registry,
		})
	}

	pvhelpers.RenderTemplate(w, r, "templates/category/list.html",
		pongo2.Context{
			"point_to": "category",
			"models":   models,
			"stub":     &CategoryLister{},
		},
		"vote")
}

func (app CategoryApp) Create(w http.ResponseWriter, r *http.Request) {
	session, _ := pvhelpers.Store.Get(r, "vote")
	category := domain.Category{}
	category.Id = "-1"
	if r.Method == "POST" {
		category.Title = strings.TrimSpace(r.FormValue("Title"))
		category.Active = strings.TrimSpace(r.FormValue("Active")) == "true"

		//TODO add validation
		category.Id = pvhelpers.GenerateUUID()
		err := app.Env.Registry.Category.Add(category)

		if err != nil {
			pvhelpers.AddFlash(w,r,pvhelpers.FlashMessage{MsgType: "danger", Msg: "Error saving contest"}, "vote")
		} else {
		pvhelpers.AddFlash(w, r, pvhelpers.FlashMessage{MsgType: "success", Msg: "Contest created successfully"}, "vote")
		session.Save(r, w)
		http.Redirect(w, r, "/voter/category", http.StatusFound)
		return
		}
	}

	session.Save(r, w)
	ctx := pongo2.Context{
		"point_to": "category",
		"model":    CategoryLister{Category: category, Registry: app.Env.Registry},
		"id":       category.Id,
	}
	pvhelpers.RenderTemplate(w, r, "templates/category/create.html", ctx, "voter")
}

func (app CategoryApp) Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryId := vars["categoryId"]
	session, _ := pvhelpers.Store.Get(r, "vote")
	category, err := app.Env.Registry.Category.GetID(categoryId)
	if err != nil {
		pvhelpers.LogErrorObject(err)
		pvhelpers.AddFlash(w, r, pvhelpers.FlashMessage{MsgType: "danger", Msg: "Could not find category"}, "vote")
		session.Save(r, w)
		http.Redirect(w, r, "/voter/category", http.StatusFound)
		return
	}
	if r.Method == "POST" {
		category.Title = strings.TrimSpace(r.FormValue("Title"))
		category.Active = strings.TrimSpace(r.FormValue("Active")) == "true"

		//TODO add validation
		//contest.Id = pvhelpers.GenerateUUID()
		err := app.Env.Registry.Category.Replace(category)

		if err != nil {
			pvhelpers.AddFlash(w,r,pvhelpers.FlashMessage{MsgType: "danger", Msg: "Error saving category"}, "vote")
		} else {
			pvhelpers.AddFlash(w, r, pvhelpers.FlashMessage{MsgType: "success", Msg: "category created successfully"}, "vote")
			session.Save(r, w)
			http.Redirect(w, r, "/voter/category", http.StatusFound)
			return
		}
	}

	session.Save(r, w)
	ctx := pongo2.Context{
		"point_to": "category",
		"model":    CategoryLister{Category: category, Registry: app.Env.Registry},
		"id":       category.Id,
	}
	pvhelpers.RenderTemplate(w, r, "templates/contest/create.html", ctx, "voter")
}

func (app CategoryApp) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryId := vars["categoryId"]
	session, _ := pvhelpers.Store.Get(r, "vote")
	category, err := app.Env.Registry.Category.GetID(categoryId)
	if err != nil {
		pvhelpers.LogErrorObject(err)
		pvhelpers.AddFlash(w, r, pvhelpers.FlashMessage{MsgType: "danger", Msg: "Could not find category"}, "vote")
		session.Save(r, w)
		http.Redirect(w, r, "/voter/category", http.StatusFound)
		return
	}

	err = app.Env.Registry.Category.Delete(category)
	if err != nil {
		pvhelpers.LogErrorObject(err)
		pvhelpers.AddFlash(w, r, pvhelpers.FlashMessage{MsgType: "danger", Msg: "Error deleting category"}, "vote")
		session.Save(r, w)
		http.Redirect(w, r, "/voter/category", http.StatusFound)
		return
	}

	session.Save(r, w)
	http.Redirect(w, r, "/voter/category", http.StatusFound)
}