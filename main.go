package main

import (
	"fmt"
	"os"

	"github.com/alexkappa/statsd/config"
	"github.com/alexkappa/statsd/daemon"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("statsd: error no configuration specified")
		os.Exit(1)
	}
	c, err := config.New(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	d, err := daemon.New(c)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = d.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}
