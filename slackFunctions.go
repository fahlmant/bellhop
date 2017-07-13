package main

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"log"
)
type responseRealTimeMessageStart struct {
	Ok    bool         `json:"ok"`
	Error string       `json:"error"`
	Url   string       `json:"url"`
	Self  responseSelf `json:"self"`
}

type responseSelf struct {
	Id string `json:"id"`
}

type Message struct {
	Id      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func startSlack(token string) (wsurl, id string, err error) {
	url := fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", token)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("API request failed with code %d", resp.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}
	var respObj responseRealTimeMessageStart
	err = json.Unmarshal(body, &respObj)
	if err != nil {
		return
	}

	if !respObj.Ok {
		err = fmt.Errorf("Slack error: %s", respObj.Error)
		return
	}

	wsurl = respObj.Url
	id = respObj.Self.Id
	return
}

func connectSlack(token string) (websock *websocket.Conn, id string) {

	websockurl, id, err := startSlack(token)
	if err != nil {
		log.Fatal(err)
	}

	websock, _, err = websocket.DefaultDialer.Dial(websockurl, nil)
	if err != nil {
		log.Fatal(err)
	}

	return websock, id

}
