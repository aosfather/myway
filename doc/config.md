#配置说明
## 系统配置
### 配置文件 config.yaml
* 格式
```$yaml
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
 * 格式说明
   * path 相关路径的配置
     * server 用于配置集群的定义
     * api 用于api接口的配置定义
     * access_log 访问日志的存放目录
     * log 日志的存放目录
     * user 鉴权的用户定义
     * role 角色定义
   * system 系统相关配置 
     * auth_redis 权限配置信息存储用的redis相关配置信息
       * address redis的IP地址
       * db 使用的数据库序号 1-32
       * expire token过期时间     