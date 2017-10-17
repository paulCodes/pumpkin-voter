package webapp

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/flosch/pongo2"
	"github.com/sebest/xff"
	"log"
	"path/filepath"
	"github.com/paulCodes/pumpkin-voter/domain"
	"github.com/paulCodes/pumpkin-voter/webapp/webtype"
	"github.com/paulCodes/pumpkin-voter/webapp/httphelpers"
	"github.com/paulCodes/pumpkin-voter/webapp/entry"
)

type ourWebApp struct {
	webtype.WebApp
}

func (app ourWebApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.Handler.ServeHTTP(w, r)
}

func RegisterFilters(app ourWebApp) {
	pongo2.RegisterFilter("appdatetime", httphelpers.CreateFilterAppSensitiveDateTime(app.WebApp))
}

func makeGetPostRoute(router *mux.Router, name string, path string, fn http.HandlerFunc) {
	router.Handle(path, handlers.MethodHandler{
		"GET":  fn,
		"POST": fn,
	}).Name(name)
}

func makeGetPostRouteForHandle(router *mux.Router, name string, path string, fn http.Handler) {
	router.Handle(path, handlers.MethodHandler{
		"GET":  fn,
		"POST": fn,
	}).Name(name)
}

func makeCreateEditRoutes(router *mux.Router, name string, create string, edit string, fn http.HandlerFunc) {
	makeGetPostRoute(router, name+"_create", create, fn)
	makeGetPostRoute(router, name, edit, fn)
}

func makeCreateEditRoutesForHandle(router *mux.Router, name string, create string, edit string, fn http.Handler) {
	makeGetPostRouteForHandle(router, name+"_create", create, fn)
	makeGetPostRouteForHandle(router, name, edit, fn)
}

func adminRoutes(app ourWebApp, adminR *mux.Router) {
	//adminR.Handle("/", handlers.MethodHandler{"GET": AdminRootHandler{app}}).Name("admin_root")
	//
	//courseApp := courses.CourseApp(app)
	////adminR.Handle("/courses", handlers.MethodHandler{"GET": http.HandlerFunc(courseApp.Courses)}).Name("courses")
	//adminR.Handle("/courses_by_cat/root", handlers.MethodHandler{"GET": http.HandlerFunc(courseApp.CoursesByCategory)}).Name("courses")
	//adminR.Handle("/courses_by_cat/{categoryId}", handlers.MethodHandler{"GET": http.HandlerFunc(courseApp.CoursesByCategory)}).Name("courses_by_category")
	//
	//adminR.Handle("/courses/choose_new_type/{categoryId}", handlers.MethodHandler{"GET": http.HandlerFunc(courseApp.CreateCourseTypeChoice)}).Name("courses_choose_type")
	//adminR.Handle("/courses/create/upload/{categoryId}", handlers.MethodHandler{"POST": http.HandlerFunc(courseApp.CreateCourseTypeLocalUpload)}).Name("course_create_upload")
	//adminR.Handle("/courses/course_delete/{courseId}", handlers.MethodHandler{"GET": http.HandlerFunc(courseApp.CourseDelete)}).Name("course_delete")
	//adminR.Handle("/courses/student_data/{courseId}", handlers.MethodHandler{"GET": http.HandlerFunc(courseApp.CourseStudentDataView)}).Name("course_student_data")
	//makeGetPostRoute(adminR, "course_create_local", "/courses/create/{serveType:local}", courseApp.CourseCreate)
	//makeGetPostRoute(adminR, "course_create_remote", "/courses/create/{serveType:remote}", courseApp.CourseCreate)
	//makeGetPostRoute(adminR, "course_edit", "/courses/{action:edit}/{courseId}", courseApp.CourseEdit)
	//adminR.Handle("/courses/view/{courseId}", handlers.MethodHandler{"GET": http.HandlerFunc(courseApp.CourseView)}).Name("course_view")
	//adminR.Handle("/courses/add_version/local/{courseId}", handlers.MethodHandler{"POST": http.HandlerFunc(courseApp.CourseAddVersionLocal)}).Name("course_add_version_local")
	//// TODO Make POST only.
	//adminR.Handle("/courses/course_delete_version/{versionId}", handlers.MethodHandler{"GET": http.HandlerFunc(courseApp.CourseDeleteVersion)}).Name("course_delete_version")
	//// TODO Make POST only.
	//adminR.Handle("/courses/make_primary_version/{versionId}", handlers.MethodHandler{"GET": http.HandlerFunc(courseApp.CourseMakePrimaryVersion)}).Name("course_make_primary")
	//adminR.Handle("/courses/download_stub/{stubId}", handlers.MethodHandler{"GET": http.HandlerFunc(courseApp.DownloadStub)}).Name("course_download_stub")
	//// TODO Make POST only.
	//makeCreateEditRoutes(adminR, "course_stub", "/courses/stubs/{action:create}/{courseId}", "/courses/stubs/{action:edit}/{courseId}/{stubId}", courseApp.StubCreateOrEdit)
	//adminR.Handle("/courses/delete_stub/{stubId}", handlers.MethodHandler{"GET": http.HandlerFunc(courseApp.DeleteStub)}).Name("course_delete_stub")
	//
	//makeCreateEditRoutes(adminR, "category", "/courses/category/{action:create}/{parentId}", "/courses/category/{action:edit}/{categoryId}", courseApp.CategoryCreateOrEdit)
	//adminR.Handle("/courses/categories/{id}", handlers.MethodHandler{"GET": http.HandlerFunc(courseApp.CategoryDisplay)}).Name("list_categories")
	//adminR.Handle("/courses/category_delete/{categoryId}", handlers.MethodHandler{"GET": http.HandlerFunc(courseApp.DeleteCategory)}).Name("category_delete")
	//adminR.Handle("/courses/categories/{id}/members", handlers.MethodHandler{
	//	"GET":  http.HandlerFunc(courseApp.CategoryMembers),
	//	"POST": http.HandlerFunc(courseApp.CategoryMembers),
	//}).Name("category_members")
	//adminR.Handle("/courses/categories/{categoryId}/members/{userId}/remove", handlers.MethodHandler{"GET": http.HandlerFunc(courseApp.RemoveCategoryMember)}).Name("category_member_remove")
	//
	//reportApp := reports.ReportApp(app)
	//adminR.Handle("/reports", handlers.MethodHandler{"GET": http.HandlerFunc(reportApp.Index)}).Name("reports")
	//adminR.Handle("/reports/usages_by_tag", handlers.MethodHandler{"GET": http.HandlerFunc(reportApp.UsagesByTag)}).Name("reports_usages_by_tag")
	//adminR.Handle("/reports/usages_by_grouping", handlers.MethodHandler{"GET": http.HandlerFunc(reportApp.UsagesByGrouping)}).Name("reports_usages_by_grouping")
	//
	//accountApp := account.AccountApp(app)
	//rcm := RoleCheckMiddleware{"admin", http.HandlerFunc(accountApp.SelectPlan), app}
	//makeGetPostRouteForHandle(adminR, "admin_account_plan", "/account/plan", rcm)
	//
	//organizationApp := organization.OrganizationApp(app)
	//rcmOrganization := RoleCheckMiddleware{"admin", http.HandlerFunc(organizationApp.Organization), app}
	//adminR.Handle("/organization", handlers.MethodHandler{"GET": rcmOrganization}).Name("organization")
	//
	//batchApp := batch.BatchApp(app)
	//adminR.Handle("/batch", handlers.MethodHandler{"GET": http.HandlerFunc(batchApp.View)}).Name("stub_batch_view")
	//makeGetPostRoute(adminR, "create_batch", "/batch/create_batch", batchApp.BatchCreateForm)
	//adminR.Handle("/batch/download_batch/{batchId}", handlers.MethodHandler{"GET": http.HandlerFunc(batchApp.DownloadBatch)}).Name("download_batch")
	//adminR.Handle("/batch/delete_batch/{batchId}", handlers.MethodHandler{"GET": http.HandlerFunc(batchApp.DeleteBatch)}).Name("delete_batch")
	//
	//usageApp := usage.UsageApp(app)
	//adminR.Handle("/usage",handlers.MethodHandler{"GET": http.HandlerFunc(usageApp.Index)}).Name("usage")
	//makeGetPostRoute(adminR, "usage_datatables","/v1/usage/datatables", usageApp.UsagesForDataTables)
	//
	//makeGetPostRoute(adminR, "admin_index", "/dashboard", app.AdminIndex)
}

