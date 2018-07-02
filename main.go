package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var SOCKET_LIST = make(map[int]*websocket.Conn)
var PLAYER_LIST = make(map[int]*Player)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Player struct {
	X          int
	Y          int
	ID         int
	PressRight bool
	PressDown  bool
	PressLeft  bool
	PressUp    bool
}

func (player *Player) updatePosition() {
	if player.PressRight {
		player.X++
	}
	if player.PressDown {
		player.Y--
	}
	if player.PressLeft {
		player.X--
	}
	if player.PressUp {
		player.Y++
	}
}

func main() {
	r := gin.Default()
	r.LoadHTMLFiles("index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)

	})

	r.GET("/socket", func(c *gin.Context) {
		websocketHandler(c.Writer, c.Request)
	})

	ticker := time.NewTicker(250 * time.Millisecond) // 	repeat 0.25 second
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				for _, player := range PLAYER_LIST {
					player.updatePosition()
					fmt.Println(*player)
					SOCKET_LIST[player.ID].WriteJSON(player)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	r.Run("localhost:12312")
}

func websocketHandler(write http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(write, req, nil)

	if err != nil {
		fmt.Println("Failed to set websocket upgrade. ", err)
		return
	}

	id := rand.Intn(1000)
	SOCKET_LIST[id] = conn

	// deletePlayer := func(num int) {
	// 	delete(SOCKET_LIST, num)
	// 	delete(PLAYER_LIST, num)
	// }

	go func(conn *websocket.Conn) {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Error while readMessage", err)
				// deletePlayer(id)
				break
			}

			var scReq map[string]interface{}
			err = json.Unmarshal([]byte(string(msg)), &scReq)
			if err != nil {
				fmt.Println("Error while unmarshal", err)
				// deletePlayer(id)
				break
			}

			switch scReq["action"] {
			case "init":
				mem := Player{
					rand.Intn(500),
					rand.Intn(500),
					id,
					false,
					false,
					false,
					false,
				}

				PLAYER_LIST[id] = &mem
				conn.WriteJSON(mem)
			case "move":
				switch scReq["direction"] {
				case "right":
					PLAYER_LIST[id].PressRight = true
				case "down":
					PLAYER_LIST[id].PressDown = true
				case "left":
					PLAYER_LIST[id].PressLeft = true
				case "up":
					PLAYER_LIST[id].PressUp = true
				default:
					break
				}
			default:
				break
			}
		}
	}(conn)
}
