package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var _upgrader_0 = websocket.Upgrader{}
var _upgrader_1 = websocket.Upgrader{}
var _msg_0 = make(chan string)
var _msg_1 = make(chan string)
var _conn_0 *websocket.Conn
var _conn_1 *websocket.Conn
var _id int = 0

func chat_0(w http.ResponseWriter, r *http.Request) {

	_conn_0, err := _upgrader_0.Upgrade(w, r, nil)
	if (err_handle(err)) { return }
	defer _conn_0.Close()

	go read(_conn_0, _msg_0)
	dispatch(_conn_0, _msg_0)
}

func chat_1(w http.ResponseWriter, r *http.Request) {

	_conn_1, err := _upgrader_1.Upgrade(w, r, nil)
	if (err_handle(err)) { return }
	defer _conn_1.Close()

	go read(_conn_1, _msg_1)
	dispatch(_conn_1, _msg_1)
}

func read(conn *websocket.Conn, chnl chan string) {

	for {

		_, _msg, err := conn.ReadMessage()
		if (err_handle(err)) { break }
		log.Println("[RECV]\t\t", _msg)
		chnl <- string(_msg)
	}
}

func dispatch(conn *websocket.Conn, chnl chan string) {

	for {
		err := conn.WriteMessage(websocket.TextMessage, []byte(<-chnl))
		log.Println("[SEND]")
		if (err_handle(err)) { break }
	}
}

func err_handle(err error) bool {

	if (err != nil) {

		log.Println("[ERROR]\t\t", err)
		return true;
	}
	return false
}

func main() {

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/chat_0", chat_0)
	http.HandleFunc("/chat_1", chat_1)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
