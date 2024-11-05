package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/meilisearch/meilisearch-go"
)

type Convention struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Module    string `json:"module"`
	Obj       string `json:"obj"`
	Line      int    `json:"line"`
	Column    int    `json:"column"`
	EndLine   int    `json:"endLine"`
	EndColumn int    `json:"endColumn"`
	Path      string `json:"path"`
	Symbol    string `json:"symbol"`
	Message   string `json:"message"`
	MessageID string `json:"message-id"`
}

var client meilisearch.ServiceManager
var index meilisearch.IndexManager

func generateID(path string, endColumn, column, endLine, line int) string {

	data := fmt.Sprintf("%s%d%d%d%d", path, endColumn, column, endLine, line)

	h := sha256.New()

	h.Write([]byte(data))

	return fmt.Sprintf("%x", h.Sum(nil))
}

func initMeiliSearchClient() {

	server := os.Getenv("MEILI_SEARCH_URL")

	key := os.Getenv("MEILI_API_KEY")

	if server == "" {
		server = "http://127.0.0.1:7700"
	}

	client = meilisearch.New(server, meilisearch.WithAPIKey(key))

	drop := os.Getenv("DROP")

	if strings.ToLower(drop) == "true" {
		task, err := client.DeleteIndex("pylint")

		if err != nil {

			log.Printf("Error: %d", err)
		}

		log.Printf("Task ID: %d", task.TaskUID)

		log.Printf("Task Status: %s", task.Status)
	}

	client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        "pylint",
		PrimaryKey: "id",
	})

	index = client.Index("pylint")
}

func uploadData(data []Convention) error {

	for i := range data {
		data[i].ID = generateID(data[i].Path, data[i].EndColumn, data[i].Column, data[i].EndLine, data[i].Line)
	}

	task, err := index.AddDocuments(data, "id")

	if err != nil {
		return err
	}

	log.Printf("Task ID: %d", task.TaskUID)

	log.Printf("Task Status: %s", task.Status)

	return nil
}

func getAllHandler(c *gin.Context) {

	var result meilisearch.DocumentsResult

	err := index.GetDocuments(&meilisearch.DocumentsQuery{
		//Fields: []string{"title", "genres", "rating", "language"},
		//Filter: "(rating > 3 AND (genres = Adventure OR genres = Fiction)) AND language = English",
	}, &result)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func uploadHandler(c *gin.Context) {

	var conventions []Convention

	if err := c.ShouldBindJSON(&conventions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := uploadData(conventions)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data uploaded successfully!"})
}

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Printf("No .env file found, proceeding without it")
	}

	initMeiliSearchClient()

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.POST("/", uploadHandler)

	router.GET("/", getAllHandler)

	port := os.Getenv("PORT")

	if port == "" {
		port = "9999"
	}

	log.Printf("Server is running on http://localhost:%s", port)

	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal(err)
	}
}
