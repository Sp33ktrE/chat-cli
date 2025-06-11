package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Sp33ktrE/chat-cli/internal/client"
)

const HOST = ""
const PORT = "5050"

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(">> Enter your name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	user, err := client.New(name, fmt.Sprintf("%s:%s", HOST, PORT))
	if err != nil {
		log.Fatal("Error connecting to the server: ", err)
	}
	user.Start()
}
