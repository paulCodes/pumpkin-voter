package webtype

import "net/http"

const WebAppFlashKey = "_webappflash"

type FlashType string

const (
	FLASH_ERROR   FlashType = "danger"
	FLASH_SUCCESS           = "success"
	FLASH_WARNING           = "warning"
	FLASH_INFO              = "info"
)

type FlashMessage struct {
	FlashType FlashType
	Msg       string
}

func (app WebApp) AddFlash(w http.ResponseWriter, r *http.Request, flashType FlashType, msg string) {
	session := app.Session(r)
	session.AddFlash(FlashMessage{FlashType: flashType, Msg: msg}, WebAppFlashKey)
	err := session.Save(r, w)
	if err != nil {
		println("Add flash error: ", err.Error())
		println("msg: ", msg)
		panic("Issue saving session.")
	}
}

func (app WebApp) AddError(w http.ResponseWriter, r *http.Request, msg string) {
	app.AddFlash(w, r, FLASH_ERROR, msg)
}

func (app WebApp) AddErrorAndLog(w http.ResponseWriter, r *http.Request, msg string, err error) {
	app.AddFlash(w, r, FLASH_ERROR, msg)
	app.Logger.Output(1, err.Error())
}

func (app WebApp) AddSuccess(w http.ResponseWriter, r *http.Request, msg string) {
	app.AddFlash(w, r, FLASH_SUCCESS, msg)
}

func (app WebApp) Flashes(r *http.Request) []FlashMessage {
	tmpFlashes := app.Session(r).Flashes(WebAppFlashKey)
	flashMsgs := make([]FlashMessage, len(tmpFlashes))
	for i, flash := range tmpFlashes {
		flashMsgs[i] = flash.(FlashMessage)
	}
	return flashMsgs
}
