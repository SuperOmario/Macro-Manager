package handlers

import (
	"MacroManager/controllers"
	"MacroManager/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//create a new recipe with one or more ingredients
func CreateRecipe(c *gin.Context) {
	var recipeRequest models.RecipeRequest
	err := c.BindJSON(&recipeRequest)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
	} else {
		recipe, err := controllers.InsertRecipe(recipeRequest.Title, recipeRequest.ServingSize, recipeRequest.Ingredients...)
		if err != nil {
			log.Fatal(err)
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
			c.IndentedJSON(http.StatusOK, recipe)
		}
	}
}

//add a single ingredient to a recipe
func AddRecipeIngredient(c *gin.Context) {
	var ingredient models.Ingredient
	recipeId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
	} else {
		err = c.BindJSON(&ingredient)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, err)
		} else {
			recipe, err := controllers.AddRecipeIngredient(recipeId, ingredient)
			if err != nil {
				c.IndentedJSON(http.StatusBadRequest, err)
			} else {
				c.IndentedJSON(http.StatusOK, recipe)
			}
		}
	}
}

func UpdateRecipe(c *gin.Context) {
	recipeId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	var recipeRequest models.RecipeUpdate
	var recipe models.Recipe

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
	} else {
		err = c.BindJSON(&recipeRequest)
		if err != nil || recipeRequest.ServingSize == 0 {
			c.IndentedJSON(http.StatusBadRequest, err)
		} else {
			err := controllers.UpdateRecipe(recipeId, recipeRequest.Title, recipeRequest.ServingSize)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
			} else {
				recipe.RecipeID = recipeId
				recipe.ServingSize = recipeRequest.ServingSize
				recipe.Title = recipeRequest.Title
				c.IndentedJSON(http.StatusOK, recipe)
			}
		}
	}
}

func RemoveIngredient(c *gin.Context) {
	recipeId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	var ingredient models.RemoveIngredient

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
	} else {
		err = c.BindJSON(&ingredient)
		if err != nil {
			log.Fatal(err)
			c.IndentedJSON(http.StatusBadRequest, err)
		} else {
			err := controllers.RemoveIngredient(recipeId, ingredient.IngredientID)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
			} else {
				c.IndentedJSON(http.StatusOK, recipeId)
			}
		}
	}
}

func DeleteRecipe(c *gin.Context) {
	recipeId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
	} else {
		err = controllers.DeleteRecipe(recipeId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
			c.IndentedJSON(http.StatusOK, recipeId)
		}
	}
}

func GetRecipesForUser(c *gin.Context) {
	recipes, err := controllers.GetRecipesForUser(1)
	if err != nil {
		c.IndentedJSON(http.StatusNoContent, err)
	} else {
		c.IndentedJSON(http.StatusOK, recipes)
	}
}

func GetRecipeById(c *gin.Context) {
	recipeId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
	} else {
		recipe, err := controllers.GetRecipeById(recipeId)
		if err != nil {
			c.IndentedJSON(http.StatusNoContent, err)
		} else {
			c.IndentedJSON(http.StatusOK, recipe)
		}
	}
}

func GetRecipeIngredientsByID(c *gin.Context) {
	recipeId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
	} else {
		ingredients, err := controllers.GetRecipeIngredientsByID(recipeId)
		if err != nil {
			c.IndentedJSON(http.StatusNoContent, err)
		} else {
			c.IndentedJSON(http.StatusOK, ingredients)
		}
	}
}

func UpdateIngredients(c *gin.Context) {
	var ingredients models.IFRRequest
	err := c.BindJSON(&ingredients)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
	} else {
		err = controllers.UpdateIngredients(ingredients)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
			c.IndentedJSON(http.StatusOK, ingredients)
		}
	}
}

func GetListedRecipes(c *gin.Context) {
	// var ids models.FoodList
	var ids []models.RecipeList
	err := c.BindJSON(&ids)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
	} else {
		foods, err := controllers.GetListedRecipes(ids[0].RecipeIDs)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(http.StatusOK, foods)
	}
}
