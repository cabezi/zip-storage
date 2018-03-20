package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	mysql := NewMySQL()
	Register(mysql.Name(), mysql)
}

//MySQL  mysql
type MySQL struct {
}

//Name return database driver name
func (mysql *MySQL) Name() string {
	return "mysql"
}

//Open return opened database
func (mysql *MySQL) Open(connStr string) (*sql.DB, error) {
	return sql.Open(mysql.Name(), connStr)
}

//NewMySQL return a MySQL object
func NewMySQL() *MySQL {
	return &MySQL{}
}
