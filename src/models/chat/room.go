package chat

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type room struct {
	// forward is channel to keep massage sending other client
	forward chan []byte
	// join is channel for client who will join the chatroom
	join chan *client
	// leave is channel fo client who will leave the chatroom
	leave chan *client
	// clients keep all joining clients
	clients map[*client]bool
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: messageBufferSize}

func NewRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) Run() {
	for {
		select {
		case client := <-r.join:
			// join
			r.clients[client] = true
		case client := <-r.leave:
			// leave
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// send message for all clients
			for client := range r.clients {
				select {
				case client.send <- msg:
					// send message
				default:
					// if send channel is closed, delete this client
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP: ", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
