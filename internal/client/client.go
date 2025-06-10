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
			fmt.Println("Received error, stop signal")
			return
		default:
			msg, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Connection from the server closed")
				close(stopCh)
				return
			}
			fmt.Println(msg)
		}
	}
}

func (client *Client) sendMessage(stopCh chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-stopCh:
			fmt.Println("Received error, stop signal")
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
		const grs = 2
		var wg sync.WaitGroup
		wg.Add(grs)
		stopCh := make(chan struct{})
		fmt.Println(msgParsed.Trailing)
		go client.readMessage(reader, stopCh, &wg)
		go client.sendMessage(stopCh, &wg)
		wg.Wait()
	} else if msgParsed.Command == "401" {
		fmt.Print(msgParsed.Trailing)
	} else {
		fmt.Println("An error connecting the server has occured: ", err)
	}
}
