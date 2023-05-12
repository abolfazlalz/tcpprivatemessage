package tcp

import "net"

type Message struct {
	Text     string
	AuthorIP string
	conn     net.Conn
}

type TCP struct {
	Port       int
	Host       string
	ln         net.Listener
	ListenerCh chan *Message
	myListener chan *Message
	quitch     chan struct{}
}
