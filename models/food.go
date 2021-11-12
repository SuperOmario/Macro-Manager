package food

//all the structs needed to take in and manipulate date from the spoonacular api and the database
type Product struct {
	Food Food `json:"product"`
}

//custom structure for db table
//handles unmarshalling of json from OpenFoodData API
type Food struct {
	FoodID    int
	PantryID  int
	Title     string `json:"product_name"`
	Nutriment struct {
		Calories     int `json:"energy-kcal_serving`
		Fat          int `json:"fat_serving"`
		Carbohydrate int `json:"carbohydrates_serving"`
		Protein      int `json:"proteins_serving"`
	} `json:"nutriments"`
	ServingSize string   `json:"serving_size"`
	Misc        []string `json:"allergens_tags"`
}
