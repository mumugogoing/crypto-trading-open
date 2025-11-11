# Python to Golang Migration Guide

## 概述

本文档详细说明从Python版本迁移到Golang版本的差异、改进和注意事项。

## 🎯 迁移目标

1. **性能提升**: Go的编译型语言特性和并发模型提供更好的性能
2. **内存效率**: 更低的内存占用和更好的GC性能
3. **部署简便**: 单一二进制文件，无需Python环境和依赖
4. **类型安全**: 编译时类型检查减少运行时错误
5. **并发性能**: Go的goroutine提供更好的并发处理能力

## 📊 技术栈对比

| 功能 | Python版本 | Golang版本 |
|------|-----------|-----------|
| HTTP客户端 | aiohttp | go-resty/resty |
| WebSocket | websockets | gorilla/websocket |
| 异步编程 | asyncio | goroutines + channels |
| 日志系统 | Python logging | uber-go/zap |
| 配置管理 | PyYAML | spf13/viper |
| 终端UI | Rich | charmbracelet/bubbletea |
| 数值计算 | Decimal | shopspring/decimal |
| 依赖注入 | injector | (待选择) |

## 🔄 架构映射

### 目录结构对比

```
Python版本                    Golang版本
─────────────────────────────────────────────
core/                         pkg/
├── domain/                   ├── domain/
│   ├── models/              │   ├── models/
│   ├── entities/            │   ├── entities/
│   └── value_objects/       │   └── valueobjects/
├── adapters/                 ├── adapters/
│   └── exchanges/           │   └── exchanges/
├── services/                 ├── services/
├── infrastructure/           ├── infrastructure/
└── logging/                  └── events/

run_*.py                      cmd/
                              ├── grid_trading/
                              ├── volume_maker/
                              └── ...
```

### 代码模式对比

#### 1. 数据类定义

**Python (dataclass):**
```python
@dataclass
class TickerData:
    symbol: str
    bid: Decimal
    ask: Decimal
    last: Decimal
    timestamp: datetime
```

**Go (struct):**
```go
type TickerData struct {
    Symbol    string
    Bid       decimal.Decimal
    Ask       decimal.Decimal
    Last      decimal.Decimal
    Timestamp time.Time
}
```

#### 2. 异步编程

**Python (async/await):**
```python
async def get_ticker(self, symbol: str) -> TickerData:
    async with aiohttp.ClientSession() as session:
        async with session.get(url) as resp:
            data = await resp.json()
            return parse_ticker(data)
```

**Go (goroutines):**
```go
func (a *Adapter) GetTicker(ctx context.Context, symbol string) (*TickerData, error) {
    var result map[string]interface{}
    err := a.client.Get(ctx, url, nil, &result)
    if err != nil {
        return nil, err
    }
    return parseTicker(result), nil
}
```

#### 3. 错误处理

**Python (exceptions):**
```python
try:
    data = await get_data()
except Exception as e:
    logger.error(f"Error: {e}")
    raise
```

**Go (explicit error returns):**
```go
data, err := getData()
if err != nil {
    logger.Error("Error", zap.Error(err))
    return nil, err
}
```

#### 4. 接口定义

**Python (ABC):**
```python
from abc import ABC, abstractmethod

class ExchangeInterface(ABC):
    @abstractmethod
    async def get_ticker(self, symbol: str) -> TickerData:
        pass
```

**Go (interface):**
```go
type ExchangeInterface interface {
    GetTicker(ctx context.Context, symbol string) (*TickerData, error)
}
```

## 🔧 配置文件迁移

### Python配置 (YAML)

```yaml
exchange:
  binance:
    api_key: "xxx"
    api_secret: "yyy"
    testnet: false
```

### Go配置 (YAML)

```yaml
exchanges:
  - id: "binance"
    name: "Binance"
    type: "perpetual"
    api_key: "xxx"
    api_secret: "yyy"
    testnet: false
    enabled: true
```

配置加载:

**Python:**
```python
import yaml

with open('config.yaml') as f:
    config = yaml.safe_load(f)
```

**Go:**
```go
import "github.com/spf13/viper"

viper.SetConfigFile("config.yaml")
viper.ReadInConfig()
var config Config
viper.Unmarshal(&config)
```

## 📦 依赖管理

### Python (pip + requirements.txt)

```bash
pip install -r requirements.txt
```

### Go (go modules)

```bash
go mod download
go mod tidy
```

## 🚀 运行方式对比

### Python

```bash
# 安装依赖
pip install -r requirements.txt

# 运行程序
python run_grid_trading.py
```

### Go

