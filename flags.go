package main

import (
	"flag"
	"time"
)

var (
	port     = flag.Int("port", 8080, "The HTTP port to listen on")
	interval = flag.Duration("interval", 1*time.Hour, "Interval between checks")
	timeout  = flag.Duration("timeout", 15*time.Second, "HTTP API call timeout")

	owners multiVar
)

type multiVar []string

func (mv *multiVar) Set(value string) error {
	*mv = append(*mv, value)
	return nil
}

func (mv *multiVar) String() string {
	all := ""
	for _, item := range *mv {
		if all != "" {
			all += ", "
		}

		all += item
	}

	return "[" + all + "]"
}

func init() {
	flag.Var(&owners, "owner", "Owners (namespaces) to list repositories for (multiple values are allowed)")

	flag.Parse()
}
