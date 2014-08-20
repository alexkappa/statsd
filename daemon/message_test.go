package daemon

import (
	"bufio"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	var tests = map[string]*Message{
		"gorets:1|c\n":        &Message{"gorets", 1, Counter, 1},
		"glork:320|ms\n":      &Message{"glork", 320, Timer, 1},
		"gaugor:333|g\n":      &Message{"gaugor", 333, Gauge, 1},
		"uniques:765|s\n":     &Message{"uniques", 765, Set, 1},
		"sampling:1|c|@0.1\n": &Message{"sampling", 1, Counter, 0.1},
	}
	for raw, expected := range tests {
		got, err := Parse([]byte(raw))
		if err != nil {
			t.Fatal(err)
		}
		if *got != *expected {
			t.Errorf("Expected %v got %v\n", expected, got)
		}
	}
}

func TestScanSplit(t *testing.T) {
	var tests = map[string][]string{
		"gorets:1|c\n":        []string{"gorets", "1", Counter},
		"glork:320|ms\n":      []string{"glork", "320", Timer},
		"gaugor:333|g\n":      []string{"gaugor", "333", Gauge},
		"uniques:765|s\n":     []string{"uniques", "765", Set},
		"sampling:1|c|@0.1\n": []string{"sampling", "1", Counter, "0.1"},
	}
	for raw, tokens := range tests {
		scanner := bufio.NewScanner(strings.NewReader(raw))
		scanner.Split(ScanStat)
		for _, token := range tokens {
			if !scanner.Scan() {
				t.Error("Scanner finished but more tokens are expected")
			}
			if token != string(scanner.Bytes()) {
				t.Errorf("Unexpected token %s. Expected %s\n", scanner.Bytes(), token)
			}
		}
	}
}
