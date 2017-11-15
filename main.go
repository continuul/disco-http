package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"continuul.io/raft-http/http"
	"continuul.io/raft-http/store"
)

// Command line defaults
const (
	DefaultHttpAddress = ":11000"
	DefaultBindAddress = ":12000"
)

var nodeName string
var bindAddress string
var httpAddress string
var joinAddress string

func init() {
	flag.StringVar(&nodeName, "node", "", "Node name")
	flag.StringVar(&bindAddress, "bind", DefaultBindAddress, "Set Raft bind address")
	flag.StringVar(&httpAddress, "client", DefaultHttpAddress, "Set the HTTP bind address")
	flag.StringVar(&joinAddress, "join", "", "Set join address, if any")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <raft-data-path> \n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "No Raft storage directory specified\n")
		os.Exit(1)
	}

	// Ensure Raft storage exists.
	raftDir := flag.Arg(0)
	if raftDir == "" {
		fmt.Fprintf(os.Stderr, "No Raft storage directory specified\n")
		os.Exit(1)
	}
	os.MkdirAll(raftDir, 0700)

	s := store.New()
	s.RaftDir = raftDir
	s.RaftBind = bindAddress
	if err := s.Open(joinAddress == "", nodeName); err != nil {
		log.Fatalf("failed to open store: %s", err.Error())
	}

	h := httpd.New(httpAddress, s)
	if err := h.Start(); err != nil {
		log.Fatalf("failed to start HTTP service: %s", err.Error())
	}

	// If join was specified, make the join request.
	if joinAddress != "" {
		if err := join(joinAddress, bindAddress, nodeName); err != nil {
			log.Fatalf("failed to join node at %s: %s", joinAddress, err.Error())
		}
	}

	log.Println("disco-http started successfully")

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
	log.Println("disco-http exiting")
}

func join(joinAddr, raftAddr, nodeID string) error {
	b, err := json.Marshal(map[string]string{"addr": raftAddr, "id": nodeID})
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/join", joinAddr), "application-type/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
