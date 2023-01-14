package main

import "strconv"

type GameRoom struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	key string
}

var rooms map[string]*GameRoom

func newGameRoom(key string) *GameRoom {
	return &GameRoom{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		key:        key,
	}
}

func (gr *GameRoom) startGame() {
	var count = 1
	for client := range gr.clients {
		select {
		case client.send <- []byte("start_game:" + strconv.Itoa(count)):
			count++
		default:
			close(client.send)
			delete(gr.clients, client)
		}
	}
}

func (gr *GameRoom) close() {
	close(gr.broadcast)
	close(gr.register)
	close(gr.unregister)
	delete(rooms, gr.key)
}

func (gr *GameRoom) run() {
	for {
		select {
		case client := <-gr.register:
			gr.clients[client] = true
			if len(gr.clients) == 2 {
				gr.startGame()
			}
		case client := <-gr.unregister:
			if _, ok := gr.clients[client]; ok {
				delete(gr.clients, client)
				close(client.send)
			}
			if len(gr.clients) == 0 {
				gr.close()
				return
			}
		case message := <-gr.broadcast:
			for client := range gr.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(gr.clients, client)
				}
			}
		}
	}
}
