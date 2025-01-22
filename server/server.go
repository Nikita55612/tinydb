package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const Logo = `
.-------------------------------------.
|  __  .__                 ._____.    |
|_/  |_|__| ____ ___.__. __| _|_ |__  |
|\   __\  |/    <   |  |/ __ | | __ \ |
| |  | |  |   |  \___  / /_/ | | \_\ \|
| |__| |__|___|  / ____\____ | |___  /|
|              \/\/         \/     \/ |
'-------------------------------------'
`

func Run() {
	println("\033[36m" + Logo + "\033[0m")
	fmt.Println("[INFO] Starting server...")
	fmt.Println("[INFO] Initializing database connection")
	InitDBConn()
	fmt.Println("[INFO] Setting up server listener and context")
	listener, ctx, stop := initialize()
	defer stop()
	addrs := listener.Addr().String()
	log.Println("[INFO] Server is running on:", addrs)
	fmt.Println("[INFO] Starting connection acceptor")
	handleAcceptConn(listener, ctx)
}

func initialize() (net.Listener, context.Context, context.CancelFunc) {
	port := flag.Uint("p", 4004, "Port")
	buff := flag.Uint64("b", 8388608, "Max conn buff size")
	flag.Parse()

	MaxBuffSize = *buff
	addrs := fmt.Sprintf(":%d", *port)
	listener, err := net.Listen("tcp", addrs)
	if err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	go func() {
		<-ctx.Done()
		log.Println("Shutting down server...")
		listener.Close()
	}()

	return listener, ctx, stop
}
