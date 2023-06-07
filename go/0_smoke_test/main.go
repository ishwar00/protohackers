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

	fmt.Printf("received new connection! from %v\n", conn.RemoteAddr())

	buff := make([]byte, 1024)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			if err != io.EOF {
				log.Println(err.Error())
			}
			return
		}
		_, err = conn.Write(buff[:n])
		if err != nil {
			log.Println("error while writing: ", err.Error())
			return
		}
	}
}

