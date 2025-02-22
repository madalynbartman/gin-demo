// Importing necessary modules
package main

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"

    _ "go-api/docs" // This is for Go Swagger generated docs
)

// Item represents an item in the inventory
type Item struct {
    Name        string  `json:"name" validate:"required"`
    Price       float64 `json:"price" validate:"required"`
    Description string  `json:"description,omitempty"`
}

// In-memory data store
var inventory = make(map[int]Item)

// Validator instance
var validate *validator.Validate

func main() {
    r := gin.Default()
    validate = validator.New()

    // Swagger documentation route
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Get item by ID
    r.GET("/get-item/:itemId", getItem)

    // Create a new item
    r.POST("/create-item/:itemId", createItem)

    // Update an existing item
    r.PUT("/update-item/:itemId", updateItem)

    // Delete an item
    r.DELETE("/delete-item/:itemId", deleteItem)

    r.Run(":8080")
}

// Get item by ID
// @Summary Get an item by ID
// @Param itemId path int true "Item ID"
// @Success 200 {object} Item
// @Failure 404 {string} string "Item ID not found"
// @Router /get-item/{itemId} [get]
func getItem(c *gin.Context) {
    itemId, err := strconv.Atoi(c.Param("itemId"))
    if err != nil || itemId < 0 {
        c.JSON(http.StatusBadRequest, "Invalid item ID")
        return
    }

    item, exists := inventory[itemId]
    if !exists {
        c.JSON(http.StatusNotFound, "Item ID not found")
        return
    }

    c.JSON(http.StatusOK, item)
}

// Create a new item
// @Summary Create a new item
// @Param itemId path int true "Item ID"
// @Param item body Item true "Item"
// @Success 200 {object} Item
// @Failure 400 {string} string "Validation error or Item ID already exists"
// @Router /create-item/{itemId} [post]
func createItem(c *gin.Context) {
    itemId, err := strconv.Atoi(c.Param("itemId"))
    if err != nil || itemId < 0 {
        c.JSON(http.StatusBadRequest, "Invalid item ID")
        return
    }

    if _, exists := inventory[itemId]; exists {
        c.JSON(http.StatusBadRequest, "Item ID already exists")
        return
    }

    var item Item
    if err := c.ShouldBindJSON(&item); err != nil {
        c.JSON(http.StatusBadRequest, "Invalid JSON")
        return
    }

    if err := validate.Struct(item); err != nil {
        c.JSON(http.StatusBadRequest, "Validation error")
        return
    }

    inventory[itemId] = item
    c.JSON(http.StatusOK, item)
}

// Update an existing item
// @Summary Update an existing item
// @Param itemId path int true "Item ID"
// @Param item body Item true "Item"
// @Success 200 {object} Item
// @Failure 404 {string} string "Item ID does not exist"
// @Failure 400 {string} string "Validation error"
// @Router /update-item/{itemId} [put]
func updateItem(c *gin.Context) {
    itemId, err := strconv.Atoi(c.Param("itemId"))
    if err != nil || itemId < 0 {
        c.JSON(http.StatusBadRequest, "Invalid item ID")
        return
    }

    if _, exists := inventory[itemId]; !exists {
        c.JSON(http.StatusNotFound, "Item ID does not exist")
        return
    }

    var item Item
    if err := c.ShouldBindJSON(&item); err != nil {
        c.JSON(http.StatusBadRequest, "Invalid JSON")
        return
    }

    if err := validate.Struct(item); err != nil {
        c.JSON(http.StatusBadRequest, "Validation error")
        return
    }

    inventory[itemId] = item
    c.JSON(http.StatusOK, item)
}

// Delete an item
// @Summary Delete an item
// @Param itemId path int true "Item ID"
// @Success 200 {string} string "Success"
// @Failure 404 {string} string "Item ID does not exist"
// @Router /delete-item/{itemId} [delete]
func deleteItem(c *gin.Context) {
    itemId, err := strconv.Atoi(c.Param("itemId"))
    if err != nil || itemId < 0 {
        c.JSON(http.StatusBadRequest, "Invalid item ID")
        return
    }

    if _, exists := inventory[itemId]; !exists {
        c.JSON(http.StatusNotFound, "Item ID does not exist")
        return
    }

    delete(inventory, itemId)
    c.JSON(http.StatusOK, "Item deleted")
}
