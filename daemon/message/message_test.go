package message

import (
	"testing"
)

var tests = map[string]*Message{
	"gorets:1|c":        &Message{"gorets", 1, Counter, 1},
	"glork:320|ms":      &Message{"glork", 320, Timer, 1},
	"gaugor:333|g":      &Message{"gaugor", 333, Gauge, 1},
	"uniques:765|s":     &Message{"uniques", 765, Set, 1},
	"sampling:1|c|@0.1": &Message{"sampling", 1, Counter, 0.1},
}

func TestParse(t *testing.T) {
	for raw, expected := range tests {
		got, err := Parse([]byte(raw))
		if err != nil {
			t.Fatal(err)
		}
		if got.Equals(expected) {
			t.Errorf("Expected %v got %v\n", expected, got)
		}
	}
}
