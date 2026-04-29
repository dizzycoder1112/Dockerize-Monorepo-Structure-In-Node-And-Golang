package handlers

import (
	"context"

	hellopb "dizzycoder1112/Dockerize-Monorepo-Structure-In-Node-And-Golang/grpc/pb/hello"
)

type Greeter struct {
	hellopb.UnimplementedGreeterServer
}

func (g *Greeter) SayHello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloReply, error) {
	return &hellopb.HelloReply{Message: "Hello from Go: " + req.Name}, nil
}
