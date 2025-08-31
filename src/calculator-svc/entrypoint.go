package calculatorsvc

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func calculateHandler(c *gin.Context) {
	var req CalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	result, err := calculate(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

func Entrypoint() {
	r := gin.Default()
	r.POST("/calculator", calculateHandler)

	log.Println("Starting calculator service on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
