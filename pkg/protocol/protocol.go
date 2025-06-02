package protocol

import (
	"fmt"
	"strings"
)

const (
	CmdNick = "NICK"
)

const (
	RplWelcome    = "001"
	RplConnected  = "100"
	RplAccepted   = "200"
	ErrServerFull = "401"
)

// message follows this structure [:sender] COMMAND [params] [:trailing]
type ProtocolMessage struct {
	Sender   string
	Command  string
	Params   []string
	Trailing string
}

func ParsePMessage(line string) (ProtocolMessage, error) {
	var fullCmd []string
	fullCmdSender := ""

	line, hasSender := strings.CutPrefix(line, ":")

	line, fullCmdMsg, _ := strings.Cut(line, ":")

	if hasSender {
		fullCmdSender, line, _ = strings.Cut(line, " ")
	}
	fullCmd = strings.Split(line, " ")
	return ProtocolMessage{
		Sender:   fullCmdSender,
		Command:  fullCmd[0],
		Params:   fullCmd[1:],
		Trailing: fullCmdMsg,
	}, nil
}

func (pmessage *ProtocolMessage) FormatPMessage() string {
	params := strings.Join(pmessage.Params, " ")
	return fmt.Sprintf(":%s %s %s:%s\n", pmessage.Sender, pmessage.Command, params, pmessage.Trailing)
}
