package models

//all the structs needed to take in and manipulate date from the spoonacular api and the database
type Product struct {
	Food Food   `json:"product"`
	Err  string `json:"status_verbose"`
}

//custom structure for db table
//handles unmarshalling of json from OpenFoodData API
type Food struct {
	UserID       int64
	IngredientID int64
	Barcode      int64
	Title        string `json:"product_name"`
	Nutriments   struct {
		Calories     float32 `json:"energy-kcal_100g"`
		Fat          float32 `json:"fat_100g"`
		Carbohydrate float32 `json:"carbohydrates_100g"`
		Protein      float32 `json:"proteins_100g"`
	} `json:"nutriments"`
	Serving_Size float32  `json:"serving_quantity,string"`
	Misc         []string `json:"allergens_tags"`
}

type FoodUpdate struct {
	Title        string
	Calories     float32
	Fat          float32
	Carbohydrate float32
	Protein      float32
	ServingSize  float32
}

type CustomFood struct {
	Title        string
	Calories     float32
	Fat          float32
	Carbohydrate float32
	Protein      float32
	Serving_Size float32
	Misc         []string
}
