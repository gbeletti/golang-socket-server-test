package server

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
)

const (
	defaultHost string = "localhost"
	defaultPort string = "8001"
	connType    string = "tcp"
)

// Start starts the server
func Start(ctx context.Context, host, port string) {
	if len(host) == 0 {
		host = defaultHost
	}
	if len(port) == 0 {
		port = defaultPort
	}
	go runServer(ctx, host, port)
}

func runServer(ctx context.Context, host, port string) {
	address := fmt.Sprintf("%s:%s", host, port)
	l, err := net.Listen(connType, address)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := l.Close(); err != nil {
			log.Printf("couldnt close server. error [%s]\n", err)
		}
	}()
	fmt.Println("server is ready and accepting connections")
	for {
		connCh := acceptConnection(l)
		select {
		case c := <-connCh:
			if c == nil {
				return
			}
			fmt.Println("client connected " + c.RemoteAddr().String())
			go handleConnection(ctx, c)
		case <-ctx.Done():
			return
		}
	}
}

func acceptConnection(l net.Listener) (connCh chan net.Conn) {
	connCh = make(chan net.Conn)
	go func() {
		c, err := l.Accept()
		if err != nil {
			log.Println("stopping accepting connection")
		}
		connCh <- c
		close(connCh)
	}()
	return
}

func handleConnection(ctx context.Context, conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("couldnt close connection. error [%s]\n", err)
		}
	}()
	for {
		msgCh := readFromConnection(conn)
		var msg []byte
		select {
		case msg = <-msgCh:
			if len(msg) == 0 { // no message, client is gone
				fmt.Println("client left on reading")
				return
			}
		case <-ctx.Done(): // server is going down
			return
		}
		if !dealWithMessage(conn, handleMsg(msg)) {
			return
		}
	}
}

func readFromConnection(conn net.Conn) (msgCh chan []byte) {
	msgCh = make(chan []byte)
	go func() {
		msg, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			log.Println("stopping reading from connection")
		}
		msgCh <- msg
		close(msgCh)
	}()
	return msgCh
}

func handleMsg(msg []byte) string {
	m := string(msg[:len(msg)-2])
	return m
}
