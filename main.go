package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"whisperd.io/whisperd/whisperd/server"
)

type StringSet map[string]struct{}

func (ss *StringSet) String() string {
	keys := make([]string, len(*ss))
	for k := range *ss {
		keys = append(keys, k)
	}
	return fmt.Sprintf("%s", keys)
}

func (ss *StringSet) Set(value string) error {
	map[string]struct{}(*ss)[value] = struct{}{}
	return nil
}

func main() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)

	waitChan := make(chan struct{})

	go func() {
		<-sigChan
		log.Printf("Caught signal\n")
		close(waitChan)
	}()

	opts := server.Opts{}
	roleSet := StringSet{}

	flag.Var(&roleSet, "role", "Server role(s)")
	flag.StringVar(&opts.Addr, "addr", "127.0.0.1:3000", "Server addresss")

	flag.StringVar(&opts.DB.Driver, "db-driver", "sqlite", "Database driver")
	flag.StringVar(&opts.DB.SQLiteFileName, "sqlite-filename", "sqlite.db", "Database file name")

	flag.Parse()

	opts.Roles = roleSet

	log.Printf("Server opts: %#v", opts)

	s, err := server.New(opts)
	if err != nil {
		log.Fatal(err)
	}
	go s.ListenAndServe()

	log.Printf("Server listening on address %s", s.Addr)
	<-waitChan

	log.Printf("Bye")
}
