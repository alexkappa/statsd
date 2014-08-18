package daemon

import (
	"fmt"
	"net"

	"github.com/alexkappa/statsd/config"
	"github.com/alexkappa/statsd/daemon/message"
)

type Daemon struct {
	address *net.UDPAddr
	mbuffer chan *message.Message
	err     chan error
}

func New(c *config.Config) (*Daemon, error) {
	address, err := net.ResolveUDPAddr("udp", c.Addr)
	if err != nil {
		return nil, err
	}
	d := &Daemon{
		address,
		make(chan *message.Message, 1000),
		make(chan error),
	}
	return d, nil
}

func (d *Daemon) Run() error {
	conn, err := net.ListenUDP("udp", d.address)
	if err != nil {
		return err
	}
	defer conn.Close()
	go d.error()
	go d.flush()
	for {
		message := make(message.Raw, 512)
		n, _, error := conn.ReadFrom(message)
		if error != nil {
			continue
		}
		go d.handle(message[0:n])
	}
}

func (d *Daemon) handle(raw message.Raw) {
	m, err := message.Parse(raw)
	if err != nil {
		d.err <- err
	} else {
		d.mbuffer <- m
	}
}

func (d *Daemon) flush() {
	for {
		m := <-d.mbuffer
		fmt.Printf("message: %s:%d|%s|@%f\n", m.Bucket, m.Value, m.Modifier, m.Sampling)
	}
}

func (d *Daemon) error() {
	for {
		fmt.Printf("error: %s\n", <-d.err)
	}
}
