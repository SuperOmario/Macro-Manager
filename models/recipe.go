package models

type Recipe struct {
	UserID      int64
	RecipeID    int64
	Title       string
	ServingSize float32
}

type RecipeRequest struct {
	Title       string
	ServingSize float32
	Ingredients []Ingredient
}

type RecipeUpdate struct {
	Title       string
	ServingSize float32
}
