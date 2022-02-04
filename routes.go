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
	router.POST("/food", handlers.CreateCustomFood)
	router.POST("/food/:upc", handlers.GetFoodProduct)

	//PATCH
	router.PATCH("food/:id", handlers.UpdateFood)

	//DELETE
	router.DELETE("/food/:id", handlers.DeleteFood)

	//recipe routes
	//GET
	router.GET("recipe", handlers.GetRecipesForUser)
	router.GET("recipe/:id", handlers.GetRecipeById)

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

	//PATCH
	router.PATCH("diary/:id", handlers.UpdateDiaryEntry)

	//DELETE
	router.DELETE("/diary/:id", handlers.DeleteDiaryEntry)

}
