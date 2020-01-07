package main

import (
	"fmt"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var _msg = make(chan string)

func main() {
	flag.Parse()
	log.SetFlags(0)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/chat_0"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	go recv(c)
	go send(c)
	scan()
}

func scan() {

	for {

		var _input string
		fmt.Scanln(&_input)
		_msg <- _input
	}
}

func recv(c *websocket.Conn) {

	for {
		_, _message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", _message)
	}
}

func send(c *websocket.Conn) {

	_done := make(chan struct{})
	defer close(_done)

	_interrupt := make(chan os.Signal, 1)
	signal.Notify(_interrupt, os.Interrupt)

	for {

		select {
		case <-_done:
			return
		case msg := <-_msg:
			err := c.WriteMessage(websocket.TextMessage,
					      []byte(msg))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-_interrupt:
			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage,
					      websocket.FormatCloseMessage(
						      websocket.CloseNormalClosure,
						      ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-_done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
