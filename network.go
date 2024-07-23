package main

import (
	"net"
	"context"
	"time"
	"log"
	"encoding/gob"
	tea "github.com/charmbracelet/bubbletea"
)

// dial the server
func ConnectToServer() net.Conn {
	address := "localhost:33183"

	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Minute)
	defer cancel()

	conn, err := d.DialContext(ctx, "tcp", address)	
	if err != nil {
		log.Fatalf("Failed to %v", err)
	}
	return conn
}

// keep the connection to the server from timing out
func KeepAlive(c net.Conn) {
	p := Message{Username: "", Text: "Ping!"}
	enc := gob.NewEncoder(c)
	err := enc.Encode(p)
	if err != nil {
		fmt.Printf("KeepAlive Encoder: %v\n", err)
		c.Close()
	}
	time.AfterFunc(30 * time.Second, func() {
		KeepAlive(c)
	})
}

type Message struct {
	Username string
	Text string
}
// send info from our model to the server 
func SendCmd(m Message, c net.Conn) tea.Cmd {
	if c == nil {
		log.Fatalf("Connection is nil")
	}

	enc := gob.NewEncoder(c)
	err := enc.Encode(m) 
	if err != nil {
		log.Fatalf("SendCmd Encoder: %v", err)
	}

	return nil
}
