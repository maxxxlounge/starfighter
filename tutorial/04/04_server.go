package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	guuid "github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var connections map[guuid.UUID]*websocket.Conn

func main() {
	l := &log.Logger{}
	connections = make(map[guuid.UUID]*websocket.Conn)
	r := mux.NewRouter()
	srv := &http.Server{
		Addr:         ":8888",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	r.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		Connect(w, r, l)
	})

	go func() {
		Execute()
	}()
	l.Fatal(srv.ListenAndServe())
}

func Connect(w http.ResponseWriter, r *http.Request, l *log.Logger) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:    4096,
		WriteBufferSize:   4096,
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		l.Fatal(err)
		return
	}
	g := guuid.New()
	connections[g] = c
	defer func(conn *websocket.Conn, g guuid.UUID) {
		delete(connections, g)
		c.Close()
	}(c, g)

	for {
		mType, m, err := c.ReadMessage()
		if err != nil {
			log.Fatal(err)
			return
		}
		if mType == websocket.TextMessage {
			fmt.Println(m)
		}
	}
}

func Execute() {
	fmt.Print("executing")
	for {
		for _, c := range connections {
			msg := "pong"
			err := c.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		time.Sleep(5 * time.Second)
	}
}
