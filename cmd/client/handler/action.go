package handler

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/dacharat/grpc-microservice-go/proto/machine"
)

func runExecute(client machine.MachineClient, instructions *machine.InstructionSet) {
	log.Printf("Executing %v", instructions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := client.Execute(ctx, instructions)
	if err != nil {
		log.Fatalf("%v.Execute(_) = _, %v: ", client, err)
	}
	log.Println(result)
}

func runExecuteStream(client machine.MachineClient, instructions *machine.InstructionSet) {
	log.Printf("Streaming %v", instructions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.ExecuteStream(ctx)
	if err != nil {
		log.Fatalf("%v.Execute(ctx) = %v, %v: ", client, stream, err)
	}
	for _, instruction := range instructions.GetInstructions() {
		if err := stream.Send(instruction); err != nil {
			log.Fatalf("%v.Send(%v) = %v: ", stream, instruction, err)
		}
	}
	result, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	log.Println(result)
}

func runServerStreamingExecute(client machine.MachineClient, instructions *machine.InstructionSet) {
	log.Printf("Executing %v", instructions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.ServerStreamingExecute(ctx, instructions)
	if err != nil {
		log.Fatalf("%v.Execute(_) = _, %v: ", client, err)
	}
	for {
		result, err := stream.Recv()
		if err == io.EOF {
			log.Println("EOF")
			break
		}
		if err != nil {
			log.Printf("Err: %v", err)
			break
		}
		log.Printf("output: %v", result.GetOutput())
	}
	log.Println("DONE!")
}

func runServerStreamingExecuteStream(client machine.MachineClient, instructions []*machine.Instruction) {
	log.Printf("Streaming %v", instructions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.ServerStreamingExecuteStream(ctx)
	if err != nil {
		log.Fatalf("%v.Execute(ctx) = %v, %v: ", client, stream, err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			result, err := stream.Recv()
			if err == io.EOF {
				log.Println("EOF")
				close(waitc)
				return
			}
			if err != nil {
				log.Printf("Err: %v", err)
			}
			log.Printf("output: %v", result.GetOutput())
		}
	}()

	for _, instruction := range instructions {
		if err := stream.Send(instruction); err != nil {
			log.Fatalf("%v.Send(%v) = %v: ", stream, instruction, err)
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err := stream.CloseSend(); err != nil {
		log.Fatalf("%v.CloseSend() got error %v, want %v", stream, err, nil)
	}
	<-waitc
}
