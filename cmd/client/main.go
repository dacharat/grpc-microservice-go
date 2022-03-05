package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/dacharat/grpc-microservice-go/cmd/client/handler"
	"github.com/dacharat/grpc-microservice-go/proto/calculator"
	"github.com/dacharat/grpc-microservice-go/proto/machine"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var (
	serverAddr  = flag.String("server_addr", "localhost:9111", "The server address in the format of host:port")
	serverAddr2 = flag.String("server_addr2", "localhost:9112", "The server address in the format of host:port")
)

func main() {
	flag.Parse()

	machineClient, close := initMachineService()
	defer close()

	calculatorClient, close2 := initCalculatorService()
	defer close2()

	h := handler.NewHandler(machineClient, calculatorClient)
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		time.Sleep(1 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	router.GET("/test", h.TestHandler)
	router.GET("/test2", h.Test2Handler)
	router.GET("/stream", h.TestStreamHandler)
	router.GET("/stream2", h.TestStream2Handler)
	router.POST("/instruction", h.InstructionHandler)
	router.POST("/calculator", h.CalculatorHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func initMachineService() (machine.MachineClient, func()) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	client := machine.NewMachineClient(conn)

	return client, func() { conn.Close() }
}

func initCalculatorService() (calculator.CalulatorClient, func()) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr2, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	client := calculator.NewCalulatorClient(conn)

	return client, func() { conn.Close() }
}
