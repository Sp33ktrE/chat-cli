package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
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
		conn.Write([]byte("ACCEPT\n"))
		reader := bufio.NewReader(conn)
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)
		fmt.Printf("[%s] Connected\n", name)
		conn.Write([]byte("OK\n"))
		server.clients[name] = conn
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
			conn.Write([]byte("FULL\n"))
			conn.Close()
		}

	}
}
