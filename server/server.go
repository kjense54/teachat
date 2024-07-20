package main

import (
	"net"
	"log"
	"fmt"
	"github.com/google/uuid"
	"bufio"
	//"encoding/gob"
)
/*
type Message struct {
	cmd string
	body string
	sender string
	receiver string
}*/

func main() {
	address := "localhost:33183"

	// listen for connections
	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to %q", err)
	}
	fmt.Printf("listening on %q\n", l.Addr().String())

	// channels
	connected := make(chan net.Conn) // new connections
	messages := make(chan string) // incoming messages
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

			go func(c net.Conn, id string) {
				r := bufio.NewReader(c)
				for {
					message, err := r.ReadString('\n')
					if err != nil {
						break
					}
					messages <- fmt.Sprintf("%s : %s", id, message)
				}
				disconnected <- id
			}(c, id)

		case m := <-messages:
			fmt.Println(m)
			for id, user := range onlineUsers {
				go func(user net.Conn, m string) {
					_, err := user.Write([]byte(m))
					if err != nil {
						disconnected <- id
					}
				}(user, m)
			}	

		case u := <-disconnected:
			fmt.Printf("user %s\n disconnected\n", u)
			onlineUsers[u].Close()
			delete(onlineUsers, u)
		}
	}
}

