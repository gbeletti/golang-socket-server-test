package server

import (
	"fmt"
	"net"
)

const (
	hello string = "hello"
	hi    string = "hi"
	help  string = "help"
	bye   string = "bye"
)

func dealWithMessage(conn net.Conn, msg string) bool {
	switch msg {
	case hello:
		return writeMessages(conn, "hi")
	case hi:
		return writeMessages(conn, "hello")
	case help:
		return sendHelp(conn)
	case bye:
		writeMessages(conn, "bye")
		return false
	}
	return writeMessages(conn, "I didnt get it. send help to see commands or say bye")
}

func sendHelp(conn net.Conn) bool {
	msgs := []string{
		"I accept the following commands:",
		"    - hi, hello, help, bye",
	}
	return writeMessages(conn, msgs...)
}

func writeMessages(conn net.Conn, msgs ...string) bool {
	for _, msg := range msgs {
		_, err := conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("client left on writing ", err)
			return false
		}
	}
	return true
}
