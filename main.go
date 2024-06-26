package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type Article struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Comment struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	ArticleId string `json:"articleId"`
}

type Playground struct {
	Id    string `json:"id"`
	Token string `json:"token"`
}

var db *sql.DB

var secret = "secret"

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
	articleID := c.Param("articleID")
	row := db.QueryRow("SELECT Id, Content, ArticleId FROM Comments WHERE Id = ? AND ArticleId = ?", commentID, articleID)

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

func createPlayground(c *gin.Context) {
	id := uuid.New().String()

	dir := "./playgrounds"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			handleError(c, err, http.StatusInternalServerError, "Failed to create directory")
			return
		}
	}

	playgroundDB, err := sql.Open("sqlite3", fmt.Sprintf("./playgrounds/%s.db", id))
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Failed to create playground database")
		return
	}
	defer playgroundDB.Close()

	expiration := time.Now().Add(time.Minute * 30).Unix()

	_, err = playgroundDB.Exec(`
	BEGIN TRANSACTION;
	CREATE TABLE IF NOT EXISTS "Info" (
		"Id"	INTEGER,
		"Expiration"	INTEGER NOT NULL,
		PRIMARY KEY("Id" AUTOINCREMENT)
	);
	INSERT INTO "Info" (Expiration) VALUES (?);
	CREATE TABLE IF NOT EXISTS "Articles" (
		"Id"	INTEGER,
		"Title"	TEXT NOT NULL,
		"Content"	TEXT NOT NULL,
		PRIMARY KEY("Id" AUTOINCREMENT)
	);
	CREATE TABLE IF NOT EXISTS "Comments" (
		"Id"	INTEGER,
		"Content"	TEXT NOT NULL,
		"ArticleId"	INTEGER NOT NULL,
		PRIMARY KEY("Id" AUTOINCREMENT)
	);
	COMMIT;
	`, expiration)
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Failed to create test table")
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["exp"] = expiration

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		handleError(c, err, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	playground := Playground{Id: id, Token: tokenString}
	c.IndentedJSON(http.StatusCreated, playground)
}

