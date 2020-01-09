package main

import (
	"fmt"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8011", "http service address")
var debug = flag.Bool("debug", false, "is debug mode")
var sensitive = []string{ "fuck", "shit" }

type Client struct {

	url			string
	conn		*websocket.Conn
	recv		chan string
	send		chan string
	interrupt	chan os.Signal
}

func mask(list []string, msg string) string {

	for _, s := range list {

		if r := strings.Index(msg, s); r != -1 {

			return msg[0:r]+ strings.Repeat("*", len(s))+ msg[r+len(s):]
		}
	}
	return msg
}

func err_detect(err error, msg string) bool {

	if err != nil {

		log.Println("[ERROR]\t\t", msg)
		return true
	}
	return false
}

func (client *Client) connect(endpoint string) {

	url				:= url.URL{ Scheme: "ws", Host: *addr, Path: endpoint }
	client.url		= url.String();
	conn, _, err	:= websocket.DefaultDialer.Dial(client.url, nil)
	if err_detect(err, "connect error") { return; }
	client.conn		= conn
	if *debug { log.Println("[CONN]\t", client.url) }
}

func (client *Client) scan() {

	var input string
	for {

		_, err := fmt.Scanln(&input)
		if err_detect(err, "scan error") { return; }
		client.send <- input
		if *debug { log.Println("[SCAN]\t", input) }
	}
}

func (client *Client) write() {

	signal.Notify(client.interrupt, os.Interrupt)
	for {

		select {
		case msg := <-client.send:

			err := client.conn.WriteMessage(websocket.TextMessage,[]byte(msg))
			if err_detect(err, "write error") { return; }
			if *debug { log.Println("[WRITE]\t", msg) }

		case <-client.interrupt:

			log.Println("interrupt")
			err := client.conn.WriteMessage(websocket.CloseMessage,
											websocket.FormatCloseMessage(
												websocket.CloseNormalClosure,
											""))
			if err_detect(err, "interrupt error") { return; }
			client.conn.Close()
			return
		}
	}
}

func (client *Client) read() {

	for {

		_, msg, err := client.conn.ReadMessage()
		if err_detect(err, "read error") { return; }
		client.recv <- string(msg)
		if *debug { log.Println("[READ]\t", string(msg)) }
	}
}

func (client *Client) show() {

	for {

		select {
		case msg := <-client.recv:

			msg = mask(sensitive, msg)
			_, err := fmt.Println(msg)
			if err_detect(err, "show error") { return; }
			if *debug { log.Println("[SHOW]\t", msg) }
		}
	}
}

func main() {

	flag.Parse()
	log.SetFlags(0)
	if *debug { log.Println("In debug mode\n") }

	client := Client {

		url:		"",
		conn:		nil,
		recv:		make(chan string),
		send:		make(chan string),
		interrupt:	make(chan os.Signal),
	}

	client.connect("/ws")
	go client.scan()
	go client.write()
	go client.read()
	client.show()
}

