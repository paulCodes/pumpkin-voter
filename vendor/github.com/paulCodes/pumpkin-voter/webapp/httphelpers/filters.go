package httphelpers

import (
	"git.jcasolutions.com/jca/traininglocker/webapp/webtype"
	"github.com/flosch/pongo2"
	"time"
)

func CreateFilterAppSensitiveDateTime(app webtype.WebApp) func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
		var t time.Time
		if in.IsNil() {
			return pongo2.AsValue(""), nil
		}
		tPtr, isTime := in.Interface().(*time.Time)
		if !isTime {
			t, isTime = in.Interface().(time.Time)
			if !isTime {
				return nil, &pongo2.Error{
					Sender:   "filter:appdate",
					ErrorMsg: "Filter input argument must be of type 'time.Time' or '*time.Time'.",
				}
			}
		} else {
			t = *tPtr
		}
		//return pongo2.AsValue(t.Format(param.String())), nil
		return pongo2.AsValue(t.Format(app.DateTimeFormat)), nil
	}
}
