package controllers

import (
	"MacroManager/models"
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

//retrieves a user by email from the db
func GetUserByEmail(email string) models.User {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		db.Close()
		log.Println(err)
	}

	rows, err := db.Query("SELECT * FROM app_user WHERE email= $1", email)
	if err != nil {
		db.Close()
		log.Println(err)
	}
	defer rows.Close()

	var user models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.FName, &user.LName, &email); err != nil {
			db.Close()
			log.Println(err)
		}
		log.Printf("id %d: name is %s %s, email is %s \n", user.UserID, user.FName, user.LName, email)
		user = models.NewUser(user.UserID, user.FName, user.LName, email)
	}
	db.Close()
	return user
}
