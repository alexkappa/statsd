package daemon

import (
	"net"
	"os"

	"github.com/alexkappa/statsd/config"
	"gopkg.in/yieldr/go-log.v0/log"
)

type Daemon struct {
	address *net.UDPAddr
	msgbuf  chan *Message
	err     chan error
	log     *log.Logger
}

func New(c *config.Config) (*Daemon, error) {
	address, err := net.ResolveUDPAddr("udp", c.Addr)
	if err != nil {
		return nil, err
	}
	logger := log.NewSimple(
		log.WriterSink(
			os.Stdout,
			log.BasicFormat,
			log.BasicFields))
	d := &Daemon{
		address,
		make(chan *Message, 1000),
		make(chan error),
		logger,
	}
	return d, nil
}

func (d *Daemon) Run() error {
	conn, err := net.ListenUDP("udp", d.address)
	if err != nil {
		return err
	}
	defer conn.Close()
	go d.Err()
	go d.Flush()
	for {
		message := make([]byte, 512)
		n, _, err := conn.ReadFrom(message)
		if err != nil {
			d.log.Error(err)
		}
		go d.handle(message[0:n])
	}
}

func (d *Daemon) handle(raw Raw) {
	m, err := Parse(raw)
	if err != nil {
		d.err <- err
	} else {
		d.msgbuf <- m
	}
}

func (d *Daemon) Flush() {
	for {
		m := <-d.msgbuf
		d.log.Infof("message: %s:%d|%s|@%f", m.Bucket, m.Value, m.Modifier, m.Sampling)
	}
}

func (d *Daemon) Err() {
	for {
		d.log.Error(<-d.err)
	}
}
