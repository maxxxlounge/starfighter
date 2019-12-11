package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func main() {
	var err error
	flag.Parse()

	u := url.URL{Scheme: "ws", Host: "localhost:8888", Path: "/connect"}
	log.Printf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()
}

func SendMessage(c *websocket.Conn, msg string) {
	fmt.Println("send message to server " + msg)
	err := c.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Println("write:", err)
		return
	}
}

func ReceiveMessage(conn *websocket.Conn) {
	mt, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}
	if mt == websocket.TextMessage {
		fmt.Printf("recv: %s", message)
	}
}
