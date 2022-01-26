package models

type DiaryEntryFood struct {
	DiaryEntryFoodID int64
	DiaryEntryID     int64
	FoodID           int64
	Servings         float32
	Calories         float32
	Fat              float32
	Carbohydrate     float32
	Protein          float32
}
