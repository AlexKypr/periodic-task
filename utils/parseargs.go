package utils

import (
	"errors"
	"flag"
)

const (
	defPort = "8080"
)

type Args struct {
	Port string
}

func (args *Args) validate() error {
	// it won't occur as we have defined default values. Nevertheless, it's a good safety measure
	// if a change in the future breaks this behavior
	if args.Port == "" {
		return errors.New("port should not be empty")
	}
	return nil
}

func ParseArgs() (*Args, error) {
	var args Args
	flag.StringVar(&args.Port, "port", defPort, "HTTP Server Listen Port\n")
	flag.Parse()

	if err := args.validate(); err != nil {
		flag.PrintDefaults()
		return nil, err
	}
	return &args, nil
}
