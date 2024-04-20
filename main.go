package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// Article struct
type Article struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Comment struct
type Comment struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	ArticleId string `json:"articleId"`
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./default.db")
	if err != nil {
		panic(err)
	}
}

func handleError(c *gin.Context, err error, status int, message string) {
	if err != nil {
		c.IndentedJSON(status, gin.H{"message": message, "error": err.Error()})
		return
	}
}

func getArticles(c *gin.Context) {
	rows, err := db.Query("SELECT Id, Title, Content FROM Articles")
	handleError(c, err, http.StatusInternalServerError, "Failed to fetch articles")
	defer rows.Close()

	articles := make([]Article, 0)
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.Id, &article.Title, &article.Content)
		handleError(c, err, http.StatusInternalServerError, "Failed to scan articles")
		articles = append(articles, article)
	}

	c.IndentedJSON(http.StatusOK, articles)
}

func getArticle(c *gin.Context) {
	articleID := c.Param("articleID")
	row := db.QueryRow("SELECT Id, Title, Content FROM Articles WHERE Id = ?", articleID)

	var article Article
	err := row.Scan(&article.Id, &article.Title, &article.Content)
	if err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Article not found"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch article", "error": err.Error()})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, article)
}

func getComments(c *gin.Context) {
	articleID := c.Param("articleID")
	rows, err := db.Query("SELECT Id, Content, ArticleId FROM Comments WHERE ArticleId = ?", articleID)
	handleError(c, err, http.StatusInternalServerError, "Failed to fetch comments")
	defer rows.Close()

	comments := make([]Comment, 0)
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.Content, &comment.ArticleId)
		handleError(c, err, http.StatusInternalServerError, "Failed to scan comments")
		comments = append(comments, comment)
	}

	c.IndentedJSON(http.StatusOK, comments)
}

func getComment(c *gin.Context) {
	commentID := c.Param("commentID")
	row := db.QueryRow("SELECT Id, Content, ArticleId FROM Comments WHERE Id = ?", commentID)

	var comment Comment
	err := row.Scan(&comment.ID, &comment.Content, &comment.ArticleId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Comment not found"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch comment", "error": err.Error()})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, comment)
}

func main() {
	router := gin.Default()
	router.GET("/articles", getArticles)
	router.GET("/articles/:articleID", getArticle)

	router.GET("/articles/:articleID/comments", getComments)
	router.GET("/articles/:articleID/comments/:commentID", getComment)

	// router.POST("/playgrounds", createPlayground)
	// router.GET("/playgrounds/:playgroundID/articles", getPlaygroundArticles)
	// router.GET("/playgrounds/:playgroundID/articles/:articleID", getPlaygroundArticle)

	// router.GET("/playgrounds/:playgroundID/articles/:articleID/comments", getPlaygroundArticleComments)
	// router.GET("/playgrounds/:playgroundID/articles/:articleID/comments/:commentID", getPlaygroundArticleComment)

	// router.POST("/playgrounds/:playgroundID/articles", createPlaygroundArticle)
	// router.POST("/playgrounds/:playgroundID/articles/:articleID/comments", createPlaygroundArticleComment)

	// router.PUT("/playgrounds/:playgroundID/articles/:articleID", updatePlaygroundArticle)
	// router.PUT("/playgrounds/:playgroundID/articles/:articleID/comments/:commentID", updatePlaygroundArticleComment)

	// router.DELETE("/playgrounds/:playgroundID/articles/:articleID", deletePlaygroundArticle)
	// router.DELETE("/playgrounds/:playgroundID/articles/:articleID/comments/:commentID", deletePlaygroundArticleComment)

	router.Run("localhost:8080")
}
