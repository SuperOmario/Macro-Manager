package controllers

// import (
// 	"MacroManager/models"
// 	"context"
// 	"database/sql"
// 	"log"
// 	"os"
// 	"time"

// 	"github.com/joho/godotenv"
// )

// I used https://www.sohamkamani.com/golang/sql-transactions/ to learn the transaction syntax
// func InsertDiaryEntryFood(foodId int64, servings float32, servingSize ...float32) {
// 	godotenv.Load()
// 	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	date := time.Now().Format("2006-01-02")

// 	ctx := context.Background()
// 	tx, err := db.BeginTx(ctx, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	row := tx.QueryRowContext(ctx, "SELECT diary_entry_id, date FROM diary_entry WHERE user_id=1 and date=$1", date)
// 	var diaryEntryPlaceHolder models.DiaryEntry
// 	err = row.Scan(&diaryEntryPlaceHolder.DiaryEntryID, &diaryEntryPlaceHolder.Date)
// 	var diaryEntryId = diaryEntryPlaceHolder.DiaryEntryID
// 	if err != nil {
// 		//must change user id to be dynamic when implementing users *TO DO*
// 		_, err = tx.ExecContext(ctx, "INSERT INTO diary_entry (user_id, date) VALUES (1, $1)", date)
// 		if err != nil {
// 			tx.Rollback()
// 			return
// 		} else {
// 			row := tx.QueryRowContext(ctx, "SELECT diary_entry_id, date FROM diary_entry WHERE user_id=1 date=$1", date)
// 			var diaryEntryPlaceHolder models.DiaryEntry
// 			err = row.Scan(&diaryEntryPlaceHolder.DiaryEntryID, &diaryEntryPlaceHolder.Date)
// 			diaryEntryId = diaryEntryPlaceHolder.DiaryEntryID
// 			if err != nil {
// 				tx.Rollback()
// 				return
// 			}
// 		}
// 	}

// 	row = tx.QueryRowContext(ctx, "SELECT calories, fat, carbohydrate, protein, serving_size FROM food where food_id=$1", foodId)
// 	var foodPlaceHolder models.Food
// 	err = row.Scan(&foodPlaceHolder.Nutriments.Calories, &foodPlaceHolder.Nutriments.Fat, &foodPlaceHolder.Nutriments.Carbohydrate, &foodPlaceHolder.Nutriments.Protein,
// 		&foodPlaceHolder.Serving_Size)
// 	if err != nil {
// 		tx.Rollback()
// 		return
// 	} else {
// 		if foodPlaceHolder.Serving_Size == 0 {
// 			if servingSize != nil {
// 				foodPlaceHolder, diaryEntryPlaceHolder = calculateNutrimentsFoodEntry(foodPlaceHolder, diaryEntryPlaceHolder, servings, servingSize[0])
// 			} else {
// 				//default serving size will go to 100g
// 				foodPlaceHolder, diaryEntryPlaceHolder = calculateNutrimentsFoodEntry(foodPlaceHolder, diaryEntryPlaceHolder, servings, 100)
// 			}

// 		} else {
// 			foodPlaceHolder, diaryEntryPlaceHolder = calculateNutrimentsFoodEntry(foodPlaceHolder, diaryEntryPlaceHolder, servings, foodPlaceHolder.Serving_Size)
// 		}
// 	}

// 	_, err = tx.ExecContext(ctx, "INSERT INTO diary_entry_food (diary_entry_id, food_id, servings, calories, fat, carbohydrate, protein) VALUES ($1, $2, $3, $4, $5, $6, $7)",
// 		diaryEntryId, foodId, servings, foodPlaceHolder.Nutriments.Calories, foodPlaceHolder.Nutriments.Fat, foodPlaceHolder.Nutriments.Carbohydrate,
// 		foodPlaceHolder.Nutriments.Protein)
// 	if err != nil {
// 		tx.Rollback()
// 		return
// 	}

// 	_, err = tx.ExecContext(ctx, "UPDATE diary_entry SET calories=$1, fat=$2, carbohydrate=$3, protein=$4 WHERE diary_entry_id=$5", diaryEntryPlaceHolder.Calories,
// 		diaryEntryPlaceHolder.Fat, diaryEntryPlaceHolder.Carbohydrate, diaryEntryPlaceHolder.Protein, diaryEntryId)
// 	if err != nil {
// 		tx.Rollback()
// 		return
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// //Helper function to modify the nutriment values based on how much of a food the user wants to enter into the diary
// func calculateNutrimentsFoodEntry(foodPlaceHolder models.Food, diaryEntryPlaceHolder models.DiaryEntry, servings float32, servingSize float32) (models.Food, models.DiaryEntry) {
// 	foodPlaceHolder.Nutriments.Calories = (foodPlaceHolder.Nutriments.Calories * (servingSize / 100)) * servings
// 	foodPlaceHolder.Nutriments.Fat = (foodPlaceHolder.Nutriments.Fat * (servingSize / 100)) * servings
// 	foodPlaceHolder.Nutriments.Carbohydrate = (foodPlaceHolder.Nutriments.Carbohydrate * (servingSize / 100)) * servings
// 	foodPlaceHolder.Nutriments.Protein = (foodPlaceHolder.Nutriments.Protein * (servingSize / 100)) * servings
// 	diaryEntryPlaceHolder.Calories += foodPlaceHolder.Nutriments.Calories
// 	diaryEntryPlaceHolder.Fat += foodPlaceHolder.Nutriments.Fat
// 	diaryEntryPlaceHolder.Carbohydrate += foodPlaceHolder.Nutriments.Carbohydrate
// 	diaryEntryPlaceHolder.Protein += foodPlaceHolder.Nutriments.Protein

// 	return foodPlaceHolder, diaryEntryPlaceHolder
// }

// // }

// func InsertDiaryEntryRecipe() {

// }
