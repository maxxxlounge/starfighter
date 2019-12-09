package main

import (
	"fmt"
	"net/http"
	"time"

	guuid "github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type CustomConn struct {
	Conn *websocket.Conn
	ID   guuid.UUID
}

var connections map[guuid.UUID]*CustomConn

func main() {
	l := &log.Logger{}
	connections = make(map[guuid.UUID]*CustomConn)
	//server
	r := mux.NewRouter()
	srv := &http.Server{
		Addr: ":8888",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	r.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		Connect(w, r, l)
	})

	go func() {
		Execute()
	}()

	l.Infof("start listening on %s", srv.Addr)
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
	cc := CustomConn{
		ID:   g,
		Conn: c,
	}
	connections[g] = &cc
	defer func(conn *websocket.Conn, g guuid.UUID) {
		delete(connections, g)
		c.Close()
	}(c, g)

	for {
		mType, m, err := cc.Conn.ReadMessage()
		if err != nil {
			l.Error(err)
			return
		}
		if mType == websocket.TextMessage {
			l.Infof(" connection %s send message %s", mType, string(m))
		}
	}
}

func Execute() {
	for {
		for u, c := range connections {
			msg := "send message to conn " + u.String()
			err := c.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(msg)
		}
		time.Sleep(1 * time.Second)
	}
}

func (c *CustomConn) ReceiveMessage() {

}
