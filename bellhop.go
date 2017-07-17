package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"strconv"
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

func listServers() (text []string) {

	text = []string{"none", "none2"}
	return
}

func getServerInfo(server string) (text []string) {

	text = []string{server}
	return
}

func reserveServers(number int) (err error) {

	return nil
}

func handleMessage(openSocket *websocket.Conn, message Message) {

	if strings.EqualFold(message.Text, "!ping") {
		go postMessage(openSocket, message, "pong")
	} else if strings.EqualFold(message.Text, "!list") {
		list := listServers()
		for i := 0; i < len(list); i += 1 {
			postMessage(openSocket, message, list[i])
		}
	} else if strings.Contains(message.Text, "!server") {
		server_name := strings.TrimLeft(message.Text, "!server")
		if server_name == "" {
			go postMessage(openSocket, message, "Error: Please provide a server name")
		} else {
			info := getServerInfo(server_name)
			for i := 0; i < len(info); i += 1 {
				postMessage(openSocket, message, info[i])
			}
		}
	} else if strings.Contains(message.Text, "!solo-reserve") {
		err := reserveServers(1)
		if err != nil {
			postMessage(openSocket, message, "Error: No servers available")
		} else {
			postMessage(openSocket, message, "Succesfully allocated 1 server for you!")
		}
	} else if strings.Contains(message.Text, "!reserve") {
		number := strings.TrimLeft(message.Text, "!reserve ")
		i, num_err := strconv.Atoi(number)
		if num_err != nil {
			postMessage(openSocket, message, "Error: Please enter a valid number")
		} else {
			err := reserveServers(i)
			if err != nil {
				postMessage(openSocket, message, "Error: No servers available")
			} else {
				postMessage(openSocket, message, "Succesfully allocated "+number+" servers for you!")
			}
		}
	} else if strings.Contains(message.Text, "!release") {

	} else if strings.Contains(message.Text, "!timer") {

	} else if strings.Contains(message.Text, "!addtime") {

	}

}

/*
 *Commands:
 * Get servers - Returns a list of the servers and their reservations
 * Get server  - Returns more detailed information about a given server
 * Reserve server - If available, reserves a specific server
 * Reserve n servers - If availabe, reserves n servers from the pool randomly
 * Release server(s) - releases the servers one has reserved
 * Query Time - Get info on time limit of server
 * More time - get more time reserved for a server reservation
 */
func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Error: Please provide a valid token\n")
		os.Exit(1)
	}

	// start a websocket-based Real Time API session
	openedWebSocket, _ := connectSlack(os.Args[1])

	for {

		message, err := getMessage(openedWebSocket)
		if err != nil {
			log.Println("read:", err)
			break
		}
		if message.Type == "message" {
			handleMessage(openedWebSocket, message)
		}
	}
}
