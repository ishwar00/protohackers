package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"

	"github/ishwar00/protohackers/go/protohacks/server"
)

type Insert struct {
	TimeStamp int32
	Price     int32
}

type Query struct {
	MinTime int32
	MaxTime int32
}

type Message struct {
	Message []byte
}

func (m Message) Type() rune {
	return rune(m.Message[0])
}

func (m Message) ProcessFormat() (any, error) {
	if len(m.Message) != 9 {
		return nil, errors.New("message is not of 9 bytes")
	}

	switch m.Type() {
	case 'I':
		timeStamp, err := DecodeInt(m.Message[1:4])
		if err != nil {
			return nil, err
		}
		price, err := DecodeInt(m.Message[5:])
		if err != nil {
			return nil, err
		}

		insert := &Insert{
			TimeStamp: timeStamp,
			Price:     price,
		}

		return insert, nil
	case 'Q':
		minTime, err := DecodeInt(m.Message[1:4])
		if err != nil {
			return nil, err
		}
		maxTime, err := DecodeInt(m.Message[5:])
		if err != nil {
			return nil, err
		}

		query := &Query{
			MinTime: minTime,
			MaxTime: maxTime,
		}

		return query, nil
	}

	return nil, errors.New("unknown format recieved")
}

func DecodeInt(buf []byte) (int32, error) {
	value, n := binary.Varint(buf)

	// local util
	Err := func(msg string) (int32, error) {
		errMsg := fmt.Sprintf("decoding failed: %s, buf = %+v\n", msg, buf)
		return 0, errors.New(errMsg)
	}
	switch {
	case n != len(buf):
		return Err("Varint() did not read all bytes of 'buf'")
	case n == 0:
		return Err("given 'buf' is too small")
	case n == 0:
		return Err("value larger than 64 bits(overflow)")
	}

	return int32(value), nil // ok?
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 9)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("server error: Read(): %s\n", err.Error())
			return
		}
		if n != 9 {
			fmt.Println("client error: could not read 9 bytes")
			return
		}

		msg := &Message{
			Message: buf,
		}

		formattedMsg, err := msg.ProcessFormat()
		if err != nil {
			fmt.Printf("ProcessFormat(): %s\n", err)
			return
		}

		switch formattedMsg.(type) {
		case Query:
		case Insert:
		default:
			panic("some thing went wrong!")
		}
	}
}

func main() {
	server := server.New("tcp", "5001")
	server.Start(handleConn)
}
