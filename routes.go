package main

import (
	"MacroManager/handlers"
)

func initialiseRoutes() {

	//admin routes
	//GET
	router.GET("/admin/food", handlers.GetAllFoodProducts)
	//router.GET("/admin/recipe")

	//food
	//GET
	router.GET("/food", handlers.GetUserFoodProducts)
	//POST
	router.POST("/food/:upc", handlers.GetFoodProduct)
	//DELETE
	router.DELETE("/food/:id", handlers.DeleteFood)

	//recipe routes
	//GET

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
	router.GET("/diary", handlers.GetAllDiaryEntriesForUser)
	router.GET("/diary/date", handlers.GetDiaryEntriesByDate)

	//POST
	router.POST("/diary", handlers.CreateDiaryEntry)

	//DELETE
	router.DELETE("/diary/:id", handlers.DeleteDiaryEntry)

}
