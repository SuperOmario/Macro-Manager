package handlers

import (
	food "MacroManager/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	//third party packages
	"github.com/gin-gonic/gin"
)

//takes in a barcode and sends a http request to the spoonacular api with the api key held in the .env file (a file which is kept hidden to store environment variables. This will not be tracked in version control)
//adapted from https://blog.logrocket.com/making-http-requests-in-go/
func ScanFood(upc string) food.SpoonacularFood {
	resp, err := http.Get("https://api.spoonacular.com/food/products/upc/" + upc + "/?apiKey=" + os.Getenv("SPOONACULAR_KEY"))
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)
	var foodProduct food.SpoonacularFood
	json.Unmarshal([]byte(sb), &foodProduct)
	return foodProduct
}

func GetFoodProduct(c *gin.Context) {
	upc := c.Param("upc")
	c.IndentedJSON(http.StatusOK, ScanFood(upc))
}
