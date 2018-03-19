package config

//Cfg
var Cfg *Config

func init() {
	Cfg = NewDefaultConfig()
}

//Config
type Config struct {
	ORIGIN   string
	WSURL    string
	RPCURL   string
	DBEngine string
	DBName   string
	DBUser   string
	DBPWD    string
	DBHost   string
	DBPort   string
	DBZone   string
}

func NewDefaultConfig() *Config {
	return &Config{
		ORIGIN:   "http://127.0.0.1:8881/",
		WSURL:    "ws://127.0.0.1:8888/",
		RPCURL:   "http://127.0.0.1:8000",
		DBEngine: "mysql",
		DBName:   "zip_storage",
		DBUser:   "root",
		DBPWD:    "root",
		DBHost:   "127.0.0.1",
		DBPort:   "3306",
		DBZone:   "Asia/Shanghai",
	}
}
