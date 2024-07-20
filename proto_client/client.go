package main

import (
	"context"
	"net"
	"log"
	"time"
) 
func main() {
	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := d.DialContext(ctx, "tcp", "localhost:33183")
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}

	// send message to server
		input := "hello there from a user\n"
		if _, err := conn.Write([]byte(string(len(input)) + input)); err != nil {
			log.Fatalf("Error writing to server: %v", err)
		}
	for {
	}
}
