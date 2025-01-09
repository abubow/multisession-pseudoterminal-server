package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	clients       map[*websocket.Conn]bool
	handleMessage func(message []byte)
}

func StartServer(handleMessage func(message []byte)) *Server {
	server := Server{
		make(map[*websocket.Conn]bool),
		handleMessage,
	}

	http.HandleFunc("/", server.echo)
	go http.ListenAndServe(":8080", nil)
	log.Println("Listening on :8080")

	return &server
}

func (server *Server) echo(w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil)

	server.clients[connection] = true

	for {
		mt, message, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			break
		}

		log.Println(string(message))

		go server.handleMessage(message)
	}

	delete(server.clients, connection)
	connection.Close()
}

func (server *Server) WriteMessage(message []byte) {
	for conn := range server.clients {
		conn.WriteMessage(websocket.TextMessage, message)
	}
}
