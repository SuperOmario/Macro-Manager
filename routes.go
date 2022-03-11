package main

import (
	"MacroManager/handlers"
)

func initialiseRoutes() {

	//admin routes
	//GET
	// router.GET("/admin/food", handlers.GetAllFoodProducts)
	//router.GET("/admin/recipe")

	//food
	//GET
	router.GET("/food", handlers.GetUserFoodProducts)
	router.GET("/food/ingredients", handlers.GetListedFoods)
	// should be a get request but Android doesn't allow GET requests with bodies
	router.POST("/food/ingredients", handlers.GetListedFoods)

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
	router.GET("recipe/ingredients/:id", handlers.GetRecipeIngredientsByID)
	router.GET("recipe/recipes", handlers.GetListedRecipes)
	// should be a get request but Android doesn't allow GET requests with bodies
	router.POST("/recipe/recipes", handlers.GetListedRecipes)

	//POST
	router.POST("/recipe", handlers.CreateRecipe)

	//PATCH
	router.PATCH("/recipe/ingredient/:id", handlers.AddRecipeIngredient)
	router.PATCH("/recipe/ingredients", handlers.UpdateIngredients)
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
