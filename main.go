package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {
	r := gin.Default()
	r.LoadHTMLFiles("index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
		// ticker := time.NewTicker(250 * time.Millisecond) // 	repeat 0.25 second
		// quit := make(chan struct{})

		// go func() {
		// 	for {
		// 		select {
		// 		case <-ticker.C:
		// 			for _, player := range SOCKET_LIST {

		// 			}
		// 		case <-quit:
		// 			ticker.Stop()
		// 			return
		// 		}
		// 	}
		// }()
	})

	r.GET("/socket", func(c *gin.Context) {
		websocketHandler(c.Writer, c.Request)
	})

	r.Run("localhost:12312")
}

var SOCKET_LIST []Player
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Player struct {
	X int
	Y int
}

func websocketHandler(write http.ResponseWriter, read *http.Request) {
	conn, err := upgrader.Upgrade(write, read, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade. ", err)
		return
	}

	go func(conn *websocket.Conn) {
		for {

			mType, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}

			var scReq map[string]interface{}
			err = json.Unmarshal([]byte(string(msg)), &scReq)
			if err != nil {
				fmt.Println("Error while unmarshal", err)
			}

			switch scReq["action"] {
			case "init":
				mem := Player{
					rand.Intn(500),
					rand.Intn(500),
				}
				SOCKET_LIST = append(SOCKET_LIST, mem)
				conn.WriteJSON(mem)
			case "move":
				conn.WriteMessage(mType, []byte("pong"))
			default:
				break
			}

			fmt.Println(string(msg))
		}
	}(conn)

}
