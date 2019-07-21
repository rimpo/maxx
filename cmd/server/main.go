package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/rimpo/maxx/pkg/max"
	"github.com/rimpo/maxx/pkg/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	port     = flag.Int("port", 8080, "server listen port e.g: 8080. (default:8080)")
	certFile = flag.String("cert_file", "", "server certificate file path.")
	keyFile  = flag.String("key_file", "", "server private key file path.")
)

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
	max.RegisterServiceServer(server, &service.Maxx{})

	//listen and serve
	if err := server.Serve(listen); err != nil {
		log.Fatalf("failed to serve. %v", err)
	}
}
