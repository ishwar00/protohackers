package main

import (
	"bufio"
	"fmt"
	"github/ishwar00/protohackers/go/protohacks/server"
	"github/ishwar00/protohackers/go/protohacks/utils"
	"net"
	"strings"
	"sync"
)

type UserConn struct {
	id       int
	conn     *net.Conn
	userName string
}

func (u *UserConn) write(msg string) error {
	_, err := (*u.conn).Write([]byte(msg))
	return err
}

type ChatRoom struct {
	mu           *sync.Mutex
	userConns    map[int]*UserConn
	lastIdOfUser int
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		mu:           &sync.Mutex{},
		userConns:    map[int]*UserConn{},
		lastIdOfUser: 0,
	}
}

func (c *ChatRoom) addUserConn(u *UserConn) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.userConns[u.id]
	if ok {
		msg := fmt.Sprintf("something is not right, found id %d to be already assigned", u.id)
		panic(msg)
	}

	msg := fmt.Sprintf("* %s has entered the romm!!\n", u.userName)
	for _, userConn := range c.userConns {
		userConn.write(msg)
	}
	c.userConns[u.id] = u
}

func (c *ChatRoom) generateNewId() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.lastIdOfUser = c.lastIdOfUser + 1
	return c.lastIdOfUser
}

func (c *ChatRoom) broadcastMessage(u *UserConn, msg string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for id, userConn := range c.userConns {
		if id != u.id {
			userConn.write(msg)
		}
	}
}

func (c *ChatRoom) getActiveUserNames() []string {
	userNames := []string{}

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, userConn := range c.userConns {
		userNames = append(userNames, userConn.userName)
	}
	return userNames
}

func (c *ChatRoom) removeUser(u *UserConn) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.userConns[u.id]
	if !ok {
		msg := fmt.Sprintf("something is not right, id %d was not found", u.id)
		panic(msg)
	}

	delete(c.userConns, u.id)
}

func HandleNewUser(conn net.Conn, c *ChatRoom) {
	scanner := bufio.NewScanner(conn)

	conn.Write([]byte("what's your name?\n"))

	if scanner.Scan() {
		name := scanner.Text()
		if len(name) > 1 && utils.IsAlhpaNumericStr(name) {
			activeUsers := c.getActiveUserNames()
			presentUsers := fmt.Sprintf("* room contains: %s\n", strings.Join(activeUsers, ","))
			conn.Write([]byte(presentUsers))
			user := &UserConn{
				id:       c.generateNewId(),
				conn:     &conn,
				userName: name,
			}

			c.addUserConn(user)

			for scanner.Scan() {
				userMessage := scanner.Text()
				chat := fmt.Sprintf("[%s] %s\n", user.userName, userMessage)
				c.broadcastMessage(user, chat)
			}

			leftMsg := fmt.Sprintf("* %s left has the room ;(\n", user.userName)
			c.broadcastMessage(user, leftMsg)
			c.removeUser(user)
		}
	}

	defer func() {
		fmt.Printf("closing connection from %s\n", conn.RemoteAddr())
		conn.Close()
	}()

	if err := scanner.Err(); err != nil {
		fmt.Printf("scanner failed with '%v'\n", err)
	}
}

func main() {
	server := server.New("tcp", "5001")

	room := NewChatRoom()
	server.Start(
		func(conn net.Conn) {
			HandleNewUser(conn, room)
		},
	)
}
