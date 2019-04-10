package orm

import (
	"com.cicidi/go-orm-compare/config"
	"com.cicidi/go-orm-compare/query"
	"database/sql"
)

// Created by cicidi / cicidi@gmail.com
// Date: 2019-04-01
// Time: 15:09

func Save(data interface{}) (sql.Result, error) {
	db := config.GetSqlx()
	return db.NamedExec(query.InsertQuery(data, query.GetTableName(data)), data)
}

// your table column id should by "id"
func FindById(data interface{}, value string) error {
	err := config.GetSqlx().Get(data, query.SelectAllQuery(data, config.ID), value)
	return err
}
