package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github/ishwar00/protohackers/go/protohacks/server"
)

func main() {
	server := server.New("tcp", "5001")
	fmt.Println("tcp server listening on port 5001...")
	server.Start(handleConn)
}

func handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	fmt.Println("received new connection!!")

	buff := make([]byte, 1024)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			return
		}
		conn.Write(buff[:n])
	}
}

