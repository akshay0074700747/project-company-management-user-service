package services

import (
	"log"
	"net"

	"github.com/akshay0074700747/projectandCompany_management_protofiles/pb/userpb"
	"google.golang.org/grpc"
)

type UserEngine struct {
	Srv userpb.UserServiceServer
}

func NewUserEngine(srv userpb.UserServiceServer) *UserEngine {
	return &UserEngine{
		Srv: srv,
	}
}
func (engine *UserEngine) Start(addr string) {

	server := grpc.NewServer()
	userpb.RegisterUserServiceServer(server, engine.Srv)

	listener, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", addr, err)
	}

	log.Printf("User Server is listening...")

	if err = server.Serve(listener); err != nil {
		log.Fatalf("Failed to listen on port %s: %v", addr, err)
	}
	
}
