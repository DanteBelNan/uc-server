package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DanteBelNan/uc-server/internal/adapters/websocket"
	"github.com/DanteBelNan/uc-server/internal/core/services"
	"github.com/gin-gonic/gin"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Inicializacion de dependencias
	roomService := services.NewRoomService()
	clipboardHandler := websocket.NewClipboardHandler(roomService)

	router := gin.Default()

	// Endpoints
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"service": "uc-server",
		})
	})

	router.GET("/ws", func(c *gin.Context) {
		clipboardHandler.HandleConnection(c)
	})

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on %s", addr)
	
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}