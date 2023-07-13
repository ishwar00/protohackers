package main

import (
	"fmt"
	"github/ishwar00/protohackers/go/protohacks/server"
	"log"
	"net"
	"strings"
	"sync"
)

type SafeConn struct {
	net.UDPConn
	mu sync.Mutex
}

func (s *SafeConn) WriteTo(key, value string, addr net.Addr) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	kv := fmt.Sprintf("%s=%s", key, value);
	return s.UDPConn.WriteTo([]byte(kv), addr)
}

type KVDatabase struct {
	mp map[string]string
	mu sync.Mutex
}

func (db *KVDatabase) read(key string) string {
	db.mu.Lock()
	defer db.mu.Unlock()

	if (key == "version") {
		return "ishwar's unstable key-value database 1.0"
	}
	return db.mp[key]
}

func (db *KVDatabase) write(key, value string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if (key != "version") {
		db.mp[key] = value
	}
}

func handleUDPConn(conn net.UDPConn) {
	safeConn := &SafeConn{UDPConn: conn}
	db := &KVDatabase{
		mp: map[string]string{},
	}

	buf := make([]byte, 1024)
	for {
		n, addr, err := safeConn.ReadFrom(buf)
		if err != nil {
			log.Printf("conn.ReadFrom(): %s", err)
			continue
		}

		go handleRequest(string(buf[:n]), addr, safeConn, db)
	}
}

func handleRequest(req string, addr net.Addr, safeConn *SafeConn, db *KVDatabase) {
	key, value, writeRequest := strings.Cut(req, "=")

	if writeRequest {
		db.write(key, value)
	} else {
		value := db.read(key)
		if _, err := safeConn.WriteTo(key, value, addr); err != nil {
			log.Printf("safeConn.WriteTo(): %s\n", err);
		}
	}
}

func main() {
	server := server.New("udp", "5001")
	server.StartUDP(handleUDPConn)
}
