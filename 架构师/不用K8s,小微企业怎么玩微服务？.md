# 不用K8s,小微企业怎么玩微服务？
​ 
​
微服务，由于独立部署、选型灵活、易扩展等特点，在行业内变得备受关注。而容器、轻量级协议，代码管理、新集成方法与工具等技术的成熟发展，更是促使互联网企业纷纷走上微服务改造的道路。相对于传统的单体架构，它具有以下几个显著的优点：

1. 具有可扩展性：与单体应用相比，微服务架构在扩展新功能时无需对之前实现的功能做过多的调整，而只需增加新的微服务节点并关联已有的微服务节点即可。
2.  技术栈灵活：容器技术使得团队可以为不同的组件使用不同的技术栈。
3.  容错能力强：系统出现故障时，微服务应用只需对有问题的部分进行修正即可，而单体应用则会导致整个服务不可用。

## 微服务不是银弹

微服务本质上是将单体架构从垂直和水平两个维度进行服务拆分。随着业务系统的复杂化，拆分的颗粒度越小，服务越来越多，随之而来的服务间的调用关系也会变得错综复杂，难以管理。因此，使用微服务架构，必然要解决服务治理的问题。下一代的微服务架构Service Mesh，只不过是将服务管理组件，独立成SideCar的机制来减少对业务系统的侵入性，本质上还是在做服务治理。再加之Istio还不够成熟，现在使用为时尚早。

## K8s看似很美好
目前常见的服务治理组件也有很多，springCloud、nacos、tars、etcd、consul等等，也有企业基于K8s实现了服务注册和治理。

众多的解决方案中，最具前景的无疑是K8s，K8s已经成为了云原生时代的操作系统。K8s容器编排平台是一项复杂而令人惊叹的技术。

然而，我们经常忽略的是，使用K8s带来的一系列的挑战。

* 首先K8s要求企业要有具备DevOps能力的运维团队。在小微企业中，连运维都不一定有专职人员，更别说是专业的运维团队了。
* 其次，小微企业的硬件资源有限，仅有的几个实例甚至完全托管于公有云，自行再搭建K8s有些多余。
* 小微企业自身的软件系统，并没有多大的并发量，并且广泛的存在单点系统架构，迁移到K8s存在困难

## 青铜时代(Docker Swarm + Consul + Registrator + Kong)

我的上家雇主就是这样一家不足50人的真小企业。在业务快速迭代，运维团队缺失的情况下，我带领着团队，学习和吸收行业的一些经验，也摸索出了一条独特的微服务架构之路。它不高大上，但一定实用，我把它叫做“青铜时代”。我们使用Docker Swarm + Consul + Registrator + Kong搭建了一套服务治理体系。

