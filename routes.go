package main

import (
	food "MacroManager/handlers"
)

func initialiseRoutes() {
	router.GET("/food/:upc", food.GetFoodProduct)
}
