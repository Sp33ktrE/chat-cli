package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/Sp33ktrE/chat-cli/pkg/protocol"
)

type Client struct {
	name string
	conn net.Conn
}

func New(name string, addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		name: name,
		conn: conn,
	}, nil
}

func (client *Client) handleReadMsg(reader *bufio.Reader) {
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection from the server closed")
			break
		}
		fmt.Println(msg)
	}
}

func (client *Client) Chat() {
	defer client.conn.Close()
	client.nickCommand()
	reader := bufio.NewReader(client.conn)
	msg, _ := reader.ReadString('\n')
	msgParsed, err := protocol.ParsePMessage(msg)
	if err != nil {
		log.Fatal(err)
	}
	if msgParsed.Command == "001" {
		fmt.Println(msgParsed.Trailing)
		go client.handleReadMsg(reader)
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Println(">> Enter your message: ")
			input, _ := reader.ReadString('\n')
			if strings.ToUpper(input) == "QUIT\n" {
				break
			}
			client.conn.Write([]byte(input))
		}
	} else if msg == "FULL\n" {
		fmt.Println("SERVER IS FULL")
	} else {
		fmt.Println("An error connecting the server has occured: ", err)
	}
}
