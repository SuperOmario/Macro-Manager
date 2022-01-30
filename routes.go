package main

import (
	"MacroManager/handlers"
)

func initialiseRoutes() {

	//admin routes
	//GET
	router.GET("/admin/food", handlers.GetAllFoodProducts)

	//food
	//GET
	router.GET("/food", handlers.GetUserFoodProducts)
	//POST
	router.POST("/food/:upc", handlers.GetFoodProduct)

	//recipe routes
	//GET
	//router.GET("/recipe")
	//POST
	router.POST("/recipe", handlers.CreateRecipe)

	//PATCH
	router.PATCH("/recipe/food/:id", handlers.AddRecipeIngredient)
	router.PATCH("/recipe/details/:id", handlers.UpdateRecipe)

	//DELETE
	router.DELETE("/recipe/ingredient/:id", handlers.RemoveIngredient)
	router.DELETE("/recipe/:id", handlers.DeleteRecipe)

	//diary routes
	//GET

	//POST
	// router.POST("/diary/foodEntry", handlers.EnterFood)

}
