package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"strings"
	"sync/atomic"
)

var (
	counter uint64
)

func getMessage(openSocket *websocket.Conn) (m Message, err error) {
	err = websocket.ReadJSON(openSocket, &m)
	return
}

func postMessage(openSocket *websocket.Conn, m Message, text string) (err error) {
	m.Text = text
	m.Id = atomic.AddUint64(&counter, 1)
	err = openSocket.WriteJSON(m)
	return
}

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Error: Please provide a valid token\n")
		os.Exit(1)
	}

	// start a websocket-based Real Time API session
	openedWebSocket, id := connectSlack(os.Args[1])

	for {

		message, err := getMessage(openedWebSocket)
		if err != nil {
			log.Println("read:", err)
			break
		}

		if message.Type == "message" && strings.Contains(message.Text, "<@"+id+">") {
			if strings.Contains(message.Text, "!ping") {
				go func(m Message) {
					m.Text = "pong"
					m.Id = atomic.AddUint64(&counter, 1)
					openedWebSocket.WriteJSON(m)
				}(message)
			}
			fmt.Println(message.Text)
		}
	}
}
