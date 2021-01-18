package config

// APIConfig will hold config vars for API
type APIConfig struct {
	IP   string
	Port string
}

// DBConfig will hold config vars for DB
type DBConfig struct {
	Driver   string
	Filename string
}
