package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/rimpo/maxx/pkg/max"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	port     = flag.Int("port", 8080, "server listen port e.g: 8080. (default:8080)")
	certFile = flag.String("cert_file", "", "server certificate file path.")
	keyFile  = flag.String("key_file", "", "server private key file path.")
)

type maxServer struct {
}

//max number will be returned in stream
func (m maxServer) FindMaxNumber(stream max.Service_FindMaxNumberServer) error {
	var maxNumber int32 = -1
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			log.Printf("stream closed. %v", err)
			return nil
		}
		if err != nil {
			log.Printf("stream error. %v", err)
			return err
		}

		log.Printf("received:%v", in.Num)
		if maxNumber < in.Num {
			maxNumber = in.Num
		}
		if err := stream.Send(&max.Response{Max: maxNumber}); err != nil {
			return err
		}
	}
}

func main() {
	flag.Parse()

	log.Printf("listen on port %v", *port)

	//List on port
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen on port:%d", *port)
	}

	if *certFile == "" || *keyFile == "" {
		log.Fatal("TLS keys not supplied.")
	}

	// Create the TLS credentials
	creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
	if err != nil {
		log.Fatalf("could not load TLS keys: %s", err)
	}
	// Create an array of gRPC options with the credentials
	opts := []grpc.ServerOption{grpc.Creds(creds)}

	// create a gRPC server object
	server := grpc.NewServer(opts...)

	//register handler
	max.RegisterServiceServer(server, &maxServer{})

	//listen and serve
	server.Serve(listen)
}
