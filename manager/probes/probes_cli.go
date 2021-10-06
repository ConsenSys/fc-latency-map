package probes

import (
	"strings"

	prompt "github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"

	"github.com/ConsenSys/fc-latency-map/manager/cli"
)

const (
	probesUpdate = "probes-update"
	probesList   = "probes-list"
)

type ProbeCommander struct {
	Handler *ProbeHandler
}

func NewProbeCommander() cli.Commander {
	return &ProbeCommander{
		Handler: NewProbeHandler(),
	}
}

// completes the input
func (cmd *ProbeCommander) Complete() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: probesUpdate, Description: "Update probes list by finding online and active probes"},
		{Text: probesList, Description: "Get probes list"},
	}
}

// executes the command
func (cmd *ProbeCommander) Execute(in string) {
	blocks := strings.Split(strings.TrimSpace(in), " ")

	switch blocks[0] {
	case probesUpdate:
		cmd.Handler.Update()
	case probesList:
		cmd.Handler.GetAllProbes()
	default:
		log.Printf("unknown command: %s\n", blocks[0])
	}
}