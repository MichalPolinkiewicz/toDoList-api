package db

import (
	"github.com/MichalPolinkiewicz/to-do-api/models"
)

func CreateTask(t *models.Task){
	db.Create(t)
}

func GetAllTasks(i *int) []models.Task {
	var tasks []models.Task
	db.Where("user_id = ?", i).Find(&tasks)

	return tasks
}

func GetTaskById(i *int, u *int) models.Task {
	var task models.Task
	db.Where("id = ?", i).Where("user_id = ?", u).Find(&task)

	return task
}

func GetTasksByStatus(s *int, u *int) *[]models.Task {
	var tasks []models.Task
	db.Where("status = ?", s).Where("user_id = ?", u).Find(&tasks)

	return &tasks
}

func GetUserFromDb(u *string, p *string) *models.User {
	var user models.User
	db.Where("username = ?", u).Where("password = ?", p).Find(&user)

	return &user
}

func CheckIfUserExistsInDb(u *string) bool{
	var user models.User
	db.Where("username = ?", u).Find(&user)

	if user.Username == *u && user.Password != "" {
		return true
	}
	return false
}

func SaveUser (u *models.User){
	db.Create(u)
}
