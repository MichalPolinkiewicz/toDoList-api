package models

//const for task status
const ToDo = 1
const InProgress = 2
const Done = 3

type Task struct {
	Id          int    `gorm:"AUTO_INCREMENT=yes;PRIMARY_KEY:yes"`
	Name        string `json:"name,omitempty" gorm:"type:varchar(100)"`
	Description string `json:"description,omitempty" gorm:"type:varchar(255)"`
	Status      int    `json:"status,omitempty" gorm:"type:int"`
	UserId      int    `gorm:"type:int"`
}

type Tasks []Task

type User struct {
	Id       int    `gorm:"AUTO_INCREMENT=yes;PRIMARY_KEY:yes"`
	Username string `json:"username" gorm:"type:varchar(25)"`
	Password string `json:"password" gorm:"type:varchar(25)"`
}

type Fail struct {
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	HTTPStatus int         `json:"-"`
}

func (t *Task) IsValidTask() bool {
	return !isEmpty(t.Name) && !isEmpty(t.Description)
}

func isEmpty(s string) bool {
	return len(s) == 0
}
