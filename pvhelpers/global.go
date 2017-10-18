package pvhelpers

import (
	"github.com/gorilla/sessions"
	"gopkg.in/gorp.v1"
)

var Cfg AppConfig
var EncryptionKey []byte
var Store = sessions.NewFilesystemStore("", []byte("byurlAnBegejAwkyoiCekOdFert"))
var GorpDb *gorp.DbMap
