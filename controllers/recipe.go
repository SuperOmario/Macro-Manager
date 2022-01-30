package controllers

import (
	"MacroManager/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InsertRecipe(title string, serving_size float32, ingredients ...models.Ingredient) (recipe models.Recipe, err error) {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	var recipeId int64
	//Must make user id dynamic when users are implemented *TO DO*
	err = tx.QueryRowContext(ctx, "INSERT INTO recipe (user_id, title, serving_size) VALUES (1, $1, $2) RETURNING recipe_id", title, serving_size).Scan(&recipeId)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return
	}
	var recipePlaceHolder models.Recipe
	recipePlaceHolder.RecipeID = recipeId
	//Must make user id dynamic when users are implemented *TO DO*
	recipePlaceHolder.UserID = 1
	recipePlaceHolder.Title = title
	recipePlaceHolder.ServingSize = serving_size

	for _, ingredient := range ingredients {
		createIngredient(ingredient, recipeId, tx, ctx)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	recipe = recipePlaceHolder
	return
}

func AddRecipeIngredient(recipeId int64, ingredient models.Ingredient) (recipe models.Recipe, err error) {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	var count int

	err = tx.QueryRowContext(ctx, "SELECT COUNT(ingredient_id) FROM recipe_ingredient WHERE recipe_id=$1 AND ingredient_id=$2", recipeId, ingredient.IngredientID).Scan(&count)
	if err != nil {
		tx.Rollback()
		return
	} else if count == 1 {
		tx.Rollback()
		err = errors.New("this ingredient is already part of this recipe")
		return
	}

	_, err = tx.Exec("INSERT INTO recipe_ingredient (ingredient_id, recipe_id, servings) VALUES ($1, $2, $3)", ingredient.IngredientID, recipeId, ingredient.Servings)
	if err != nil {
		tx.Rollback()
		return
	}

	err = tx.QueryRowContext(ctx, "SELECT * FROM recipe WHERE recipe_id=$1", recipeId).Scan(&recipe.UserID, &recipe.RecipeID, &recipe.Title, &recipe.ServingSize)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	return
}

func UpdateRecipe(recipeId int64, title string, servingSize float32) error {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(title, servingSize, recipeId)

	_, err = db.Exec("UPDATE recipe SET title=$1, serving_size=$2 WHERE recipe_id=$3", title, servingSize, recipeId)
	return err
}

func RemoveIngredient(recipeId int64, ingredientId int64) error {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("DELETE FROM recipe_ingredient WHERE recipe_id=$1 AND ingredient_id=$2", recipeId, ingredientId)
	return err
}

func DeleteRecipe(recipeId int64) error {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("DELETE FROM recipe WHERE recipe_id=$1", recipeId)
	return err
}

//creates a recipe ingredient entry in the database
func createIngredient(ingredient models.Ingredient, recipeId int64, tx *sql.Tx, ctx context.Context, servingSize ...float32) {
	fmt.Println(ingredient.IngredientID)
	_, err := tx.ExecContext(ctx, "INSERT INTO recipe_ingredient (ingredient_id, recipe_id, servings) VALUES ($1, $2, $3)", ingredient.IngredientID, recipeId, ingredient.Servings)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
	}
}

//calculates the nutriments of a recipe ingredient for entry into the recipe
// func calculateNutrimentsRecipeIngredient(foodPlaceHolder models.Food, servings float32, servingSize float32) models.Food {
// 	foodPlaceHolder.Nutriments.Calories = (foodPlaceHolder.Nutriments.Calories * (servingSize / 100)) * servings
// 	foodPlaceHolder.Nutriments.Fat = (foodPlaceHolder.Nutriments.Fat * (servingSize / 100)) * servings
// 	foodPlaceHolder.Nutriments.Carbohydrate = (foodPlaceHolder.Nutriments.Carbohydrate * (servingSize / 100)) * servings
// 	foodPlaceHolder.Nutriments.Protein = (foodPlaceHolder.Nutriments.Protein * (servingSize / 100)) * servings

// 	return foodPlaceHolder
// }

//calculates the nutriments of a recipe
// func calculateNutrimentsRecipe(foodPlaceHolder models.Food, recipePlaceHolder models.Recipe) models.Recipe {
// 	recipePlaceHolder.Calories += foodPlaceHolder.Nutriments.Calories
// 	recipePlaceHolder.Fat += foodPlaceHolder.Nutriments.Fat
// 	recipePlaceHolder.Carbohydrate += foodPlaceHolder.Nutriments.Carbohydrate
// 	recipePlaceHolder.Protein += foodPlaceHolder.Nutriments.Protein
// 	recipePlaceHolder.Misc = append(recipePlaceHolder.Misc, foodPlaceHolder.Misc...)

// 	return recipePlaceHolder
// }

//in development

// func calculateNutrimentsRecipeRemove(foodPlaceHolder models.Food, recipePlaceHolder models.Recipe) models.Recipe {
// 	recipePlaceHolder.Calories -= foodPlaceHolder.Nutriments.Calories
// 	recipePlaceHolder.Fat -= foodPlaceHolder.Nutriments.Fat
// 	recipePlaceHolder.Carbohydrate -= foodPlaceHolder.Nutriments.Carbohydrate
// 	recipePlaceHolder.Protein -= foodPlaceHolder.Nutriments.Protein
// 	for i, misc := range recipePlaceHolder.Misc {
// 		for _, misc2 := range foodPlaceHolder.Misc {
// 			if misc == misc2 {
// 				recipePlaceHolder.Misc[i] = recipePlaceHolder.Misc[len(recipePlaceHolder.Misc)-1]
// 			}
// 		}
// 	}
// 	return recipePlaceHolder
// }
