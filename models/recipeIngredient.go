package models

type RecipeIngredient struct {
	RecipeIngredientID int
	FoodID             int
	RecipeID           int
	Servings           float32
}
