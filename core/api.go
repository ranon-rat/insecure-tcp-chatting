package core

import "net"

type Message struct {
	Author  string   `json:"author"`
	Msg     string   `json:"msg"`
	Channel string   `json:"channel"`
	Conn    net.Conn `json:"-"`
}
type Joining struct {
	Author  string `json:"author"`
	Channel string `json:"channel"`
}
