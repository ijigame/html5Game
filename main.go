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

type Player struct {
	x int
	y int
}

type SocketRequest struct {
	Action string `json:"action"`
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

			ret := []byte(string(msg))

			fmt.Println("ret ", ret)
			var scReq SocketRequest

			err = json.Unmarshal(ret, &scReq)
			if err != nil {
				fmt.Println("There was error while unmarshal", err)
			}

			switch scReq.Action {
			case "init":
				fmt.Println("init")
				conn.WriteJSON(Player{
					x: rand.Intn(500),
					y: rand.Intn(500),
				})
			case "move":
				fmt.Println("move")
				conn.WriteMessage(mType, []byte("pong"))
			default:
				break
			}

			fmt.Println(string(msg))
		}
	}(conn)

}
