package main

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/ranon-rat/tcp-msg-service/core"
)

var connections = map[string]map[net.Conn]bool{}
var msgPackages = make(chan core.Message)

func readMessages(conn net.Conn, channel string) {
	reader := json.NewDecoder(conn)
	for {
		var msg core.Message
		if reader.Decode(&msg) != nil {
			delete(connections[channel], conn)
			conn.Close()
			break
		}
		msg.Conn = conn
		msgPackages <- msg
	}
}
func handleMessages() {
	for {
		msg := <-msgPackages
		for c := range connections[msg.Channel] {
			if c == msg.Conn {
				continue
			}
			if json.NewEncoder(c).Encode(msg) != nil {
				delete(connections[msg.Channel], c)

			}
		}
		if len(connections[msg.Channel]) == 0 {
			delete(connections, msg.Channel)
		}
	}

}
func main() {
	server, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer server.Close()
	fmt.Println("starting server")
	go handleMessages()
	for {
		conn, err := server.Accept()
		if err != nil {
			continue
		}
		var joining core.Joining
		if err := json.NewDecoder(conn).Decode(&joining); err != nil {
			conn.Close()
			continue
		}

		if _, e := connections[joining.Channel]; !e {
			connections[joining.Channel] = map[net.Conn]bool{}
		}
		msgPackages <- core.Message{Author: "server", Msg: "hey, we have a new member named:" + joining.Author + "\n", Channel: joining.Channel, Conn: conn}
		connections[joining.Channel][conn] = true
		go readMessages(conn, joining.Channel)

	}
}
