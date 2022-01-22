package handlers

import (
	"MacroManager/controllers"
	"MacroManager/models"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	//third party packages
	"github.com/gin-gonic/gin"
)

//takes in a barcode and sends a http request to the OpenFoodFacts API https://world.openfoodfacts.org/data
//adapted from https://blog.logrocket.com/making-http-requests-in-go/
func ScanFood(upc string) (models.Food, error) {

	resp, err := http.Get("https://world.openfoodfacts.org/api/v0/product/" + upc)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	//maps returned json to food model
	sb := string(body)
	var foodProduct models.Product
	json.Unmarshal([]byte(sb), &foodProduct)

	//error handling for barcodes which are not in the openfoodfacts API
	if foodProduct.Err == "product not found" {
		return foodProduct.Food, errors.New(foodProduct.Err)
	}

	//converts barcode from string to an integer to give food its unique ID
	foodId, err := strconv.ParseInt(upc, 0, 64)
	if err != nil {
		log.Fatal(err)
	} else {
		foodProduct.Food.Barcode = int(foodId)
	}

	foodProduct.Food.PantryID = 1
	return foodProduct.Food, nil
}

func GetFoodProduct(c *gin.Context) {
	upc := c.Param("upc")
	food, err := ScanFood(upc)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, err)
	} else {
		c.IndentedJSON(http.StatusOK, food)
		controllers.InsertFood(food, upc)
	}
}

func GetAllFoodProducts(c *gin.Context) {
	foods := controllers.GetAllFood()
	fmt.Println(foods)
	c.IndentedJSON(http.StatusOK, foods)
}

func GetPantry(c *gin.Context) {
	foods := controllers.GetPantry()
	fmt.Println(foods)
	c.IndentedJSON(http.StatusOK, foods)
}
