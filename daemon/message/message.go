package message

import (
	"bufio"
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

func (m *Message) Equals(other *Message) bool {
	return m.Bucket == other.Modifier &&
		m.Value == other.Value &&
		m.Modifier == other.Modifier &&
		m.Sampling == other.Sampling
}

func ConsumeBucketToken(m *Message, b []byte) error {
	m.Bucket = string(b)
	return nil
}

func ConsumeValueToken(m *Message, b []byte) (err error) {
	m.Value, err = strconv.Atoi(string(b))
	return err
}

func ConsumeModifierToken(m *Message, b []byte) error {
	m.Modifier = string(b)
	return nil
}

func ConsumeSamplingToken(m *Message, b []byte) error {
	f, err := strconv.ParseFloat(string(b), 32)
	if err != nil {
		m.Sampling = 1
		return err
	}
	m.Sampling = float32(f)
	return err
}

func ScanStat(data []byte, atEOF bool) (advance int, token []byte, err error) {
	l := len(data)
	if atEOF && l == 0 {
		return 0, nil, nil
	}
	for i := 0; i < l; i++ {
		switch data[i] {
		case ':', '|':
			return i + 1, data[0:i], nil
		}
	}
	return l + 1, data, nil
}

// The consumer type defines a function signature that can be used when
// consuming message tokens.
type Consumer func(*Message, []byte) error

var consumers = map[int]Consumer{
	0: ConsumeBucketToken,
	1: ConsumeValueToken,
	2: ConsumeModifierToken,
	3: ConsumeSamplingToken,
}

func Parse(r Raw) (*Message, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(ScanStat)
	i := 0
	m := new(Message)
	for scanner.Scan() {
		err := consumers[i](m, scanner.Bytes())
		if err != nil {
			return nil, err
		}
		i++
	}
	return m, nil
}
