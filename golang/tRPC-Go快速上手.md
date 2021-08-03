# trpc-go的设计

做后端开发的同学肯定接触过不少的框架，常见的框架比如：gin beego echo iris martini。Python也有著名的Django，PHP中的Laravel以及nodejs的koa、express。它们都有个共同的名字，叫“Web框架”。Web框架主要是面向web开发，业务处理的默认都是HTTP请求。而trpc-go则有些不同，它是一个rpc框架。由于客观历史原因，客户端和服务器的直接交互大多还是http请求，但是系统内部大量的微服务之间的交互则并不局限于http协议。从性能和可读性的角度来说，thrift dubbo等协议显然是比http更好的选择。trpc-go作为一个全新的开发框架，如果只锚定了http，显然受众非常有限。因此它的设计中最重要也是最核心的一点，就是支持多种协议。trpc-go中很大一部分数据结构抽象都是围绕着多协议支持这个目的来的，搞清楚这一点可以让你更容易地理解trpc-go。

回想一下我们的日常开发工作，我们在开发业务逻辑的时候其实并不关心底层的通信协议，一般来说业务逻辑最理想都是从如下函数开始的：

```
func handleLogin(req LoginReq) (LoginRsp, error) {

  // your logic

}
```

那么站在框架设计者的角度，如果我们要给用户提供这样的开发体验，我们就必须要做以下工作：

