# API网关
轻量级的API网关。
具备接口转发、报文转换、安全防护、流量控制、熔断等特性
并支持流行的注册中心
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
### version 1.0.0
* 具有访问日志 --走新的日志（logrus）--ok
* 支持iphash 负载策略 --ok
* 支持基于tag切换的负载策略 --ok
* 支持access的黑白名单 --ok
* 支持从配置文件中注册api、server接口 --ok
* 支持console指令:shutdown  --ok
* 支持accesstoken方式的鉴权（存储走redis) --ok
* 支持权限校验和拦截 --ok
* 支持权限配置 --ok
* 支持console指令:reload 重新装载 api定义 --ok
* 支持reload方法来重新装载权限配置 --ok
* 完整的配置  --ok


### feature pool
* 支持分布式任务调度（workflow)用于组合api    --  developing
* 实现动态注册api、server接口、动态权限接口
* 支持熔断
* 支持开关规则拦截请求，并返回对应的错误报文
* 支持 lua扩展
* 支持lua脚本，可以自定义路由规则
* 支持自动生成trace id对系统无侵入
* 支持模仿nginx错误返回，例如 404 500 等
* 支持模拟apache错误返回 
* 支持网关集群(paxos或Raft--接入ETCD)
* 支持 ETCD 注册中心 
* 支持 prometheus
* 支持服务上线、验证 状态，支持服务平滑升级