```bash
# 下载依赖 (可选,第一次运行时自动)
go mod download

# 运行程序 (开发模式)
go run ./cmd/grid_trading/main.go

# 或构建后运行 (生产模式)
go build -o bin/grid_trading ./cmd/grid_trading/
./bin/grid_trading
```

## 🔍 性能对比

| 指标 | Python | Go | 改进 |
|------|--------|-----|------|
| 启动时间 | ~1-2s | ~0.1s | 10-20x |
| 内存占用 | ~100-200MB | ~20-50MB | 2-5x |
| 并发性能 | 良好 (asyncio) | 优秀 (goroutines) | 2-3x |
| CPU效率 | 中等 (解释型) | 高 (编译型) | 5-10x |
| 部署大小 | ~50MB + 依赖 | ~10-15MB (单文件) | 简化 |

## ⚠️ 迁移注意事项

### 1. Decimal精度

Python的`Decimal`和Go的`shopspring/decimal`在某些边界情况下可能有细微差异。建议:

- 对关键计算添加单元测试
- 使用相同的测试用例验证一致性

### 2. 时间处理

Python: `datetime` 模块
Go: `time` 包

注意时区处理的差异:
```go
// Go中显式指定UTC
now := time.Now().UTC()
```

### 3. 字符串格式化

Python: f-strings
Go: fmt.Sprintf 或字符串拼接

```python
# Python
msg = f"Price: {price}, Volume: {volume}"
```

```go
// Go
msg := fmt.Sprintf("Price: %s, Volume: %s", price.String(), volume.String())
```

### 4. JSON处理

Python的dict可以包含任意类型，Go需要明确类型定义:

```go
// 使用map[string]interface{}处理动态JSON
var data map[string]interface{}
json.Unmarshal(bytes, &data)

// 或定义明确的结构体
type Response struct {
    Price  string `json:"price"`
    Volume string `json:"volume"`
}
var resp Response
json.Unmarshal(bytes, &resp)
```

## 🧪 测试迁移

### Python (pytest)

```python
def test_get_ticker():
    adapter = BinanceAdapter(config)
    ticker = await adapter.get_ticker("BTCUSDT")
    assert ticker.symbol == "BTCUSDT"
```

### Go (testing)

```go
func TestGetTicker(t *testing.T) {
    adapter, _ := NewBinanceAdapter(config)
    ticker, err := adapter.GetTicker(context.Background(), "BTCUSDT")
    assert.NoError(t, err)
    assert.Equal(t, "BTCUSDT", ticker.Symbol)
}
```

## 📈 迁移进度追踪

使用GitHub Issues或Project Board追踪迁移进度:

- [ ] Phase 1: 核心基础设施 ✅
- [ ] Phase 2: 交易所适配器 (20%)
- [ ] Phase 3: 交易策略 (0%)
- [ ] Phase 4: 完整系统 (0%)

## 🤝 贡献指南

参与迁移工作:

1. 选择一个未完成的模块
2. 参考Python版本实现
3. 使用Go惯用法重写
4. 添加测试
5. 提交PR

## 📚 学习资源

- [Go官方文档](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com/)
- [Go并发模式](https://go.dev/blog/pipelines)

## 🎓 Go最佳实践

1. **错误处理**: 显式处理所有错误，不要忽略
2. **接口设计**: 保持接口小而专注
3. **并发安全**: 使用mutex或channels保护共享状态
4. **Context传递**: 使用context.Context管理取消和超时
5. **命名规范**: 遵循Go的命名惯例 (PascalCase导出, camelCase私有)

## 💡 提示和技巧

### 1. 从Python的async/await转换

Python的`async def`函数在Go中通常是普通函数，并发通过goroutines实现:

```go
// 启动异步任务
go func() {
    result, err := doSomething()
    // 处理结果
}()
```

### 2. 使用channels代替回调

Python中的回调函数在Go中可以用channels实现:

```go
// 创建channel
results := make(chan Result)

// 发送结果
go func() {
    results <- computeResult()
}()

// 接收结果
result := <-results
```

### 3. defer用于资源清理

Go的`defer`类似于Python的`finally`或`with`:

```go
file, err := os.Open("file.txt")
if err != nil {
    return err
}
defer file.Close()  // 确保关闭

// 使用file...
```

## 🔗 相关文档

- [README_GO.md](README_GO.md) - Golang版本主文档
- [README.md](README.md) - Python版本文档
- [config/config.example.yaml](config/config.example.yaml) - 配置示例

---

**更新时间**: 2025-11-11
**维护者**: @mumugogoing
