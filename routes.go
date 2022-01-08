package main

import (
	"MacroManager/handlers"
)

func initialiseRoutes() {

	//food routes
	router.GET("/food", handlers.GetAllFoodProducts)
	router.GET("/food/:upc", handlers.GetFoodProduct)

}
