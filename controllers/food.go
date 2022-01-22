package controllers

import (
	"MacroManager/models"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

//checks if food is already in the current users pantry and if not inserts it into the database
func InsertFood(food models.Food, upc string) {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	// must change pantry id to be dynamic when implementing that feature *TO DO*
	row := db.QueryRow("SELECT * FROM food WHERE barcode=$1 and pantry_id=1", upc)
	if err != nil {
		log.Fatal(err)
	} else {
		var foodPlaceHolder models.Food
		err := row.Scan(&foodPlaceHolder.PantryID, &foodPlaceHolder.FoodID, &foodPlaceHolder.Barcode, &foodPlaceHolder.Title, &foodPlaceHolder.Nutriments.Calories,
			&foodPlaceHolder.Nutriments.Fat, &foodPlaceHolder.Nutriments.Carbohydrate, &foodPlaceHolder.Nutriments.Protein, &foodPlaceHolder.Serving_Size, pq.Array(&foodPlaceHolder.Misc))
		if err != nil {
			rows, err := db.Query("INSERT INTO food(pantry_id, barcode, title, calories, fat, carbohydrate, protein, serving_size, misc) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
				food.PantryID, upc, food.Title, food.Nutriments.Calories, food.Nutriments.Fat, food.Nutriments.Carbohydrate, food.Nutriments.Protein, food.Serving_Size, pq.Array(food.Misc))
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
		} else {
			fmt.Println("Food already saved for this user")
		}
	}
}

func GetAllFood() []models.Food {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	// must change pantry id to be dynamic when implementing that feature *TO DO*
	rows, err := db.Query("SELECT * FROM food")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var foods []models.Food
	for rows.Next() {
		var foodPlaceHolder models.Food
		err := rows.Scan(&foodPlaceHolder.PantryID, &foodPlaceHolder.FoodID, &foodPlaceHolder.Barcode, &foodPlaceHolder.Title, &foodPlaceHolder.Nutriments.Calories,
			&foodPlaceHolder.Nutriments.Fat, &foodPlaceHolder.Nutriments.Carbohydrate, &foodPlaceHolder.Nutriments.Protein, &foodPlaceHolder.Serving_Size, pq.Array(&foodPlaceHolder.Misc))
		if err != nil {
			log.Fatal(err)
		}
		foods = append(foods, foodPlaceHolder)
	}
	return foods
}

func GetPantry() []models.Food {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	// must change pantry id to be dynamic when implementing that feature *TO DO*
	rows, err := db.Query("SELECT * FROM food WHERE pantry_id=1")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var foods []models.Food
	for rows.Next() {
		var foodPlaceHolder models.Food
		err := rows.Scan(&foodPlaceHolder.PantryID, &foodPlaceHolder.FoodID, &foodPlaceHolder.Barcode, &foodPlaceHolder.Title, &foodPlaceHolder.Nutriments.Calories,
			&foodPlaceHolder.Nutriments.Fat, &foodPlaceHolder.Nutriments.Carbohydrate, &foodPlaceHolder.Nutriments.Protein, &foodPlaceHolder.Serving_Size, pq.Array(&foodPlaceHolder.Misc))
		if err != nil {
			log.Fatal(err)
		}
		foods = append(foods, foodPlaceHolder)
	}
	return foods
}
