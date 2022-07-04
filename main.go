package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Message struct {
	UserId  string `json:"userid"`
	Message string `json:"message"`
	Time    int64  `json:"time"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var messages []Message
var clients map[*websocket.Conn]bool

func sendMessages(conn *websocket.Conn) {
	for _, m := range messages {
		conn.WriteJSON(&m)
	}
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			conn.Close()
			delete(clients, conn)
			return
		}

		var message Message
		err = json.Unmarshal(p, &message)
		if err != nil {
			log.Println(err)
			continue
		}

		messages = append(messages, message)
		log.Println(message)
		for c, _ := range clients {
			if err := c.WriteMessage(messageType, p); err != nil {
				log.Println(err)
				conn.Close()
				delete(clients, conn)
			}
		}

	}
}

func ws(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	clients[ws] = true
	log.Println("Client successfully connected..")
	sendMessages(ws)
	go reader(ws)
}

func routes() {
	http.HandleFunc("/ws", ws)
	http.Handle("/", http.FileServer(http.Dir("./spa")))
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
}

func main() {
	clients = make(map[*websocket.Conn]bool)
	messages = make([]Message, 0)
	routes()
	fmt.Println("Server: http://localhost:8000/")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
