package server_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/gbeletti/golang-socket-server-test/server"
)

func TestStart(t *testing.T) {
	host, port := "localhost", "8001"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	server.Start(ctx, host, port)
	t.Log("sleeping one second for server to listen")
	time.Sleep(time.Second * 1)

	address := fmt.Sprintf("%s:%s", host, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		t.Errorf("couldnt connect to server. error [%s]", err)
		return
	}
	err = conn.Close()
	if err != nil {
		t.Errorf("couldnt close connection to server. error [%s]", err)
		return
	}
}

func TestHandleMsg(t *testing.T) {
	tcases := []struct {
		Name     string
		Input    []byte
		Expected string
	}{
		{
			Name:     "01 - new line",
			Input:    []byte("hi\n "),
			Expected: "hi",
		},
		{
			Name:     "02 - new line",
			Input:    []byte("hello\n "),
			Expected: "hello",
		},
	}
	for _, tcase := range tcases {
		t.Run(tcase.Name, func(t *testing.T) {
			testHandleMsg(t, tcase.Name, tcase.Input, tcase.Expected)
		})
	}
}
func testHandleMsg(t *testing.T, name string, input []byte, expected string) {
	got := server.HandleMsg(input)
	if got != expected {
		t.Errorf("test [%s] failed. expected [%s] got [%s]", name, expected, got)
	}
}
