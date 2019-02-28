###grpc-gateway

grpc-gateway是protoc的一个插件。它读取gRPC服务定义，并生成一个反向代理服务器，将RESTful JSON API转换为gRPC。此服务器是根据gRPC定义中的自定义选项生成的。

```$xslt
    go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
    go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
    
    进入 github.com/grpc-ecosystem/grpc-gateway/rotoc-gen-grpc-gateway 
    执行go build 然后安装 go install
    进入 github.com/grpc-ecosystem/grpc-gateway/pprotoc-gen-swagger 
    执行go build 然后安装 go install
    
    **
    需要额外下载
    go get github.com/ghodss/yaml
    go get google/api/annotations.proto （被墙使用：github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis）
    git clone https://github.com/protocolbuffers/protobuf
    **
```
#####grpc-gateway-proto文件例子

```gotemplate
    //需要在编译的时候 指定包的路径，并且该包依赖 https://github.com/protocolbuffers/protobuf
    import "google/api/annotations.proto";   
    
    service Menu {
        rpc Save (MenuModel) returns (Res) {
            option (google.api.http) = {
              post: "/api/menu"
              body: "*"
            };
        }

        rpc View (RepMenuView) returns (Res) {
            option (google.api.http) = {
              get: "/api/menu"
            };
        }
    }

```
#####编译
```      
    protoc --go_out=plugins=grpc:. role_manage.proto
    
    protoc -I$GOPATH/src \
      -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
      -I$GOPATH/src/protobuf/src/ \
      -I$GOPATH/src/vsiapi/proto \
      --grpc-gateway_out=logtostderr=true:. \
      role_manage_gw.proto
```

###Golang/glog  简介

####log 等级

基本的 log 等级分为４类：Info, Warning, Error, Fatal.

```

    glog.Info("Test.")
    glog.Flush()

```

####Flush

log 产生后，会暂存在内存的buffer中。只有显示的调用 glog.Flush(), 数据才会真正被写入文件。glog package 的 init 函数中启动了一个 go routine 用来周期性的调用 glog.Flush() 来保证数据被写入文件, 默认的 Flush 周期为30 秒。

当程序运行至 glog.Fatal() 时, glog package 中保证了在退出前程序前会将所有缓存中的log写入文件。但是对于 Info, Warning 以及 Error, 如果程序正常退出，那么在程序退出前 30 秒的 log 就会丢失。defer 可以被用来防止这种情况的发生。

```
    defer glog.Flush()
    glog.Info("Test.")
```

####调用测试
```gotemplate
   交叉编译
    set GOARCH=amd64
    
    set GOOS=linux
    
    go build xx.go


```
$ curl -d '{"ID":1}' -X GET 'http://localhost:3344/api/menu'

#Dokcer 部署
  
  #### https://github.com/felixyao1985/go-grpc-gateway-examples/blob/master/docker.md