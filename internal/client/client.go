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
	reader := bufio.NewReader(client.conn)
	msg, err := reader.ReadString('\n')
	if msg == "ACCEPT\n" {
		client.conn.Write([]byte(client.name + "\n"))
		msg, _ := reader.ReadString('\n')
		if msg == "OK\n" {
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
		}
	} else if msg == "FULL\n" {
		fmt.Println("SERVER IS FULL")
	} else {
		fmt.Println("An error connecting the server has occured: ", err)
	}
}
