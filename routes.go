package main

import (
	"MacroManager/handlers"
)

func initialiseRoutes() {

	//food routes
	//GET
	router.GET("/food", handlers.GetAllFoodProducts)

	//POST
	router.POST("/food/:upc", handlers.GetFoodProduct)

	//diary routes
	//GET

	//POST
	router.POST("/diary/foodEntry", handlers.EnterFood)

	//pantry routes
	//GET
	router.GET("/pantry", handlers.GetPantry)

}
