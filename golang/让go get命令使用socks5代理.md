
不知道是因为什么鬼打墙的原因，GFW把go的官网给禁了，也许因为Go出身Google吧，真是*了狗了。不仅是 golang.org 被墙了，现在Go的包管理站 http://gopkg.in 也被墙了。直接用 go get gopkg.in/package 半天都没反应，只好用代理了。

这里只说明Mac OS的设置。

### 科学上网必备 shadowsocks

配置自己的服务端。客户端使用ShadowsocksX-NG

偏好设置--> 开启HTTP代理服务器，端口号1080

### 设置 git 使用代理
使用代理，终端中输入:

```bash
// 1080 是
git config --global http.proxy 127.0.0.1:1080 
```
如果不使用git代理了，执行命令：

```bash
git config --global --unset-all http.proxy
```

### 设置 go 使用代理
编辑 ~/.bash_profile , 可以执行： vi ~/.bash_profile， 加入:

```bash
alias go='http_proxy=127.0.0.1:1080 go'
```


请随便执行 go get 命令， have fun!.
