package tlserve

import (
	"database/sql"
	"encoding/gob"
	"errors"
	"fmt"
	"git.jcasolutions.com/jca/traininglocker/domain"
	"git.jcasolutions.com/jca/traininglocker/interface/relationaldb"
	"git.jcasolutions.com/jca/traininglocker/payment"
	"git.jcasolutions.com/jca/traininglocker/webapp"
	"git.jcasolutions.com/jca/traininglocker/webapp/webtype"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/flosch/pongo2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/stripe/stripe-go"
	"gopkg.in/gcfg.v1"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var cfg AppConfig
var logFileHandle *os.File
var loggerHandle *log.Logger

func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		loggerHandle.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}
	return http.HandlerFunc(fn)
}

func checkConfig(cfg AppConfig) bool {
	if cfg.Deploy.PublicContentHttpsUrl == "" ||
		!strings.HasPrefix(cfg.Deploy.PublicContentHttpsUrl, "https://") ||
		strings.HasSuffix(cfg.Deploy.PublicContentHttpsUrl, "/") {
		println("publicContentHttpsUrl must be non-empty and of the form: https://<domain>(:port) .  No trailing slash.")
		return false
	}
	if cfg.Deploy.PublicContentHttpUrl == "" ||
		!strings.HasPrefix(cfg.Deploy.PublicContentHttpUrl, "http://") ||
		strings.HasSuffix(cfg.Deploy.PublicContentHttpUrl, "/") {
		println("publicContentHttpUrl must be non-empty and of the form: http://<domain>(:port) .  No trailing slash.")
		return false
	}

	binaryExists := func(path string) bool {
		_, err := os.Stat(path)
		return !os.IsNotExist(err)
	}

	if cfg.System.ZipBinary == "" ||
		!binaryExists(cfg.System.ZipBinary) {
		println("ZipBinary must be non-empty and must exist.")
		return false
	}

	if cfg.System.ZipFixType == "" {
		println("ZipFixType must be non-empty.")
		return false
	}

	return true
}

func LoadConfig(configDir string) (cfg AppConfig, err error) {
	err = gcfg.ReadFileInto(&cfg, filepath.Join(configDir, "config.gcfg"))
	if err != nil {
		println("Configuration file (config.gcfg) could not be found in directory: ("+configDir+").", err.Error())
		return
	}

	configValid := checkConfig(cfg)
	if !configValid {
		err = errors.New("Configuration invalid, exiting.")
		println(err.Error())
		return
	}
	return
}

func initStore(cfg AppConfig, storeLogger *log.Logger) (storeRegistry domain.StoreRegistry, err error) {
	var db *sql.DB
	if cfg.Database.DBType == "mysql" {
		sqlConn := cfg.Database.Username + ":" + cfg.Database.Password + "@" + cfg.Database.Host + "/" + cfg.Database.DB + "?parseTime=true&loc=Local&timeout=5s"
		db, err = sql.Open("mysql", sqlConn)
	} else if cfg.Database.DBType == "mssql" {
		sqlConn := "server=" + cfg.Database.Host + ";user id=" + cfg.Database.Username + ";password=" + cfg.Database.Password + ";database=" + cfg.Database.DB
		db, err = sql.Open("mssql", sqlConn)
	} else {
		err = errors.New("Unrecognized DB type: " + cfg.Database.DBType)
		return
	}

	if err != nil {
		println("Unable to instantiate database connection.  Details: ", err.Error())
		println("The server has failed to start.  Please check your database settings and try again.")
		return
	}

	// Actually test if the DB is there.
	err = db.Ping()
	if err != nil {
		println("Unable to open database connection.  Details: ", err.Error())
		println("The server has failed to start.  Please check your database settings and try again.")
		return
	}
	wrappedDB := relationaldb.WrapDB(db, storeLogger, cfg.Database.Debug)

	rel := relationaldb.NewRelationalStore(wrappedDB, cfg.Database.DBType)
	storeRegistry = domain.StoreRegistry{
		OrganizationStore:     rel.OrganizationStore(),
		AccessTokenStore:      rel.AccessTokenStore(),
		CourseStore:           rel.CourseStore(),
		LicenseStore:          rel.LicenseStore(),
		LearnerScormDataStore: rel.LearnerScormDataStore(),
		PaymentStore:          rel.PaymentStore(),
		PlanStore:             rel.PlanStore(),
		StripeCustomerStore:   rel.StripeCustomerStore(),
		StubStore:             rel.StubStore(),
		UsageStore:            rel.UsageStore(),
		UserStore:             rel.UserStore(),
	}
	return
}

