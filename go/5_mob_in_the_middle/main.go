package main

import (
	"bufio"
	"fmt"
	"github/ishwar00/protohackers/go/protohacks/server"
	"log"
	"net"
	"strings"
)

const tonyAddr = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"

// reference: https://github.com/adityathebe/protohackers/blob/3c5db03ed3014e7f5c082e0a8bf5bbe8dd2ed723/5.mob_in_the_middle/main.go#L67
func isBogusCoinAddress(s string) bool {
	return len(s) >= 26 && len(s) <= 35 && s[0] == '7'
}

func WriteBoguscoinAddress(msg string) string {
	splits := strings.Split(msg, " ")
	for i, s := range splits {
		if isBogusCoinAddress(s) {
			splits[i] = tonyAddr
		}
	}

	return strings.Join(splits, " ")
}

func proxy(from *net.Conn, to *net.Conn) {
	fromScnr := bufio.NewScanner(*from)

	for fromScnr.Scan() {
		serverMsg := fromScnr.Text()
		bogusMsg := WriteBoguscoinAddress(serverMsg) + "\n"
		_, err := (*to).Write([]byte(bogusMsg))
		if err != nil {
			log.Printf("to.Write(): %s", err)
			return
		}
	}

	if err := fromScnr.Err(); err != nil {
		log.Printf("fromScnr.Err(): %s\n", fromScnr.Err())
	}
}

func handleConn(clientConn net.Conn) {
	addr := fmt.Sprintf("chat.protohackers.com:16963")
	serverConn, err := net.Dial("tcp", addr)
	defer serverConn.Close()
	if err != nil {
		log.Fatalf("Server.Start(): %v\n", err)
	}

	go proxy(&clientConn, &serverConn)
	proxy(&serverConn, &clientConn)
}

func main() {
	server := server.New("tcp", "5001")
	server.Start(handleConn)
}

