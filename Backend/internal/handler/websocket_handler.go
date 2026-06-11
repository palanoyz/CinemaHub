package handler

import (
	"github.com/palanoyz/cinemahub/internal/websocket"
	"net/http"

	"github.com/gin-gonic/gin"
	gorilla "github.com/gorilla/websocket"
)

var upgrader = gorilla.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWS(hub *websocket.Hub) gin.HandlerFunc {

	return func(c *gin.Context) {
		showtimeID := c.Param("id")
		if showtimeID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Showtime ID required"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		hub.AddClient(showtimeID, conn)

		// Keep connection alive and listen for closure
		defer hub.RemoveClient(showtimeID, conn)

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}
}
