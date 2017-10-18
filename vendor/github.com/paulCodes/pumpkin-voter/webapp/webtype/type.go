package webtype

import (
	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"github.com/paulCodes/pumpkin-voter/domain"
)

type WebApp struct {
	AppDir                string
	Router                *mux.Router
	Handler               http.Handler
	TemplateSet           *pongo2.TemplateSet
	SessionStore          sessions.Store
	SessionKey            string
	StoreRegistry         domain.StoreRegistry
	DateTimeFormat        string
	ContentStoragePath    string
	ContentSessionStore   sessions.Store
	ContentSessionKey     string
	PublicContentHttpsUrl string
	PublicContentHttpUrl  string
	Logger                log.Logger
}

func (app WebApp) Session(r *http.Request) *sessions.Session {
	session, _ := app.SessionStore.Get(r, app.SessionKey)
	return session
}
