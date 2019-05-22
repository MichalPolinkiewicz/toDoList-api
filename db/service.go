package db

import (
	"github.com/MichalPolinkiewicz/to-do-api/models"
)

func CreateTask(t *models.Task){
	db.Create(*t)
}

func GetAllTasks(i *int) []models.Task {
	var tasks []models.Task
	db.Where("user_id = ?", *i).Find(&tasks)

	return tasks
}

func GetTaskById(i *int) models.Task {
	var task models.Task
	db.Where("id = ?", *i).Find(&task)

	return task
}

func GetTasksByStatus(s *int) *[]models.Task {
	var tasks []models.Task
	db.Where("status = ?", *s).Find(&tasks)

	return &tasks
}

func GetUserByUsernameAndPassword(u *string, p *string) *models.User {
	var user models.User
	db.Where("username = ?", *u).Where("password = ?", *p).Find(&user)

	return &user
}
