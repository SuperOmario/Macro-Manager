package models

//all the structs needed to take in and manipulate date from the spoonacular api and the database
type Product struct {
	Food Food   `json:"product"`
	Err  string `json:"status_verbose"`
}

//custom structure for db table
//handles unmarshalling of json from OpenFoodData API
type Food struct {
	PantryID   int
	FoodID     int
	Barcode    int
	Title      string `json:"product_name"`
	Nutriments struct {
		Calories     float32 `json:"energy-kcal_100g"`
		Fat          float32 `json:"fat_100g"`
		Carbohydrate float32 `json:"carbohydrates_100g"`
		Protein      float32 `json:"proteins_100g"`
	} `json:"nutriments"`
	Serving_Size float32  `json:"serving_quantity,string"`
	Misc         []string `json:"allergens_tags"`
}
