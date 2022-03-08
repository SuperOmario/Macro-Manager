package models

type RecipeIngredient struct {
	RecipeIngredientID int64
	RecipeID           int64
	Ingredient         Ingredient
}

type Ingredient struct {
	IngredientID int64
	Servings     float32
}

type RemoveIngredient struct {
	IngredientID int64
}

type FoodList struct {
	IngredientIDs []int `json:"IDs"`
}
