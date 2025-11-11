# 多交易所策略自动化系统 - Golang版本

## 🚀 项目状态

这是原Python系统向Golang的完全重写。目前处于**活跃开发中**。

### 完成进度

- ✅ **Phase 1: 核心基础设施** (100%)
  - Go项目结构和依赖管理
  - 领域模型和实体
  - 交易所接口定义
  - 事件系统
  - 日志系统 (uber-go/zap)
  - 配置管理 (spf13/viper)
  - HTTP客户端工具
  - WebSocket客户端工具

- 🔄 **Phase 2: 交易所适配器** (20%)
  - ✅ Binance适配器 (REST API完成, WebSocket待实现)
    - 市场数据 (ticker, orderbook, trades, OHLCV)
    - 交易操作 (创建/取消订单)
    - 账户数据 (余额, 持仓)
    - 持仓管理 (杠杆, 保证金模式)
  - ⏳ Hyperliquid适配器
  - ⏳ 其他交易所 (Backpack, Lighter, OKX, EdgeX)

- ⏳ **Phase 3: 交易策略** (0%)
  - 网格交易系统
  - 刷量交易
  - 套利监控
  - 价格提醒

- ⏳ **Phase 4: 完整系统** (0%)
  - 终端UI (charmbracelet/bubbletea)
  - 主应用程序
  - 文档
  - 测试

## 🏗️ 项目结构

```
crypto-trading-open/
├── cmd/                          # 可执行程序入口
│   ├── test_binance/            # Binance测试程序
│   ├── grid_trading/            # 网格交易 (待实现)
│   ├── volume_maker/            # 刷量交易 (待实现)
│   ├── arbitrage_monitor/       # 套利监控 (待实现)
│   └── price_alert/             # 价格提醒 (待实现)
├── pkg/                          # 公共库
│   ├── domain/                   # 领域层
│   │   ├── models/              # 领域模型
│   │   ├── entities/            # 实体
│   │   └── valueobjects/        # 值对象
│   ├── adapters/                 # 适配器层
│   │   └── exchanges/           # 交易所适配器
│   │       ├── binance/         # Binance适配器
│   │       ├── hyperliquid/     # Hyperliquid适配器 (待实现)
│   │       └── ...
│   ├── services/                 # 服务层 (待实现)
│   │   ├── grid/                # 网格交易服务
│   │   ├── volume_maker/        # 刷量服务
│   │   ├── arbitrage/           # 套利服务
│   │   └── price_alert/         # 价格提醒服务
│   ├── infrastructure/          # 基础设施层
│   │   ├── config/              # 配置管理
│   │   ├── logging/             # 日志系统
│   │   ├── http/                # HTTP客户端
│   │   ├── websocket/           # WebSocket客户端
│   │   └── di/                  # 依赖注入 (待实现)
│   └── events/                   # 事件系统
├── internal/                     # 内部包
│   └── terminal/                # 终端UI (待实现)
├── config/                       # 配置文件
└── go.mod                        # Go模块定义
```

## 📦 技术栈

### 核心依赖