func getPlaygroundArticles(c *gin.Context) {
	playgroundID := c.Param("id")
	playgroundDB, err := sql.Open("sqlite3", fmt.Sprintf("./playgrounds/%s.db", playgroundID))
	handleError(c, err, http.StatusInternalServerError, "Failed to open playground database")
	defer playgroundDB.Close()

	rows, err := playgroundDB.Query("SELECT Id, Title, Content FROM Articles")
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

func getPlaygroundArticle(c *gin.Context) {
	playgroundID := c.Param("id")
	articleID := c.Param("articleID")

	playgroundDB, err := sql.Open("sqlite3", fmt.Sprintf("./playgrounds/%s.db", playgroundID))
	handleError(c, err, http.StatusInternalServerError, "Failed to open playground database")
	defer playgroundDB.Close()

	row := playgroundDB.QueryRow("SELECT Id, Title, Content FROM Articles WHERE Id = ?", articleID)

	var article Article
	err = row.Scan(&article.Id, &article.Title, &article.Content)
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

func getPlaygroundArticleComments(c *gin.Context) {
	playgroundID := c.Param("id")
	articleID := c.Param("articleID")

	playgroundDB, err := sql.Open("sqlite3", fmt.Sprintf("./playgrounds/%s.db", playgroundID))
	handleError(c, err, http.StatusInternalServerError, "Failed to open playground database")
	defer playgroundDB.Close()

	rows, err := playgroundDB.Query("SELECT Id, Content, ArticleId FROM Comments WHERE ArticleId = ?", articleID)
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

func getPlaygroundArticleComment(c *gin.Context) {
	playgroundID := c.Param("id")
	commentID := c.Param("commentID")
	articleID := c.Param("articleID")

	playgroundDB, err := sql.Open("sqlite3", fmt.Sprintf("./playgrounds/%s.db", playgroundID))
	handleError(c, err, http.StatusInternalServerError, "Failed to open playground database")
	defer playgroundDB.Close()

	row := playgroundDB.QueryRow("SELECT Id, Content, ArticleId FROM Comments WHERE Id = ? AND ArticleId = ?", commentID, articleID)

	var comment Comment
	err = row.Scan(&comment.ID, &comment.Content, &comment.ArticleId)
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

func getPlaygroundComment(c *gin.Context) {
	playgroundID := c.Param("id")
	commentID := c.Param("commentID")

	playgroundDB, err := sql.Open("sqlite3", fmt.Sprintf("./playgrounds/%s.db", playgroundID))
	handleError(c, err, http.StatusInternalServerError, "Failed to open playground database")
	defer playgroundDB.Close()

	row := playgroundDB.QueryRow("SELECT Id, Content, ArticleId FROM Comments WHERE Id = ?", commentID)

	var comment Comment
	err = row.Scan(&comment.ID, &comment.Content, &comment.ArticleId)
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

func createPlaygroundArticle(c *gin.Context) {
	playgroundID := c.Param("id")
	playgroundDB, err := sql.Open("sqlite3", fmt.Sprintf("./playgrounds/%s.db", playgroundID))
	handleError(c, err, http.StatusInternalServerError, "Failed to open playground database")
	defer playgroundDB.Close()

	var article Article
	err = c.BindJSON(&article)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	result, err := playgroundDB.Exec("INSERT INTO Articles (Title, Content) VALUES (?, ?)", article.Title, article.Content)
	handleError(c, err, http.StatusInternalServerError, "Failed to insert article")

	id, err := result.LastInsertId()
	handleError(c, err, http.StatusInternalServerError, "Failed to get last insert ID")

	article.Id = fmt.Sprintf("%d", id)
	c.IndentedJSON(http.StatusCreated, article)
}

func createPlaygroundArticleComment(c *gin.Context) {
	playgroundID := c.Param("id")
	articleID := c.Param("articleID")
	playgroundDB, err := sql.Open("sqlite3", fmt.Sprintf("./playgrounds/%s.db", playgroundID))
	handleError(c, err, http.StatusInternalServerError, "Failed to open playground database")
	defer playgroundDB.Close()

	var comment Comment
	err = c.BindJSON(&comment)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	result, err := playgroundDB.Exec("INSERT INTO Comments (Content, ArticleId) VALUES (?, ?)", comment.Content, articleID)
	handleError(c, err, http.StatusInternalServerError, "Failed to insert comment")

	id, err := result.LastInsertId()
	handleError(c, err, http.StatusInternalServerError, "Failed to get last insert ID")

	comment.ID = fmt.Sprintf("%d", id)
	comment.ArticleId = articleID
	c.IndentedJSON(http.StatusCreated, comment)
}

func updatePlaygroundArticle(c *gin.Context) {
	playgroundID := c.Param("id")
	articleID := c.Param("articleID")
	playgroundDB, err := sql.Open("sqlite3", fmt.Sprintf("./playgrounds/%s.db", playgroundID))
	handleError(c, err, http.StatusInternalServerError, "Failed to open playground database")
	defer playgroundDB.Close()

	var article Article
	err = c.BindJSON(&article)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	_, err = playgroundDB.Exec("UPDATE Articles SET Title = ?, Content = ? WHERE Id = ?", article.Title, article.Content, articleID)
	handleError(c, err, http.StatusInternalServerError, "Failed to update article")

	article.Id = articleID
	c.IndentedJSON(http.StatusOK, article)
}

func updatePlaygroundArticleComment(c *gin.Context) {
	playgroundID := c.Param("id")
	articleID := c.Param("articleID")
	commentID := c.Param("commentID")
	playgroundDB, err := sql.Open("sqlite3", fmt.Sprintf("./playgrounds/%s.db", playgroundID))
	handleError(c, err, http.StatusInternalServerError, "Failed to open playground database")
	defer playgroundDB.Close()

	var comment Comment
	err = c.BindJSON(&comment)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	_, err = playgroundDB.Exec("UPDATE Comments SET Content = ? WHERE Id = ? AND ArticleId = ?", comment.Content, commentID, articleID)
	handleError(c, err, http.StatusInternalServerError, "Failed to update comment")

	comment.ID = commentID
	comment.ArticleId = articleID
	c.IndentedJSON(http.StatusOK, comment)
}

func deletePlaygroundArticle(c *gin.Context) {
	playgroundID := c.Param("id")
	articleID := c.Param("articleID")
	playgroundDB, err := sql.Open("sqlite3", fmt.Sprintf("./playgrounds/%s.db", playgroundID))
	handleError(c, err, http.StatusInternalServerError, "Failed to open playground database")
	defer playgroundDB.Close()

	_, err = playgroundDB.Exec("DELETE FROM Articles WHERE Id = ?", articleID)
	handleError(c, err, http.StatusInternalServerError, "Failed to delete article")

	c.IndentedJSON(http.StatusNoContent, gin.H{})
}

func deletePlaygroundArticleComment(c *gin.Context) {
	playgroundID := c.Param("id")
	articleID := c.Param("articleID")
	commentID := c.Param("commentID")
	playgroundDB, err := sql.Open("sqlite3", fmt.Sprintf("./playgrounds/%s.db", playgroundID))
	handleError(c, err, http.StatusInternalServerError, "Failed to open playground database")
	defer playgroundDB.Close()

	_, err = playgroundDB.Exec("DELETE FROM Comments WHERE Id = ? AND ArticleId = ?", commentID, articleID)
	handleError(c, err, http.StatusInternalServerError, "Failed to delete comment")

	c.IndentedJSON(http.StatusNoContent, gin.H{})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil // Use the same secret key used to generate the token
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Check if the ID in the token matches the ID in the URL
		id := c.Param("id")
		tokenID, ok := claims["id"].(string)
		if !ok || tokenID != id {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token does not match ID"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func deleteExpiredPlaygrounds() {
	files, err := os.ReadDir("./playgrounds")
	if err != nil {
		fmt.Println("Failed to read playgrounds directory:", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		db, err := sql.Open("sqlite3", fmt.Sprintf("./playgrounds/%s", file.Name()))
		if err != nil {
			fmt.Println("Failed to open playground database:", err)
			continue
		}
		defer db.Close()

		row := db.QueryRow("SELECT Expiration FROM Info")
		var expiration int64
		err = row.Scan(&expiration)
		if err != nil {
			fmt.Println("Failed to fetch expiration:", err)
			continue
		}

		if time.Now().Unix() > expiration {
			err = os.Remove(fmt.Sprintf("./playgrounds/%s", file.Name()))
			if err != nil {
				fmt.Println("Failed to delete playground:", err)
			}
			log.Printf("Deleted expired playground: %s", file.Name())
		}
	}
}

func main() {
	go func() {
		for {
			deleteExpiredPlaygrounds()
			time.Sleep(time.Minute)
		}
	}()
	router := setupRouter()
	router.Run("localhost:8080")
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/articles", getArticles)
	router.GET("/articles/:articleID", getArticle)

	router.GET("/articles/:articleID/comments", getComments)
	router.GET("/articles/:articleID/comments/:commentID", getComment)

	router.POST("/playgrounds", createPlayground)

	playgrounds := router.Group("/playgrounds/:id")
	playgrounds.Use(authMiddleware())
	playgrounds.GET("/articles", getPlaygroundArticles)
	playgrounds.GET("/articles/:articleID", getPlaygroundArticle)

	playgrounds.GET("/articles/:articleID/comments", getPlaygroundArticleComments)
	playgrounds.GET("/articles/:articleID/comments/:commentID", getPlaygroundArticleComment)
	playgrounds.GET("/comments/:commentID", getPlaygroundComment)

	playgrounds.POST("/articles", createPlaygroundArticle)
	playgrounds.POST("/articles/:articleID/comments", createPlaygroundArticleComment)

	playgrounds.PUT("/articles/:articleID", updatePlaygroundArticle)
	playgrounds.PUT("/articles/:articleID/comments/:commentID", updatePlaygroundArticleComment)

	playgrounds.DELETE("/articles/:articleID", deletePlaygroundArticle)
	playgrounds.DELETE("/articles/:articleID/comments/:commentID", deletePlaygroundArticleComment)

	return router
}
