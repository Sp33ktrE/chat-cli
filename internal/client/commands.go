package client

import (
	"github.com/Sp33ktrE/chat-cli/pkg/protocol"
)

func (client *Client) nickCommand() {
	command := protocol.New("", protocol.CmdNick, []string{client.name}, "")
	client.conn.Write([]byte(command.FormatPMessage()))
}
