package service

import (
	"context"
	"io"
	"log"
	"net"
	"testing"
	"time"

	"github.com/rimpo/maxx/pkg/max"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var lis *bufconn.Listener

func bufDialer(string, time.Duration) (net.Conn, error) {
	return lis.Dial()
}

func init() {
	lis = bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	max.RegisterServiceServer(server, &Maxx{})
	go func() {
		log.Print("server started")

		if err := server.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
		log.Print("Server stopped")
	}()
}

func setupClient() (max.ServiceClient, func(), error) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial bufnet: %v", err)
		return nil, nil, err
	}

	client := max.NewServiceClient(conn)

	return client, func() { conn.Close() }, nil

}
func TestMaxer_FindMaxNumber(t *testing.T) {

	cases := []struct {
		name     string
		send     []int32
		expected []int32
	}{
		{
			"greatest_first",
			[]int32{95, 93, 66, 12, 59},
			[]int32{95},
		},
		{
			"increasing_order",
			[]int32{1, 2, 3, 4, 5, 6, 7},
			[]int32{1, 2, 3, 4, 5, 6, 7},
		},
		{
			"random_case",
			[]int32{21, 81, 73, 94, 54},
			[]int32{21, 81, 94},
		},
		{
			"negative_numbers",
			[]int32{-20, -19, -18, -16, -1, -10},
			[]int32{-20, -19, -18, -16, -1},
		},
	}

	for _, c := range cases {

		client, shutdown, err := setupClient()
		if err != nil {
			log.Fatalf("failed to setupClient. %v", err)
		}
		defer shutdown()

		stream, err := client.FindMaxNumber(context.Background())
		if err != nil {
			log.Fatalf("failed to get stream. %v", err)
		}

		waitc := make(chan struct{})

		//start recv
		var recv []int32
		go func() {
			for {
				in, err := stream.Recv()
				if err == io.EOF {
					close(waitc)
					return
				}
				recv = append(recv, in.Max)
			}
		}()

		//start send
		for _, s := range c.send {
			if err := stream.Send(&max.Request{Num: s}); err != nil {
				t.Errorf("failed to send.")
			}
		}

		stream.CloseSend()
		<-waitc

		if len(recv) != len(c.expected) {
			t.Errorf("case %s failed expected:%v got:%v", c.name, c.expected, recv)
		}

		for i := range recv {
			if recv[i] != c.expected[i] {
				t.Errorf("case %s failed expected:%v got:%v", c.name, c.expected, recv)
				break
			}
		}
	}
}
