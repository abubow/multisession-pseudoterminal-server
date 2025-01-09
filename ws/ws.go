package ws

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/creack/pty"
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
	c := exec.Command("bash")
	terminal, err := pty.Start(c)
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			mt, message, err := connection.ReadMessage()

			if err != nil || mt == websocket.CloseMessage {
				break
			}

			terminal.Write([]byte(string(message) + "\n"))

			go server.handleMessage(message)
		}
	}()
	func() {
		reader := bufio.NewReader(terminal)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			connection.WriteMessage(websocket.TextMessage, []byte(line))
		}
	}()

	io.Copy(os.Stdout, terminal)
	delete(server.clients, connection)
	connection.Close()
}
