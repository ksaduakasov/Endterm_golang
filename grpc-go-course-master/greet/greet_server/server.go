package main

import (
	"google.golang.org/grpc"
	"io"
	"kalbek/greet/greetpb"
	"log"
	"net"
	"time"
)

type Server struct {
	greetpb.UnimplementedGreetServiceServer
}
type CalculatorService struct{
	greetpb.UnimplementedCalculatorServiceServer
}


func (s *CalculatorService) PrimeNumberDecomposition(req *greetpb.IntRequest, stream greetpb.CalculatorService_PrimeNumberDecompositionServer) (error) {
	num := req.GetNumber()
	primes := FindPrimeNumber(num)
	for i := 0; i < len(primes); i++ {
		res := &greetpb.IntResponse{Result: primes[i]}
		if err := stream.Send(res); err != nil {
			log.Fatalf("error with responses: %v", err.Error())
		}
		time.Sleep(time.Second)
	}
	return nil
}

func FindPrimeNumber(n int32) []int32 {
	var arr []int32
	for{
		arr = append(arr, 2)
		n /= 2
		if n % 2 != 0{
			break
		}
	}
	var i int32 = 0
	for i = 3; i <= n*n; i+=2{
		for {
			arr = append(arr, 3)
			n /= i
			if n % i != 0 {
				break
			}
		}
	}
	if n > 2 {
		arr = append(arr, n)
	}

	return arr
}

func (s *CalculatorService) ComputeAverage(stream greetpb.CalculatorService_ComputeAverageServer) error {
	var sum int32 = 0
	var count int32 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			response := &greetpb.AvgResponse{Result: float64(sum / count)}
			return stream.SendAndClose(response)
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		sum += req.GetNumber()
		count++
	}
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil{
		log.Fatalf("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterCalculatorServiceServer(s, &CalculatorService{})
	if err := s.Serve(l); err != nil{
		log.Fatalf("Failed to serve: %v", err)
	}
}

