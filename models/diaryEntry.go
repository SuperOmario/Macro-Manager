package models

type DiaryEntry struct {
	DiaryEntryID int64
	DiaryID      int64
	Date         string
	Calories     float32
	Fat          float32
	Carbohydrate float32
	Protein      float32
}
