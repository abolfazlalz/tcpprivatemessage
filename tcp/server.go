package tcp

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	*TCP
}

func NewServer(listener chan *Message, myListener chan *Message, port int, host string) *Server {
	return &Server{
		TCP: new(listener, myListener, port, host),
	}
}

func (t *Server) Listen() error {
	addr := fmt.Sprintf("%s:%d", t.Host, t.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	defer ln.Close()
	t.ln = ln

	go t.acceptLoop()

	<-t.quitch
	return nil
}

func (t *Server) listenForSend(conn net.Conn) {
	for message := range t.myListener {
		if _, err := conn.Write([]byte(message.Text)); err != nil {
			log.Printf("Error in write message: %v\n", err)
		}
	}
}

func (t *Server) acceptLoop() {
	for {
		conn, err := t.ln.Accept()
		go t.listenForSend(conn)
		log.Printf("A user connected %s", conn.RemoteAddr())
		if err != nil {
			log.Printf("accept error: %v", err)
			conn.Close()
			continue
		}

		go t.readLoop(conn)
	}
}

func (t *Server) readLoop(conn net.Conn) {
	buf := make([]byte, 2048)
	defer conn.Close()
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Error in read message from a tcp client: %v", err)
			return
		}

		t.ListenerCh <- &Message{
			Text:     strings.TrimSpace(string(buf[:n])),
			AuthorIP: conn.RemoteAddr().String(),
			conn:     conn,
		}
	}
}

func (c *Message) Send(text string) {
	c.conn.Write([]byte(text + "\n"))
}