![](https://tva1.sinaimg.cn/large/007S8ZIlgy1gjwqiays1vj30k00900v0.jpg)

### 配置
IP地址 | 组件分配
--- | ---
192.168.10.100 | Registrator，consul，swarm manager
192.168.10.101 |  Registrator，consul，swarm worker, docker registry
192.168.10.102 | Registrator，consul，swarm worker
192.168.10.103 | Registrator，kong, konga, postgres

### Docker Network

Docker Network是一种虚拟网络，能够分离容器与真实外界网络的通信，有助于保护系统安全。Docker Network分类众多，本文使用的是Overlay网络，它是多个Docker主机之间的分布式网络，通常用于多主机通信，同一个Overlay网络中的所有节点都是可以相互通信的，却不能对外通信，如果需要的话应该让节点主动开放对外的端口。

创建一个Overlay网络，网络名称为micro

```
docker network create -d overlay --attachable micro
```
###  Docker Swarm
Docker Swarm是Docker内置的一种容器编排工具，它最实用的功能就是一条指令开启多个容器。在Swarm中，节点有Manager和Node两种角色，其中Manager负责容器调度，Node负责听从Manager指挥开启/关闭容器。

在192.168.10.100主机中创建一个Swarm，该节点为Swarm Manager：

```
$ docker swarm init --advertise-addr 192.168.10.100
Swarm initialized: current node (dxn1zf6l61qsb1josjja83ngz) is now a manager.

To add a worker to this swarm, run the following command:

    docker swarm join \
    --token SWMTKN-1-49nj1cmql0jkz5s954yi3oex3nedyz0fb0xx14ie39trti4wxv-8vxv8rssmk743ojnwacrr2e7c \
    192.168.10.100:2377

To add a manager to this swarm, run 'docker swarm join-token manager' and follow the instructions.
```
> Note:指令运行后的输出包含其他主机加入该Swarm的指令
> 

在192.168.10.101、192.168.10.102、192.168.10.103执行命令加入Swarm：

```
docker swarm join \
    --token SWMTKN-1-49nj1cmql0jkz5s954yi3oex3nedyz0fb0xx14ie39trti4wxv-8vxv8rssmk743ojnwacrr2e7c \
    192.168.10.100:2377
```

在Swarm Manager中查看Swarm节点情况:

```
docker node ls 
```

### Docker Registry

Docker Registry是Docker的一款开源项目，用于搭建一个私有的镜像服务器，配合Swarm能够方便各主机之间镜像的分发，同时也能保密系统内的代码。

在192.168.10.101启动registry私有镜像服务器:

```
docker run -d -p 5000:5000 registry 
```

上传本地镜像到registry私有镜像服务器:

```
docker image tag <本地镜像名称> 192.168.10.101:5000/<本地镜像名称>
docker push 192.168.10.101:5000/<本地镜像名称>
```

> Note:docker image tag命令为本地镜像创建一个别名，当执行docker push命令时，docker将自动读取该别名前的ip+port，并将镜像上传到ip+port对应的镜像服务器中。
> 
> Note:执行上述命令若出现http: server gave HTTP response to HTTPS client这种错误信息，则在/etc/docker/daemon.json中追加（若该文件不存在则创建一个同名文件）：
"insecure-registries" : ["<registry容器地址>:5000"]
并使用systemctl restart docker重启docker，强烈建议只在测试环境中使用该方法，生产环境中请使用CA证书，详情请见https://docs.docker.com/registry

拉取registry私有镜像服务器中的镜像:

```
docker pull 192.168.10.101:5000/<镜像名称> 
```
### consul

Consul是一个服务发现和服务配置的解决方案，它是一个一个分布式的，高度可用的系统，而且开发使用都很简便。

在192.168.10.100启动第一个consul实例

```
docker run -d -p 8500:8500 -p 8600:8600 -p 8600:8600/udp -p 8300:8300 -p 8301:8301 -p 8301:8301/udp -p 8302:8302 -p 8302:8302/udp consul agent -server -bootstrap -ui -client=0.0.0.0 -advertise=192.168.10.100
```
在192.168.10.101和192.168.10.102启动另外两个consul实例:

```
(192.168.10.101)$ docker run -d -p 8500:8500 -p 8600:8600 -p 8600:8600/udp -p 8300:8300 -p 8301:8301 -p 8301:8301/udp -p 8302:8302 -p 8302:8302/udp consul agent -server -join=192.168.10.100 -client=0.0.0.0 -advertise=192.168.10.101

(192.168.10.102)$ docker run -d -p 8500:8500 -p 8600:8600 -p 8600:8600/udp -p 8300:8300 -p 8301:8301 -p 8301:8301/udp -p 8302:8302 -p 8302:8302/udp consul agent -server -join=192.168.10.100 -client=0.0.0.0 -advertise=192.168.10.102

```

在浏览器中打开http://192.168.10.100:8500，能够看到控制台页面，则consul集群搭建成功

### Registrator

Registrator是一款开源的自动化服务注册与服务注销工具，使得应用无需编写额外的代码进行服务注册，在应用下线时，也能自动完成服务注销。

在Swarm Manager上为所有Swarm主机上启动Registrator:

```
docker service create --name=registrator --mount type=bind,src=/var/run/docker.sock,dst=/tmp/docker.sock --mode global --network=micro gliderlabs/registrator -deregister="always" -internal consul://192.168.10.100:8500
```

> Note:以下是命令中相关参数的说明
> 
> --name设定该服务的名称是registrator
> 
> --mount设定每个容器挂载主机的目录，这是Registrator启动所要求的
> 
> --mode设定global表示每个Swarm主机都要启动一个Registrator实例
> 
> --network设定该服务内所有实例都属于micro这个Overlay网络
> 
> -deregister设定"always"表示服务下线时自动执行服务注销
> 
> -internal设定Registrator可以自动注册网络内部的实例，方便网络内实例互相调用
> 
> consul://192.168.10.100:8500指明了Consul的服务地址，写consul集群中其他实例的地址也可以

### Kong
kong是一款基于nginx的API网关，与nginx不同的是，kong可以查询DNS信息动态获取服务信息并将请求分发到正在活跃的服务节点上。用户可以在kong的配置中加入用户鉴权、流量控制、日志收集等超实用的功能。

在192.168.10.103启动Kong依赖的数据库postgres：

```
docker run -d --name kong-database \
               --network=micro \
               -p 5432:5432 \
               -e "POSTGRES_USER=kong" \
               -e "POSTGRES_DB=kong" \
               -e "POSTGRES_PASSWORD=kong" \
               postgres:9.6
```

在192.168.10.103查询postgres在Overlay网络中的地址:

```
$ docker container inspect kong-database

...
"Networks": {
  "micro": {
    "IPAMConfig": {
      "IPv4Address": "10.0.2.174"  <-- 这个就是postgres在micro网络中的ip
  },
...
```
在192.168.10.103初始化Kong的数据库:

```
docker run --rm \
     --network=micro \
     -e "KONG_DATABASE=postgres" \
     -e "KONG_PG_HOST=kong-database" \
     -e "KONG_PG_USER=kong" \
     -e "KONG_PG_PASSWORD=kong" \
     -e "KONG_CASSANDRA_CONTACT_POINTS=kong-database" \
     kong:latest kong migrations bootstrap
```
在192.168.10.103启动Kong:

```
docker run -d --name kong\
     --network=micro \
     -e "KONG_DATABASE=postgres" \
     -e "KONG_PG_HOST=10.0.2.174" \
     -e "KONG_PG_USER=kong" \
     -e "KONG_PG_PASSWORD=kong" \
     -e "KONG_CASSANDRA_CONTACT_POINTS=kong-database" \
     -e "KONG_PROXY_ACCESS_LOG=/dev/stdout" \
     -e "KONG_ADMIN_ACCESS_LOG=/dev/stdout" \
     -e "KONG_PROXY_ERROR_LOG=/dev/stderr" \
     -e "KONG_ADMIN_ERROR_LOG=/dev/stderr" \
     -e "KONG_ADMIN_LISTEN=0.0.0.0:8001, 0.0.0.0:8444 ssl" \
     -e "KONG_DNS_RESOLVER=192.168.10.100:8600,192.168.10.101:8600,192.168.10.102:8600" \
     -e "KONG_DNS_ORDER=SRV,LAST,A,CNAME" \
     -p 8000:8000 \
     -p 8443:8443 \
     -p 8001:8001 \
     -p 8444:8444 \
     kong:latest
```

> Note:以下是命令中需要配置的参数说明
> 
> KONG\_PG_HOST是上一步中查询到的postgres在Overlay网络中的地址
> 
> KONG\_DNS_RESOLVER是DNS解析地址列表，本文中开启了三个consul服务，其DNS解析端口默认是8600
> 

测试一下Kong是否启动，若如下命令输出一大堆信息则表示启动成功:

```
curl -i http://192.168.10.103:8001/
```

### Konga

Konga是管理kong的webUI。

在192.168.10.103启动Konga

```
docker run -d -p 1337:1337 --network micro pantsel/konga
```

打开浏览器，进入Konga管理界面：http://192.168.10.103:1337


## 微服务部署示例

### 启动多服务实例
打包制作开发好的应用镜像，并成功上传到Docker Registry以后，我们可以在Swarm Manager启动服务集群

```
docker service create --name=tornado-srv --network=micro --replicas=3 192.168.10.101:5000/tornado-srv:v1 python3 /root/main.py

```
> Note:以下是命令中相关参数的说明
> 
> --name设定服务名称为"tornado-srv"
> 
> --network设定服务网络为“micro”
> 
> --replicas设定服务有3个副本
> 
> 192.168.10.101:5000/tornado-srv:v1是私有镜像名称
> 
> python3 /root/main.py是第一次启动容器实例时执行的命令，这里我们要启动web服务

进入http://192.168.10.100:8500就能看到部署的实例，这里后缀指明了服务的端口，事实上这里我们只要使用8000端口就可以了，即nodejs-srv-8000

![](https://tva1.sinaimg.cn/large/007S8ZIlgy1gjwrca1l1yj30dz08dglw.jpg)

### Kong路由配置

进入konga控制台，首先需要连接Kong，需要填写Name和Kong Admin URL

![](https://tva1.sinaimg.cn/large/007S8ZIlgy1gjwrdpfy6oj30ia0djmy0.jpg)

连接成功后点击add service，填写相关信息

![](https://tva1.sinaimg.cn/large/007S8ZIlgy1gjwrf0s4n0j30k00hyq4d.jpg)

> Note:以下是信息中相关参数的说明
> 
> Name：服务的名称
> 
> Protocol：服务的协议，这里我们只有http
> 
> Host：由于我们配置了DNS自动解析，因此这里我们只需要填写DNS解析域名即可，格式为<服务在consul中的名称>.service.consul
> 
> Port：服务的端口，这里设定容器的开放服务端口为8000
> 

为service添加route

![](https://tva1.sinaimg.cn/large/007S8ZIlgy1gjwrg5808pj30ex0bc74z.jpg)

> Note:以下是信息中相关参数的说明
> 
> Name：route的命名
> 
> Paths：route的路径，每个service可以有多个route
> 

使用Postman测试请求

![](https://tva1.sinaimg.cn/large/007S8ZIlgy1gjwrhl7nnoj30f70ab74t.jpg)


至此，我们的整个微服务架构搭建完成