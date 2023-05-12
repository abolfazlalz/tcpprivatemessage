package tcp

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Client struct {
	*TCP
}

func NewClient(c chan *Message, myListener chan *Message, port int, host string) *Client {
	return &Client{
		TCP: new(c, myListener, port, host),
	}
}

func (c *Client) Listen() error {
	dial, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port))
	if err != nil {
		return err
	}

	go c.sendChan(dial)
	go c.readLoop(dial)
	if err != nil {
		return err
	}
	return nil
}

func (t *Client) readLoop(conn net.Conn) {
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

func (c *Client) sendChan(dial net.Conn) {
	for message := range c.myListener {
		_, err := dial.Write([]byte(message.Text))
		if err != nil {
			return
		}
	}
}
