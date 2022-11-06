package utils

import (
	"errors"
	"flag"
)

const (
	defAddress = "127.0.0.1"
	defPort    = "8080"
)

type Args struct {
	Address string
	Port    string
}

func (args *Args) validate() error {
	// it won't occur as we have defined default values. Nevertheless, it's a good safety measure
	// if a change in the future breaks this behavior
	if args.Address == "" || args.Port == "" {
		return errors.New("address and port should not be empty")
	}
	return nil
}

func ParseArgs() (*Args, error) {
	var args Args
	flag.StringVar(&args.Address, "addr", defAddress, "HTTP Server Listen Address\n")
	flag.StringVar(&args.Port, "port", defPort, "HTTP Server Listen Port\n")
	flag.Parse()

	if err := args.validate(); err != nil {
		flag.PrintDefaults()
		return nil, err
	}
	return &args, nil
}
