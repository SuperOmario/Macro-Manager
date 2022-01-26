package handlers

import (
	"MacroManager/controllers"
	"MacroManager/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateRecipe(c *gin.Context) {
	fmt.Println("Request received")
	var recipeRequest models.RecipeRequest
	fmt.Println("variable initialised")
	err := c.BindJSON(&recipeRequest)
	fmt.Println("bind called")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
	} else {
		recipe, err := controllers.InsertRecipe(recipeRequest.Title, recipeRequest.Ingredients...)
		if err != nil {
			log.Fatal(err)
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
			c.IndentedJSON(http.StatusOK, recipe)
		}
	}
}
