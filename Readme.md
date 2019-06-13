# API网关
轻量级的网关，支持自有协议，自动注册，作为一个小的注册中心。
我们认为的一个小型的微服务环境：一个网关、一个cache、一个MQ
支持spring cloud的注册中心，完成服务发现，并自动注册进入网关中。
## 模型
* API -服务接口
  * url 访问的地址
  * serverUrl 对应真实服务的地址
  * Namespace 命名空间-模块，不同的namespace url可以重复
  * domain 所属的服务器域名，可以作为区分不同的环境或兼容以前的旧有的基于域名做区分的设计
  * cluster 集群

* Cluster -集群
集群是由服务器组成
  * Name 集群名称，唯一
  * servers 服务器列表

* Server -服务器
  * IP   地址
  * port 端口
  * tag 标签

* 特性
  * auth   鉴权
  * access 访问控制
  * loadbalance 负载均衡
  * 熔断
  * 限速
  
## RoadMap
### version 0.5.0
* 具有访问日志 --走新的日志（logrus）--ok
* 支持iphash 负载策略 --ok
* 支持基于tag切换的负载策略 --ok
* 支持access的黑白名单 --ok

### version 0.6.0
* 支持从配置文件中注册api、server接口 --ok
* 支持console指令:shutdown  --ok


### version 0.7.0
* 支持accesstoken方式的鉴权（存储走redis) --ok
* 支持权限校验和拦截 --ok

### version 0.8.0
* 支持权限配置 --ok
* 支持console指令:reload 重新装载 api定义 --doing

### version 0.9.0
* 支持静态资源
* 支持默认获取accesstoken的实现 --planing 
类似于微信 method：GET ，参数：grandtype ==access_token ,appid==分配的应用id,secret==秘钥

### version 1.0.0
* 支持reload方法来重新装载权限配置
* 完整的配置
* 支持模仿nginx错误返回，例如 404 500 等

### version 1.1.0
* 支持开关规则拦截请求，并返回对应的错误报文


### version 1.2.0
* 实现动态注册api、server接口、动态权限接口
* 支持网关集群
* 支持熔断

### version 1.5.0
* 支持模拟apache错误返回
* 支持lua扩展

### version 2.0.0
* 支持Eureka 注册中心 
* 支持 Consul 注册中心 
