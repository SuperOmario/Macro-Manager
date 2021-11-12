package handlers

import (
	food "MacroManager/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	//third party packages
	"github.com/gin-gonic/gin"
)

//takes in a barcode and sends a http request to the OpenFoodData API https://world.openfoodfacts.org/data
//adapted from https://blog.logrocket.com/making-http-requests-in-go/
func ScanFood(upc string) food.Food {

	resp, err := http.Get("https://world.openfoodfacts.org/api/v0/product/" + upc)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)
	var foodProduct food.Product
	json.Unmarshal([]byte(sb), &foodProduct)
	return foodProduct.Food
}

func GetFoodProduct(c *gin.Context) {
	upc := c.Param("upc")
	foodProduct := ScanFood(upc)

	c.IndentedJSON(http.StatusOK, foodProduct)
}
