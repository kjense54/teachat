package main

import (
	"context"
	"log"
	"time"
	"net"
	tea "github.com/charmbracelet/bubbletea"
)

// dial the server
func ConnectToServer() *net.Conn {
	address := "localhost:33183"
	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Minute)
	defer cancel()

	conn, err := d.DialContext(ctx, "tcp", address)	
	if err != nil {
		log.Fatal("Failed to %v", err)
	}
	return &conn
}

// keep the connection to the server from timing out
func KeepAlive(c net.Conn) {
	_, err := c.Write([]byte("ping!"))
	if err != nil {
		log.Fatalf("Connection broken: %v", err)
	}
	time.AfterFunc(30 * time.Second, func() {
		KeepAlive(c)
	})
}

func SendCmd(m model, c net.Conn) tea.Cmd {
	if c == nil {
		log.Fatalf("Connection is nil")
	}
	if m.messageToSend == "" {
		log.Fatalf("messageToSend is nil")
	}
	_, err := c.Write([]byte(m.messageToSend + "\n"))
	if err != nil {
		log.Fatalf("Failed to send message. %v", err)
	}
	return nil
}



// fix incorrect wrapping in viewport by manually resizing strings
func (m model) ChopText(text string, size int) []string {
	if len(text) == 0 {
		return nil
	}
	if len(text) < size {
		return []string{text}
	}
	var chopped []string = make([]string, 0, (len(text)-1)/size+1)
	currentLen := 0
	currentStart := 0
	for i := range text {
		if currentLen == size {
			chopped = append(chopped, text[currentStart:i])
			currentLen = 0
			currentStart = i
		} 
		currentLen++
	}
	// add extra bits at end
	chopped = append(chopped, text[currentStart:])
	return chopped
}