| 库 | 版本 | 用途 |
|---|---|---|
| [shopspring/decimal](https://github.com/shopspring/decimal) | v1.4.0 | 精确的十进制数学运算 |
| [gorilla/websocket](https://github.com/gorilla/websocket) | v1.5.3 | WebSocket客户端 |
| [go-resty/resty](https://github.com/go-resty/resty) | v2.16.5 | HTTP客户端 |
| [uber-go/zap](https://github.com/uber-go/zap) | v1.27.0 | 高性能结构化日志 |
| [spf13/viper](https://github.com/spf13/viper) | v1.21.0 | 配置管理 |
| [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) | v1.2.4 | 终端UI框架 |

## 🚀 快速开始

### 前置要求

- Go 1.24+ 
- 交易所API密钥 (Binance, Hyperliquid等)

### 安装

```bash
# 克隆仓库
git clone https://github.com/mumugogoing/crypto-trading-open.git
cd crypto-trading-open

# 下载依赖
go mod download

# 构建测试程序
go build -o bin/test_binance ./cmd/test_binance/
```

### 配置

复制配置模板并填写API密钥:

```bash
cp config/config.example.yaml config/config.yaml
# 编辑 config/config.yaml 填写您的API密钥
```

配置示例:

```yaml
app:
  name: "Crypto Trading System"
  version: "2.0.0-go"
  environment: "production"

exchanges:
  - id: "binance"
    name: "Binance"
    type: "perpetual"
    api_key: "YOUR_API_KEY"
    api_secret: "YOUR_API_SECRET"
    testnet: false
    enabled: true

trading:
  default_leverage: 1
  default_margin_mode: "cross"
  max_position_size: 10000
  risk_per_trade: 0.02

logging:
  level: "info"
  output_path: "stdout"
  colored: true
```

### 运行测试

```bash
# 测试Binance适配器
./bin/test_binance

# 构建其他程序 (待实现)
# go build -o bin/grid_trading ./cmd/grid_trading/
# go build -o bin/volume_maker ./cmd/volume_maker/
```

## 🔧 开发

### 构建

```bash
# 构建所有程序
make build

# 构建特定程序
go build -o bin/test_binance ./cmd/test_binance/
```

### 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./pkg/adapters/exchanges/binance/...
```

### 代码格式化

```bash
# 格式化代码
go fmt ./...

# 静态检查
go vet ./...

# 使用golangci-lint (推荐)
golangci-lint run
```

## 📚 API文档

### Binance适配器

```go
import (
    "context"
    "github.com/mumugogoing/crypto-trading-open/pkg/adapters/exchanges"
    "github.com/mumugogoing/crypto-trading-open/pkg/adapters/exchanges/binance"
)

// 创建适配器
config := exchanges.NewExchangeConfig(
    "binance", "Binance", exchanges.ExchangeTypePerpetual,
    "API_KEY", "API_SECRET",
)
adapter, _ := binance.NewBinanceAdapter(config)

// 连接
ctx := context.Background()
adapter.Connect(ctx)

// 获取行情
ticker, _ := adapter.GetTicker(ctx, "BTCUSDT")
fmt.Printf("BTC价格: %s\n", ticker.Last.String())

// 获取订单簿
orderbook, _ := adapter.GetOrderBook(ctx, "BTCUSDT", 10)
fmt.Printf("最佳买价: %s\n", orderbook.Bids[0][0].String())

// 创建订单
order, _ := adapter.CreateOrder(
    ctx, "BTCUSDT",
    exchanges.OrderSideBuy,
    exchanges.OrderTypeLimit,
    decimal.NewFromFloat(0.001),
    decimal.NewFromFloat(50000),
    nil,
)
```

## 🗺️ 路线图

### 短期目标 (1-2周)

- [ ] 完成Binance WebSocket实现
- [ ] 实现Hyperliquid适配器
- [ ] 开始网格交易系统实现
- [ ] 添加单元测试

### 中期目标 (1-2月)

- [ ] 完成所有交易所适配器
- [ ] 实现所有交易策略
- [ ] 添加终端UI
- [ ] 完善文档和示例

### 长期目标 (3-6月)

- [ ] 性能优化
- [ ] 添加回测系统
- [ ] Web管理界面
- [ ] 移动端支持

## 🤝 贡献

欢迎贡献！请遵循以下步骤:

1. Fork本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 详见 LICENSE 文件

## 🙏 致谢

- 原Python版本的所有贡献者
- Go社区的优秀开源项目
- 各交易所的API文档和支持

## 📞 联系方式

- GitHub Issues: [问题追踪](https://github.com/mumugogoing/crypto-trading-open/issues)
- 项目主页: [crypto-trading-open](https://github.com/mumugogoing/crypto-trading-open)

---

**注意**: 本项目仅供学习和研究使用。加密货币交易存在风险，请谨慎使用并自行承担风险。
