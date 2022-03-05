package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/dacharat/grpc-microservice-go/cmd/server/handler"
	"github.com/dacharat/grpc-microservice-go/proto/machine"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 9111, "Port on which gRPC server should listen TCP conn.")
)

func main() {
	flag.Parse()
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	machine.RegisterMachineServer(grpcServer, &handler.MachineServer{})

	errChan := make(chan error)
	stopChan := make(chan os.Signal)

	// bind OS events to the signal channel
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	// run blocking call in a separate goroutine, report errors via channel
	go func() {
		if err := grpcServer.Serve(listen); err != nil {
			errChan <- err
		}
	}()

	// terminate your environment gracefully before leaving main function
	defer func() {
		grpcServer.GracefulStop()
	}()

	// block until either OS signal, or server fatal error
	select {
	case err := <-errChan:
		log.Printf("Fatal error: %v", err)
	case <-stopChan:
	}
}
