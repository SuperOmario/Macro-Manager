package controllers

import (
	"MacroManager/models"
	"database/sql"
	"log"
)

//retrieves a user by email from the db
func GetUserByEmail(email string, db *sql.DB) models.User {
	rows, err := db.Query("SELECT * FROM app_user WHERE email= $1", email)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var user models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.FName, &user.LName, &email); err != nil {
			log.Fatal(err)
		}
		log.Printf("id %d: name is %s %s, email is %s \n", user.UserID, user.FName, user.LName, email)
		user = models.NewUser(user.UserID, user.FName, user.LName, email)
	}
	return user
}
