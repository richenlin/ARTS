# OpenResty从入门到开发一个网关服务

## 简介

OpenResty（也称为 ngx_openresty）是一个全功能的 Web 应用服务器。它打包了标准的 Nginx 核心，很多的常用的第三方模块，以及它们的大多数依赖项。

通过揉和众多设计良好的 Nginx 模块，OpenResty 有效地把 Nginx 服务器转变为一个强大的 Web 应用服务器，基于它开发人员可以使用 Lua 编程语言对 Nginx 核心以及现有的各种 Nginx C 模块进行脚本编程，构建出可以处理一万以上并发请求的极端高性能的 Web 应用。

ngx_openresty 目前有两大应用目标：

1. 通用目的的 web 应用服务器。在这个目标下，现有的 web 应用技术都可以算是和 OpenResty 或多或少有些类似，比如 Nodejs, PHP 等等。ngx_openresty 的性能（包括内存使用和 CPU 效率）算是最大的卖点之一。
2. Nginx 的脚本扩展编程，用于构建灵活的 Web 应用网关和 Web 应用防火墙。有些类似的是 NetScaler。其优势在于 Lua 编程带来的巨大灵活性。

引用自：[OpenResty最佳实践](https://moonbingbing.gitbooks.io/openresty-best-practices/content/base/intro.html)

## 入门教程

* [OpenResty 最佳实践](https://moonbingbing.gitbooks.io/openresty-best-practices/content/)
* [跟我学OpenResty(Nginx+Lua)开发](https://www.iteye.com/blog/jinnianshilongnian-2190344)

目前已经有大佬写了很完整的教程，没有必要重复造轮子了，我觉得这两个是最好的。更多openresty相关资料可以看这里[https://blog.fengjx.com/awesome/#lua-openresty](https://blog.fengjx.com/awesome/#lua-openresty)

## 网关服务开发

### 整体架构

