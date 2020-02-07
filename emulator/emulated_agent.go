package emulator

import (
	"context"
	"fmt"
	"os"

	"github.com/choria-io/go-choria/choria"
	"github.com/choria-io/go-choria/providers/agent/mcorpc"
	"github.com/choria-io/go-choria/server/agents"
)

type GenerateRequest struct {
	Size int `json:"size"`
}

type GenerateReply struct {
	Message string `json:"message"`
}

func NewEmulatedAgent(fw *choria.Framework, count int) agents.Agent {
	metadata := &agents.Metadata{
		Name:        fmt.Sprintf("emulated%d", count),
		Description: "Emulated Agent",
		Author:      "R.I.Pienaar <rip@devco.net>",
		Version:     "1.0.0",
		License:     "Apache-2.0",
		Timeout:     5,
		URL:         "http://choria.io",
	}

	agent := mcorpc.New(metadata.Name, metadata, fw, fw.Logger(metadata.Name))
	agent.MustRegisterAction("generate", generateAction)
	agent.MustRegisterAction("exit_emulator", exitAction)

	return agent
}

func exitAction(ctx context.Context, req *mcorpc.Request, reply *mcorpc.Reply, agent *mcorpc.Agent, conn choria.ConnectorInfo) {
	os.Exit(0)
}

func generateAction(ctx context.Context, req *mcorpc.Request, reply *mcorpc.Reply, agent *mcorpc.Agent, conn choria.ConnectorInfo) {
	genreq := &GenerateRequest{}
	if !mcorpc.ParseRequestData(genreq, req, reply) {
		return
	}

	reply.Data = &GenerateReply{
		Message: randomString(genreq.Size),
	}
}

func randomString(strlen int) string {
	chars := "01234567890"
	result := ""

	for i := 0; i < strlen; i++ {
		start := i % 10
		result += chars[start : start+1]
	}

	return result
}
