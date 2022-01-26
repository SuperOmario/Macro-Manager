package controllers

import (
	"MacroManager/models"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

func InsertRecipe(title string, ingredients ...models.Ingredient) (recipe models.Recipe, err error) {
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
	//Must make pantry id and user id dynamic when users are implemented *TO DO*
	err = tx.QueryRowContext(ctx, "INSERT INTO recipe (user_id, pantry_id, title) VALUES (1,1, $1) RETURNING recipe_id", title).Scan(&recipeId)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return
	}

	var foodPlaceHolder models.Food
	var foods []models.Food
	var recipePlaceHolder models.Recipe
	recipePlaceHolder.RecipeID = recipeId
	//Must make pantry id dynamic when users are implemented *TO DO*
	recipePlaceHolder.PantryID = 1
	recipePlaceHolder.Title = title

	for _, ingredient := range ingredients {
		foodPlaceHolder = createIngredient(ingredient, recipeId, tx, ctx)
		if foodPlaceHolder.Nutriments.Calories != 0 {
			foods = append(foods, foodPlaceHolder)
		} else {
			return
		}
	}

	for _, food := range foods {
		recipePlaceHolder = calculateNutrimentsRecipe(food, recipePlaceHolder)
	}

	_, err = tx.ExecContext(ctx, "UPDATE recipe SET calories=$1, fat=$2, carbohydrate=$3, protein=$4, misc=$5 WHERE recipe_id=$6", recipePlaceHolder.Calories,
		recipePlaceHolder.Fat, recipePlaceHolder.Carbohydrate, recipePlaceHolder.Protein, pq.Array(recipePlaceHolder.Misc), recipeId)
	if err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	recipe = recipePlaceHolder
	return
}

func createIngredient(ingredient models.Ingredient, recipeId int64, tx *sql.Tx, ctx context.Context, servingSize ...float32) models.Food {
	var emptyFood models.Food
	_, err := tx.ExecContext(ctx, "INSERT INTO recipe_ingredient (food_id, recipe_id, servings) VALUES ($1, $2, $3)", ingredient.FoodID, recipeId, ingredient.Servings)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
	} else {
		err = tx.QueryRowContext(ctx, "SELECT calories, fat, carbohydrate, protein, serving_size, misc FROM food WHERE food_id=$1",
			ingredient.FoodID).Scan(&emptyFood.Nutriments.Calories, &emptyFood.Nutriments.Fat,
			&emptyFood.Nutriments.Carbohydrate, &emptyFood.Nutriments.Protein, &emptyFood.Serving_Size, pq.Array(&emptyFood.Misc))
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
		} else {
			fmt.Println(emptyFood.Misc)
			if emptyFood.Serving_Size == 0 {
				if servingSize != nil {
					emptyFood = calculateNutrimentsRecipeIngredient(emptyFood, ingredient.Servings, servingSize[0])
				} else {
					//default serving size will go to 100g
					emptyFood = calculateNutrimentsRecipeIngredient(emptyFood, ingredient.Servings, 100)
				}
			} else {
				emptyFood = calculateNutrimentsRecipeIngredient(emptyFood, ingredient.Servings, emptyFood.Serving_Size)
			}
		}
	}
	return emptyFood
}

func calculateNutrimentsRecipeIngredient(foodPlaceHolder models.Food, servings float32, servingSize float32) models.Food {
	foodPlaceHolder.Nutriments.Calories = (foodPlaceHolder.Nutriments.Calories * (servingSize / 100)) * servings
	foodPlaceHolder.Nutriments.Fat = (foodPlaceHolder.Nutriments.Fat * (servingSize / 100)) * servings
	foodPlaceHolder.Nutriments.Carbohydrate = (foodPlaceHolder.Nutriments.Carbohydrate * (servingSize / 100)) * servings
	foodPlaceHolder.Nutriments.Protein = (foodPlaceHolder.Nutriments.Protein * (servingSize / 100)) * servings

	return foodPlaceHolder
}

func calculateNutrimentsRecipe(foodPlaceHolder models.Food, recipePlaceHolder models.Recipe) models.Recipe {
	recipePlaceHolder.Calories += foodPlaceHolder.Nutriments.Calories
	recipePlaceHolder.Fat += foodPlaceHolder.Nutriments.Fat
	recipePlaceHolder.Carbohydrate += foodPlaceHolder.Nutriments.Carbohydrate
	recipePlaceHolder.Protein += foodPlaceHolder.Nutriments.Protein
	recipePlaceHolder.Misc = append(recipePlaceHolder.Misc, foodPlaceHolder.Misc...)

	return recipePlaceHolder
}
