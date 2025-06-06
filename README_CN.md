# 🃏 《二十一点》智能游戏 (Intelligent Blackjack Game)

**🌐 Language / 语言选择**: [English](README.md) | [中文](README_CN.md)

一个用Go语言编写的功能完整的二十一点游戏，集成了**蒙特卡洛概率分析**、**凯利公式资金管理**和**智能决策建议**系统。

## ✨ 核心特性

### 🎯 **智能概率分析**
- **实时概率计算**: 基于蒙特卡洛模拟(10,000次)计算获胜概率
- **操作胜率对比**: 分析要牌、停牌、加倍等操作的期望胜率
- **最优策略推荐**: 基于基本策略自动推荐最佳决策

### 💰 **凯利公式资金管理**
- **智能下注建议**: 基于资金状况提供科学的投注金额建议
- **风险评估**: 实时评估当前资金状况和风险等级
- **加倍决策分析**: 评估加倍的期望ROI和风险回报比

### 🧠 **决策支持系统**
- **基本策略集成**: 内置专业的blackjack基本策略
- **实时数据分析**: 显示爆牌概率、21点概率等关键指标
- **娱乐成本预估**: 帮助玩家了解预期的娱乐成本

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
1. **下注阶段** → 选择下注金额 + 资金管理建议
2. **发牌阶段** → 玩家庄家各2张牌(庄家1张暗牌)
3. **玩家回合** → 选择要牌/停牌/加倍 + 概率分析
4. **庄家回合** → 庄家按规则自动行动
5. **结算阶段** → 比较点数并结算筹码

## 📊 智能分析系统

### 🎯 **实时概率分析**
游戏会在每一回合显示详细的概率分析：

```
📊 当前获胜概率分析
────────────────────────────────────────
🟢 玩家获胜概率: 45.6%
🔴 庄家获胜概率: 16.7%
🟡 平局概率:     37.6%

📈 详细分析:
   💥 玩家爆牌概率: 85.1%
   💥 庄家爆牌概率: 22.7%
   🎯 玩家21点概率: 8.5%
   🎯 庄家21点概率: 6.3%

🎯 操作胜率对比:
   ✋ 停牌: 45.8% ⭐ (推荐)
   👆 要牌: 41.9%
   ⚡ 加倍: 36.3%
```

### 💰 **凯利公式资金管理**

#### **下注阶段建议**
```
💰 资金管理建议:
📊 建议下注: 15 筹码 (1.5% 资金)
💡 您的资金状况良好，可以适度下注
🟢 风险状况: 资金充足，风险可控
🎮 预期娱乐成本: 0.06% 每局
```

#### **加倍决策分析**
```
💰 凯利公式加倍分析:
⚡ 推荐加倍 (期望ROI: 18.0%)
🔴 加倍风险等级: High (凯利比例: 0.180)
```

### 🧮 **分层资金管理策略**

| 资金水平   | 建议比例 | 风险等级 | 策略描述             |
| ---------- | -------- | -------- | -------------------- |
| ≥ 1000筹码 | 1.5%     | 🟢 Low    | 资金充足，可适度下注 |
| ≥ 500筹码  | 1.0%     | 🟢 Low    | 保守下注，控制风险   |
| ≥ 200筹码  | 0.5%     | 🟡 Medium | 更加保守，最小下注   |
| < 200筹码  | 最小     | 🔴 High   | 建议谨慎或离开游戏   |

## 🧠 概率计算原理

### 🎲 **蒙特卡洛模拟方法**
游戏采用蒙特卡洛模拟来计算各种概率：

1. **模拟次数**: 每次分析进行10,000次模拟
2. **策略模拟**: 玩家使用基本策略，庄家按固定规则
3. **剩余牌堆**: 基于当前已知牌面和剩余牌堆
4. **统计结果**: 统计各种结果的出现频率

### 📈 **相比直接计算的优势**

| 方面           | 蒙特卡洛模拟       | 直接计算             |
| -------------- | ------------------ | -------------------- |
| **复杂度管理** | ✅ 处理复杂策略简单 | ❌ 组合爆炸问题       |
| **策略集成**   | ✅ 自然集成基本策略 | ❌ 难以处理策略变化   |
| **扩展性**     | ✅ 易于添加新规则   | ❌ 需重写计算逻辑     |
| **多变量处理** | ✅ 自然处理多个因素 | ❌ 维度增加复杂度激增 |
| **动态适应**   | ✅ 自动适应牌堆变化 | ❌ 需要重新推导公式   |

### 💡 **凯利公式应用**

#### **基本公式**
```
f* = (bp - q) / b
```
- `f*`: 最优投注比例
- `b`: 赔率 (1:1 或 3:2)
- `p`: 获胜概率
- `q`: 失败概率

#### **实际应用策略**
1. **保守系数**: 使用25%的凯利比例以降低风险
2. **分层管理**: 根据资金状况采用不同策略
3. **娱乐导向**: 重点是资金管理而非严格期望收益

## 🏗️ 架构设计

