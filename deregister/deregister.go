package main

import "github.com/hashicorp/consul/api"

func main() {
	//1，初始化consul配置
	consultConfig := api.DefaultConfig()

	//2,创建consul对象
	consulClient, _ := api.NewClient(consultConfig)

	//3，注销服务
	consulClient.Agent().ServiceDeregister("1")

}
