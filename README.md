# Hyperliquid Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/funcblock-quant/hyperliquid-go-sdk.svg)](https://pkg.go.dev/github.com/funcblock-quant/hyperliquid-go-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/funcblock-quant/hyperliquid-go-sdk)](https://goreportcard.com/report/github.com/funcblock-quant/hyperliquid-go-sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Hyperliquid Go SDK 是一个功能完整的 Go 语言 SDK，为 [Hyperliquid](https://hyperliquid.xyz) 去中心化永续期货交易所提供全面的 API 支持。

## 📋 目录

- [特性](#特性)
- [安装](#安装)
- [快速开始](#快速开始)
- [功能说明](#功能说明)
- [API 参考](#api-参考)
- [示例](#示例)
- [环境配置](#环境配置)
- [测试](#测试)
- [贡献](#贡献)


## ✨ 特性

### 🔌 核心功能

- **多环境支持**: 支持主网和测试网
- **HTTP 客户端**: 高性能 HTTP 客户端，支持超时控制
- **WebSocket 客户端**: 实时数据订阅，支持自动重连
- **类型安全**: 完整的类型定义和结构体映射
- **错误处理**: 完善的错误处理和类型化错误响应

### 📊 市场数据

- **实时价格**: 获取所有交易对的中间价格
- **订单簿**: L2 深度数据查询
- **K线数据**: 支持多种时间间隔的历史 K 线数据
- **交易历史**: 用户交易记录和成交明细
- **资金费率**: 历史资金费率和用户资金费用记录

### 💼 交易功能

- **订单管理**: 限价单、市价单、触发单
- **批量操作**: 批量下单、取消、修改订单
- **杠杆设置**: 调整交易对杠杆倍数（全仓/逐仓）
- **保证金管理**: 逐仓保证金调整
- **减仓功能**: 安全的仓位管理

### 🔐 安全与签名

- **以太坊签名**: 完整的 EIP-712 签名支持
- **私钥管理**: 安全的私钥处理和签名生成
- **Nonce 管理**: 自动的 nonce 管理和时间戳同步

### 💰 资产管理

- **用户状态**: 查询账户余额、持仓和保证金信息
- **转账功能**: USD 转账和 Vault 操作
- **资金历史**: 充值、提现和资金费用记录
- **现货交易**: 支持现货市场操作

### 📈 高级功能

- **实时订阅**: WebSocket 实时数据订阅
- **推荐系统**: 推荐状态查询
- **质押功能**: 质押奖励和委托管理
- **子账户**: 多账户管理支持

## 🚀 安装

确保您的 Go 版本为 1.24.0 或更高版本：

```bash
go mod init your-project
go get github.com/funcblock-quant/hyperliquid-go-sdk
```

### 依赖要求

- Go 1.24.0+
- `github.com/ethereum/go-ethereum` - 以太坊相关功能
- `github.com/gorilla/websocket` - WebSocket 支持
- `github.com/vmihailenco/msgpack/v5` - 消息序列化

## 🚀 快速开始

### 基本市场数据查询

```go
package main

import (
    "fmt"
    "log"
    
    sdk "github.com/funcblock-quant/hyperliquid-go-sdk"
)

func main() {
    // 创建信息客户端
    info, err := sdk.NewInfo(sdk.MainnetAPIURL)
    if err != nil {
        log.Fatal(err)
    }
    
    // 获取所有交易对价格
    mids, err := info.AllMids()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("BTC 价格: %s\n", mids["BTC"])
    fmt.Printf("ETH 价格: %s\n", mids["ETH"])
}
```

### 交易示例

```go
package main

import (
    "log"
    
    "github.com/ethereum/go-ethereum/crypto"
    sdk "github.com/funcblock-quant/hyperliquid-go-sdk"
)

func main() {
    // 从私钥创建签名器
    privateKey, err := crypto.HexToECDSA("your_private_key_without_0x_prefix")
    if err != nil {
        log.Fatal(err)
    }
    signer, err := sdk.NewSignerFromPrivateKey(privateKey)
    if err != nil {
        log.Fatal(err)
    }
    
    // 获取元数据
    info, err := sdk.NewInfo(sdk.MainnetAPIURL)
    if err != nil {
        log.Fatal(err)
    }
    meta, err := info.Meta()
    if err != nil {
        log.Fatal(err)
    }
    
    // 创建交易客户端
    exchange := sdk.NewExchange(sdk.MainnetAPIURL, nil, meta, signer)
    
    // 下限价单
    orderReq := sdk.OrderRequest{
        Coin:    "BTC",
        IsBuy:   true,
        Size:    0.01,
        LimitPx: 40000.0,
        OrderType: sdk.OrderType{
            Limit: &sdk.LimitOrderType{Tif: sdk.TifGtc},
        },
    }
    
    result, err := exchange.Order(orderReq, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("订单结果: %+v", result)
}
```

### WebSocket 订阅

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    sdk "github.com/funcblock-quant/hyperliquid-go-sdk"
)

func main() {
    // 创建 WebSocket 客户端
    ws := sdk.NewWebsocketClient(sdk.MainnetAPIURL)
    
    // 连接
    ctx := context.Background()
    if err := ws.Connect(ctx); err != nil {
        log.Fatal(err)
    }
    defer ws.Close()
    
    // 订阅 BTC 交易数据
    sub := sdk.Subscription{
        Type: sdk.SubTypeTrades,
        Coin: "BTC",
    }
    
    _, err := ws.Subscribe(sub, func(data interface{}) {
        fmt.Printf("BTC 交易数据: %+v\n", data)
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // 保持连接
    select {}
}
```

## 📚 功能说明

### Info 客户端 - 市场数据和账户信息

```go
info, err := sdk.NewInfo(sdk.MainnetAPIURL)

// 市场数据
mids, err := info.AllMids()                    // 所有交易对价格
meta, err := info.Meta()                       // 交易对元数据
l2Book, err := info.L2Snapshot("BTC")          // 订单簿快照
candles, err := info.CandlesSnapshot("BTC", "1h", startTime, endTime)

// 用户数据
userState, err := info.UserState("0x...")     // 用户状态
openOrders, err := info.OpenOrders("0x...")   // 未成交订单
fills, err := info.UserFills("0x...")         // 成交记录
fundingHistory, err := info.UserFundingHistory("0x...", startTime, nil)
```

### Exchange 客户端 - 交易操作

```go
exchange := sdk.NewExchange(apiURL, vaultAddr, meta, signer)

// 订单操作
result, err := exchange.Order(orderReq, nil)           // 下单
cancelResult, err := exchange.Cancel(cancelReq)        // 取消订单
modifyResult, err := exchange.ModifyOrder(modifyReq)   // 修改订单

// 批量操作
results, err := exchange.BulkOrders(orders, nil)       // 批量下单
cancelResults, err := exchange.BulkCancel(cancelReqs)  // 批量取消

// 杠杆和保证金
err = exchange.UpdateLeverage("BTC", true, 10)         // 更新杠杆
err = exchange.UpdateIsolatedMargin("BTC", 1000.0)     // 调整逐仓保证金
```

### WebSocket 客户端 - 实时数据

```go
ws := sdk.NewWebsocketClient(sdk.MainnetAPIURL)

// 市场数据订阅
ws.Subscribe(sdk.Subscription{Type: sdk.SubTypeTrades, Coin: "BTC"}, callback)
ws.Subscribe(sdk.Subscription{Type: sdk.SubTypeL2Book, Coin: "ETH"}, callback)

// 用户数据订阅
ws.Subscribe(sdk.Subscription{Type: sdk.SubTypeUserFills, User: "0x..."}, callback)
ws.Subscribe(sdk.Subscription{Type: sdk.SubTypeOrderUpdates, User: "0x..."}, callback)
```

## 🔧 环境配置

### 环境变量

创建 `.env` 文件：

```bash
# 私钥 (不包含 0x 前缀)
HL_PRIVATE_KEY=your_private_key_here

# 可选: Vault 地址
HL_VAULT_ADDRESS=0x...
```

### 网络配置

```go
// 主网
info, err := sdk.NewInfo(sdk.MainnetAPIURL)

// 测试网
info, err := sdk.NewInfo(sdk.TestnetAPIURL)
```

## 🧪 测试

### 运行测试

```bash
# 加载环境变量
source .env

# 运行所有测试
go test ./examples/ -v

# 运行特定测试
go test ./examples/ -run TestGetAllTokens -v
```

### 测试示例

项目包含丰富的测试示例：

- `examples/all_tokens_test.go` - 基础信息查询
- `examples/order_test.go` - 订单操作
- `examples/trade_test.go` - 交易功能
- `examples/websocket_test.go` - WebSocket 订阅
- `examples/candles_test.go` - K线数据
- `examples/leverage_test.go` - 杠杆操作

## 📖 API 参考

### 常用订单类型

```go
// 限价单
orderType := sdk.OrderType{
    Limit: &sdk.LimitOrderType{Tif: sdk.TifGtc},
}

// 只做市单
orderType := sdk.OrderType{
    Limit: &sdk.LimitOrderType{Tif: sdk.TifAlo},
}

// 立即成交或取消
orderType := sdk.OrderType{
    Limit: &sdk.LimitOrderType{Tif: sdk.TifIoc},
}
```

### 错误处理

```go
result, err := exchange.Order(orderReq, nil)
if err != nil {
    if apiErr, ok := err.(*sdk.APIError); ok {
        log.Printf("API 错误: %s", apiErr.Message)
    } else {
        log.Printf("其他错误: %v", err)
    }
}
```

## 🤝 贡献

欢迎贡献代码！请按照以下步骤：

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 🔗 相关链接

- [Hyperliquid 官网](https://hyperliquid.xyz)
- [Hyperliquid API 文档](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api)
- [官方 Python SDK](https://github.com/hyperliquid-dex/hyperliquid-python-sdk)

## ⚠️ 免责声明

本 SDK 为非官方实现，使用前请充分测试。交易有风险，请谨慎使用。

---

**如有问题或建议，欢迎提交 Issue 或 Pull Request！**