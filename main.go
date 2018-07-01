package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {
	r := gin.Default()
	r.LoadHTMLFiles("index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/socket", func(c *gin.Context) {
		websocketHandler(c.Writer, c.Request)
	})

	r.Run("localhost:12312")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func websocketHandler(write http.ResponseWriter, read *http.Request) {
	conn, err := upgrader.Upgrade(write, read, nil)

	if err != nil {
		fmt.Println("Failed to set websocket upgrade. ", err)
		return
	}

	go func(conn *websocket.Conn) {
		for {
			mType, _, err := conn.ReadMessage()

			if err != nil {
				break
			}

			conn.WriteMessage(mType, []byte("pong"))
		}
	}(conn)

}
