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
	const cap = 2
	sem := make(chan bool, cap)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		select {
		case sem <- true:
			{
				go func() {
					defer conn.Close()
					{
						conn.Write([]byte("ACCEPT\n"))
						fmt.Println("Client connected")
						reader := bufio.NewReader(conn)
						for {
							msg, err := reader.ReadString('\n')
							if err != nil {
								fmt.Println("Connection closed or err: ", err)
								break
							}
							fmt.Println("Message recieved: ", msg)
						}
					}
					<-sem
				}()
			}
		default:
			conn.Write([]byte("FULL\n"))
			conn.Close()
		}

	}
}
