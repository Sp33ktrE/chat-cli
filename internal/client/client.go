package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct {
	name string
	conn net.Conn
}

func New() (*Client, error) {
	conn, err := net.Dial("tcp", ":5050")
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(">> Enter your name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	return &Client{
		name: name,
		conn: conn,
	}, nil
}

func (client *Client) Chat() {
	defer client.conn.Close()
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println(">> Enter your message: ")
		input, _ := reader.ReadString('\n')
		if strings.ToUpper(input) == "QUIT\n" {
			break
		}
		client.conn.Write([]byte(input))
	}
}
