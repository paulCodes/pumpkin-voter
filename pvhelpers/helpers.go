package pvhelpers

import (
	"fmt"
	"github.com/flosch/pongo2"
	"net"
	"net/http"
	"os"
	"github.com/twinj/uuid"
)

func RenderTemplate(rw http.ResponseWriter, req *http.Request, tplFilename string, ctx pongo2.Context, sessionGroup string) {
	session, _ := Store.Get(req, sessionGroup)

	if _, ok := ctx["session"]; ok {
		panic("In RenderTemplate, noticed 'session' already exists.  'session' key is not allowed.")
	}
	ctx["session"] = session
	ctx["STATIC"] = Cfg.Url.StaticPath
	ctx["BASEURL"] = Cfg.Url.BaseUrl
	ctx["BASE_API_URL"] = session.Values["BASE_API_URL"]
	ctx["LOGIN_TIMEOUT_MINUTES"] = Cfg.App.LoginTimeoutMinutes
	ctx["ONDEMANDBASEURL"] = Cfg.Url.OnDemandBaseURL
	ctx["FlashesSingleType"] = true
	flashes := session.Flashes()
	ctx["SessionFlashes"] = flashes
	if len(flashes) > 0 {
		msgType := flashes[0].(*FlashMessage).MsgType
		ctx["FlashesSingleTypeType"] = msgType
		for _, x := range flashes[1:] {
			msg := x.(*FlashMessage)
			if msg.MsgType != msgType {
				ctx["FlashesSingleType"] = false
				break
			}
			msg.MsgType = msgType
		}

	}

	tpl, err := pongo2.FromFile(tplFilename)
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

	// TODO: This is currently only here to consume flashes, so saving every time is overkill.
	// Can we make save smarter, and only do it if there were actual flashes to consume?
	session.Save(req, rw)
}

func RenderTemplateFromString(rw http.ResponseWriter, tplString string, ctx pongo2.Context) string {
	tpl, err := pongo2.FromString(tplString)
	out, err := tpl.Execute(ctx)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	return out
}

func AddFlashes(rw http.ResponseWriter, r *http.Request, flashes []FlashMessage, sessionGroup string) {
	reqSession, _ := Store.Get(r, sessionGroup)
	for i := 0; i < len(flashes); i++ {
		reqSession.AddFlash(&flashes[i])
	}
	reqSession.Save(r, rw)
}

func AddFlash(rw http.ResponseWriter, r *http.Request, flash FlashMessage, sessionGroup string) {
	reqSession, _ := Store.Get(r, sessionGroup)
	reqSession.AddFlash(&flash)
	reqSession.Save(r, rw)
}

func PrintServerConnectionDetails() {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		LogError(fmt.Sprintln(err))
		os.Exit(1)
	}

	fmt.Printf("Your server URL will be one of the following:\n\n")
	for _, address := range addrs {

		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Printf("    http://%s%s\n", ipnet.IP.String(), Cfg.Deploy.Listen)
			}

		}
	}
	fmt.Printf("\n")
}

func GenerateUUID() string {
	uuid.SwitchFormat(uuid.FormatCanonical)
	return uuid.NewV4().String()
}