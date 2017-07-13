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

func getMessage(ws *websocket.Conn) (m Message, err error) {
	err = websocket.ReadJSON(ws, &m)
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
		messagetype, message, err := openedWebSocket.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		if messagetype == websocket.TextMessage && strings.Contains(string(message[:]), "<@"+id+">") {
		/*	if strings.Contains(string(message[:]), "!ponging") {
				err = postMessage(openedWebSocket, "pong")
				if err != nil {
					log.Println("write:", err)
					break
				}
			}*/
			fmt.Println(string(message))
		}
	}

}
