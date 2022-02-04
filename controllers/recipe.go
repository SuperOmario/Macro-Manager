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
	"github.com/lib/pq"
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

func GetRecipesForUser(userId int64) (recipes []models.RecipeDetails, err error) {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return
	}

	rows, err := db.Query("SELECT recipe_id, title, serving_size FROM recipe WHERE user_id=$1", userId)
	if err != nil {
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			var recipeDetails models.RecipeDetails
			var recipe models.Recipe
			rows.Scan(&recipe.RecipeID, &recipe.Title, &recipe.ServingSize)
			recipeDetails, err = getTotalNutrimentsRecipe(db, recipe.RecipeID)
			if err != nil {
				return
			}
			recipes = append(recipes, recipeDetails)
		}
	}
	return
}

func GetRecipeById(recipeId int64) (recipe models.RecipeDetails, err error) {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return
	}

	err = db.QueryRow("SELECT recipe_id, title, serving_size FROM recipe WHERE user_id=$1", recipeId).Scan(&recipe.RecipeID, &recipe.Title, &recipe.ServingSize)
	if err != nil {
		return
	} else {
		recipe, err = getTotalNutrimentsRecipe(db, recipe.RecipeID)
	}
	return
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

//Helper function to grab all of the ingredients involved in a specific diary entry and run calculations for total nutriments
func getTotalNutrimentsRecipe(db *sql.DB, recipeId int64) (recipe models.RecipeDetails, err error) {

	rows, err := db.Query(
		"SELECT calories, fat, carbohydrate, protein, serving_size, misc FROM ingredient LEFT JOIN recipe_ingredient ON ingredient.ingredient_id=recipe_ingredient.ingredient_id WHERE recipe_ingredient.recipe_id=$1",
		recipeId)
	if err != nil {
		log.Print(err)
		return
	}
	if !rows.Next() {
		deleteEmptyRecipe(db, recipeId)
		err = errors.New("deleted recipe due to no ingredients")
		return
	}
	defer rows.Close()
	for rows.Next() {
		var food models.Food
		err := rows.Scan(&food.Nutriments.Calories, &food.Nutriments.Fat, &food.Nutriments.Carbohydrate, &food.Nutriments.Protein, &food.Serving_Size, pq.Array(&food.Misc))
		if err != nil {
			log.Print(err)
		}
		recipe = calculateNutrimentsRecipe(food, recipe)
	}
	return
}

//calculates the nutriments of a recipe ingredient for entry into the recipe
// func calculateNutrimentsRecipeIngredient(foodPlaceHolder models.Food, servings float32, servingSize float32) models.Food {
// 	foodPlaceHolder.Nutriments.Calories = (foodPlaceHolder.Nutriments.Calories * (servingSize / 100)) * servings
// 	foodPlaceHolder.Nutriments.Fat = (foodPlaceHolder.Nutriments.Fat * (servingSize / 100)) * servings
// 	foodPlaceHolder.Nutriments.Carbohydrate = (foodPlaceHolder.Nutriments.Carbohydrate * (servingSize / 100)) * servings
// 	foodPlaceHolder.Nutriments.Protein = (foodPlaceHolder.Nutriments.Protein * (servingSize / 100)) * servings

// 	return foodPlaceHolder
// }

// helper function to calculate the nutriments of a recipe
func calculateNutrimentsRecipe(foodPlaceHolder models.Food, recipePlaceHolder models.RecipeDetails) models.RecipeDetails {
	recipePlaceHolder.Calories += foodPlaceHolder.Nutriments.Calories
	recipePlaceHolder.Fat += foodPlaceHolder.Nutriments.Fat
	recipePlaceHolder.Carbohydrate += foodPlaceHolder.Nutriments.Carbohydrate
	recipePlaceHolder.Protein += foodPlaceHolder.Nutriments.Protein
	recipePlaceHolder.Misc = append(recipePlaceHolder.Misc, foodPlaceHolder.Misc...)

	return recipePlaceHolder
}

func deleteEmptyRecipe(db *sql.DB, recipeId int64) error {
	_, err := db.Exec("DELETE FROM recipe WHERE recipe_id=$1", recipeId)
	return err
}

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
