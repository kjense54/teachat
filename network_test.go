package main

import (
	"testing"
	"reflect"
	"bytes"
	"encoding/gob"
)

type testData struct {
	Username string
	Message string
}

// test encoding/decoding structs with gob
func TestGobEncode(t *testing.T) { 
	var network bytes.Buffer
	enc := gob.NewEncoder(&network) // writes to network
	dec := gob.NewDecoder(&network) // reads from network

	sent := testData{Username: "lemon28", Message: "hello there!"}
	err := enc.Encode(sent)
	if err != nil {
		t.Errorf("Encoder error: %v", err)
	}

	var received testData
	err = dec.Decode(&received)
	if err != nil {
		t.Errorf("Decoder error: %v", err)
	}

if !reflect.DeepEqual(sent, received) {
	t.Errorf("received %q, sent %q", received, sent)
	}
}