### 📁 项目结构
```
go-blackjack/
├── cmd/                         # 🚀 程序入口
│   └── main.go
├── internal/                    # 🔒 内部模块
│   ├── domain/                  # 🎯 领域层 - 核心业务逻辑
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
│   │   │   ├── game.go          # 游戏应用服务
│   │   │   └── probability.go   # 概率计算服务
│   │   └── dtos/
│   │       └── game.go          # 数据传输对象
│   └── interfaces/              # 🖥️ 接口层 - 用户交互
│       └── cli/
│           ├── game.go          # 命令行处理器
│           └── display.go       # 显示服务
├── go.mod                       # 📦 依赖管理
├── go.sum
└── README.md                    # 📖 项目文档
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
- **包含**: 应用服务、概率计算、DTO对象

```go
// 游戏应用服务 - 编排游戏用例
type GameApplicationService struct {
    game            *entities.Game
    probabilityCalc *ProbabilityCalculator
}

// 概率计算器 - 蒙特卡洛模拟
type ProbabilityCalculator struct {
    trials int // 模拟次数
    rng    *rand.Rand
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
用户输入 → CLI Handler → Application Service → Probability Calculator
                    ↓                              ↓
Domain Entities ← Application Service ← Monte Carlo Simulation
                    ↓
CLI Display ← Kelly Formula Analysis
```

## 🧪 核心特性

### 🎮 **状态机管理**
```go
type GameState int
const (
    StateWaitingToBet GameState = iota
    StatePlayerTurn
    StateDealerTurn
    StateGameOver
)
```

### 🃏 **智能牌值计算**
自动计算A牌的最优值(1或11)，确保玩家获得最佳点数。

### 📊 **概率分析引擎**
```go
type ProbabilityResult struct {
    PlayerWinProbability  float64
    DealerWinProbability  float64
    PlayerBustProbability float64
    ActionAnalysis        *ActionAnalysis
}
```

### 💰 **凯利公式决策**
```go
type KellyRecommendation struct {
    RecommendedBetAmount   int
    RecommendedBetFraction float64
    ShouldDouble           bool
    RiskLevel              string
}
```

### 💡 **业务规则验证**
在领域层进行严格的业务规则校验，确保游戏逻辑的正确性。

### 🔧 **可扩展设计**
- **多界面支持**: 易于添加Web、移动端界面
- **功能扩展**: 可轻松添加分牌、投降等功能
- **算法优化**: 支持不同的概率计算方法

## 🎯 游戏截图示例

### 下注阶段 - 资金管理建议
```
💰 当前筹码: 1000
请选择下注金额:
1. 10 筹码
2. 25 筹码
3. 50 筹码
4. 100 筹码
5. 200 筹码
6. 退出游戏

💰 资金管理建议:
📊 建议下注: 15 筹码 (1.5% 资金)
💡 您的资金状况良好，可以适度下注
🟢 风险状况: 资金充足，风险可控
🎮 预期娱乐成本: 0.06% 每局
```

### 游戏回合 - 智能分析
```
👨 庄家手牌 (第一张牌隐藏):
🂠 🃏A

👨 玩家手牌 (点数: 9):
🃏5 🃏4

────────────────────────────────────────
📊 当前获胜概率分析
────────────────────────────────────────
🟢 玩家获胜概率: 59.1%
🔴 庄家获胜概率: 34.0%
🟡 平局概率:     6.8%

📈 详细分析:
   💥 玩家爆牌概率: 0.0%
   💥 庄家爆牌概率: 44.4%
   🎯 玩家21点概率: 0.0%
   🎯 庄家21点概率: 8.7%

🎯 操作胜率对比:
   ✋ 停牌: 44.2%
   👆 要牌: 56.4%
   ⚡ 加倍: 59.0% ⭐ (推荐)

🏆 最优策略期望胜率: 59.0%

💰 凯利公式加倍分析:
   ⚡ 推荐加倍 (期望ROI: 18.0%)
   🔴 加倍风险等级: High (凯利比例: 0.180)
────────────────────────────────────────
```

## 🎯 未来计划

### 短期目标
- [ ] 🧪 完善单元测试覆盖
- [ ] 📈 添加历史统计功能
- [ ] 🎮 实现分牌(Split)功能
- [ ] 🔄 添加投降(Surrender)选项

### 中期目标
- [ ] 🌐 开发Web界面版本
- [ ] 💾 实现游戏数据持久化
- [ ] 📱 移动端适配
- [ ] 🤖 AI对手模式

### 长期目标
- [ ] 🏆 多人对战功能
- [ ] 📊 详细统计报表
- [ ] 🎯 自定义规则设置
- [ ] 🌍 国际化支持

## 🤝 贡献指南

欢迎提交Issues和Pull Requests！

### 开发环境设置
```bash
git clone https://github.com/yourusername/go-blackjack.git
cd go-blackjack
go mod tidy
go test ./...
```

### 代码风格
- 遵循Go官方代码规范
- 使用有意义的变量和函数名
- 添加适当的注释和文档

## 📄 许可证

本项目采用MIT许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

---

⭐ 如果这个项目对您有帮助，请给我们一个星标！
