package controllers

import (
	"MacroManager/models"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

//checks if food is already in the current users pantry and if not inserts it into the database
func InsertFood(food models.Food, upc string) (id string, err error) {
	id = ""
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println(err)
		db.Close()
		return
	}

	// must change user id to be dynamic when implementing that feature *TO DO*
	row := db.QueryRow("SELECT * FROM ingredient WHERE barcode=$1 and user_id=1", upc)
	if err != nil {
		log.Println(err)
		db.Close()
		return
	} else {
		var foodPlaceHolder models.Food
		err = row.Scan(&foodPlaceHolder.UserID, &foodPlaceHolder.IngredientID, &foodPlaceHolder.Barcode, &foodPlaceHolder.Title, &foodPlaceHolder.Nutriments.Calories,
			&foodPlaceHolder.Nutriments.Fat, &foodPlaceHolder.Nutriments.Carbohydrate, &foodPlaceHolder.Nutriments.Protein, &foodPlaceHolder.Serving_Size, pq.Array(&foodPlaceHolder.Misc))
		if err != nil {
			err = db.QueryRow("INSERT INTO ingredient(user_id, barcode, title, calories, fat, carbohydrate, protein, serving_size, misc) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING ingredient_id",
				food.UserID, upc, food.Title, food.Nutriments.Calories, food.Nutriments.Fat, food.Nutriments.Carbohydrate, food.Nutriments.Protein, food.Serving_Size, pq.Array(food.Misc)).Scan(&id)
			if err != nil {
				log.Println(err)
				db.Close()
				return
			} else {
				ingredient := models.Ingredient{IngredientID: foodPlaceHolder.IngredientID, Servings: 1}
				InsertRecipe(food.Title, food.Serving_Size, ingredient)
				return
			}
		} else {
			id = strconv.Itoa(int(foodPlaceHolder.IngredientID))
			fmt.Println("Food already saved for this user")
			db.Close()
			return
		}
	}
}

func InsertCustomFood(food models.CustomFood) (err error) {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Print(err)
		db.Close()
		return
	}
	//must make userID dynamic *TO DO*
	_, err = db.Exec("INSERT INTO ingredient(user_id, title, calories, fat, carbohydrate, protein, serving_size, misc) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		1, food.Title, food.Calories, food.Fat, food.Carbohydrate, food.Protein, food.Serving_Size, pq.Array(food.Misc))
	if err != nil {
		log.Print(err)
		db.Close()
		return
	}
	db.Close()
	return
}

func DeleteFood(foodId int64) error {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println(err)
		db.Close()
		return err
	}

	_, err = db.Exec("DELETE FROM ingredient WHERE ingredient_id=$1", foodId)
	db.Close()
	return err
}

func GetAllFood() []models.Food {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		db.Close()
		log.Println(err)
	}
	// must change pantry id to be dynamic when implementing that feature *TO DO*
	rows, err := db.Query("SELECT * FROM ingredient")
	if err != nil {
		db.Close()
		log.Println(err)
	}
	defer rows.Close()
	var foods []models.Food
	for rows.Next() {
		var foodPlaceHolder models.Food
		err := rows.Scan(&foodPlaceHolder.IngredientID, &foodPlaceHolder.Barcode, &foodPlaceHolder.Title, &foodPlaceHolder.Nutriments.Calories,
			&foodPlaceHolder.Nutriments.Fat, &foodPlaceHolder.Nutriments.Carbohydrate, &foodPlaceHolder.Nutriments.Protein, &foodPlaceHolder.Serving_Size, pq.Array(&foodPlaceHolder.Misc))
		if err != nil {
			log.Println(err)
		}
		foods = append(foods, foodPlaceHolder)
	}
	db.Close()
	return foods
}

func GetPantry() []models.Food {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		db.Close()
		log.Println(err)
	}
	// must change user id to be dynamic when implementing that feature *TO DO*
	rows, err := db.Query("SELECT * FROM ingredient WHERE user_id=1")
	if err != nil {
		db.Close()
		log.Println(err)
	}
	defer rows.Close()
	var foods []models.Food
	for rows.Next() {
		var foodPlaceHolder models.Food
		err := rows.Scan(&foodPlaceHolder.UserID, &foodPlaceHolder.IngredientID, &foodPlaceHolder.Barcode, &foodPlaceHolder.Title, &foodPlaceHolder.Nutriments.Calories,
			&foodPlaceHolder.Nutriments.Fat, &foodPlaceHolder.Nutriments.Carbohydrate, &foodPlaceHolder.Nutriments.Protein, &foodPlaceHolder.Serving_Size,
			pq.Array(&foodPlaceHolder.Misc))
		if err != nil {
			db.Close()
			log.Print(err)
		}
		foods = append(foods, foodPlaceHolder)
	}
	db.Close()
	return foods
}

func UpdateFood(ID int64, food models.FoodUpdate) error {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		db.Close()
		log.Print(err)
	}

	var barcode string = ""
	_ = db.QueryRow("SELECT barcode FROM ingredient WHERE ingredient_id=$1", ID).Scan(&barcode)
	if barcode != "" && food.ServingSize != 0 {
		fmt.Print("Case 1 ", barcode, "    ", food.ServingSize)
		_, err = db.Exec("UPDATE ingredient SET title=$1, serving_size=$2 WHERE ingredient_id=$3", food.Title, food.ServingSize, ID)
		if err != nil {
			db.Close()
			log.Print(err)
			return err
		}
	} else if barcode == "" && food.ServingSize != 0 {
		fmt.Print("Case 2 ", barcode, "    ", food.ServingSize)
		_, err = db.Exec("UPDATE ingredient SET calories=$1, fat=$2, carbohydrate=$3, protein=$4, serving_size=$5, title=$6 WHERE ingredient_id=$7", food.Calories,
			food.Fat, food.Carbohydrate, food.Protein, food.ServingSize, food.Title, ID)
		if err != nil {
			db.Close()
			log.Print(err)
			return err
		}
	} else {
		fmt.Print("Case 3 ", barcode, "    ", food.ServingSize)
		_, err = db.Exec("UPDATE ingredient SET calories=$1, fat=$2, carbohydrate=$3, protein=$4, title=$5 WHERE ingredient_id=$6", food.Calories,
			food.Fat, food.Carbohydrate, food.Protein, food.Title, ID)
		if err != nil {
			db.Close()
			log.Print(err)
			return err
		}
	}
	db.Close()
	return err
}
