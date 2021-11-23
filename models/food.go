package models

//all the structs needed to take in and manipulate date from the spoonacular api and the database
type Product struct {
	Food Food   `json:"product"`
	Err  string `json:"status_verbose"`
}

//custom structure for db table
//handles unmarshalling of json from OpenFoodData API
type Food struct {
	FoodID     int
	PantryID   int
	Title      string `json:"product_name"`
	Nutriments struct {
		Calories     float32 `json:"energy-kcal_serving"`
		Fat          float32 `json:"fat_serving"`
		Carbohydrate float32 `json:"carbohydrates_serving"`
		Protein      float32 `json:"proteins_serving"`
	} `json:"nutriments"`
	ServingSize string   `json:"serving_size"`
	Misc        []string `json:"allergens_tags"`
}
