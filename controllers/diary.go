package controllers

import (
	"MacroManager/models"
	"context"
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

	rows, err := db.Query("SELECT * FROM diary_entry WHERE user_id=$1", 1)
	if err != nil {
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			diaryEntries, err = createDiaryEntryArray(rows)
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
			diaryEntries, err = createDiaryEntryArray(rows)
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
		log.Fatal(err)
	}

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := tx.QueryContext(ctx,
		"SELECT calories, fat, carbohydrate, protein, serving_size, misc FROM ingredient LEFT JOIN recipe_ingredient ON ingredient.ingredient_id=recipe_ingredient.ingredient_id WHERE recipe_ingredient.recipe_id=$1",
		recipeId)
	if err != nil {
		log.Fatal(err)
		tx.Rollback()
		return
	}
	defer rows.Close()
	var diaryEntry models.DiaryEntry
	for rows.Next() {
		var food models.Food
		err := rows.Scan(&food.Nutriments.Calories, &food.Nutriments.Fat, &food.Nutriments.Carbohydrate, &food.Nutriments.Protein, &food.Serving_Size, pq.Array(&food.Misc))
		if err != nil {
			log.Fatal(err)
		}
		diaryEntry = calculateNutrimentsDiaryEntry(food, diaryEntry, servings, food.Serving_Size)
	}

	err = tx.QueryRowContext(ctx, "INSERT INTO diary_entry (user_id, recipe_id, date, meal, calories, fat, carbohydrate, protein, servings, misc) VALUES (1, $1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING diary_entry_id",
		recipeId, date, meal, diaryEntry.Calories, diaryEntry.Fat, diaryEntry.Carbohydrate, diaryEntry.Protein, servings, pq.Array(diaryEntry.Misc)).Scan(&diaryEntryId)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return
	}

	return
}

func DeleteDiaryEntry(diaryEntryID int64) {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("DELETE FROM diary_entry WHERE diary_entry_id=$1", diaryEntryID)
	if err != nil {
		log.Fatal(err)
	}
}

//Helper function to modify the nutriment values based on how much of a food the user wants to enter into the diary
func calculateNutrimentsDiaryEntry(foodPlaceHolder models.Food, diaryEntryPlaceHolder models.DiaryEntry, servings float32, servingSize float32) models.DiaryEntry {
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

func createDiaryEntryArray(rows *sql.Rows) (diaryEntries models.DiaryEntries, err error) {
	var diaryEntry models.DiaryEntry
	err = rows.Scan(&diaryEntry.UserID, &diaryEntry.DiaryEntryID, &diaryEntry.RecipeID, &diaryEntry.Date, &diaryEntry.Meal,
		&diaryEntry.Calories, &diaryEntry.Fat, &diaryEntry.Carbohydrate, &diaryEntry.Protein, &diaryEntry.Servings, pq.Array(&diaryEntry.Misc))
	diaryEntries = append(diaryEntries, diaryEntry)
	if err != nil {
		return
	}
	return
}
