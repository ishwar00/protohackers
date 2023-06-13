package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github/ishwar00/protohackers/go/protohacks/server"
)

type request struct {
	Method string   `json:"method"`
	Number *float64 `json:"number"`
}

type response struct {
	Method  string `json:"method"`
	IsPrime bool   `json:"prime"`
}

func (r *request) isMalformed() error {
	if r.Method != "isPrime" {
		return errors.New("bad request: 'method' is not `isPrime`")
	}

	if r.Number == nil {
		return errors.New("bad request: field 'number' is not present")
	}

	return nil
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("received new connection from %v\n", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Printf("debugInfo: %v\n", scanner.Text())

		var req request
		err := json.Unmarshal(scanner.Bytes(), &req)
		if err != nil {
			fmt.Printf("bad request: Unmarshal() failed on '%v' with error '%v'\n", scanner.Text(), err)
			conn.Write([]byte("{}"))
			return
		}

		if err := req.isMalformed(); err != nil {
			fmt.Println(err.Error())
			conn.Write([]byte("{}"))
			return
		}

		res := response{
			Method:  "isPrime",
			IsPrime: false,
		}

		if *req.Number == float64(int64(*req.Number)) {
			integer := int64(*req.Number)
			res.IsPrime = isPrime(integer)
		}

		jsonRes, err := json.Marshal(res)
		if err != nil {
			fmt.Printf("server error: failed to form response with error '%v'\n", err)
			conn.Write([]byte("{}"))
			return
		}

		_, err = conn.Write(append(jsonRes, byte('\n')))

		if err != nil {
			fmt.Printf("server error: failed to write to the connection with error '%v'\n", err)
			conn.Write([]byte("{}"))
			return
		}
		fmt.Printf("server served successfully '%+v'\n", string(jsonRes))
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("scanner failed with '%v'\n", err)
	}
}

func isPrime(n int64) bool {
	for i := int64(2); i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}

	return n >= 2
}

func main() {
	server := server.New("tcp", "5001")
	server.Start(handleConn)
}
