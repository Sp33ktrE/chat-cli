package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func Chat() {
	conn, err := net.Dial("tcp", ":5050")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter your message:")
		input, _ := reader.ReadString('\n')
		if strings.ToUpper(input) == "QUIT\n" {
			break
		}
		conn.Write([]byte(input))
	}
}
