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

type HandleTCPConn func(conn net.Conn)
type HandleUDPConn func(conn net.UDPConn)

func (s *Server) StartUDP(handleConn HandleUDPConn) {
	addr := fmt.Sprintf("127.0.0.1:%v", s.port)
	network := "udp"
	laddr, err := net.ResolveUDPAddr(network, addr)
	if err != nil {
		log.Panicf("ResolveUDPAddr(): %s", err)
	}
	udpConn, err := net.ListenUDP(s.network, laddr)
	if err != nil {
		log.Fatalf("net.ListenUDP(): %v", err)
	}

	fmt.Printf("%s server listening on %s...\n", s.network, addr)

	defer udpConn.Close()

	handleConn(*udpConn)
}

func (s *Server) Start(handleConn HandleTCPConn) {
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