func NewWebApp(appDir string, storeRegistry domain.StoreRegistry, templateSet *pongo2.TemplateSet,
	sessionStore sessions.Store, sessionKey string, contentStoragePath string,
	contentSessionStore sessions.Store, contentSessionKey string, publicContentHttpsUrl string, publicContentHttpUrl string,
	logger log.Logger) ourWebApp {
	router := mux.NewRouter()
	router.StrictSlash(true)
	smx := http.NewServeMux()
	xff, err := xff.New(xff.Options{Debug: true})
	if err != nil {
		panic(err.Error())
	}

	app := ourWebApp{
		webtype.WebApp{
			AppDir:                appDir,
			Router:                router,
			Handler:               xff.Handler(smx),
			TemplateSet:           templateSet,
			SessionStore:          sessionStore,
			SessionKey:            sessionKey,
			StoreRegistry:         storeRegistry,
			DateTimeFormat:        "2006-01-02 15:04:05",
			ContentStoragePath:    contentStoragePath,
			ContentSessionStore:   contentSessionStore,
			ContentSessionKey:     contentSessionKey,
			PublicContentHttpsUrl: publicContentHttpsUrl,
			PublicContentHttpUrl:  publicContentHttpUrl,
			Logger:                logger,
		},
	}
	RegisterFilters(app)

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(appDir, "static/")))))

	entryApp := entry.EntryApp(app)
	makeGetPostRoute(router, "entry", "/entry/list", entryApp.Entries)
	makeGetPostRoute(router, "entry", "/entry/create", entryApp.EntryCreate)


	smx.Handle("/", router)

	return app
}

type RoleCheckMiddleware struct {
	RequiresRole string
	http.Handler
	app ourWebApp
}

func (mh RoleCheckMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO Put role check here.

	mh.Handler.ServeHTTP(w, r)
}
