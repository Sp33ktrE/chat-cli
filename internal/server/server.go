package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	host string
	port string
}

func New(host string, port string) *Server {
	return &Server{
		host: host,
		port: port,
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
		for {
			msg, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Connection closed or err: ", err)
				break
			}
			fmt.Printf("[%s]: %s\n", name, msg)
		}
	}
	<-ch
}

func (server *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("***Chat Server Started***")
	const cap = 2
	sem := make(chan bool, cap)
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
