package models

type Recipe struct {
	UserID       int64
	PantryID     int64
	RecipeID     int64
	Title        string
	Calories     float32
	Fat          float32
	Carbohydrate float32
	Protein      float32
	ServingSize  float32
	Misc         []string
}

type RecipeRequest struct {
	Title       string
	Ingredients []Ingredient
}
