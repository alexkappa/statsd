package daemon

import (
	"bufio"
	"bytes"
	"strconv"
)

const (
	Gauge   string = "g"
	Counter string = "c"
	Timer   string = "ms"
	Set     string = "s"
)

// The Raw type is a wrapper for byte slice.
type Raw []byte

// To satisfy the `io.Reader` interface.
func (r Raw) Read(p []byte) (n int, err error) {
	for i := 0; i < len(r); i++ {
		p[i] = r[i]
	}
	return len(r), nil
}

// The Message type represents a parsed message.
type Message struct {
	Bucket   string
	Value    int
	Modifier string
	Sampling float32
}

// Consumes a token into the messages bucket property.
func ConsumeBucketToken(m *Message, b []byte) error {
	m.Bucket = string(b)
	return nil
}

// Consumes a token into the messages value property.
func ConsumeValueToken(m *Message, b []byte) (err error) {
	m.Value, err = strconv.Atoi(string(b))
	return err
}

// Consumes a token into the messages modifier property.
func ConsumeModifierToken(m *Message, b []byte) error {
	m.Modifier = string(b)
	return nil
}

// Consumes a token into the messages sampling property.
func ConsumeSamplingToken(m *Message, b []byte) error {
	f, err := strconv.ParseFloat(string(b[1:]), 32)
	if err != nil {
		m.Sampling = 1
		return err
	}
	m.Sampling = float32(f)
	return err
}

func isSplitChar(b byte) bool {
	switch b {
	case ':', '|', '@', '\n':
		return true
	}
	return false
}

// Used to split a byte slice into tokens with `bufio.Scanner.Split`.
func ScanStat(data []byte, atEOF bool) (advance int, token []byte, err error) {
	length := len(data)
	start := 0
	for ; start < length; start++ {
		if !isSplitChar(data[start]) {
			break
		}
	}
	if atEOF && length == 0 {
		return 0, nil, nil
	}
	for i := start; i < length; i++ {
		if isSplitChar(data[i]) {
			return i + 1, data[start:i], nil
		}
	}
	if atEOF && length > start {
		return length, data[start:], nil
	}
	return 0, nil, nil
}

// The consumer type defines a function signature that can be used when
// consuming message tokens.
type Consumer func(*Message, []byte) error

// We create this map to represent the order in which tokens should be consumed.
// The scanner increments the index every time it scans a token, calling a
// different consumer each time.
var consumers = map[int]Consumer{
	0: ConsumeBucketToken,
	1: ConsumeValueToken,
	2: ConsumeModifierToken,
	3: ConsumeSamplingToken,
}

// Parses a `Raw` message to a `Message`.
func Parse(b []byte) (*Message, error) {
	scanner := bufio.NewScanner(bytes.NewReader(b))
	scanner.Split(ScanStat)
	i := 0
	m := new(Message)
	m.Sampling = 1
	for scanner.Scan() {
		err := consumers[i](m, scanner.Bytes())
		if err != nil {
			return nil, err
		}
		i++
	}
	return m, nil
}
