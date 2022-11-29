package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/xid"

	"github.com/gin-gonic/gin"
)

type Recipe struct {
	ID string `json:"id"`

	Name string `json:"name"`

	Tags []string `json:"tags"`

	Ingredients []string `json:"ingredients"`

	Instructions []string `json:"instructions"`

	PublishedAt time.Time `json:"publishedAt"`
}

var recipes []Recipe

func init() {
	recipes = make([]Recipe, 0)
	file, _ := ioutil.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(file), &recipes)
}

func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	err := c.ShouldBindJSON(&recipe)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)
}

func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"recipes": recipes,
	})
}

func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	err := c.ShouldBindJSON(&recipe)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	idx := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			idx = i
		}
	}
	if idx == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found"})
		return
	}
	recipes[idx] = recipe
	c.JSON(http.StatusOK, recipe)

}

func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	idx := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found"})
		return
	}
	recipes = append(recipes[:idx], recipes[idx+1:]...)
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been deleted"})

}

func main() {
	r := gin.Default()
	r.POST("/recipes", NewRecipeHandler)
	r.GET("/recipes", ListRecipesHandler)
	//TODO : handle case of empty ID after put execution.
	r.PUT("/recipes/:id", UpdateRecipeHandler)
	r.DELETE("/recipes/:id", DeleteRecipeHandler)
	r.Run()
}
