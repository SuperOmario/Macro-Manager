package food

type Food struct {
	FoodID       int
	PantryID     int
	Title        string
	Calories     int
	Fat          int
	Carbohydrate int
	Protein      int
	ServingSize  int
	Misc         []string
}

type SpoonacularFood struct {
	Id               int          `json:"id"`
	Title            string       `json:"title"`
	Badges           []string     `json:"badges"`
	ImportantBadges  []string     `json:"importantBadges"`
	Breadcrumbs      []string     `json:"breadcrumbs"`
	GeneratedText    string       `json:"generatedText"`
	ImageType        string       `json:"imageType"`
	IngredientCount  int          `json:"ingredientCount"`
	IngredientList   string       `json:"ingredientList"`
	Ingredients      []Ingredient `json:"ingredients"`
	Likes            int          `json:"likes"`
	Nutrition        Nutrition    `json:"nutrition"`
	Price            float32      `json:"price"`
	Servings         Serving      `json:"servings`
	SpoonacularScore float32      `json:"spoonacularScore"`
}

type Ingredient struct {
	Description  string `json:"description"`
	Name         string `json:"name"`
	Safety_level string `json:"safety_level"`
}

type Nutrition struct {
	Nutrients        []Nutrient       `json:"nutrients`
	CaloricBreakdown CaloricBreakdown `json:"caloricBreakdown`
}

type Nutrient struct {
	Name                string  `json:"name"`
	Amount              int     `json:"amount"`
	Unit                string  `json:"unit"`
	PercentOfDailyNeeds float32 `json:"percentOfDailyNeeds"`
}

type CaloricBreakdown struct {
	PercentProtein float32 `json:"percentProtein`
	PercentFat     float32 `json:"percentFat"`
	PercentCarbs   float32 `json:"percentCarbs"`
}
type Serving struct {
	Number int    `json:"number"`
	Size   int    `json:"size"`
	Unit   string `json:"unit"`
}
