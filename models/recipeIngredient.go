package models

type RecipeIngredient struct {
	RecipeIngredientID int64
	RecipeID           int64
	Ingredient         Ingredient
}

type Ingredient struct {
	FoodID   int64
	Servings float32
}
