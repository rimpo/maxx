package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/rimpo/maxx/pkg/max"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8080, "server listen port e.g: 8080. (default:8080)")
)

type maxServer struct {
}

//max number will be returned in stream
func (m maxServer) FindMaxNumber(stream max.Service_FindMaxNumberServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("received:%v", in.Num)

		if err := stream.Send(&max.Response{Max: in.Num}); err != nil {
			return err
		}
	}
}

func main() {

	log.Printf("listen on port %v", *port)

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen on port:%d", *port)
	}

	server := grpc.NewServer()

	//register handler
	max.RegisterServiceServer(server, &maxServer{})

	//listen and serve
	server.Serve(listen)
}
