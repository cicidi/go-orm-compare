package config

//Created by cicidi / cicidi@gmail.com
//Date: 2019-04-06
//Time: 11:09

import (
	"com.cicidi/go-orm-compare/model"
	"com.cicidi/go-orm-compare/query"
	"com.cicidi/go-orm-compare/util"
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gopkg.in/gorp.v1"
	"log"
	"time"
)

var gormConnection *gorm.DB
var sqlConnection *sql.DB
var sqlxConnection *sqlx.DB

//db config parameter
var User = "postgres"
var Password = "password"
var Url = "localhost"
var Port = "5432"
var Dialect = "postgres"
var Database = "postgres"
var JdbcPattern = "%s://%s:%s@%s:%s/%s?sslmode=disable"
var SqlDriverName = "org.postgresql.Driver"

var maxConnection = 10
var maxIdle = 10
var maxLifetime = time.Duration(0)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	gormConnection = initGorm()
	sqlConnection = initSql()
	sqlxConnection = initSqlx()
	createTable()
}

func initGorm() *gorm.DB {

	log.Println("initGorm")

	if gormConnection != nil {
		return gormConnection
	}

	db, err := gorm.Open(Dialect, fmt.Sprintf(JdbcPattern, Dialect, User, Password, Url, Port, Database))

	util.CheckErr(err, "gorp open failed")

	db.DB().SetMaxOpenConns(maxConnection)
	db.DB().SetMaxIdleConns(maxIdle)
	db.DB().SetConnMaxLifetime(maxLifetime)
	db.AutoMigrate(&model.People{})
	gormConnection = db
	return gormConnection
}

func initSql() *sql.DB {
	log.Println("initSql")

	if sqlConnection != nil {
		return sqlConnection
	}
	// has to use sql + gorm open here
	db, err := sql.Open(Dialect, fmt.Sprintf(JdbcPattern, Dialect, User, Password, Url, Port, Database))

	util.CheckErr(err, "sql.Open failed")

	db.SetMaxOpenConns(maxConnection) // The default is 0 (unlimited)
	db.SetMaxIdleConns(maxIdle)       // defaultMaxIdleConns = 2
	db.SetConnMaxLifetime(maxLifetime)
	sqlConnection = db
	return sqlConnection
}

func initSqlx() *sqlx.DB {
	log.Println("initSqlx")
	if sqlxConnection != nil {
		return sqlxConnection
	}
	db, err := sqlx.Connect(Dialect, fmt.Sprintf(JdbcPattern, Dialect, User, Password, Url, Port, Database))

	util.CheckErr(err, "sqlx open failed")
	db.SetMaxOpenConns(maxConnection) // The default is 0 (unlimited)
	db.SetMaxIdleConns(maxIdle)       // defaultMaxIdleConns = 2
	db.SetConnMaxLifetime(maxLifetime)
	sqlxConnection = db
	return sqlxConnection
}

// I dont like either the way how gorp or gorm create table
//  gorp  use field create to create folumn name without understore which is not postgres standard
//  gorm create table name  by default add "s" to it
func createTable() {
	// construct a gorp DbMap
	if sqlConnection == nil {
		initSql()
	}
	dbmap := &gorp.DbMap{Db: sqlConnection, Dialect: gorp.PostgresDialect{}}

	// since this is only one time job, don't keep the connection
	//defer sqlConnection.Close()
	//add table with name
	dbmap.AddTableWithName(model.People{}, query.GetTableName(model.People{})).SetKeys(false, ID)

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err := dbmap.CreateTablesIfNotExists()
	fmt.Print(err)
	util.CheckErr(err, "Create tables failed")
}

func GetGorm() *gorm.DB {
	if gormConnection == nil {
		initGorm()
	}
	return gormConnection
}

func GetSqlx() *sqlx.DB {
	if sqlxConnection == nil {
		initSqlx()
	}
	return sqlxConnection
}

func GetSql() *sql.DB {
	if sqlConnection == nil {
		sqlConnection = initSql()
	}
	return sqlConnection
}
