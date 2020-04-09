package main

import (
	"context"
	"fmt"

	proto "github.com/locpham24/go-authentication/proto"
	"github.com/micro/go-micro/v2"
)

/*

Example usage of top level service initialisation

*/

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	rsp.Greeting = "Hello " + req.Name
	return nil
}

func main() {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(micro.Name("greeter"))

	// Init will parse the command line flags. Any flags set will
	// override the above settings. Options defined here will
	// override anything set on the command line.
	service.Init()

	// By default we'll run the server unless the flags catch us

	// Setup the server

	// Register handler
	proto.RegisterGreeterHandler(service.Server(), new(Greeter))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
