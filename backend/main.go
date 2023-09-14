package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ShortenedUrl struct {
	gorm.Model
	OriginalUrl string `gorm:"unique"`
	ShortUrl    string `gorm:"unique"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=Asia/Bangkok"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&ShortenedUrl{})

	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/shorten", func(c *gin.Context) {
		var data struct {
			Url string `json:"url" binding:"required"`
		}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var link ShortenedUrl
		result := db.Where("original_url = ?", data.Url).First(&link)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				shortUrl := generateShortUrl()
				link = ShortenedUrl{OriginalUrl: data.Url, ShortUrl: shortUrl}
				result = db.Create(&link)
				if result.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
					return
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{"short_url": link.ShortUrl})
	})

	r.GET("/:shortUrl", func(c *gin.Context) {
		shortUrl := c.Param("shortUrl")
		var link ShortenedUrl
		result := db.Where("short_url = ?", shortUrl).First(&link)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			}
			return
		}
		c.Redirect(http.StatusMovedPermanently, link.OriginalUrl)
	})
	r.Run(":8000")
}

func generateShortUrl() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())

	var shortUrl string
	for i := 0; i < 6; i++ {
		shortUrl += string(chars[rand.Intn(len(chars))])
	}

	return shortUrl
}
