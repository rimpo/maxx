package main

import (
	"context"
	"flag"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/rimpo/maxx/pkg/max"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:8080", "The server address in the format of host:port")
	count      = flag.Int("count", 10, "count of random numbers. (default: 10)")
	certFile   = flag.String("cert_file", "", "server certificate file for TLS")
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	flag.Parse()

	if *certFile == "" {
		log.Fatal("TLS certificate not supplied.")
	}

	// Create the client TLS credentials
	creds, err := credentials.NewClientTLSFromFile(*certFile, "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}

	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(creds))
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
	for i := 0; i < *count; i++ {
		//send random number from 1 to 100
		if err := stream.Send(&max.Request{Num: int32(rand.Intn(100) + 1)}); err != nil {
			log.Fatalf("Failed to send a note: %v", err)
		}
	}
	//close the stream
	stream.CloseSend()

	//wait for close
	<-waitc
}
