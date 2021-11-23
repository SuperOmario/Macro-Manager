package controllers

import (
	"MacroManager/models"
	"database/sql"
	"log"
	"os"

	"github.com/lib/pq"
)

func SaveFood(food models.Food, upc string) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("INSERT INTO food(pantry_id, barcode, title, calories, fat, carbohydrate, protein, serving_size, misc) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", food.PantryID, upc, food.Title, food.Nutriments.Calories, food.Nutriments.Fat, food.Nutriments.Carbohydrate, food.Nutriments.Protein, food.ServingSize, pq.Array(food.Misc))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
}
