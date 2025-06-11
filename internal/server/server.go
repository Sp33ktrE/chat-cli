package server

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/Sp33ktrE/chat-cli/pkg/protocol"
)

type Server struct {
	host      string
	port      string
	broadcast chan string
	clients   map[string]net.Conn
}

func New(host string, port string) *Server {
	return &Server{
		host:      host,
		port:      port,
		broadcast: make(chan string),
		clients:   make(map[string]net.Conn),
	}
}

func (server *Server) handleClient(ch chan bool, conn net.Conn) {
	defer conn.Close()
	{
		reader := bufio.NewReader(conn)
		clientNickCommand, _ := reader.ReadString('\n')
		fmt.Println(clientNickCommand)
		pMessage, err := protocol.ParsePMessage(clientNickCommand)
		if err != nil {
			fmt.Println(err)
			return
		}
		clientName := pMessage.Params[0]
		fmt.Printf("[%s] Connected\n", clientName)
		acceptMessage := protocol.New("server", protocol.RplWelcome, []string{}, "Connected to CLI chat, welcome!!")
		conn.Write([]byte(acceptMessage.FormatPMessage()))
		server.clients[clientName] = conn
		for {
			msg, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Connection closed or err: ", err)
				break
			}
			//fmt.Printf("[%s]: %s\n", name, msg)
			server.broadcast <- msg
		}
	}
	<-ch
}

func (server *Server) handleBroadcast() {
	for {
		msg := <-server.broadcast
		fmt.Println(msg)
		for _, conn := range server.clients {
			conn.Write([]byte(msg + "\n"))
		}
		// TODO: i should add client map so i can broadcast to all their conns here!!
	}
}

func (server *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("***Chat Server Started***")
	const cap = 2
	sem := make(chan bool, cap)
	go server.handleBroadcast()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		select {
		case sem <- true:
			go server.handleClient(sem, conn)
		default:
			fullErrMessage := protocol.New("server", protocol.ErrServerFull, []string{}, "Server is full, try again later!!")
			conn.Write([]byte(fullErrMessage.FormatPMessage()))
			conn.Close()
		}
	}
}
