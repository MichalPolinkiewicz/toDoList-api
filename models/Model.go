package models

//const for task status
const ToDo = 1
const InProgress = 2
const Done = 3

type Task struct {
	ID          int    `gorm:"AUTO_INCREMENT=yes;PRIMARY_KEY:yes"`
	Name        string `json:"name,omitempty" gorm:"type:varchar(100)"`
	Description string `json:"description,omitempty" gorm:"type:varchar(255)"`
	Status      int    `json:"status,omitempty" gorm:"type:int"`
}

func (t *Task) IsValidTask() bool {
	return !isEmpty(t.Name) && !isEmpty(t.Description) && StatusIsValid(t.Status)
}

func StatusIsValid(s int)bool{
	return s >= 1 && s <=3
}

func isEmpty(s string) bool {
	return len(s) == 0
}

type User struct {
	Login    string
	Password string
	Key      string //login + password + datetime now + szyfr
	IsLogged bool
}
