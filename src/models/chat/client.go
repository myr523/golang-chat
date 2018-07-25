/*	client.go
 * クライアントとチャットルームを定義
 */

package chat

import (
	"github.com/gorilla/websocket"
	"fmt"
)

type client struct {
	// socket for client
	socket *websocket.Conn
	// send is channel for sending massage
	send chan []byte
	// room is chatroom client joined
	room *room
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			fmt.Println(msg)
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
