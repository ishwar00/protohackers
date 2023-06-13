package server

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	network string
	port    string
}

type HandleConn func(conn net.Conn)

func (s *Server) Start(handleConn HandleConn) {
	addr := fmt.Sprintf("127.0.0.1:%v", s.port)
	listener, err := net.Listen(s.network, addr)
	if err != nil {
		log.Fatalf("Server.Start(): %v", err)
	}

	fmt.Printf("%s server listening on %s...\n", s.network, addr)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("listener.Accept(): %v", err)
			continue
		}

		go handleConn(conn)
	}
}

func New(network, port string) *Server {
	return &Server{
		network: network,
		port:    port,
	}
}
