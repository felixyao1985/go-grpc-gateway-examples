package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	gw "go-grpc-gateway-examples/proto"
)

var (
	//常规服务器通信
	echoEndpoint = flag.String("hello_endpoint", "localhost:50051", "endpoint of YourService")

	//基于docker
	//echoEndpoint = flag.String("hello_endpoint", "felixserver:50051", "endpoint of YourService")
)

func run() error {


	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterMenuHandlerFromEndpoint(ctx, mux, *echoEndpoint, opts)
	if err != nil {
		return err
	}
	log.Printf("listen: %v ....", 3344)
	return http.ListenAndServe(":3344", mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}