package main

import (
	_ "embed"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/wailsapp/wails"
)

func basic(word string) string {
	return "World!" + " " + word
}

// go:embed frontend-v12/dist/frontend-v12/main.js
var js string

//go:embed frontend-v12/dist/frontend-v12/styles.css
var css string

//DB Connect
var db *sqlx.DB
var tables []string

type MySQLConnectionEnv struct {
	Host     string
	Port     string
	User     string
	DBName   string
	Password string
}

type MySQLColumn struct {
	Name string
	Type string
}

func newMysqlDbConnection(host string, port string, dbName string, username string, password string) *MySQLConnectionEnv {
	return &MySQLConnectionEnv{
		Host:     host,
		Port:     port,
		User:     username,
		DBName:   dbName,
		Password: password,
	}
}

func (mc *MySQLConnectionEnv) ConnectDB() (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", mc.User, mc.Password, mc.Host, mc.Port, mc.DBName)
	return sqlx.Open("mysql", dsn)
}

func connectToDatabase(host string, port string, dbName string, username string, password string) {
	var err error
	db, err = newMysqlDbConnection(host, port, dbName, username, password).ConnectDB()
	if err != nil {
		fmt.Println("DB connection failed : ", err)
	}
	res, err := db.Query(`SHOW TABLES`)
	if err != nil {
		fmt.Println("DB connection failed : ", err)
	}

	tables = []string{}
	for res.Next() {
		var tableName string
		res.Scan(&tableName)
		tables = append(tables, tableName)
	}
}

func tablesNames() []string {
	return tables
}

func columnsNames(tableName string) []string {
	res, err := db.Query(`SELECT * FROM ` + tables[0] + ` LIMIT 1`)
	fmt.Println(err)
	columnsName, _ := res.Columns()
	columnsType, _ := res.ColumnTypes()
	columns := []string{}
	for i, t := range columnsType {
		json_s, _ := json.Marshal(MySQLColumn{Name: columnsName[i], Type: t.DatabaseTypeName()})
		columns = append(columns, string(json_s))
	}
	return columns
}

func main() {
	app := wails.CreateApp(&wails.AppConfig{
		Width:  1024,
		Height: 768,
		Title:  "MySQL-auto-generate",
		JS:     js,
		CSS:    css,
		Colour: "#131313",
	})
	app.Bind(basic)
	app.Bind(connectToDatabase)
	app.Bind(tablesNames)
	app.Bind(columnsNames)
	app.Run()
}
