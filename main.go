package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Playground struct
type Playground struct {
	ID         string    `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
	Articles   []Article `json:"articles"`
}

// Article struct
type Article struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Comments []Comment `json:"comments"`
}

// Comment struct
type Comment struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
