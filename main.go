package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Book struct {
	Id     int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int32  `json:"year"`
}

type DbConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type ApiConfig struct {
	ApiHost string
	ApiPort string
}

func GetAllBookHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var books []Book
	if err := db.Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}

	c.JSON(http.StatusOK, books)
}

func PostBookHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := db.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func initDb() (*gorm.DB, error) {
	dbConfig := DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name)
	db, err := gorm.Open(postgres.Open(psqlConn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db = db.Debug()

	err = db.AutoMigrate(&Book{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	router := gin.Default()

	db, err := initDb()
	if err != nil {
		panic(err)
	}

	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	router.GET("/books", GetAllBookHandler)
	router.POST("/books", PostBookHandler)

	apiPort := os.Getenv("API_PORT")
	host := fmt.Sprintf(":%s", apiPort)
	router.Run(host)
}
