package emulator

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type EmulatedAgent struct {
	license     string
	author      string
	timeout     int
	name        string
	version     string
	url         string
	description string
	request     *RPCRequestBody
}

type RPCRequestBody struct {
	Agent  string          `json:"agent"`
	Action string          `json:"action"`
	Data   *RPCRequestData `json:"data"`
}

type RPCRequestData struct {
	Size int `json:"size"`
}

type RPCReply struct {
	Statuscode int           `json:"statuscode"`
	Statusmsg  string        `json:"statusmsg"`
	Data       *RPCReplyData `json:"data"`
}

type RPCReplyData struct {
	Message string `json:"message"`
}

func (self *EmulatedAgent) Name() string {
	return self.name
}

func (self *EmulatedAgent) Init() error {
	self.license = "ASL-2.0"
	self.author = "R.I.Pienaar <rip@devco.net>"
	self.timeout = 2
	self.version = "1.0.0"
	self.url = "http://choria.io"
	self.description = "Emulated Agent"

	return nil
}

func (self *EmulatedAgent) newReply() *RPCReply {
	reply := &RPCReply{
		Statuscode: 0,
		Statusmsg:  "OK",
		Data:       &RPCReplyData{},
	}

	return reply
}

func (self *EmulatedAgent) requestFromMsg(msg string) (*RPCRequestBody, error) {
	r := &RPCRequestBody{}

	err := json.Unmarshal([]byte(msg), r)
	if err != nil {
		return nil, fmt.Errorf("Could not parse incoming request: %s", err.Error())
	}

	return r, nil
}

func (self *EmulatedAgent) HandleAgentMsg(msg string) (*[]byte, error) {
	reply := self.newReply()
	request, err := self.requestFromMsg(msg)
	if err != nil {
		return nil, err
	}

	switch request.Action {
	case "generate":
		self.generateAction(request, reply)
	default:
		reply.Statuscode = 2
		reply.Statusmsg = fmt.Sprintf("Unknown action %s", request.Action)
	}

	j, err := json.Marshal(&reply)
	if err != nil {
		log.Errorf("Could not marshall JSON reply: %s", err.Error())
	}

	return &j, nil
}

func (self *EmulatedAgent) generateAction(request *RPCRequestBody, reply *RPCReply) {
	reply.Data.Message = self.randomString(request.Data.Size)
}

func (self *EmulatedAgent) randomString(strlen int) string {
	chars := "01234567890"
	result := ""

	for i := 0; i < strlen; i++ {
		start := i % 10
		result += chars[start : start+1]
	}

	return result
}
