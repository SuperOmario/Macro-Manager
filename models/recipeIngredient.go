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

type IFRRequest []IngredientForRecipe

type IngredientForRecipe struct {
	RecipeIngredientID int64
	Title              string
	IngredientID       int64
	ServingSize        float32
	Servings           float32
}

type RemoveIngredient struct {
	IngredientID int64
}

type FoodListRequest []FoodList

type FoodList struct {
	IngredientIDs []int `json:"IDs"`
}
