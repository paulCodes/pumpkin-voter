package httphelpers

import (
	"net/http"

	. "git.jcasolutions.com/jca/traininglocker/webapp/webtype"
	"github.com/flosch/pongo2"
)

func RenderTemplate(app WebApp, rw http.ResponseWriter, req *http.Request, tplFilename string, ctx pongo2.Context) {
	session := app.Session(req)

	if _, ok := ctx["_session"]; ok {
		panic("In RenderTemplate, noticed 'session' already exists.  'session' key is not allowed.")
	}
	if _, ok := ctx["Reverse"]; ok {
		panic("In RenderTemplate, noticed 'Reverse' already exists.  'Reverse' key is not allowed.")
	}
	ctx["_session"] = session

	userId, ok := session.Values["user"].(string)
	if ok {
		user, err := app.StoreRegistry.UserStore.FindById(userId)

		if err != nil {
			panic("In RenderTemplate, could not find user by ID.")
		} else {
			ctx["_userRole"] = user.Role
			ctx["_username"] = user.Username
		}
	}

	ctx["Reverse"] = app.Reverse
	ctx["_flashesSingleType"] = ""
	flashes := app.Flashes(req)
	ctx["_flashes"] = flashes
	if len(flashes) > 0 {
		ctx["_flashesSingleType"] = flashes[0].FlashType
		for _, msg := range flashes[1:] {
			if msg.FlashType != flashes[0].FlashType {
				ctx["_flashesSingleType"] = ""
				break
			}
		}
		// Consume the flashes.
		session.Save(req, rw)
	}

	//tpl, err := app.TemplateSet.FromFile(tplFilename)
	tpl, err := app.TemplateSet.FromCache(tplFilename)
	if err != nil {
		// A major template error occurred.  Fall out immediately.
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tpl.ExecuteWriter(ctx, rw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
