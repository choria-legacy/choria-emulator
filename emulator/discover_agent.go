package emulator

import (
	"fmt"
	"strings"
)

type DiscoveryAgent struct {
	license     string
	author      string
	timeout     int
	name        string
	version     string
	url         string
	description string
}

func (self *DiscoveryAgent) Init() error {
	self.license = "ASL-2.0"
	self.author = "R.I.Pienaar <rip@devco.net>"
	self.timeout = 2
	self.name = "discovery"
	self.version = "1.0.0"
	self.url = "http://choria.io"
	self.description = "Discovery Agent"

	return nil
}

func (self *DiscoveryAgent) Name() string {
	return self.name
}

func (self *DiscoveryAgent) HandleAgentMsg(msg string) (*[]byte, error) {
	if strings.Contains(msg, "ping") {
		r := []byte("pong")
		return &r, nil
	}

	return nil, fmt.Errorf("unknown request: %s", msg)
}