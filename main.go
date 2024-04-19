package main

import (
	"encoding/json"
	"net/http"
	"os"
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

// Get Articals from file
func getArticlesFromFile() ([]Article, error) {
	file, err := os.Open("articles.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	articles := []Article{}
	err = json.NewDecoder(file).Decode(&articles)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

// GetArticles
func getArticles(c *gin.Context) {
	articles, err := getArticlesFromFile()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, articles)
}

// GetArticle by ID
func getArticle(c *gin.Context) {
	articleID := c.Param("articleID")
	articles, err := getArticlesFromFile()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	for _, article := range articles {
		if article.ID == articleID {
			c.IndentedJSON(http.StatusOK, article)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Article not found"})
}

// GetComments by article ID
func getComments(c *gin.Context) {
	articleID := c.Param("articleID")
	articles, err := getArticlesFromFile()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	for _, article := range articles {
		if article.ID == articleID {
			c.IndentedJSON(http.StatusOK, article.Comments)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Article not found"})
}

// GetComment by article ID and comment ID
func getComment(c *gin.Context) {
	articleID := c.Param("articleID")
	commentID := c.Param("commentID")
	articles, err := getArticlesFromFile()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	for _, article := range articles {
		if article.ID == articleID {
			for _, comment := range article.Comments {
				if comment.ID == commentID {
					c.IndentedJSON(http.StatusOK, comment)
					return
				}
			}
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Comment not found"})
}

func main() {
	router := gin.Default()
	router.GET("/articles", getArticles)
	router.GET("/articles/:articleID", getArticle)

	router.GET("/articles/:articleID/comments", getComments)
	router.GET("/articles/:articleID/comments/:commentID", getComment)

	router.Run("localhost:8080")
}
