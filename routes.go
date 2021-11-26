package main

import (
	"MacroManager/handlers"
)

func initialiseRoutes() {
	router.GET("/food/:upc", handlers.GetFoodProduct)
}
