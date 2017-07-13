package main

import (
	"fmt"
	//"github.com/gorilla/websocket"
	"log"
	"os"
	//"strings"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: bellhop slack-bot-token\n")
		os.Exit(1)
	}

	// start a websocket-based Real Time API session
	openedWebSocket, id := connectSlack(os.Args[1])
	fmt.Println("bellhop ready")
	fmt.Println(id)

	for {
		_, message, err := openedWebSocket.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		fmt.Println(string(message))
		/*if messagetype == websocket.TextMessage && strings.Contains(string(message[:]), "<@"+id+">") {
			if strings.Contains(string(message[:]), "!ponging") {
				err = postMessage(openedWebSocket, "pong")
				if err != nil {
					log.Println("write:", err)
					break
				}
			}
		}*/
	}

}
