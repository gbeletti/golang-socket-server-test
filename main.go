package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gbeletti/golang-socket-server-test/server"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	cancel := start()
	defer shutdown(cancel)
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-sigc
	fmt.Printf("got signal [%s]. Time to go\n", s)
}

func start() context.CancelFunc {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	ctx, cancel := context.WithCancel(context.Background())
	server.Start(ctx, host, port)
	return cancel
}

func shutdown(cancel context.CancelFunc) {
	cancel()
	fmt.Println("shutting down...")
	time.Sleep(time.Second * 2)
	fmt.Println("bye")
}
