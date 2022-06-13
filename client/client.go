package main

import (
	"context"
	"fmt"
	"go-grpc/proto"
	"log"
	"strconv"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/hashicorp/consul/api"

	"google.golang.org/grpc"
)

func main() {
	//1,初始化consul 配置
	consultConfig := api.DefaultConfig()

	//2，创建consul对象  --(可以重新指定consul属性：IP/Port，也可以使用默认)
	consulClient, err := api.NewClient(consultConfig)
	if err != nil {
		return
	}

	//3，服务发现，从consul上获取健康的服务
	services, _, err := consulClient.Health().Service("grpc and consul", "grpc", true, nil)
	if err != nil {
		return
	}

	// services[0].Service.Port
	addr := services[0].Service.Address + ":" + strconv.Itoa(services[0].Service.Port)
	addrs := string(addr)
	//1 配置grpc服务端的端口作为客户端的监听
	//使用服务发现consul上的IP/Port 来与服务建立连接
	conn, err := grpc.Dial(addrs, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("正在监听服务端 : %v\n", err)
	}

	defer conn.Close()

	//2 实例化 UserInfoService 服务的客户端
	client := proto.NewUserInfoServiceClient(conn)

	//3 调用grpc服务
	req := new(proto.UserRequest)
	req.Name = "YMX"
	resp, err := client.GetUserInfo(context.Background(), req)
	if err != nil {
		log.Fatalf("请求错误 : %v\n", err)
	}
	fmt.Printf("响应内容 : %v\n", resp)
}
