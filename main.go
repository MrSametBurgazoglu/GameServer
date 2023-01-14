package main

import (
	"fmt"
	"log"
	"net/http"
)

func createGameRoom(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	roomID := GetRoomID(3)
	gameRoom := newGameRoom(roomID)
	go gameRoom.run()
	rooms[roomID] = gameRoom
	fmt.Fprintf(w, "room_id:%s", roomID)
}

func main() {
	rooms = make(map[string]*GameRoom)
	http.HandleFunc("/", createGameRoom)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		joinGameRoom(w, r)
	})
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
