package main

import (
	"fmt"
	"log"

	"github.com/abolfazlalz/tcpprivatemessage/tcp"
	"resenje.org/cipher/xor"
)

func messengerAsServer(c chan *tcp.Message, myMessages chan *tcp.Message, port int) {
	srv := tcp.NewServer(c, myMessages, port, "")

	if err := srv.Listen(); err != nil {
		log.Fatal(err)
	}
}

func messengerAsClient(c chan *tcp.Message, myMessages chan *tcp.Message, host string, port int) {
	srv := tcp.NewClient(c, myMessages, port, host)

	if err := srv.Listen(); err != nil {
		log.Fatal("Error: ", err)
	}
}

func main() {
	c := make(chan *tcp.Message)
	myMessages := make(chan *tcp.Message)

	var (
		key        string
		host       string
		port       int
		serverType string
	)

	fmt.Print("Server type: ")
	fmt.Scan(&serverType)
	fmt.Print("Enter key: ")
	fmt.Scan(&key)

	cipher := xor.New([]byte(key))

	go func() {
		for message := range c {
			text, _ := cipher.DecryptString(message.Text)
			log.Println("Message from: ", text)
		}
	}()

	if serverType == "server" {
		fmt.Print("Enter Port: ")
		fmt.Scanf("%d", &port)
		go messengerAsServer(c, myMessages, port)
	} else {
		fmt.Print("Enter Host Port: ")
		fmt.Scanf("%s %d", &host, &port)
		go messengerAsClient(c, myMessages, host, port)
	}

	for {
		var message string
		if _, err := fmt.Scanln(&message); err != nil {
			log.Printf("Can't read user input: %v", err)
			continue
		}

		text, _ := cipher.EncryptString(message)
		myMessages <- &tcp.Message{Text: text, AuthorIP: ""}
	}
}
