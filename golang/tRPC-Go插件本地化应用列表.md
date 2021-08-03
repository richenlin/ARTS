# tRPC-Go插件本地化应用列表

| 组件                    | 官方组件                                                | 是否需要自研         |  自研项目 | 备注 |
| ----------------------- | ---------------------------- | -------------------------------- | ---- | ---- | 
mysql | [git.code.oa.com/trpc-go/trpc-database/mysql](https://git.code.oa.com/trpc-go/trpc-database/mysql) | 否 |  [git.code.oa.com/TGS/tgs-ext/tgs-mysql](https://git.code.oa.com/TGS/tgs-ext/tgs-mysql) | 封装了sqlbuilder
redis | [git.code.oa.com/trpc-go/trpc-database/redis](https://git.code.oa.com/trpc-go/trpc-database/redis) | 是 |[ git.code.oa.com/TGS/tgs-ext/tgs-redis](https://git.code.oa.com/TGS/tgs-ext/tgs-redis) |官方组件基于redigo，不支持sentinel
cli | [git.code.oa.com/trpc-go/trpc-cli](https://git.code.oa.com/trpc-go/trpc-cli) | 否 | [git.code.oa.com/TGS/tgs-cli](https://git.code.oa.com/TGS/tgs-cli) | 基于官方cli工具，简化封装；项目目录结构有区别
mongo | [git.code.oa.com/trpc-go/trpc-database/mongodb](https://git.code.oa.com/trpc-go/trpc-database/mongodb) | 否 |  | 可以直接使用，但是部分API不完整
cousul | [git.code.oa.com/trpc-go/trpc-naming-consul](https://git.code.oa.com/trpc-go/trpc-naming-consul) |  |  | 待验证
etcd | [git.code.oa.com/trpc-go/trpc-config-etcd](https://git.code.oa.com/trpc-go/trpc-config-etcd)  |  |  | 待验证
kafka | [git.code.oa.com/trpc-go/trpc-database/kafka](https://git.code.oa.com/trpc-go/trpc-database/kafka) |  |  |  待验证
rabbitmq | [git.code.oa.com/trpc-go/trpc-database/rabbitmq](https://git.code.oa.com/trpc-go/trpc-database/tree/master/rabbitmq) |  |  |  待验证
skywalking | [git.code.oa.com/trpc-go/trpc-opentracing-skywalking](https://git.code.oa.com/trpc-go/trpc-opentracing-skywalking) |  |  | 待验证
prometheus | [git.code.oa.com/trpc-go/trpc-metrics-prometheus](https://git.code.oa.com/trpc-go/trpc-metrics-prometheus) |  |  | 待验证
cos（文件上传协议）| [git.code.oa.com/trpc-go/trpc-database/cos](https://git.code.oa.com/trpc-go/trpc-database/cos) |  |  | 待验证
熔断限流 | [git.code.oa.com/cooperyan/trpc-filter/tree/master/hystrix](https://git.code.oa.com/cooperyan/trpc-filter/tree/master/hystrix) |  |  | 待验证，版本非stable
重试策略 | [iwiki.woa.com/pages/viewpage.action?pageId=429400811](https://iwiki.woa.com/pages/viewpage.action?pageId=429400811)  |  |  | 待验证
定时器 | [git.code.oa.com/trpc-go/trpc-database/tree/master/timer](https://git.code.oa.com/trpc-go/trpc-database/tree/master/timer) |  |  | 待验证
mock | [git.code.oa.com/cooperyan/trpc-filter/tree/master/mock](https://git.code.oa.com/cooperyan/trpc-filter/tree/master/mock)|  |  | 待验证