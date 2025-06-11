package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

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

func (client *Client) readMessage(reader *bufio.Reader, stopCh chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-stopCh:
			fmt.Println("Connection from the server closed")
			return
		default:
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Connection from the server closed")
				close(stopCh)
				return
			}
			fmt.Println(message)
		}
	}
}

func (client *Client) sendMessage(stopCh chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-stopCh:
			fmt.Println("Connection from the server closed")
			return
		default:
			//os.Stdin read blocks even if we receive an error from another goroutine, fixable but needs some adjsuments
			reader := bufio.NewReader(os.Stdin)
			fmt.Println(">> Enter your message: ")
			input, _ := reader.ReadString('\n')
			if strings.ToUpper(input) == "QUIT\n" {
				close(stopCh)
				return
			}
			client.conn.Write([]byte(input))
		}
	}
}

func (client *Client) Chat(reader *bufio.Reader) {
	const grs = 2
	var wg sync.WaitGroup
	wg.Add(grs)
	stopCh := make(chan struct{})
	go client.readMessage(reader, stopCh, &wg)
	go client.sendMessage(stopCh, &wg)
	wg.Wait()
}

func (client *Client) Start() {
	defer client.conn.Close()
	// Init server handshake
	client.nickCommand()
	reader := bufio.NewReader(client.conn)
	serverReply, _ := reader.ReadString('\n')
	serverReplyParsed, err := protocol.ParsePMessage(serverReply)
	if err != nil {
		log.Fatal("Error parsing server response: ", err)
	}
	switch serverReplyParsed.Command {
	case "001":
		fmt.Print(serverReplyParsed.Trailing)
		client.Chat(reader)
	case "401":
		fmt.Println(serverReplyParsed.Trailing)
	default:
		fmt.Println("Error connecting to the server")
	}
}
