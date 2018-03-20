package config

//Cfg
var Cfg *Config

func init() {
	Cfg = NewDefaultConfig()
}

//Config
type Config struct {
	ORIGIN   string `yaml:"origin"`
	WSURL    string `yaml:"wsurl"`
	RPCURL   string `yaml:"rpcurl"`
	DBEngine string `yaml:"engine"`
	DBName   string `yaml:"name"`
	DBUser   string `yaml:"user"`
	DBPWD    string `yaml:"password"`
	DBHost   string `yaml:"host"`
	DBPort   string `yaml:"port"`
	DBZone   string `yaml:"zone"`
}

func NewDefaultConfig() *Config {
	return &Config{
		ORIGIN:   "http://127.0.0.1:8888/",
		WSURL:    "ws://127.0.0.1:8888/",
		RPCURL:   "http://127.0.0.1:8881",
		DBEngine: "mysql",
		DBName:   "zip_storage",
		DBUser:   "root",
		DBPWD:    "root",
		DBHost:   "127.0.0.1",
		DBPort:   "3306",
		DBZone:   "Asia/Shanghai",
	}
}
