package client

import (
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

func RunWebSocketClient() {
	ws, err := websocket.Dial("ws://localhost:8080/ws", "", "http://localhost/")
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	for {
		var availableSpots map[string]int
		if err := websocket.JSON.Receive(ws, &availableSpots); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Available spots: %v\n", availableSpots)
	}
}
