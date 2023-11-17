# go-sdk-core

### 主要功能

* cache：cache接口用于缓存SDK中所需要的token等信息
* log：定义了一个简单的日志接口，允许插入不同的日志实现
* rest：设计了一个通用的 HTTP 客户端，支持基本的 HTTP 请求和响应处理，支持中间件
* token：定义了一个用于获取和管理访问token的 `TokenProvider` 接口和 从远程服务器获取token的 `TokenFetcher`

### 许可证

此项目根据 [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0) 许可证授权。
