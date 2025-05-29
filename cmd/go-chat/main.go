package main

import (
	"github.com/Sp33ktrE/chat-cli/internal/server"
)

const HOST = ""
const PORT = "5050"

func main() {
	server := server.New(HOST, PORT)
	server.Run()
}
