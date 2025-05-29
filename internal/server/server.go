package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
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

func (server *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("***Chat Server Started***")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		go func() {
			fmt.Println("-Client Connected-")
			defer conn.Close()
			reader := bufio.NewReader(conn)
			for {
				msg, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Connection closed or err: ", err)
					break
				}
				fmt.Println("Message recieved: ", msg)
			}
		}()
	}
}
