# Comparing Sql vs Sqlx vs Gorm
[![godoc v2](https://img.shields.io/badge/godoc-v2-375EAB.svg)](https://godoc.org/gopkg.in/gorp.v2)
[![GoDoc](https://godoc.org/github.com/jinzhu/gorm?status.svg)](https://godoc.org/github.com/jinzhu/gorm)

### Before reading

First, I want to say thank you to all contributors who are / were contributing to both project.
You guys are awesome, make Golang community much better, and make my work easier.

Second, I start using Golang No 2018, I am still new , and learning. I am create this small test, *NOT to compare which library is better*. This is just
for my person interest. I want to compare which is a better fit to my project. The result doesn't not
apply to any other project.

Third, if you dont agree with my code / comparison , don't hesitate to email me at walterchen.ca@gmail.com
I am happy to make change.


## Introduction

Currently, my project requires a lot of `WRITE` operation to postgres database. At meantime,
the project has a lot of domain models and tables. So an ORM tool is needed. I found top 2 Golng ORM project on
Github, which is `github.com/jinzhu/gorm` and `github.com/jmoiron/sqlx`. but I still want to know which one is a better fit to my project.
I know that Golang reflect is kind slow, and ORM library use a lot of reflection. That might cause
use performance issue. So I create a test injest record thousands of time with multi-threads. And comparing which is a faster.

I am comparing
- database/sql
- github.com/jmoiron/sqlx
- github.com/jinzhu/gorm

## Configuration

- local postgres instance with `max_connections` = 100
- running the test on local Macbook.

## Result

Is seems that `github.com/jmoiron/sqlx` is the fastest, `database/sql` is second, and `github.com/jinzhu/gorm` is the 3rd.


```
GOROOT=/usr/local/go #gosetup
GOPATH=/Users/walter/go #gosetup
/usr/local/go/bin/go test -c -o /private/var/folders/r2/ydgw6d9j1czcvhwrxzvnp_p80000gn/T/___Test_Compare_in_com_cicidi_go_orm_compare_test com.cicidi/go-orm-compare/test #gosetup
/usr/local/go/bin/go tool test2json -t /private/var/folders/r2/ydgw6d9j1czcvhwrxzvnp_p80000gn/T/___Test_Compare_in_com_cicidi_go_orm_compare_test -test.v -test.run ^Test_Compare$ #gosetup
2019/04/09 18:56:28 database.go:49: initGorm
2019/04/09 18:56:28 database.go:68: initSql
2019/04/09 18:56:28 database.go:86: initSqlx
<nil>=== RUN   Test_Compare
2019/04/09 18:56:28 gorm_sqlx_test.go:107:  Insert and read 1000 records
2019/04/09 18:56:28 gorm_sqlx_test.go:129: sql use 226 milliseconds
2019/04/09 18:56:29 gorm_sqlx_test.go:147: sqlx use 170 milliseconds
2019/04/09 18:56:29 gorm_sqlx_test.go:164: gorm use 427 milliseconds
2019/04/09 18:56:29 gorm_sqlx_test.go:165: sqlx.DB spend -24.778761 % (more/less) time than sql.DB
2019/04/09 18:56:29 gorm_sqlx_test.go:166: gorm.DB spend 88.938053 % (more/less) time than sql.DB
2019/04/09 18:56:29 gorm_sqlx_test.go:167: gorm.DB spend 151.176471 % (more/less) time than sqlx.DB
PASS
```


## My experince


- Sql is the official JDBC library. But it need extra work for table, column mapping.
- Sqlx provides a better coding experience, we can do column mapping by using table, and also auto fulfill your "object", however, Sqlx doesn't support

    - create table
    - insert "object" with type as `interface{}`. I need to create diffent `save` func for different type of "object"
- Gorm has the best coding experience, it is like Spring data in Java, you dont need worry about anything. The tradeoff is the performance. Gorm supports

    - create table func
    - insert an instance as  ` interface{} `
    - transactional



## Enhancement
In the end , I choose Sqlx in project, and I did some enhancement in `github.com/jmoiron/sqlx`.
   - CREATE TABLE:  To create table, I use `github.com/go-gorp/gorp` , sqlx is a derived from this project, they share the same `tagName`.
   - ORM wrapper: During project initiation, I load all table name and fieldName into a sync.Map, so the reflection is a one-time job. So far, I don't find any issue of this design.


`GetFieldNames`  create / return all field names of a type
```go
func GetFieldNames(instance interface{}) sync.Map {
	typeName := structs.Name(instance)
	if _, ok := typeFieldWithTag.Load(typeName); !ok {
		t := reflect.TypeOf(instance)
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
```

`GetTableName`  create / return table Names of a type

```go
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

```
