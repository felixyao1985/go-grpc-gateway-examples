package main

import (
	"fmt"
	gw "go-grpc-gateway-examples/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"vsiapi/camera"
)

// 业务实现方法的容器
type server struct{}

// 为server定义 DoMD5 方法 内部处理请求并返回结果
// 参数 (context.Context[固定], *test.Req[相应接口定义的请求参数])
// 返回 (*test.Res[相应接口定义的返回参数，必须用指针], error)
func (s *server) List(ctx context.Context, in *gw.RepMenuList) (*gw.Res, error) {

	Meuns := []*gw.MenuModel{}

	return &gw.Res{Code:1,Msg:"GRPC Menu List调用成功",DataInfo:Meuns}, nil
}

func (s *server) View(ctx context.Context, in *gw.RepMenuView) (*gw.Res, error) {
	fmt.Println("Menu View 参数:",in)
	newRow0 := camera.MENU{}

	id := int(in.ID)
	List ,_:= newRow0.View(id)
	Menu := &gw.MenuModel{}
	Menu.ID = int32(List.ID)
	Menu.CODE = List.CODE
	Menu.PID = int32(List.PID)
	Menu.URL = List.URL

	Meuns := []*gw.MenuModel{}
	Meuns = append(Meuns,Menu)

	fmt.Println("列表 List1:",List)
	return &gw.Res{Code:1,Msg:"GRPC Menu View调用成功",DataInfo:Meuns}, nil
}

func (s *server) Save(ctx context.Context, in *gw.MenuModel) (*gw.Res, error) {
	fmt.Println("Menu Save 参数:",in)

	cMENU := camera.MENU{}
	cMENU.CODE = in.CODE
	cMENU.Insert()

	Meuns := []*gw.MenuModel{}

	Menu := &gw.MenuModel{CODE:cMENU.CODE}
	Meuns = append(Meuns,Menu)

	return &gw.Res{Code:1,Msg:"GRPC Menu Save调用成功",DataInfo:Meuns}, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051") //监听所有网卡50051端口的TCP连接

	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}
	s := grpc.NewServer() //创建gRPC服务

	/**注册接口服务
	 * 以定义proto时的service为单位注册，服务中可以有多个方法
	 * (proto编译时会为每个service生成Register***Server方法)
	 * 包.注册服务方法(gRpc服务实例，包含接口方法的结构体[指针])
	 */
	gw.RegisterMenuServer(s, &server{})
	/**如果有可以注册多个接口服务,结构体要实现对应的接口方法
	 * user.RegisterLoginServer(s, &server{})
	 * minMovie.RegisterFbiServer(s, &server{})
	 */
	// 在gRPC服务器上注册反射服务

	reflection.Register(s)
	// 将监听交给gRPC服务处理
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	log.Printf("listen: %v ....", 50051)
}
