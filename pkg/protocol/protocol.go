package protocol

import (
	"fmt"
	"strings"
)

const (
	CmdNick    = "NICK"
	CmdJoin    = "JOIN"
	CmdPrivMsg = "PRIVMSG"
	CmdQuit    = "QUIT"
	CmdPing    = "PING"
	CmdPong    = "PONG"
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

func New(sender string, command string, params []string, trailing string) *ProtocolMessage {
	return &ProtocolMessage{
		Sender:   sender,
		Command:  command,
		Params:   params,
		Trailing: trailing,
	}
}

func (pmessage *ProtocolMessage) FormatPMessage() string {
	params := strings.Join(pmessage.Params, " ")
	message := fmt.Sprintf(" %s %s", pmessage.Command, params)
	if pmessage.Sender != "" {
		message = fmt.Sprintf(":%s%s", pmessage.Sender, message)
	}
	if pmessage.Trailing != "" {
		message = fmt.Sprintf("%s:%s", message, pmessage.Trailing)
	}
	return message + "\n"
}

func ParsePMessage(line string) (*ProtocolMessage, error) {
	var fullCmd []string
	fullCmdSender := ""

	line, hasSender := strings.CutPrefix(line, ":")

	line, fullCmdMsg, _ := strings.Cut(line, ":")

	if hasSender {
		fullCmdSender, line, _ = strings.Cut(line, " ")
	}
	fullCmd = strings.Split(line, " ")
	return New(
		fullCmdSender, fullCmd[0], fullCmd[1:], fullCmdMsg,
	), nil
}
