package main

import (
	"fmt"
	"multisession-pseudoterminal-server/ws"
	"time"
)

func main() {
	server := ws.StartServer(messageHandler)

	for {
		server.WriteMessage([]byte("Hello"))
		time.Sleep(2 * time.Second)
	}
}

func messageHandler(message []byte) {
	fmt.Println(string(message))
}
