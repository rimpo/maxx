package service

import (
	"io"
	"log"
	"math"

	"github.com/rimpo/maxx/pkg/max"
)

//Maxx - service sending max number
type Maxx struct {
}

//FindMaxNumber - sends max number on stream whenever max number is encountered
func (m Maxx) FindMaxNumber(stream max.Service_FindMaxNumberServer) error {
	var maxNumber int32 = math.MinInt32
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			log.Print("client disconnected.")
			return nil
		}
		if err != nil {
			log.Printf("stream error. %v", err)
			return err
		}

		if maxNumber < in.Num {
			maxNumber = in.Num
			//send when max
			if err := stream.Send(&max.Response{Max: maxNumber}); err != nil {
				return err
			}
		}

	}
}
