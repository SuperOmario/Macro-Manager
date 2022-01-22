package models

type DiaryEntryFood struct {
	DiaryEntryFoodID int
	DiaryEntryID     int
	FoodID           int
	Servings         float32
	Calories         float32
	Fat              float32
	Carbohydrate     float32
	Protein          float32
}
