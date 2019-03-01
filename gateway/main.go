package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go-grpc-gateway-examples/jwt"
	gw "go-grpc-gateway-examples/proto"
	"go-study/lib/negroni"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

var (
	//常规服务器通信
	echoEndpoint1 = flag.String("hello_endpoint1", "localhost:50051", "endpoint of YourService1")
	//PS：name不可重复
	echoEndpoint2 = flag.String("hello_endpoint2", "localhost:50052", "endpoint of YourService2")

	//基于docker
	//echoEndpoint = flag.String("hello_endpoint", "felixserver:50051", "endpoint of YourService")
)


func run() error {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterMenuHandlerFromEndpoint(ctx, mux, *echoEndpoint1, opts)

	//测试一下 同时注册多个grpc端口是否可行。结论是：可以通信
	err = gw.RegisterMenu2HandlerFromEndpoint(ctx, mux, *echoEndpoint2, opts)
	if err != nil {
		return err
	}
	log.Printf("listen: %v ....", 3344)

	fmt.Println(mux)
	/*
		negroni.Recovery - 异常（恐慌）恢复中间件
		negroni.Logging - 请求 / 响应 log 日志中间件
		negroni.Static - 静态文件处理中间件，默认目录在 "public" 下.
	*/

	n := negroni.Classic()

	//n.UseHandler(NewRouter())
	n.Use(negroni.Wrap(jwt.ValidateTokenMiddleware()(mux)))

	//n.UseHandler(negroni.New(negroni.Wrap(NewRouter()),negroni.Wrap(jwt.ValidateTokenMiddleware()(mux))))


	//n.UseHandler(mux)
	n.Run(":3344")
	return nil

	//return http.ListenAndServe(":3344", mux)
}


func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}