package pvhelpers

type FlashMessage struct {
	MsgType string
	Msg     string
	FieldId string
}

type AppConfig struct {
	Database struct {
		Host         string
		Username     string
		Password     string
		RunMigration bool
		DB           string
		Debug        bool
	}
	Url struct {
		StaticPath      string
		BaseUrl         string
		OnDemandBaseURL string
		PrepApiBaseURL  string
		RegApiBaseURL   string
	}
	App struct {
		LogLevel            string // Verbose, Info, Warn, Error
		LoginTimeoutMinutes int64
	}
	Deploy struct {
		Listen string
	}
}

type DbIndices int
