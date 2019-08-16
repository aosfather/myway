#配置说明
## 系统配置
### 配置文件 config.yaml
* 格式
```yaml
path:
  server: config/clusters
  api: config/apis
  access_log:
  log:
  user: config/users
  role: config/roles

system:
  auth_redis:
        address: 127.0.0.1
        db: 1
        expire: 7200
```
 * ### 格式说明
   * path 相关路径的配置
     * server 用于配置集群的定义。一个集群一个文件，yaml格式。
     * api 用于api接口的配置定义。一个文件定义多个api，建议一个namespace一个文件。
     * access_log 访问日志的存放目录。
     * log 日志的存放目录。
     * user 鉴权的用户定义。一个用户一个文件，yaml格式。
     * role 角色定义。一个角色一个文件，yaml格式。
   * system 系统相关配置 。
     * auth_redis 权限配置信息存储用的redis相关配置信息。
       * address redis的IP地址。
       * db 使用的数据库序号 1-32。
       * expire token过期时间。
       
 * server 配置文件格式说明
 ```yaml
 id: test          #集群id
 name: "测试服务器"  #集群名称
 balance: 2        #负载策略
 balance_config: "test" #负载的配置
 servers:   #服务器列表
   - {id: 100 ,ip: "127.0.0.1",port: 8990,tag: "test,dev"}
   # id 服务器的标识 ip 服务器地址 port 服务端口 tag 服务标签
```

* api 配置文件格式说明
```yaml
namespace: test   #命名空间
apis:             # api列表
  - {url: "bb",max_qps: 2,cluster: "test",server_url: "bb",auth: true}
  - {url: "cc",max_qps: 2,cluster: "-",server_url: "bb"}
  - {url: "dd",max_qps: 5,cluster: ".",server_url: "tokens"}
  # url api的url max_qps 允许的最大的QPS cluster 对应的集群 server_url 对应的服务的url地址 auth 是否启用鉴权
```

* user 配置文件格式说明
```yaml
name: "test1"   #用户名
pwd:  "xxxxxx"  #用户密码
role: test1     #对应的角色，只能有一个角色
```

* role 配置文件格式说明
```yaml
name: test1  #角色名称
super:       #父级角色名称，会继承父级角色的权限
white:       # 白名单，允许访问的url列表
   - {url: "/aa/aa/a"}
black:      # 黑名单，不允许访问的url列表
   -
```         