![1615362722422rgdWCL](https://upic-1258482165.cos.ap-chengdu.myqcloud.com/2021-03-10/1615362722422rgdWCL.png)

其实包括标准库在内几乎所有的web框架都是这个处理流程。

## 协议识别


如果要考虑支持多协议的话，事情就变得麻烦了。因为要支持多种协议(http、http2、http3、udp、trpc、grpc)，首先面临的问题就是，框架怎么从socket的字节流中识别出本次请求的协议类型呢？

但是你从实际应用想一个问题，真的有客户端在一个连接中一会用协议A发送数据一会儿用协议B吗？显然对于业务逻辑来说，这是个伪需求！没人会这么干。但是请求服务A用A协议，服务B用B协议是完全可能存在的。所以trpc-go进行了第一个抽象——service。一个trpc-go服务可以包含多个service，每个service监听一个端口，使用一种协议！所以对于trpc-go来说是不需要做协议识别的，在实例化service时就可以确定协议类型，后续都用该协议进行编解码。



## Codec

协议解析是必不可少的一个环节。对于不同的协议，其解析方式都是完全不同的。于是在框架层可以考虑抽象出一个interface（Codec），不同协议的编解码只要实现这个interface就行。有了这层抽象，用户就可以在配置文件中设置服务使用哪个协议进行编解码，然后框架初始化时为该service加载对应协议的编解码模块。这个设计看起来很美好，既能够封装足够多的细节，又能够提供快速扩展的能力。

trpc-go的Codec定义如下：

```go
// Codec 业务协议打解包接口, 业务协议分成包头head和包体body
// 这里只解析出二进制body，具体业务body结构体通过serializer来处理,
// 一般body都是 pb json jce等，特殊情况可由业务自己注册serializer
type Codec interface {
	Encode(message Msg, body []byte) (buffer []byte, err error)
	Decode(message Msg, buffer []byte) (body []byte, err error)
}
```



## Transport

Codec只是负责对协议进行编解码，除此之外还有个Transport的概念也很重要。http2相对http1.1是一个巨大的进步，受限于TCP协议的限制，如果某个请求发生丢包会进行重传，重传就会阻塞其它请求。这种一荣俱荣一损俱损的特点也算是http2的一个不小的瑕疵，因此又有了基于UDP的HTTP3。http3并不是基于TCP的协议，而是基于UDP的。

一般的TCP服务都有一个AcceptLoop，比如：

```go
func ListenAndServe(addr) {
	fd, _ := net.Listen("tcp", addr)
	for {
		conn, err := fd.Accept()
		setting(conn)
		go serveConn(conn)
	}
}
```

在TCP服务中，由于有“连接的概念”，所以我们每次需要在监听句柄上获取新的连接。于是有了上面的AcceptLoop。但是UDP没有连接的概念，监听之后直接读数据，因此它是没有AceeptLoop的。因此trpc-go引入了transport的抽象，不同协议(udp tcp)有一个自己的transport，各自实现自己的ListenAndServe:

```go
// ServerTransport server通讯层接口
type ServerTransport interface {
	ListenAndServe(ctx context.Context, opts ...ListenServeOption) error
}
```



通过Transport和Codec这两层抽象，trpc-go就基本可以做到适配各种不同的网络通信协议，不论是UDP还是TCP，是否multiplexing，这些都能支持。



## Serializer

通信协议主要是为了更好地传输数据，但是具体的数据是啥呢？目前为止我们都是用[]byte这么一个字节流来保存。实际上业务数据的编码方式是各种各样，仅在http协议中，常见的编码方式就有：

- multipart/form-data
- application/x-www-form-urlencoded
- application/json

除了http的这几种序列化方式，我们常用的还包括:

- protobuf
- thrift
- jce

Serializer的作用就是把[]byte反序列化成一个特定的结构体对象，方便业务处理。而协议的payload具体使用什么方式进行序列化和反序列化，一般都是在请求帧的头部指定。因此Codec需要在解析请求时设置Msg的SerializeType，后续处理通过SerializeType找到对应的序列化工具，进行真正的序列化。



## 处理流程

![161536387323053cxO1](https://upic-1258482165.cos.ap-chengdu.myqcloud.com/2021-03-10/161536387323053cxO1.png)



## Filters

Filter类似各种Web框架中的middleware模式。只不过trpc-go在设计时要考虑pb生成的service handler函数签名的不确定性，跟middleware有一些区别。

![1615368142142HGlw06](https://upic-1258482165.cos.ap-chengdu.myqcloud.com/2021-03-10/1615368142142HGlw06.png)

请求先被第一个filter（middleware）处理，然后依次交给后续的filter，最后交给业务逻辑处理。处理完之后又从内层的filter一层一层向外层返回，最后返回给客户端。这种模式其实就是trpc-go文档中说的支持在业务逻辑前和后进行hook，看这个图就能理解，请求先一层一层地进入洋葱，再一层一层地出来，这两个地方都可以写逻辑，也就是所谓的hook。



## Plugin

Plugin和trpc-go是比较松耦合，或者可以说，“基本没啥关系”。正是通过各种plugin，我们的trpc-go服务才能无缝地和公司的各种系统对接。

trpc-go约定了一种配置格式，插件的配置都按那种格式配置到trpc.yaml中，然后trpc-go根据插件名找到对应插件，然后把该插件的配置分离出来，用这部分小配置去实例化插件。这样对插件唯一的约束就是插件必须要支持通过yaml格式进行配置，至于yaml配置的格式完全可以自定义。

# 开发示例



## 环境准备



 protoc：

```sh
brew install protoc
```

protoc-gen-go：

```shell
brew install protoc-gen-go		
```

protoc-gen-secv：

```shell
go get git.code.oa.com/devsec/protoc-gen-secv
```

trpc-go-cmdline：

```sh
go get git.code.oa.com/trpc-go/trpc-go-cmdline/trpc
```



## 服务选型

### 标准HTTP服务

用**“net/http”**标准库方式来开发HTTP服务，同时也能复用框架的服务治理能力，如自动上报监控，模调，调用链等关键信息。支持http，https，http2 和http3协议。

直接注册接口的URL和处理函数的方式：

```go
// URL注册函数：pattern 为http请求的URL， handler 为路由处理函数
func HandleFunc(pattern string, handler func(w http.ResponseWriter, r *http.Request) error)
// 注册HTTP标准服务
func RegisterDefaultService(s server.Service)
```

示例代码：

```go
package main
import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "git.code.oa.com/trpc-go/trpc-go/log"
    trpc "git.code.oa.com/trpc-go/trpc-go"
    thttp "git.code.oa.com/trpc-go/trpc-go/http"
)
// Data 请求报文数据
type Data struct {
    Msg string
}
func handle(w http.ResponseWriter, r *http.Request) error {
    // 获取请求报文头里的 "request" 字段
    reqHead := r.Header.Get("request")
    // 获取请求报文中的数据
    msg, _ := ioutil.ReadAll(r.Body)
    log.Infof("data is %s, request head is %s\n", msg, reqHead)
    // 为响应报文设置Cookie
    cookie := &http.Cookie{Name: "sample", Value: "sample", HttpOnly: false}
    http.SetCookie(w, cookie)
    // 注意: 使用ResponseWriter回包时，Set/WriteHeader/Write这三个方法必须严格按照以下顺序调用
    w.Header().Set("Content-type", "application/json")
    // 为响应报文头添加 “reply” 字段
    w.Header().Add("reply", "tested")
    // 为响应报文设置HTTP状态码
    // w.WriteHeader(403)
    // 为响应报文设置Body
    rsp, _ := json.Marshal(&Data{Msg: "Hello, World!"})
    w.Write(rsp)
    return nil
}
func main() {
    s := trpc.NewServer()
    // 路由注册
    thttp.HandleFunc("/v1/hello", handle)
    // 服务注册
    thttp.RegisterDefaultService(s)
    s.Serve()
}
```

### HTTP RPC服务

泛HTTP RPC服务是RPC服务，服务调用接口由PB文件定义，可以由工具生产桩代码。而泛HTTP标准服务是一个普通的HTTP服务，不使用PB文件定义服务，用户需要自己编写代码定义服务接口，注册URL，组包和解包HTTP报文。

#### 定义服务接口

tRPC采用protobuf来描述一个服务，我们用protobuf定义服务方法，请求参数和响应参数。helloworld.proto:

```protobuf
syntax = "proto3";
package trpc.demo.hello;
option go_package = "git.code.oa.com/trpcprotocol/linyyyang/demo/hello";

service Hello {
    rpc SayHello(HelloRequest) returns (HelloReply) {};
}
// 请求参数
message HelloRequest {
    string msg = 1;
}
// 响应参数
message HelloReply {
    string msg = 1;
}
```



这里我们定义了一个`Greeter`服务，这个服务里面有个`SayHello`方法，接收一个包含`msg`字符串的`HelloRequest`参数，返回`HelloReply`数据。
这里需要注意以下几点：

- `syntax`必须是`proto3`，tRPC都是基于proto3通信的。
- `package`内容格式推荐为`trpc.{app}.{server}`，以trpc为固定前缀，标识这是一个trpc服务协议，app为你的应用名，server为你的服务进程名。
- `package`后面必须有`option go_package="git.code.oa.com/trpcprotocol/{app}/{server}";`，指明你的pb.go生成文件的git存放地址，协议与服务分离，方便其他人直接引用，git地址用户可以自己随便定，也可以使用tRPC-Go提供的公用group：`trpcprotocol`。更加详细的接口管理见[这里](https://iwiki.woa.com/pages/viewpage.action?pageId=99485686)。
- 定义`rpc`方法时，一个`server`（服务进程）可以有多个`service`（对`rpc`逻辑分组），一般是一个`server`一个`service`，一个`service`中可以有多个`rpc`调用。
- 编写protobuf时必须遵循谷歌[官方规范](https://developers.google.com/protocol-buffers/docs/style)。

#### 生成服务代码

首次使用，用该命令生成完整工程，当前目录下不要出现跟pb同名的目录名，如pb名为helloworld.proto，则当前目录不要出现helloworld的目录名：

```shell
trpc create -p helloworld.proto -o . --protocol=http
```



### TRPC服务

 tRPC服务使用tRPC私有协议。这是与HTTP RPC服务的区别。

#### 生成服务代码

```shell
trpc create --protofile=helloworld.proto
```

如果在现有的rRPC-Go项目进行新增功能，则只生成rpcstub代码：

```shell
trpc create -p helloworld.proto --rpconly -o
```



## 项目结构

tRPC-Go推荐每个服务server以一个repo为粒度放在具体的app下面来组织管理服务结构，每个服务独立仓库。在微服务领域里面无可厚非，但是实际业务场景中，现有的项目结构往往很难遵循，码客上对此也有争论：

https://mk.woa.com/q/261629?ADTAG=rtx

个人觉得适当的按照业务领域进行拆分，并非一定要遵循颗粒度这么细的拆分模式。且tRPC-Go该条规范是推荐性而非强制

## 调试

Trpc-cli :

```shell
trpc-cli -func /trpc.demo.test.Hello/SayHello -target ip://127.0.0.1:8000 -body '{"msg":"hello"}'
```



## 错误处理

tRPC-Go推荐在写服务端业务逻辑时，使用tRPC-Go封装的errors.New()来返回业务错误码，这样框架能自动上报业务错误码到监控系统。如果业务自定义error的话，就只能靠业务主动调用Metrics SDK来上报错误码。关于错误码的API使用， 请参考 [这里](http://godoc.oa.com/git.code.oa.com/tRPC-Go/tRPC-Go/errs)。

tRPC-Go对错误码的数据类型和含义都做了规划， 对于常见错误码的问题定位也都做了解释。具体请参考 [tRPC-Go错误码手册](https://iwiki.woa.com/pages/viewpage.action?pageId=276029299)。

## URL 自定义

有两种办法定义URL， 一是通过命令参数alias：

```shell
trpc create -p helloworld.proto --protocol=http --alias=/api/hello
```

也可以在proto文件内定义：

```protobuf
package hello;
import "trpc.proto";
service HelloSvr {
  rpc Hello(Req) returns(Rsp) {option (trpc.alias) = "/api/hello"; };
}
```



## 参数验证

**SECV**是一款protoc插件，用于RPC程序的自动化数据校验（Validation）。

安装 https://git.code.oa.com/devsec/protoc-gen-secv

```shell
# 下载插件源代码至$GOPATH
go get -d git.code.oa.com/devsec/protoc-gen-secv

# 执行命令，将SECV安装至$GOPATH/bin
make build
```

编写proto文件，在扩展字段设置Validation规则。例如：

```protobuf
string email = 2 [(validate.rules).string.email = true];
```

示例项目

 https://git.code.oa.com/yuyangzhou/trpc_validate_demo

也可以使用参数校验拦截器

https://git.code.oa.com/cooperyan/trpc-filter/tree/master/validation

## 超时控制

RPC-Go的超时控制全部通过配置文件指定即可。
注意：以下设置的均是当前服务自身的超时配置，不是上游对自己的超时配置。

### 链路超时

超时时间默认会从最源头服务一直通过协议字段透传下去，用户可以自己配置开关开启或关闭。

```yaml
server:  service:    
  - name: trpc.app.server.service      
    disable_request_timeout: true  #默认false，默认超时时间会继承上游的设置时间，配置true则禁用，表示忽略上游服务调用我时协议传递过来的超时时间
```

### 消息超时

每个服务启动时都可配置该服务所有请求的最长处理超时时间， 该时间只用于调用下游服务时的单个调用超时。如果服务端处理调用请求的时间超过了消息超时时间，处理协程不会立马结束。

```yaml
server:  service:    
  - name: trpc.app.server.service      
    timeout: 1000  #单位ms，每个接收到的请求最多允许1000ms的执行时间，所以要注意权衡当前请求内的所有串行rpc调用的超时时间分配
```

### 调用超时

每个rpc后端调用都可以配置当次调用请求的最大超时时间，如果代码里面有设置`WithTimeout Option`，则`调用超时以代码为准，该配置不生效`，代码不够灵活，建议不要在代码里面设置`WithTimeout Option`。

```yaml
client:  service:    
  - name: trpc.app.server.service  #后端服务协议文件的service name，格式为：pbpackagename.pbservicename      
    timeout: 500  #单位ms，每个发起的请求最多允许500ms的超时时间
```

每次rpc请求会取 `链路超时` `消息超时` `调用超时` 的最小值来调用后端，当前消息的最长处理超时时间会实时计算剩余时间。



## 自定义序列化类型

如果框架自带的序列化类型不满足业务需求，业务可以自定义序列化类型。

```go
import (
    thttp "git.code.oa.com/trpc-go/trpc-go/http"
)
// ExampleSerialization 
type ExampleSerialization struct {
}
// Unmarshal 反序列
func (s *ExampleSerialization) Unmarshal(in []byte, body interface{}) error {
    // 业务需要实现把in反序列化的数据写到body中
    ...
}
// Marshal 序列化
func (s *ExampleSerialization) Marshal(body interface{}) ([]byte, error) {
    // 业务需要实现把body的数据序列化，并返回
    ...
}
func init() {
    thttp.RegisterSerializer("application/x-example", 1101, &ExampleSerialization{})
}
```



## 自定义返回处理函数

典型场景：服务端使用HTTP RPC模式开发，但客户端使用HTTP请求，并且要求HTTP错误响应报文遵循以下格式：

```
{
    "retcode": 10000,
    "retmsg": "服务器超载"
}
```

自定义错误拦截:

```go
thttp.DefaultServerCodec.ErrHandler = func(w http.ResponseWriter, r *http.Request, e *errs.Error) {
        // 填充指定格式错误信息到HTTP Body
        w.Write([]byte(fmt.Sprintf(`{"retcode":%d, "retmsg":"%s"}`, e.Code, e.Msg)))
    }
```

自定义返回拦截:

```go
type Response struct {
    Code    int32           `json:"code"`
    Message string          `json:"message"`
    Data    json.RawMessage `json:"data"`
}

thttp.DefaultServerCodec.RspHandler = func(w http.ResponseWriter, r *http.Request, rspbody []byte) (err error) {
        if len(rspbody) == 0 {
            return nil
        }
        bs, _ := json.Marshal(&Response{Code: 0, Message: "OK", Data: rspbody})
        _, err = w.Write(bs)
        return
    }
```



## 协议代理

服务作为一个Proxy，接收标准HTTP服务请求，然后转化成tRPC协议格式向后端的tRPC服务发送请求。 示例代码如下：

```go
package main
import (
    "context"
    "io/ioutil"
    "net/http"
    "git.code.oa.com/trpc-go/trpc-go/client"
    pb "git.code.oa.com/trpcprotocol/test/helloworld"
    trpc "git.code.oa.com/trpc-go/trpc-go"
    thttp "git.code.oa.com/trpc-go/trpc-go/http"
)
func handle(w http.ResponseWriter, r *http.Request) error {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "can't read body", http.StatusBadRequest)
        return nil
    }
    proxy := pb.NewGreeterClientProxy()
    req := &pb.HelloRequest{Msg: string(body[:])}
    // 向tRPC服务请求
    rsp, err := proxy.SayHello(context.Background(), req, client.WithTarget("ip://127.0.0.1:8001"))
    if err != nil {
        http.Error(w, "call fails！", http.StatusBadRequest)
        return nil
    }
    // 回响应给HTTP客户端
    w.Header().Set("Content-type", "application/text")
    w.Write([]byte(rsp.Msg))
    return nil
}
func main() {
    s := trpc.NewServer()
    // 路由注册
    thttp.HandleFunc("/v1/hello", handle)
    // 服务注册
    thttp.RegisterDefaultService(s)
    s.Serve()
}
```



## 测试以及Mock

todo

## 调用RPC服务

引入pb协议:

```go
import pb "git.code.oa.com/trpcprotocol/linyyyang/demo/hello" // 引入pb协议
```

rpc调用

```go
testProxy := testData.NewGreeterClientProxy()
	opts := []client.Option{
		client.WithNamespace("Development"),
		client.WithTarget("ip://127.0.0.1:8989"), //如果使用了注册发现，则此处无需配置
	}
	t, err := testProxy.SayHello(ctx, &testData.HelloRequest{
		Msg: req.Msg,
	}, opts...)
	if err != nil {
		log.Errorf("SayTest fail:%v", err)
		return err
	}
```



## 使用mysql

编辑项目配置，引入mysql client：

```yaml
client:                                            #客户端调用的后端配置
  timeout: 1000                                    #针对所有后端的请求最长处理时间
  namespace: Development                           #针对所有后端的环境
  filter:                                          #针对所有后端调用函数前后的拦截器列表
  service:                                         #针对单个后端的配置
     - name: trpc.mysql.demo.openplatform
       target: dsn://root:123456@tcp(127.0.0.1:3306)/db?timeout=1s&parseTime=true&interpolateParams=true
```



修改proto文件定义协议:

```protobuf

service Hello {
    rpc SayHello(HelloRequest) returns (HelloReply) {};
    rpc GetDB(GetDBRequest) returns (GetDBReply) {}; //增加GetDB
}

// 请求参数
message HelloRequest {
    string msg = 1;
}
// 响应参数
message HelloReply {
    string msg = 1;
}
message GetDBRequest {
    int64 id = 1;
}
message GetDBReply {
    int64 id    = 1;
    string name = 2;
}
```

重新生成代码:

```shell
trpc create -p proto/hello.proto --rpconly -o stub/git.code.oa.com/trpcprotocol/linyyyang/demo/hello -f
```

编写repo:

```go
package repository

import (
	"context"

	"git.code.oa.com/trpc-go/trpc-database/mysql"
	"git.code.oa.com/trpc-go/trpc-go/log"
	pb "git.code.oa.com/trpcprotocol/linyyyang/demo/hello"
	"go.uber.org/zap"
)

// GetUserById
func GetUserById(ctx context.Context, id int64) (*pb.GetDBReply, error) {
	var users []*pb.GetDBReply

	proxy := mysql.NewClientProxy("trpc.mysql.demo.openplatform") // 必须要跟client配置的name一致
	// 读取数据，select字段尽量只select自己关心的字段，不要用*，以下只是简单示例
	err := proxy.QueryToStructs(ctx, &users, "select id, name from user where id=? limit 1", id)
	log.DebugContext(ctx, "GetUserById",
		zap.String("sql", "select id, name from user where id=? limit 1"),
		zap.Int64("id", id),
	)
	if err != nil || len(users) == 0 {
		log.ErrorContext(ctx, "GetUserById",
			zap.String("sql", "select id, name from user where id=? limit 1"),
			zap.Int64("id", id),
			zap.Error(err),
		)
		return nil, err
	}
	return users[0], nil
}

```

实现服务逻辑:

```go
// GetDB
func (s *HelloServiceImpl) GetDB(ctx context.Context, req *pb.GetDBRequest, rsp *pb.GetDBReply) error {
	u, err := repository.GetUserById(ctx, req.Id)
	if err != nil {
		return errs.New(404, "user not found")
	}
	rsp.Id = u.Id
	rsp.Name = u.Name
	return nil
}

```



访问测试:

```shell
trpc-cli -func /trpc.demo.test.Hello/GetDB -target ip://127.0.0.1:8000 -body '{"id":5}'
```



## 使用Redis

编辑项目配置，引入redis client：

```yaml
client: #客户端调用的后端配置
  timeout: 1000 #针对所有后端的请求最长处理时间
  namespace: Development #针对所有后端的环境
  filter: #针对所有后端调用函数前后的拦截器列表
  service: #针对单个后端的配置
    - name: trpc.redis.tgs-demo.test
      namespace: Development
      target: ip://host:port
      password:
      timeout: 800 #当前这个请求最长处理时间
```

使用：

```go
import (
	"context"
	"time"

	"git.code.oa.com/trpc-go/trpc-database/redis"
)

const redisService = "trpc.redis.tgs-demo.test"
// Get
func Get(ctx context.Context, key string) (string, error) {
	c := redis.NewClientProxy(redisService)
	rep, err := redis.String(c.Do(ctx, "GET", key))
	if err != nil {
		return "", err
	}
	return rep, nil
}
```





## 服务发现consul

Todo

## 配置中心etcd

Todo

## 调用链SkyWalking/zipkin

Todo

## 日志收集elk

Todo

# 资料
* [tRPC-Go](https://iwiki.oa.tencent.com/pages/viewpage.action?pageId=279550562)
* [trpc-go新人落地step by step](http://km.oa.com/articles/show/472111)
* [trpc-go入门篇-step by step](http://km.oa.com/articles/show/470544?kmref=search&from_page=1&no=2)
* [内容安全部trpc最佳实践标准](https://iwiki.woa.com/pages/viewpage.action?pageId=447432816)
* [trpc-go 内容安全部实践标准](http://km.oa.com/group/46279/articles/show/439640?kmref=search&from_page=1&no=3)
* [trpc-go字段校验规则](https://git.code.oa.com/devsec/protoc-gen-secv/wikis/%E6%A0%A1%E9%AA%8C%E8%A7%84%E5%88%99/)
* [腾讯云微服务治理平台(TSF)现已支持tRPC-go框架](http://km.oa.com/group/33359/articles/show/453310?kmref=search&from_page=1&no=5)
* [TSF平台部署tRPC-Go服务Step By Step](http://km.oa.com/group/17746/articles/show/458194?kmref=search&from_page=1&no=3)
* [redigo支持单机，哨兵，集群](https://segmentfault.com/a/1190000017879129)
* [tRPC快速接入TSF](https://git.code.oa.com/qcloud_middleware/tsf-go-trpc/blob/master/doc/TRPC.md)





