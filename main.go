package main

import (
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
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("connection is good")
	}

	// var ctx context.Context

	// db.ExecContext(ctx,
	// 	"INSERT INTO User (FName, LName, Email) VALUES ($1, $2, $3)", "Omar", "Sallam", "omarsallamala@gmail.com")

	godotenv.Load()

	router = gin.Default()

	initialiseRoutes()

	router.Run(":8080")
}
