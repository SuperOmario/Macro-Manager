package handlers

import (
	"MacroManager/controllers"
	"MacroManager/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllDiaryEntriesForUser(c *gin.Context) {
	diaryEntries, err := controllers.GetAllDiaryEntriesForUser()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
	} else {
		c.IndentedJSON(http.StatusOK, diaryEntries)
	}
}

func GetDiaryEntriesByDate(c *gin.Context) {
	var date models.DiaryDate
	err := c.BindJSON(&date)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
	} else {
		diaryEntries, err := controllers.GetDiaryEntriesByDate(date.Date)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
			fmt.Println(diaryEntries)
			c.IndentedJSON(http.StatusOK, diaryEntries)
		}

	}
}

func CreateDiaryEntry(c *gin.Context) {
	var diaryRequest models.DiaryRequest
	err := c.BindJSON(&diaryRequest)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
	} else {
		diaryEntryID := controllers.InsertDiaryEntry(diaryRequest.RecipeID, diaryRequest.Servings, diaryRequest.Date, diaryRequest.Meal)
		if diaryEntryID == 0 {
			c.IndentedJSON(http.StatusInternalServerError, diaryEntryID)
		} else {
			c.IndentedJSON(http.StatusOK, diaryEntryID)
		}
	}
}

func UpdateDiaryEntry(c *gin.Context) {
	var diaryUpdate models.DiaryUpdate
	ID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
	} else {
		err := c.BindJSON(&diaryUpdate)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, err)
		} else {
			diaryEntry := controllers.UpdateDiaryEntry(ID, diaryUpdate.Servings)
			c.IndentedJSON(http.StatusOK, diaryEntry)
		}
	}
}

func DeleteDiaryEntry(c *gin.Context) {
	dairyEntryId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	fmt.Println(dairyEntryId)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
	} else {
		controllers.DeleteDiaryEntry(dairyEntryId)
		c.IndentedJSON(http.StatusOK, dairyEntryId)
	}
}
