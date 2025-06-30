package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"zg6zy5/apiway/inits"
	"zg6zy5/service/srvhandler"

	"google.golang.org/grpc"
	__ "zg6zy5/service/protobuf"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.

func main() {
	inits.ExampleClient()
	inits.InitMysql()
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	__.RegisterUserServer(s, &srvhandler.Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
