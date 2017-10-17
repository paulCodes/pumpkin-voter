package tlserve

type AppConfig struct {
	Database struct {
		DBType   string
		Host     string
		Username string
		Password string
		DB       string
		Debug    bool
	}
	Deploy struct {
		Listen                string
		PublicContentHttpsUrl string
		PublicContentHttpUrl  string
		LogFile               string
	}
	Payment struct {
		StripePublishableKey string
		StripeSecretKey      string
	}
	System struct {
		ZipFixType string
		ZipBinary  string
	}
}
