package model

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/cabezi/zip-storage/config"
	"github.com/cabezi/zip-storage/model/database"
	"github.com/cabezi/zip-storage/model/table"
)

var DB *sql.DB

func init() {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=%s&parseTime=true",
		config.Cfg.DBUser, config.Cfg.DBPWD, config.Cfg.DBHost, config.Cfg.DBPort, config.Cfg.DBName, url.QueryEscape(config.Cfg.DBZone))
	if tdb, err := database.Open(config.Cfg.DBEngine, connStr); err != nil {
		panic(fmt.Sprintf("open database error connStr=%s err=%v", connStr, err))
	} else {
		DB = tdb
	}
	DB.SetMaxOpenConns(2000)
	DB.SetMaxIdleConns(1000)
	if err := DB.Ping(); err != nil {
		panic(err)
	}
	table.InitTables(DB)
}

func ClearDB() {
	table.DropTables(DB)
}
