package test

import (
	"com.cicidi/go-orm-compare/config"
	"com.cicidi/go-orm-compare/model"
	"com.cicidi/go-orm-compare/orm"
	"com.cicidi/go-orm-compare/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"sync"
	"testing"
	"time"
)

func createPeopleList(total int) []model.People {
	list := make([]model.People, 0)
	for i := 0; i < total; i++ {
		id, _ := uuid.NewRandom()
		people := model.People{
			Id:               id,
			CreatedTimestamp: currentMilliseconds(),
			IsStudent:        true,
			Address:          "Santa_Clara",
			Gender:           "male",
			Race:             "Asian"}
		list = append(list, people)
	}
	return list
}

func saveBySql(t *testing.T, wg *sync.WaitGroup, people *model.People) {
	defer wg.Done()

	//save

	stmt, err := config.GetSql().Prepare("INSERT INTO people(id,created_timestamp,is_student,address,gender,race) VALUES($1, $2, $3, $4, $5, $6)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(people.Id, people.CreatedTimestamp, people.IsStudent, people.Address, people.Gender, people.Race)
	//lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	//rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("affected = %d\n", rowCnt)
	util.CheckErr(err, "save by sql failed")

	//read
	rows, err := config.GetSql().Query("select * from people where id = $1", people.Id.String())
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var savedPeople model.People
	for rows.Next() {
		err := rows.Scan(&savedPeople.Id, &savedPeople.CreatedTimestamp, &savedPeople.IsStudent, &savedPeople.Address, &savedPeople.Gender, &savedPeople.Race)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	//compare
	assert.Equal(t, people.Id.String(), savedPeople.Id.String())

}
func saveBySqlx(t *testing.T, wg *sync.WaitGroup, people *model.People) {
	defer wg.Done()

	//save
	_, err := orm.Save(*people)
	util.CheckErr(err, "save by sqlx failed")
	//read
	var savedPeople model.People
	orm.FindById(&savedPeople, people.Id.String())
	//compare
	assert.Equal(t, people.Id.String(), savedPeople.Id.String())

}
func saveByGorm(t *testing.T, wg *sync.WaitGroup, people *model.People) {
	defer wg.Done()

	//save
	// gorm is by default use tableName as {objectName}s
	config.GetGorm().Save(&people)

	//read
	var savedPeople model.People
	config.GetGorm().Where("id = ?", people.Id.String()).First(&savedPeople)

	//compare
	assert.Equal(t, people.Id.String(), savedPeople.Id.String(), "save by gorm")

}

func Test_Compare(t *testing.T) {

	//totally times of insert
	var total = 1000
	log.Printf(" Insert and read %d records", total)
	list_sql := createPeopleList(total)
	list_sqlx := createPeopleList(total)
	list_gorm := createPeopleList(total)

	//sqlx
	sqlStart := currentMilliseconds()
	sqlDone := make(chan bool, total)
	var wg sync.WaitGroup

	wg.Add(total)
	for i := 0; i < total; i++ {
		go func(i int) {
			saveBySql(t, &wg, &list_sql[i])
			sqlDone <- true
		}(i)
	}
	wg.Wait()
	close(sqlDone)

	sqlEnd := currentMilliseconds()
	var durationSql = sqlEnd - sqlStart
	log.Printf("sql use %d milliseconds ", durationSql)

	//sqlx
	sqlxStart := currentMilliseconds()
	sqlxDone := make(chan bool, total)

	wg.Add(total)
	for i := 0; i < total; i++ {
		go func(i int) {
			saveBySqlx(t, &wg, &list_sqlx[i])
			sqlxDone <- true
		}(i)
	}
	wg.Wait()
	close(sqlxDone)

	sqlxEnd := currentMilliseconds()
	var durationSqlx = sqlxEnd - sqlxStart
	log.Printf("sqlx use %d milliseconds ", durationSqlx)

	//Grom
	gormDone := make(chan bool, total)
	gormStart := currentMilliseconds()
	wg.Add(total)
	for i := 0; i < total; i++ {
		go func(i int) {
			saveByGorm(t, &wg, &list_gorm[i])
			gormDone <- true
		}(i)
	}
	wg.Wait()
	close(gormDone)

	gormEnd := currentMilliseconds()
	var durationGorm = gormEnd - gormStart
	log.Printf("gorm use %d milliseconds ", durationGorm)
	log.Printf("sqlx.DB spend %f %% (more/less) time than sql.DB", float64(durationSqlx-durationSql)*100/float64(durationSql))
	log.Printf("gorm.DB spend %f %% (more/less) time than sql.DB", float64(durationGorm-durationSql)*100/float64(durationSql))
	log.Printf("gorm.DB spend %f %% (more/less) time than sqlx.DB", float64(durationGorm-durationSqlx)*100/float64(durationSqlx))
}

func currentMilliseconds() int64 {
	return time.Now().UnixNano() / 1000000
}
