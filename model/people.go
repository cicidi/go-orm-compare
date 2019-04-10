package model

//Created by cicidi / cicidi@gmail.com
//Date: 2019-04-06
//Time: 11:09
import "github.com/google/uuid"

type People struct {
	Id uuid.UUID `db:"id", gorm:"primary_key"`

	CreatedTimestamp int64 `db:"created_timestamp"`

	IsStudent bool `db:"is_student"`

	Address string `db:"address"`

	Gender string `db:"gender"`

	Race string `db:"race"`
}
