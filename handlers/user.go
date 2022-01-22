package handlers

import (
	"MacroManager/controllers"
	"MacroManager/models"
)

func LoginUser(email string) models.User {
	user := controllers.GetUserByEmail(email)
	return user
}
