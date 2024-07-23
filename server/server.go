package main

import (
	"net"
	"log"
	"fmt"
	"github.com/google/uuid"
	"encoding/gob"
)
type Message struct {
	Username string
	Text string
}

func main() {
	gob.Register(Message{})
	address := "localhost:33183"

	// listen for connections
	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to %q", err)
	}
	fmt.Printf("listening on %q\n", l.Addr().String())

	// channels
	connected := make(chan net.Conn) // new connections
	messages := make(chan Message) // incoming messages
	disconnected := make(chan string) // disconected connections

	// accept connections indefinitely 
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Fatalf("Failed to %q", err)
			}
			connected <- conn	
		}
	}()
	
	// handle connections 
	onlineUsers := make(map[string]net.Conn)

	for {
		select {
		case c := <-connected:
			fmt.Printf("connection received\n")
			id := uuid.New().String()
			onlineUsers[id] = c

			// read messages from connection
			go func(c net.Conn, id string) {
				for {
					dec := gob.NewDecoder(c)
					var m Message
					err := dec.Decode(&m)
					if err != nil {
						fmt.Println("Decode: %v", err)
						break
					}
					messages <- m
				}
				disconnected <- id
			}(c, id)

		case m := <-messages:
			switch m.Text {
			case "ping!":
				// do nothing
				default:	
				fmt.Printf("%s: %s", m.Username, m.Text)
				/*
				for id, user := range onlineUsers {
					go func(user net.Conn, m string) {
						_, err := user.Write([]byte(m))
						if err != nil {
							disconnected <- id
						}
					}(user, m)
				}	
				*/
			}

		case u := <-disconnected:
			fmt.Printf("user %s\n disconnected\n", u)
			onlineUsers[u].Close()
			delete(onlineUsers, u)
		}
	}
}

