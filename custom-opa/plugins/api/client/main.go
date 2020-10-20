package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/patrick-east/kubecon-na-2020/custom-opa/plugins/api"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(os.Args[1], grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := api.NewAuthorizerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Authz(ctx, &api.AuthzRequest{
		Jwt: os.Args[2],
	})
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Allow: %t\n", r.Allow)
	os.Exit(0)
}