func MakeLogger(cfg AppConfig) (makeLogger *log.Logger, err error) {
	prefix := "wlog: "
	logFlags := log.Ldate | log.Ltime | log.Lmicroseconds
	if cfg.Deploy.LogFile == "" {
		println("no logfile")
		makeLogger = log.New(os.Stdout, prefix, logFlags)
	} else {
		println("logfile is ", cfg.Deploy.LogFile)
		if logFileHandle == nil {
			logFileHandle, err = os.OpenFile(cfg.Deploy.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				println("Log file could not be opened, exiting.")
				return
			}
		}
		makeLogger = log.New(logFileHandle, prefix, logFlags)
		//log.SetOutput(logFileHandle)
	}
	return
}

func Serve(appDir string) {
	var err error
	cfg, err = LoadConfig(appDir)
	if err != nil {
		println("Error occurred, exiting: " + err.Error())
		return
	}
	loggerHandle, err = MakeLogger(cfg)
	if err != nil {
		println("Could not open log file: ", err.Error())
		return
	}
	loggerHandle.Println("serve() called...")

	storeRegistry, err := initStore(cfg, loggerHandle)
	if err != nil {
		loggerHandle.Println("Error occurred instantiating store: " + err.Error())
		return
	}

	//sessionStore, err := mysqlstore.NewMySQLStore(sqlConn, "session_app", "/", 86400, []byte("CiHecFeDotivarAyWeoWhouj0On"))
	//if err != nil {
	//	println("Error creating sessionStore: ", err.Error())
	//}
	sessionStore := sessions.NewFilesystemStore("", []byte("CiHecFeDotivarAyWeoWhouj0On"))

	gob.Register(webtype.FlashMessage{})
	var mainSessionKey = "traininglockerSessionKey"
	localTemplateLoader, err := pongo2.NewLocalFileSystemLoader(filepath.Join(appDir, "templates/"))
	if err != nil {
		loggerHandle.Println("Local template loader failed to open: " + err.Error())
		return
	}
	templateSet := pongo2.NewSet("main template set", localTemplateLoader)
	templateSet.Debug = true
	contentStoragePath := filepath.Join(appDir, "..", "_content/")

	//var sessionContentStore, _ = mysqlstore.NewMySQLStore(sqlConn, "session_content", "/", 86400, []byte("dografEchnedwuadEydirenthIsDesumJu"))
	sessionContentStore := sessions.NewFilesystemStore("", []byte("CiHecFeDotivarAyWeoWhouj0On"))

	gob.Register(webtype.FlashMessage{})
	var contentSessionKey = "traininglockerContentKey"

	var stripeConfig = payment.StripeConfig{
		PublishableKey: cfg.Payment.StripePublishableKey,
		SecretKey:      cfg.Payment.StripeSecretKey,
	}
	stripe.Key = cfg.Payment.StripeSecretKey

	app := webapp.NewWebApp(appDir, storeRegistry, templateSet, sessionStore, mainSessionKey,
		contentStoragePath, sessionContentStore, contentSessionKey, stripeConfig, cfg.System.ZipFixType,
		cfg.System.ZipBinary, cfg.Deploy.PublicContentHttpsUrl, cfg.Deploy.PublicContentHttpUrl, *loggerHandle)
	loggerHandle.Println("Serving requests on " + cfg.Deploy.Listen)
	log.Fatalf("%v", http.ListenAndServe(cfg.Deploy.Listen, RecoverWrap(loggingHandler(app))))
}

func RecoverWrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				w.WriteHeader(http.StatusInternalServerError)
				stack := make([]byte, 1024*8)
				stack = stack[:runtime.Stack(stack, false)]

				f := "PANIC: %s\n\n%s"
				log.Printf(f, err, stack)

				fmt.Fprintf(w, f, err, stack)
				fmt.Printf(f, err, stack)

				//sendMeMail(err)
				//http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}
