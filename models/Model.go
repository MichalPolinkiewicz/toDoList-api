package models

const ToDo  = 1
const InProgress  = 2
const Done  = 3

type Task struct {
	Id string
	Name string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Status int `json:"name:status,omitempty"`
}

var Tasks []Task
