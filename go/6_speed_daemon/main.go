package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	Server "github/ishwar00/protohackers/go/protohacks/server"
	"net"
)

func decodeInt(buf []byte) (int32, error) {
	var intBuf int32
	err := binary.Read(bytes.NewReader(buf), binary.BigEndian, &intBuf)
	if err != nil {
		fmt.Printf("binary.Read(): failed with: %+v\n", err)
		return 0, err
	}

	return intBuf, err
}

func splitByMessages(data []byte, atEOF bool) (int, []byte, error) {

	if atEOF && len(data) == 0 {
		return 0, data, bufio.ErrFinalToken
	}

	switch data[0] {
	case 0x20:
		// Plate
		// plate: str
		// timestamp: u32
		if len(data) > 1 {
			strLen, err := decodeInt(data[1:2])
			if err != nil {
				return 0, data, err
			}
			plateDataLen := 2 + strLen + 4
			if len(data) >= int(plateDataLen) {
				return int(plateDataLen), data[:plateDataLen], nil
			}
		}
	case 0x40:
		// WantHeartBeat
		// interval: u32 (deciseconds)
		intervalLen := 1 + 4
		if len(data) >= intervalLen {
			return intervalLen, data[:intervalLen], nil
		}
	case 0x80:
		// IAmCamera
		// road: u16
		// mile: u16
		// limit: u16 (miles per hour)
		iAmCameraLen := 3*2 + 1
		if len(data) >= iAmCameraLen {
			return iAmCameraLen, data[:iAmCameraLen], nil
		}
	case 0x81:
		// IAmDispatcher
		// numroads: u8
		// roads: [u16] (array of u16)
		if len(data) > 1 {
			numroads, err := decodeInt(data[1:2])
			if err != nil {
				return 0, data, err
			}

			iAmDispatcherLen := 1 + 1 + numroads*2
			if len(data) >= int(iAmDispatcherLen) {
				return int(iAmDispatcherLen), data[:iAmDispatcherLen], nil
			}
		}
	default:
		return 0, data, fmt.Errorf("CLIENT_BAD_DATA: no valid message found.")
	}

	if !atEOF {
		return 0, nil, nil
	}

	return 0, data, bufio.ErrFinalToken
}

func ErrorRes(errMessage string) []byte {
	message := []byte{0x10}
	return append(message, []byte(errMessage)...)
}

func handleConn(clientConn net.Conn) {
	fmt.Printf("handling connection from %v\n", clientConn.RemoteAddr().String())

	scanner := bufio.NewScanner(clientConn)
	scanner.Split(splitByMessages)
	for scanner.Scan() {
		message := scanner.Bytes()
		switch message[0] {
		case 0x81:
		case 0x80:
		default:
			defer clientConn.Close()
			clientConn.Write(ErrorRes("Tell me who the F*ck are you, camera or dispatcher?"))
			return
		}
	}

	if err := scanner.Err(); err != nil {
		defer clientConn.Close()
		fmt.Printf("scanner failed with '%v'\n", err)
		clientConn.Write(ErrorRes(err.Error()))
	}
}

func main() {
	server := Server.New("tcp", "5001")
	server.Start(handleConn)
}
