package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"kalbek/greet/greetpb"
	"log"
	"time"
)

func main()  {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect: %v", err)
	}
	defer conn.Close()
	temp := greetpb.NewCalculatorServiceClient(conn)
	//FindPrimeNumber(240, temp)
	FindAverage([]int32{2, 4, 6, 8}, temp)
}

func FindPrimeNumber(n int32, temp greetpb.CalculatorServiceClient)  {
	req := &greetpb.IntRequest{Number: n}
	stream, err := temp.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("error with server stream RPC %v", err)
	}
	defer stream.CloseSend()

LOOP:
	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break LOOP
			}
			log.Fatalf("error with response from server stream RPC %v", err)
		}
		log.Printf(fmt.Sprint(res.GetResult(), " "))
	}
}

func FindAverage(numbers []int32, temp greetpb.CalculatorServiceClient) {
	ctx := context.Background()
	stream, err := temp.ComputeAverage(ctx)
	if err != nil {
		log.Fatalf("error related to function ComputeAverage: %v", err)
	}
	for _, n := range numbers {
		stream.Send(&greetpb.IntRequest{Number: n})
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error related to function ComputeAverage: %v", err)
	}
	fmt.Printf("The average of this array is: %v\n", res.GetResult())
}



