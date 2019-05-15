package models

//const for task status
const ToDo  = 1
const InProgress  = 2
const Done  = 3

type Task struct {
	ID int `gorm:"AUTO_INCREMENT=yes;PRIMARY_KEY:yes"`
	Name        string `json:"name,omitempty" gorm:"type:varchar(100)"`
	Description string `json:"description,omitempty" gorm:"type:varchar(255)"`
	Status      int    `json:"name:status,omitempty" gorm:"type:int"`
}

//var Tasks []Task

//type User struct {
//	Login string
//	Password string
//	Token string
//	IsLogged bool
//}
