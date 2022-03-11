package models

type Recipes []Recipe

type Recipe struct {
	UserID      int64
	RecipeID    int64
	Title       string
	ServingSize float32
}

type RecipeDetails struct {
	RecipeID     int64
	Title        string
	ServingSize  float32
	Calories     float32
	Fat          float32
	Carbohydrate float32
	Protein      float32
	Misc         []string
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

type RecipeListRequest []FoodList

type RecipeList struct {
	RecipeIDs []int `json:"IDs"`
}
