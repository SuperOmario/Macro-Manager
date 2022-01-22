package models

type Recipe struct {
	UserID       int
	PantryID     int
	RecipeID     int
	Title        string
	Calories     float32
	Fat          float32
	Carbohydrate float32
	Protein      float32
	ServingSize  float32
	Misc         []string
}
