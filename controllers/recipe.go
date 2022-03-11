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
		db.Close()
		log.Println(err)
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		db.Close()
		log.Println(err)
	}

	var recipeId int64
	//Must make user id dynamic when users are implemented *TO DO*
	err = tx.QueryRowContext(ctx, "INSERT INTO recipe (user_id, title, serving_size) VALUES (1, $1, $2) RETURNING recipe_id", title, serving_size).Scan(&recipeId)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		db.Close()
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
		db.Close()
		log.Println(err)
	}

	db.Close()
	recipe = recipePlaceHolder
	return
}

func AddRecipeIngredient(recipeId int64, ingredient models.Ingredient) (recipe models.Recipe, err error) {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		db.Close()
		log.Println(err)
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		db.Close()
		log.Println(err)
	}

	var count int

	err = tx.QueryRowContext(ctx, "SELECT COUNT(ingredient_id) FROM recipe_ingredient WHERE recipe_id=$1 AND ingredient_id=$2", recipeId, ingredient.IngredientID).Scan(&count)
	if err != nil {
		tx.Rollback()
		db.Close()
		return
	} else if count == 1 {
		tx.Rollback()
		db.Close()
		err = errors.New("this ingredient is already part of this recipe")
		return
	}

	_, err = tx.Exec("INSERT INTO recipe_ingredient (ingredient_id, recipe_id, servings) VALUES ($1, $2, $3)", ingredient.IngredientID, recipeId, ingredient.Servings)
	if err != nil {
		tx.Rollback()
		db.Close()
		return
	}

	err = tx.QueryRowContext(ctx, "SELECT * FROM recipe WHERE recipe_id=$1", recipeId).Scan(&recipe.UserID, &recipe.RecipeID, &recipe.Title, &recipe.ServingSize)
	if err != nil {
		db.Close()

		return
	}
	err = tx.Commit()
	if err != nil {
		db.Close()
		return
	}
	db.Close()
	return
}

func UpdateRecipe(recipeId int64, title string, servingSize float32) error {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		db.Close()
		log.Println(err)
	}

	fmt.Println(title, servingSize, recipeId)

	_, err = db.Exec("UPDATE recipe SET title=$1, serving_size=$2 WHERE recipe_id=$3", title, servingSize, recipeId)
	db.Close()
	return err
}

func RemoveIngredient(recipeId int64, ingredientId int64) error {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		db.Close()
		log.Println(err)
	}

	_, err = db.Exec("DELETE FROM recipe_ingredient WHERE recipe_id=$1 AND ingredient_id=$2", recipeId, ingredientId)
	db.Close()
	return err
}

func DeleteRecipe(recipeId int64) error {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		db.Close()
		log.Println(err)
	}

	_, err = db.Exec("DELETE FROM recipe WHERE recipe_id=$1", recipeId)
	db.Close()
	return err
}

func GetRecipesForUser(userId int64) (recipes []models.RecipeDetails, err error) {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		db.Close()
		return
	}

	//*TODO* make user id dynamic
	rows, err := db.Query("SELECT recipe_id, title, serving_size FROM recipe WHERE user_id=$1", userId)
	if err != nil {
		db.Close()
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			var recipe models.RecipeDetails
			rows.Scan(&recipe.RecipeID, &recipe.Title, &recipe.ServingSize)
			recipe, err = getTotalNutrimentsRecipe(db, recipe.RecipeID, recipe)
			if err != nil {
				db.Close()
				return
			}
			recipes = append(recipes, recipe)
		}
	}
	db.Close()
	return
}

func GetRecipeById(recipeId int64) (recipe models.RecipeDetails, err error) {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		db.Close()
		return
	}
	//must make user id dynamic *TO DO*
	err = db.QueryRow("SELECT recipe_id, title, serving_size FROM recipe WHERE user_id=1 AND recipe_id=$1", recipeId).Scan(&recipe.RecipeID, &recipe.Title, &recipe.ServingSize)
	if err != nil {
		db.Close()
		return
	} else {
		recipe, err = getTotalNutrimentsRecipe(db, recipe.RecipeID, recipe)
	}
	db.Close()
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
func getTotalNutrimentsRecipe(db *sql.DB, recipeId int64, recipe models.RecipeDetails) (updatedRecipe models.RecipeDetails, err error) {

	rows, err := db.Query(
		"SELECT calories, fat, carbohydrate, protein, serving_size, misc FROM ingredient LEFT JOIN recipe_ingredient ON ingredient.ingredient_id=recipe_ingredient.ingredient_id WHERE recipe_ingredient.recipe_id=$1",
		recipeId)
	if err != nil {
		log.Print(err)
		return
	} else {
		counter := 0
		defer rows.Close()
		for rows.Next() {
			var food models.Food
			err := rows.Scan(&food.Nutriments.Calories, &food.Nutriments.Fat, &food.Nutriments.Carbohydrate, &food.Nutriments.Protein, &food.Serving_Size, pq.Array(&food.Misc))
			if err != nil {
				log.Print(err)
			}
			updatedRecipe = calculateNutrimentsRecipe(food, recipe)
			counter++
		}
		if counter == 0 {
			deleteEmptyRecipe(db, recipeId)
			err = errors.New("deleted recipe due to no ingredients")
			return
		}
	}
	return
}

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

func GetRecipeIngredientsByID(recipeID int64) (ingredients []models.IngredientForRecipe, err error) {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		db.Close()
		return
	}

	rows, err := db.Query("SELECT recipe_ingredient_id, ingredient.title, ingredient.ingredient_id, ingredient.serving_size, recipe_ingredient.servings FROM ingredient LEFT JOIN recipe_ingredient ON ingredient.ingredient_id = recipe_ingredient.ingredient_id WHERE recipe_id=$1", recipeID)
	if err != nil {
		db.Close()
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			var ingredient models.IngredientForRecipe
			rows.Scan(&ingredient.RecipeIngredientID, &ingredient.Title, &ingredient.IngredientID, &ingredient.ServingSize, &ingredient.Servings)
			ingredients = append(ingredients, ingredient)
		}
	}
	db.Close()
	return
}

func UpdateIngredients(ingredients models.IFRRequest) error {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		db.Close()
		log.Println(err)
	}

	for _, ingredient := range ingredients {
		_, err = db.Exec("UPDATE recipe_ingredient SET servings=$1 WHERE recipe_ingredient_id=$2", ingredient.Servings, ingredient.RecipeIngredientID)
	}

	db.Close()
	return err
}

func GetListedRecipes(ids []int) (recipes []models.RecipeDetails, err error) {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		db.Close()
		log.Println(err)
	}
	var userId int64 = 1
	// must change user id to be dynamic when implementing that feature *TO DO*
	rows, err := db.Query("SELECT recipe_id, title, serving_size FROM recipe WHERE user_id=$1 AND ingredient_id = ANY($2)", userId, pq.Array(ids))
	if err != nil {
		db.Close()
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			var recipe models.RecipeDetails
			rows.Scan(&recipe.RecipeID, &recipe.Title, &recipe.ServingSize)
			recipe, err = getTotalNutrimentsRecipe(db, recipe.RecipeID, recipe)
			if err != nil {
				db.Close()
				return
			}
			recipes = append(recipes, recipe)
		}
	}
	db.Close()
	return
}
