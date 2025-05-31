# 🃏 Go 黑杰克游戏 (Go Blackjack Game)

一个用Go语言编写的完整二十一点游戏，支持终端命令行交互。

## 🚀 快速开始

### 环境要求
- Go 1.22+ 
- 支持UTF-8的终端

### 安装依赖
```bash
go mod tidy
```

### 运行游戏
```bash
# 方式1: 直接运行
go run ./cmd

# 方式2: 编译后运行
go build -o blackjack ./cmd
./blackjack
```

## 🎮 游戏操作

### 基本操作
- `h` / `hit` - 要牌
- `s` / `stand` - 停牌  
- `d` / `double` / `doubledown` - 加倍下注
- `q` / `quit` - 退出游戏
- `y` / `yes` - 继续游戏
- `n` / `no` - 结束游戏

### 菜单选项
- `1` - 开始游戏
- `2` - 查看游戏规则
- `3` - 退出程序

## 🃏 游戏规则

### 🎯 游戏目标
让手牌点数尽可能接近21点，但不能超过，同时要比庄家的点数更高。

### 🎴 牌面点数
- **数字牌(2-10)**: 按牌面数字计算
- **花牌(J,Q,K)**: 每张都是10点  
- **A**: 智能计算为1点或11点(选择对玩家最有利的值)

### 💰 下注系统
- **初始筹码**: 1000
- **下注选项**: 10, 25, 50, 100, 200筹码
- **智能下注**: 筹码不足时自动提供全押选项
- **赔率规则**:
  - 🏆 普通获胜: 1:1
  - 🌟 Blackjack 获胜: 3:2 (非加倍情况)
  - 🤝 平局: 返还原金额

### ⚡ 加倍功能  
- **触发条件**: 前两张牌时可选择加倍
- **筹码要求**: 当前筹码≥当前下注金额
- **游戏规则**: 加倍后只能再拿一张牌
- **赔率调整**: 加倍后 Blackjack 按1:1计算

### 🎲 游戏流程
1. **下注阶段** → 选择下注金额
2. **发牌阶段** → 玩家庄家各2张牌(庄家1张暗牌)
3. **玩家回合** → 选择要牌/停牌/加倍
4. **庄家回合** → 庄家按规则自动行动
5. **结算阶段** → 比较点数并结算筹码

## 🏗️ 架构设计

### 📁 项目结构
```
go-blackjack/
├── cmd/                           # 🚀 程序入口
│   └── main.go                   
├── internal/                      # 🔒 内部模块
│   ├── domain/                   # 🎯 领域层 - 核心业务逻辑
│   │   └── entities/             
│   │       ├── game.go          # 游戏聚合根
│   │       ├── player.go        # 玩家实体  
│   │       ├── dealer.go        # 庄家实体
│   │       ├── card.go          # 卡牌实体
│   │       ├── deck.go          # 牌堆实体
│   │       ├── hand.go          # 手牌实体
│   │       └── types.go         # 类型定义
│   ├── application/             # 🔄 应用层 - 用例编排
│   │   ├── services/           
│   │   │   └── game_service.go  # 游戏应用服务
│   │   └── dtos/               
│   │       └── game_dtos.go     # 数据传输对象
│   └── interfaces/              # 🖥️ 接口层 - 用户交互
│       └── cli/                
│           └── game_handler.go  # 命令行处理器
├── go.mod                        # 📦 依赖管理
├── go.sum                       
└── README.md                     # 📖 项目文档
```

### 🎯 架构层次

#### **Domain Layer (领域层)**
- **职责**: 核心业务逻辑和规则
- **特点**: 无外部依赖，纯业务代码
- **包含**: 游戏实体、业务规则、状态管理

```go
// 游戏聚合根 - 统一管理游戏状态
type Game struct {
    ID          string
    Player      *Player  
    Dealer      *Dealer
    Deck        *Deck
    State       GameState // 状态机管理
    RoundNumber int
}
```

#### **Application Layer (应用层)**  
- **职责**: 用例编排，协调领域对象
- **特点**: 处理业务流程，数据转换
- **包含**: 应用服务、DTO对象

```go
// 游戏应用服务 - 编排游戏用例
type GameApplicationService struct {
    game *entities.Game
}
```

#### **Interface Layer (接口层)**
- **职责**: 用户交互，输入输出处理  
- **特点**: 可扩展多种界面实现
- **包含**: CLI处理器、显示服务

```go
// 命令行游戏处理器
type GameHandler struct {
    gameService *services.GameApplicationService
    display     *DisplayService
}
```

### 🔄 数据流向
```
用户输入 → CLI Handler → Application Service → Domain Entities → Application Service → CLI Display
```

## 🧪 核心特性

### 🎮 状态机管理
```go
type GameState int
const (
    StateWaitingToBet GameState = iota
    StatePlayerTurn  
    StateDealerTurn
    StateGameOver
)
```

### 🃏 智能牌值计算
自动计算A牌的最优值(1或11)，确保玩家获得最佳点数。

### 💡 业务规则验证
在领域层进行严格的业务规则校验，确保游戏逻辑的正确性。

### 🔧 可扩展设计
- **多界面支持**: 易于添加Web、移动端界面
- **功能扩展**: 可轻松添加分牌、投降等功能
- **数据持久化**: 支持添加数据库存储

## 🎯 未来计划

### 短期目标
- [ ] 🧪 完善单元测试覆盖
- [ ] 🌐 添加Web界面支持
- [ ] 💾 实现游戏数据持久化
- [ ] 📊 添加游戏统计功能

### 长期目标  
- [ ] 🃏 支持分牌(Split)功能
- [ ] 👥 多玩家游戏支持
- [ ] 🤖 AI策略提示系统
- [ ] 📱 移动端适配
- [ ] 🔧 配置文件支持

## 🤝 贡献指南

1. Fork 本项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启 Pull Request

### 贡献者
感谢所有为这个项目做出贡献的开发者！

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 🎉 致谢

感谢开源社区的支持，让我们能够构建更好的软件！

---

**🎮 享受游戏，远离赌博！Happy Coding! 🚀** 