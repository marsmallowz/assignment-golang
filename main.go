package main

import (
	model "assignment-golang/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dbURI := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	database, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	database.AutoMigrate(&model.Order{}, &model.Item{})
	db = database
}

func main() {
	initDB()
	r := gin.Default()

	// CRUD operations for Order
	r.POST("/orders", createOrder)
	r.GET("/orders", getOrder)
	r.PUT("/orders/:id", updateOrder)
	r.DELETE("/orders/:id", deleteOrder)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}

func createOrder(c *gin.Context) {
	var newOrder model.Order
	if err := c.ShouldBindJSON(&newOrder); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if newOrder.OrderedAt.IsZero() {
		newOrder.OrderedAt = time.Now()
	}
	db.Create(&newOrder)
	c.JSON(200, newOrder)
}

func getOrder(c *gin.Context) {
	var orders []model.Order
	if err := db.Preload("Items").Find(&orders).Error; err != nil {
		c.JSON(404, gin.H{"error": "Error while fetching orders"})
		return
	}
	c.JSON(200, orders)
}

func updateOrder(c *gin.Context) {
	id := c.Params.ByName("id")
	var updatedOrder model.Order
	if err := db.First(&updatedOrder, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Order not found"})
		return
	}

	if err := c.ShouldBindJSON(&updatedOrder); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db.Save(&updatedOrder)
	c.JSON(200, updatedOrder)
}

func deleteOrder(c *gin.Context) {
	id := c.Params.ByName("id")
	var order model.Order
	if err := db.First(&order, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Order not found"})
		return
	}
	db.Delete(&order)
	c.JSON(200, gin.H{"success": "Order deleted"})
}
