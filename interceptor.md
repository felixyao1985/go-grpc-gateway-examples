引用  [带入gRPC：Unary and Stream interceptor](https://studygolang.com/articles/15389?fr=sidebar"With a Title"). 

## interceptor （拦截器）

在 grpc 中进行通信前 与 后的逻辑处理

在 gRPC 中，大类可分为两种 RPC 方法，与拦截器的对应关系是：

- 普通方法：一元拦截器（grpc.UnaryInterceptor）

- 流方法：流拦截器（grpc.StreamInterceptor）

常规实现方式：
```gotemplate
var opts []grpc.ServerOption
// 注册interceptor
var interceptor grpc.UnaryServerInterceptor
interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

    /*
        一些中间逻辑代码
    */
    return handler(ctx, req)
}
opts = append(opts, grpc.UnaryInterceptor(interceptor))

```

grpc 本身只能设置一个拦截器，为了避免所有逻辑都写在可以 可以使用(go-grpc-middleware)来处理

```gotemplate

import "github.com/grpc-ecosystem/go-grpc-middleware"

//RPC 方法的入参出参的日志输出
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    log.Printf("gRPC method: %s, %v", info.FullMethod, req)
    resp, err := handler(ctx, req)
    log.Printf("gRPC method: %s, %v", info.FullMethod, resp)
    return resp, err
}

//RPC 方法的异常保护和日志输出
func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
    defer func() {
        if e := recover(); e != nil {
            debug.PrintStack()
            err = status.Errorf(codes.Internal, "Panic err: %v", e)
        }
    }()

    return handler(ctx, req)
}

    
myServer := grpc.NewServer(
    grpc.Creds(c), //安全通信,处理证书认证相关
    grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
        ...
    )),
    grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
       ...
    )),
)

or 

opts := []grpc.ServerOption{
    grpc.Creds(c),
    grpc_middleware.WithUnaryServerChain(
        RecoveryInterceptor,
        LoggingInterceptor,
    ),
}

```