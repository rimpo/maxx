package main

import (
	"context"
	"flag"
	"io"
	"log"

	"github.com/rimpo/maxx/pkg/max"
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:8080", "The server address in the format of host:port")
)

func main() {

	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect. %v", err)
	}

	defer conn.Close()

	client := max.NewServiceClient(conn)

	stream, err := client.FindMaxNumber(context.Background())
	if err != nil {
		log.Fatalf("failed to get stream. %v", err)
	}

	waitc := make(chan struct{})

	//receive from stream
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}

			log.Printf("current max:%v", in.Max)
		}
	}()

	//send on stream
	for i := 0; i < 10; i++ {
		if err := stream.Send(&max.Request{Num: int32(i)}); err != nil {
			log.Fatalf("Failed to send a note: %v", err)
		}
	}
	//close the stream
	stream.CloseSend()

	//wait for close
	<-waitc
}
