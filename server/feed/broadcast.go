package feed

import (
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Client struct {
	Closed bool
	mu     sync.Mutex
}

type broadcastPostType struct {
	Action string `json:"action"`
	Post   *Post  `json:"post"`
}

var (
	broadcast  = make(chan broadcastPostType)
	clients    = make(map[*websocket.Conn]*Client)
	register   = make(chan *websocket.Conn)
	unregister = make(chan *websocket.Conn)
)

func broadcastPost(msg broadcastPostType) {
	broadcast <- msg
}

func listenToPostBroadcast() {
	for {
		select {
		case connection := <-register:
			clients[connection] = &Client{}
			remoteAddr := connection.RemoteAddr().String()
			localAddr := connection.LocalAddr().String()
			log.Printf("New WebSocket connection from %s to %s", remoteAddr, localAddr)

		case msg := <-broadcast:
			for connection, client := range clients {
				go func(connection *websocket.Conn, client *Client) {
					client.mu.Lock()
					defer client.mu.Unlock()

					if client.Closed {
						return
					}

					err := connection.WriteJSON(msg)
					if err != nil {
						client.Closed = true
						log.Printf("write: %s\n", err)
						connection.Close()
						unregister <- connection
					}
				}(connection, client)
			}

		case connection := <-unregister:
			delete(clients, connection)
			remoteAddr := connection.RemoteAddr().String()
			localAddr := connection.LocalAddr().String()
			log.Printf("Connection disconnected from %s to %s", remoteAddr, localAddr)

		}
	}
}

func BroadcastHandler(c *websocket.Conn) {
	defer func() {
		unregister <- c
		c.Close()
	}()

	// register the client
	register <- c

	select {}
}

func init() {
	// listen for messages on the broadcast channel
	go listenToPostBroadcast()
}
