package query

import (
	"github.com/fatih/structs"
	"github.com/iancoleman/strcase"
	"log"
	"reflect"
	"strings"
	"sync"
)

//Created by cicidi / cicidi@gmail.com
//Date: 2019-03-29
//Time: 14:51

var typeFieldWithTag sync.Map
var tableNameMap sync.Map

func GetFieldNames(instance interface{}) sync.Map {
	typeName := structs.Name(instance)
	if _, ok := typeFieldWithTag.Load(typeName); !ok {
		// this is to get type of a pointer
		t := reflect.TypeOf(instance).Elem()
		pairs := Pairs{}
		for _, fieldName := range structs.Names(instance) {
			var tagName string
			if f, ok := t.FieldByName(string(fieldName)); !ok {
				log.Fatalln("field : {} not found", fieldName)
			} else {
				if tagName, ok = f.Tag.Lookup("db"); !ok {
					// this is following gorp
					// gorp.go line 745
					tagName = f.Name
				}
				pairs = append(pairs, Pair{fieldName, tagName})
			}
		}
		typeFieldWithTag.Store(typeName, pairs)
	}
	return typeFieldWithTag
}

func InsertQuery(instanceType interface{}, tableName string) string {
	var sb strings.Builder
	var secondSb strings.Builder
	sb.WriteString("INSERT INTO " + tableName + " ( ")
	secondSb.WriteString(" VALUES ( ")
	typeName := structs.Name(instanceType)
	GetFieldNames(instanceType)
	x, _ := typeFieldWithTag.Load(typeName)
	pairs := x.(Pairs)
	//pairs := Pairs{}
	for i, pair := range pairs {
		sb.WriteString(pair.TagName)
		if i != len(pairs)-1 {
			sb.WriteString(" , ")
		} else {
			sb.WriteString(" ) ")
		}

		secondSb.WriteString(":" + pair.TagName)
		if i != len(pairs)-1 {
			secondSb.WriteString(" , ")
		} else {
			secondSb.WriteString(" ) ")
		}
	}
	sb.WriteString(secondSb.String())
	return sb.String()
}

func SelectAllQuery(data interface{}, fieldName string) string {
	return "SELECT * FROM " + GetTableName(data) + " WHERE " + fieldName + " =$1"
}

func GetTableName(instance interface{}) string {
	typeName := structs.Name(instance)
	if val, ok := tableNameMap.Load(typeName); !ok {
		tableName := strcase.ToSnake(typeName)
		tableNameMap.Store(typeName, tableName)
		return tableName
	} else {
		return val.(string)
	}
}

type Pair struct {
	FieldName string
	TagName   string
}
type Pairs []Pair
