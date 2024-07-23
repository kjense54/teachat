package main

import (
	"testing"
	"reflect"
)

func TestChopText(t *testing.T) {
	var m model
	t.Run("no text", func(t *testing.T) {
		got := m.ChopText("", 10)
		var want []string 

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})
	t.Run("text shorter than wrap width", func (t *testing.T) {
		got := m.ChopText("123456", 10)
		want := []string{"123456"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})
	t.Run("text longer than wrap width", func (t * testing.T) {
		got := m.ChopText("This is a longer test sentence which will be chopped.", 10)
		want := []string{"This is a ", "longer tes", "t sentence", " which wil", "l be chopp", "ed."}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
