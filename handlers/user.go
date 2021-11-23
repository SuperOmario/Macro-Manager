package handlers

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
		var (
			user_id int
			fname   string
			lname   string
			email   string
		)
		if err := rows.Scan(&user_id, &fname, &lname, &email); err != nil {
			log.Fatal(err)
		}
		log.Printf("id %d: name is %s %s, email is %s \n", user_id, fname, lname, email)
		user = models.NewUser(user_id, fname, lname, email)
	}
	return user
}
