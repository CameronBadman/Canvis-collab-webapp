package pkg

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for now. In production, restrict this.
		},
	}
	clients      = make(map[string]map[*websocket.Conn]bool)
	clientsMutex = &sync.Mutex{}
)

func HandleWebSocket(c *gin.Context) {
	canvasID := c.Param("canvasId")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}

	addClient(canvasID, conn)
	defer removeClient(canvasID, conn)

	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error reading json: %v", err)
			break
		}

		switch msg["type"] {
		case "draw":
			handleDrawMessage(canvasID, msg)
			broadcastToCanvas(canvasID, msg)
		}
	}
}

func addClient(canvasID string, conn *websocket.Conn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	if clients[canvasID] == nil {
		clients[canvasID] = make(map[*websocket.Conn]bool)
	}
	clients[canvasID][conn] = true
}

func removeClient(canvasID string, conn *websocket.Conn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	if _, ok := clients[canvasID]; ok {
		delete(clients[canvasID], conn)
		conn.Close()
	}
}

func handleDrawMessage(canvasID string, msg map[string]interface{}) {
	ctx := context.Background()
	canvasJSON, err := rdb.Get(ctx, "canvas:"+canvasID).Result()
	if err != nil {
		log.Printf("Failed to get canvas from Redis: %v", err)
		return
	}

	var canvas Canvas
	json.Unmarshal([]byte(canvasJSON), &canvas)
	canvas.Vectors = append(canvas.Vectors, msg["vector"].(Vector))
	canvas.UpdatedAt = time.Now()

	updatedCanvasJSON, _ := json.Marshal(canvas)
	rdb.Set(ctx, "canvas:"+canvasID, updatedCanvasJSON, 0)
}

func broadcastToCanvas(canvasID string, msg map[string]interface{}) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	for client := range clients[canvasID] {
		err := client.WriteJSON(msg)
		if err != nil {
			log.Printf("error broadcasting message: %v", err)
			client.Close()
			delete(clients[canvasID], client)
		}
	}
}

func SyncRedisToMongo() {
	for {
		ctx := context.Background()
		keys, err := rdb.Keys(ctx, "canvas:*").Result()
		if err != nil {
			log.Printf("Failed to get keys from Redis: %v", err)
			continue
		}

		for _, key := range keys {
			canvasID := key[7:] // Remove "canvas:" prefix
			canvasJSON, err := rdb.Get(ctx, key).Result()
			if err != nil {
				continue
			}

			var canvas Canvas
			json.Unmarshal([]byte(canvasJSON), &canvas)

			update := bson.M{
				"$set": bson.M{
					"vectors":   canvas.Vectors,
					"updatedAt": canvas.UpdatedAt,
				},
			}

			_, err = canvasesCollection.UpdateOne(
				ctx,
				bson.M{"canvasId": canvasID},
				update,
			)
			if err != nil {
				log.Printf("Failed to update MongoDB: %v", err)
			}
		}

		time.Sleep(5 * time.Minute)
	}
}
