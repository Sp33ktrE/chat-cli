package main

import (
	"log"

	"github.com/Sp33ktrE/chat-cli/internal/client"
)

func main() {
	client, err := client.New()
	if err != nil {
		log.Fatal(err)
	}
	client.Chat()
}
