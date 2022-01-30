package models

type DiaryEntry struct {
	DiaryEntryID int64
	UserID       int64
	RecipeID     int64
	Date         string
	Meal         string
	Calories     float32
	Fat          float32
	Carbohydrate float32
	Protein      float32
	Servings     float32
	Misc         []string
}
