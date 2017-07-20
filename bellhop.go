package main

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
	"log"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
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

func reserveServer(server string, cli *clientv3.Client) (err error) {

	key := server + "/reserved"
	_, err = cli.Put(context.TODO(), key, "true")
	return nil
}

func releaseServer(server string, cli *clientv3.Client) (err error) {

	key := server + "/reserved"
	_, err = cli.Put(context.TODO(), key, "false")
	return nil
}

func getTimer(server string, cli *clientv3.Client) (time int, err error) {

	key := server + "/time"
	resp, err := cli.Get(context.TODO(), key)
	time_str := resp.Kvs[0].Value
	time, _ = strconv.Atoi(string(time_str))
	return time, nil
}

func addTime(server string, amount int64, cli *clientv3.Client) (err error) {

	resp, err := cli.Grant(context.TODO(), amount)
	_, err = cli.Put(context.TODO(), server, "time", clientv3.WithLease(resp.ID))
	return nil
}

func handleMessage(openSocket *websocket.Conn, message Message, cli *clientv3.Client) {

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
	} else if strings.Contains(message.Text, "!reserve") {
		server_name := strings.TrimLeft(message.Text, "!reserve ")
		if server_name == "" {
			go postMessage(openSocket, message, "Error: Please provide a server name")
		} else {
			err := reserveServer(server_name, cli)
			if err != nil {
				postMessage(openSocket, message, "Something went wrong. Is the server name valid?")
			} else {
				postMessage(openSocket, message, "Succesfully reserved "+server_name)
			}
		}
	} else if strings.Contains(message.Text, "!release") {
		server_name := strings.TrimLeft(message.Text, "!release ")
		if server_name == "" {
			go postMessage(openSocket, message, "Error: Please provide a server name")
		} else {
			err := releaseServer(server_name, cli)
			if err != nil {
				postMessage(openSocket, message, "Something went wrong. Do you own the server? Is the server name valid?")
			} else {
				postMessage(openSocket, message, "Succesfully released "+server_name)
			}
		}
	} else if strings.Contains(message.Text, "!timer") {
		server_name := strings.TrimLeft(message.Text, "!timer ")
		if server_name == "" {
			go postMessage(openSocket, message, "Error: Please provide a server name")
		} else {
			time_left, _ := getTimer(server_name, cli)
			postMessage(openSocket, message, ""+server_name+" has "+strconv.Itoa(time_left)+" minutes left")
		}
	} else if strings.Contains(message.Text, "!addtime") {
		args := strings.Split(message.Text, " ")
		time, _ := strconv.Atoi(args[2])
		time64 := int64(time)
		err := addTime(args[1], time64, cli)
		if err != nil {
			postMessage(openSocket, message, "Something went wrong.")
		} else {
			postMessage(openSocket, message, "Succesfully added "+args[2]+" minutes to "+args[1])
		}
	}

}

/*
 *Commands:
 * !list                     - Returns a list of the servers and their reservations
 * !server   <name>          - Returns more detailed information about a given server
 * !reserve  <num>           - If availabe, reserves n servers from the pool randomly
 * !release  <name>          - releases the server reserved if requestor is the owner
 * !timer    <name>          - Get info on time limit of server
 * !addtime  <name> <amount> - get more time reserved for a server reservation
 */
func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Error: Please provide a valid token\n")
		os.Exit(1)
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// start a websocket-based Real Time API session
	openedWebSocket, _ := connectSlack(os.Args[1])
	for {

		message, err := getMessage(openedWebSocket)
		if err != nil {
			log.Println("read:", err)
			break
		}
		if message.Type == "message" {
			handleMessage(openedWebSocket, message, cli)
		}
	}

	defer cli.Close()
}
