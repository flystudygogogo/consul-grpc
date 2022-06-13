package main

import (
	"context"
	"fmt"

	"go-grpc/proto"
	"log"
	"net"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

//定义服务端 实现 约定的接口
type UserInfoService struct{}

var u = UserInfoService{}

//实现 interface
func (s *UserInfoService) GetUserInfo(ctx context.Context, req *proto.UserRequest) (resp *proto.UserResponse, err error) {
	name := req.Name
	if name == "YMX" {
		resp = &proto.UserResponse{
			Id:    1,
			Name:  name,
			Age:   22,
			Title: []string{"Java", "Go"},
		}
	}
	err = nil
	return
}
func main() {
	//把grpc项目，注册到consul
	//1,初始化consul配置
	consultConfig := api.DefaultConfig()

	//2,创建consul对象
	consulClient, err := api.NewClient(consultConfig)
	if err != nil {
		fmt.Println("api.newClient err:", err)
		return
	}
	//3,告诉consul，即将注册的服务的配置信息
	registerService := api.AgentServiceRegistration{
		ID:      "1",
		Tags:    []string{"grpc", "consul"},
		Name:    "grpc and consul",
		Address: "127.0.0.1",
		Port:    8800,
		Check: &api.AgentServiceCheck{
			CheckID:  "consul grpc test",
			TCP:      "127.0.0.1:8800",
			Timeout:  "1s",
			Interval: "5s",
		},
	}

	//4,注册服务到consul上
	consulClient.Agent().ServiceRegister(&registerService)

	//1 添加监听的端口
	port := ":8800"
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("端口监听错误 : %v\n", err)
	}
	fmt.Printf("正在监听： %s 端口\n", port)
	//2 启动grpc服务
	s := grpc.NewServer()
	//3 将UserInfoService服务注册到gRPC中
	// 注意第二个参数 UserInfoServiceServer 是接口类型的变量，需要取地址传参
	proto.RegisterUserInfoServiceServer(s, &u)
	s.Serve(l)
}
