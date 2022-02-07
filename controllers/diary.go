package controllers

import (
	"MacroManager/models"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

func GetAllDiaryEntriesForUser() (diaryEntries models.DiaryEntries, err error) {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return
	}

	//must change user id to be dynamic *TO DO*
	rows, err := db.Query("SELECT * FROM diary_entry WHERE user_id=1")
	if err != nil {
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			var RecipeID, ID, UserID int64
			var date, meal string
			var servings float32
			err = rows.Scan(&ID, &UserID, &RecipeID, &date, &meal, &servings)
			if err != nil {
				return
			}
			diaryEntry := getTotalNutrimentsDiary(db, RecipeID, servings)
			diaryEntry.UserID = UserID
			diaryEntry.Date = date
			diaryEntry.Meal = meal
			diaryEntry.DiaryEntryID = ID
			diaryEntry.RecipeID = RecipeID
			diaryEntries = append(diaryEntries, diaryEntry)
			if err != nil {
				return
			}
		}
		return
	}
}

func GetDiaryEntriesByDate(date string) (diaryEntries models.DiaryEntries, err error) {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return
	}

	//must make user id dynamic *TO DO*
	rows, err := db.Query("SELECT * FROM diary_entry WHERE date=$1 AND user_id=1", date)
	if err != nil {
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			var RecipeID, ID, UserID int64
			var date, meal string
			var servings float32
			err = rows.Scan(&ID, &UserID, &RecipeID, &date, &meal, &servings)
			if err != nil {
				return
			}
			diaryEntry := getTotalNutrimentsDiary(db, RecipeID, servings)
			diaryEntry.UserID = UserID
			diaryEntry.Date = date
			diaryEntry.Meal = meal
			diaryEntry.DiaryEntryID = ID
			diaryEntry.RecipeID = RecipeID
			diaryEntries = append(diaryEntries, diaryEntry)
			if err != nil {
				return
			}
		}
		return
	}
}

// I used https://www.sohamkamani.com/golang/sql-transactions/ to learn the transaction syntax
func InsertDiaryEntry(recipeId int64, servings float32, date string, meal string) (diaryEntryId int64) {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Print(err)
	}

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	err = db.QueryRow("INSERT INTO diary_entry (user_id, recipe_id, date, meal, servings) VALUES (1, $1, $2, $3, $4) RETURNING diary_entry_id",
		recipeId, date, meal, servings).Scan(&diaryEntryId)
	if err != nil {
		log.Print(err)
	}

	return
}

func UpdateDiaryEntry(ID int64, servings float32) error {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Print(err)
	}
	_, err = db.Exec("UPDATE diary_entry SET servings=$1 WHERE diary_entry_id=$2", servings, ID)

	return err
}

func DeleteDiaryEntry(diaryEntryID int64) {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Print(err)
	}

	_, err = db.Exec("DELETE FROM diary_entry WHERE diary_entry_id=$1", diaryEntryID)
	if err != nil {
		log.Print(err)
	}
}

//Helper function to modify the nutriment values based on how much of a food the user wants to enter into the diary
func calculateNutrimentsDiary(foodPlaceHolder models.Food, diaryEntryPlaceHolder models.Diary, servings float32, servingSize float32) models.Diary {
	foodPlaceHolder.Nutriments.Calories = (foodPlaceHolder.Nutriments.Calories * (servingSize / 100)) * servings
	foodPlaceHolder.Nutriments.Fat = (foodPlaceHolder.Nutriments.Fat * (servingSize / 100)) * servings
	foodPlaceHolder.Nutriments.Carbohydrate = (foodPlaceHolder.Nutriments.Carbohydrate * (servingSize / 100)) * servings
	foodPlaceHolder.Nutriments.Protein = (foodPlaceHolder.Nutriments.Protein * (servingSize / 100)) * servings
	diaryEntryPlaceHolder.Calories += foodPlaceHolder.Nutriments.Calories
	diaryEntryPlaceHolder.Fat += foodPlaceHolder.Nutriments.Fat
	diaryEntryPlaceHolder.Carbohydrate += foodPlaceHolder.Nutriments.Carbohydrate
	diaryEntryPlaceHolder.Protein += foodPlaceHolder.Nutriments.Protein

	return diaryEntryPlaceHolder
}

//Helper function to grab all of the ingredients involved in a specific diary entry and run calculations for total nutriments
func getTotalNutrimentsDiary(db *sql.DB, recipeId int64, servings float32) (diaryEntry models.Diary) {

	rows, err := db.Query(
		"SELECT calories, fat, carbohydrate, protein, serving_size, misc FROM ingredient LEFT JOIN recipe_ingredient ON ingredient.ingredient_id=recipe_ingredient.ingredient_id WHERE recipe_ingredient.recipe_id=$1",
		recipeId)
	if err != nil {
		log.Print(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var food models.Food
		err := rows.Scan(&food.Nutriments.Calories, &food.Nutriments.Fat, &food.Nutriments.Carbohydrate, &food.Nutriments.Protein, &food.Serving_Size, pq.Array(&food.Misc))
		if err != nil {
			log.Print(err)
		}
		diaryEntry = calculateNutrimentsDiary(food, diaryEntry, servings, food.Serving_Size)
	}
	return
}
