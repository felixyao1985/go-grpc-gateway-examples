##同一主机下的Docker 部署

```
生成docker 镜像 ./docker
docker build -t go-gw .
docker build -t go_server .
```

- 首先，对于Gateway指定的 server要改成 "container_name:port" 模式

	因为docker下 容器间无法直接通信。需要指定物理IP，或者指定容器网络。通常选择容器名称，因为微服务下容器IP 会随时变

- 其次docker 部署需要指定 --network 这里我推荐两种方式

    - Docker DNS Server：
        - docker daemon 实现了一个内嵌的 DNS server，使容器可以直接通过“容器名”通信。方法很简单，只要在启动时用 --name 为容器命名就可以了
	
	    - 使用 docker DNS 有个限制：只能在 user-defined 网络中使用。也就是说，默认的 bridge 网络是无法使用 DNS 的。
	    
    ```
    案例：(代码中使用了此方式demo)
    先创建网络
    docker network create d4_ghost
    
    docker run -it --name=felixgw --net=d4_ghost -p 3344:3344 -d go-gw

    docker run -it --name=felixserver --net=d4_ghost -p 50051:50051 -d go_server

    ```
	- joined 容器 :
	    
	    - joined 容器非常特别，它可以使两个或多个容器共享一个网络栈，共享网卡和配置信息，joined 容器之间可以通过 127.0.0.1 直接通信。
	    
	    - joined 容器非常适合以下场景
	    
	        - 不同容器中的程序希望通过 loopback 高效快速地通信，比如 web server 与 app server。
	        
	        - 希望监控其他容器的网络流量，比如运行在独立容器中的网络监控程序。    
    ```
        docker run -it --name=felixserver  -p 50051:50051 -d go_server
        
        ##指定 jointed 容器为 felixserver
        docker run -it --name=felixgw --network=container:felixserver  -p 3344:3344 -d go-gw
    ```
    
 #切记给代码文件加可执行权限!