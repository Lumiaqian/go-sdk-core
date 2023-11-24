# Go-SDK-Core

## 项目简介

Go-SDK-Core 是一个为 Go 语言开发的核心软件开发工具包（SDK）。它提供了一系列的核心功能和接口，以便于开发者在构建特定于不同第三方服务（如 WeChat、DingTalk 等）的 Go SDK 时能够重用这些功能。

## 主要功能

- `cache`：cache 接口用于缓存 SDK 中所需要的 token 等信息。
- `log`：定义了一个简单的日志接口，允许插入不同的日志实现。
- `rest`：设计了一个通用的 HTTP 客户端，支持基本的 HTTP 请求和响应处理，支持中间件。
- `token`：定义了一个用于获取和管理访问 token 的 `TokenProvider` 接口，以及从远程服务器获取 token 的 `TokenFetcher`。

## 安装

```bash
go get github.com/Lumiaqian/go-sdk-core
```

### 许可证

此项目根据 [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0) 许可证授权。
