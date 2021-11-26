package main

import (
	"MacroManager/controllers"
	"fmt"
	"log"
	"os"

	"database/sql"
	//third party packages
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

//taken and adapted from this site https://semaphoreci.com/community/tutorials/building-go-web-applications-and-microservices-using-gin
var router *gin.Engine

func main() {
	//loads environment variables sets up port and db connection
	godotenv.Load()
	port := os.Getenv("PORT")
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	//verifies db connection
	status := "up"
	if err := db.Ping(); err != nil {
		status = "down"
	}
	log.Println(status)

	user := controllers.GetUserByEmail("fake_email@gmail.com", db)

	fmt.Println(user)

	router = gin.Default()

	initialiseRoutes()

	router.Run(":" + port)
}
