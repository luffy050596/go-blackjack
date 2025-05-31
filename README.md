# 二十一点游戏 (Blackjack Game)

一个用Go语言编写的完整二十一点游戏，支持终端命令行交互和完整的下注系统。

## 🚀 快速开始

### 运行游戏
```bash
go run ./cmd
```

### 游戏操作
- `h` 或 `hit` - 要牌
- `s` 或 `stand` - 停牌
- `d` 或 `double` 或 `doubledown` - 加倍(仅前两张牌时可用)  
- `q` 或 `quit` - 退出游戏
- `y` 或 `yes` - 再来一局
- `n` 或 `no` - 结束游戏

## 🃏 游戏规则

### 基本规则
- **目标**: 让手牌点数尽可能接近21点，但不能超过
- **牌面点数**:
  - 数字牌(2-10): 按牌面数字计算
  - 花牌(J,Q,K): 每张都是10点
  - A: 可以是1点或11点(自动选择最优)

### 💰 下注系统
- **初始筹码**: 1000
- **下注选项**: 10, 25, 50, 100, 200 筹码(从预设选项中选择)
- **智能选项**: 筹码不足时自动提供全押选项
- **赔率系统**:
  - 普通获胜: 1:1 赔率
  - Blackjack获胜: 3:2 赔率
  - 平局: 返还下注金额
- **筹码管理**: 筹码用完时可选择重新开始

### 游戏流程
1. **下注阶段**: 从预设选项中选择下注金额
2. 玩家和庄家各发2张牌
3. 庄家有一张暗牌(背面朝上)
4. 玩家选择要牌、停牌或加倍
5. 庄家小于17点必须要牌，17点以上必须停牌
6. 比较点数决定胜负并结算筹码

### ⚡ 加倍功能
- **使用条件**: 只能在拿到前两张牌时使用
- **筹码要求**: 需要有足够筹码将下注翻倍
- **游戏规则**: 加倍后只能再拿一张牌，然后必须停牌
- **赔率调整**: 加倍后的Blackjack按1:1赔率计算(非3:2)

### 特殊情况
- **Blackjack**: 前两张牌就是21点(A+10点牌)，非加倍时3:2赔率
- **爆牌**: 点数超过21点立即失败
- **平局**: 双方点数相同，返还下注金额

## 🏗️ 代码结构

### 文件组织

```
blackjack/
├── main.go      # 程序入口点
├── card.go      # 扑克牌相关定义
├── deck.go      # 牌堆管理
├── hand.go      # 手牌逻辑
├── player.go    # 玩家和庄家结构(含下注系统)
├── game.go      # 游戏主逻辑
├── menu.go      # 菜单系统
├── help.go      # 帮助和规则说明
├── strategy.go  # 策略相关功能
├── go.mod       # Go模块定义
└── README.md    # 项目说明
```

### 核心模块

#### 1. `card.go` - 扑克牌基础
```go
// 花色和牌面枚举
type Suit int    // ♥♦♣♠
type Rank int    // A,2-10,J,Q,K

// 扑克牌结构
type Card struct {
    Suit Suit
    Rank Rank
}
```

#### 2. `deck.go` - 牌堆管理
```go
type Deck struct {
    Cards []Card
}

// 主要功能
- NewDeck()     // 创建并洗牌
- Shuffle()     // Fisher-Yates洗牌算法
- Deal()        // 发牌
```

#### 3. `hand.go` - 手牌逻辑
```go
type Hand struct {
    Cards []Card
}

// 主要功能
- Value()       // 智能点数计算
- IsBlackjack() // Blackjack检测
- IsBust()      // 爆牌检测
- AddCard()     // 添加牌
```

#### 4. `player.go` - 参与者与下注系统
```go
type Player struct {
    Name  string
    Hand  Hand
    Bet   int // 当前下注金额
    Chips int // 玩家筹码总数
}

type Dealer struct {
    Hand Hand
}

// 下注系统功能
- NewPlayer()   // 创建玩家
- CanBet()      // 检查下注能力
- PlaceBet()    // 下注
- WinBet()      // 赢得下注
- LoseBet()     // 输掉下注
- PushBet()     // 平局返还
- HasChips()    // 检查筹码状态
```

#### 5. `game.go` - 游戏引擎
```go
type Game struct {
    Deck    *Deck
    Player  *Player
    Dealer  *Dealer
    Scanner *bufio.Scanner
}

// 主要功能
- placeBet()        // 下注阶段
- playerTurn()      // 玩家回合
- dealerTurn()      // 庄家回合
- determineWinner() // 胜负判断与筹码结算
- showMenu()        // 用户界面
- showChipsStatus() // 筹码状态显示
```

#### 6. `main.go` - 程序入口
```go
func main() {
    game := NewGame()
    game.showMenu()
}
```

## 📋 系统要求

- Go 1.23.0 或更高版本
- 支持ANSI转义序列的终端(用于清屏)

## 🎯 未来改进方向

- [x] 添加下注系统
- [ ] 支持分牌(Split)功能
- [x] 支持加倍(Double Down)功能
- [ ] 添加多玩家支持
- [ ] 实现基本策略提示
- [ ] 添加游戏统计功能
- [ ] 支持自定义牌堆数量
- [ ] 添加单元测试
- [ ] 实现配置文件支持

## 📝 许可证

MIT License

## 🤝 贡献

欢迎提交Issue和Pull Request来改进这个游戏！

---

**远离赌博，享受游戏！🎉** 