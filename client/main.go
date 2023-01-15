package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"

	"github.com/ranon-rat/tcp-msg-service/core"
)

var sender *json.Encoder
var reader *json.Decoder

func main() {

	addrs := flag.String("server", "localhost:8080", "")
	//key := flag.String("key", "hello world", "")
	username := flag.String("username", "guest"+strconv.Itoa(rand.Int()), "")
	channel := flag.String("channel", "public channel", "")

	flag.Parse()
	server, err := net.Dial("tcp", *addrs)
	if err != nil {
		panic("the server is closed or something idk")

	}
	sender, reader = json.NewEncoder(server), json.NewDecoder(server)
	sender.Encode(core.Joining{Author: *username, Channel: *channel})

	go func() {
		for {
			fmt.Print(">")
			r := bufio.NewReader(os.Stdin)
			msg, err := r.ReadString('\n')
			if err != nil {
				continue
			}
			sender.Encode(core.Message{
				Author:  *username,
				Msg:     msg,
				Channel: *channel,
			})

		}
	}()
	for {
		var msg core.Message
		if reader.Decode(&msg) != nil {
			continue
		}
		fmt.Printf("\r%s>%s\r>", msg.Author, msg.Msg)

	}

}
