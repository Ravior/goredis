GoRedis 是一个用 Go 语言实现的 Redis 服务器。

关键功能:
- 支持 string, list, hash, set, sorted set, bitmap 数据结构
- 自动过期功能(TTL)
- 发布订阅
- 地理位置
- AOF 持久化及 AOF 重写
- 加载和导出 RDB 文件
- Multi 命令开启的事务具有`原子性`和`隔离性`. 若在执行过程中遇到错误, godis 会回滚已执行的命令
- 内置集群模式. 集群对客户端是透明的, 您可以像使用单机版 redis 一样使用 godis 集群
    - `MSET`, `MSETNX`, `DEL`, `Rename`, `RenameNX`  命令在集群模式下原子性执行, 允许 key 在集群的不同节点上
    - Multi 命令开启的事务在集群模式下支持在同一个 slot 内执行
- 并行引擎, 无需担心您的操作会阻塞整个服务器.