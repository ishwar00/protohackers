package main

import (
	"bufio"
	"bytes"
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

type Buffer struct {
	Buf []Insert
}

func (b *Buffer) insert(i Insert) {
	b.Buf = append(b.Buf, i)
}

func (b Buffer) query(q Query) int32 {
	price_sum := 0
	price_count := 0
	for _, msg := range b.Buf {
		if msg.TimeStamp <= q.MaxTime && q.MinTime <= msg.TimeStamp {
			fmt.Printf("debubInfo: msg: %+v\n", msg)
			price_sum += int(msg.Price)
			price_count++
		}
	}

	if price_count != 0 {
		return int32(price_sum / price_count)
	}

	return 0
}

func processFormat(buf []byte) (any, error) {
	if len(buf) != 9 {
		return nil, errors.New("message is not of 9 bytes")
	}

	typeOfMsg, first4, second4, err := decodeMessage(buf)
	if err != nil {
		return nil, err
	}

	switch typeOfMsg {
	case 'I':
		insert := Insert{
			TimeStamp: first4,
			Price:     second4,
		}

		return insert, nil
	case 'Q':
		query := Query{
			MinTime: first4,
			MaxTime: second4,
		}

		return query, nil
	}

	return nil, errors.New("unknown format recieved")
}

func decodeMessage(buf []byte) (rune, int32, int32, error) {
	typeOfMsg := rune(buf[0])
	first4, err := decodeInt(buf[1:5])
	if err != nil {
		return ' ', 0, 0, err
	}

	second4, err := decodeInt(buf[5:])
	if err != nil {
		return ' ', 0, 0, err
	}

	return typeOfMsg, first4, second4, nil
}

func decodeInt(buf []byte) (int32, error) {
	var intBuf int32
	err := binary.Read(bytes.NewReader(buf), binary.BigEndian, &intBuf)
	if err != nil {
		fmt.Printf("binary.Read(): failed with: %+v\n", err)
		return 0, err
	}

	return intBuf, err
}

func encodeInt(value int32) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, &value)
	if err != nil {
		fmt.Printf("binary.Write(): failed with: %+v\n", err)
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("handling connection from %v\n", conn.RemoteAddr().String())

	buffer := Buffer{
		Buf: []Insert{},
	}

	split := func(data []byte, atEOF bool) (int, []byte, error) {
		if len(data) >= 9 {
			return 9, data[:9], nil
		}

		if !atEOF {
			return 0, nil, nil
		}

		return 0, data, bufio.ErrFinalToken
	}

	scanner := bufio.NewScanner(conn)
	scanner.Split(split)
	for scanner.Scan() {
		formattedMsg, err := processFormat(scanner.Bytes())
		if err != nil {
			fmt.Printf("processFormat(): %s\n", err)
			conn.Write([]byte{0, 0, 0, 0})
			continue
		}

		switch formattedMsg.(type) {
		case Query:
			q := formattedMsg.(Query)

			if q.MinTime > q.MaxTime {
				conn.Write([]byte{0, 0, 0, 0})
				continue
			}

			mean := buffer.query(q)
			response, err := encodeInt(mean)
			if err != nil {
				conn.Write([]byte{0, 0, 0, 0})
				continue
			}
			conn.Write(response[:4])
		case Insert:
			i := formattedMsg.(Insert)
			buffer.insert(i)
		default:
			conn.Write([]byte{0, 0, 0, 0})
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("scanner failed with '%v'\n", err)
	}
}

func main() {
	server := server.New("tcp", "5001")
	server.Start(handleConn)
}
