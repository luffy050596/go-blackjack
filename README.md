# 二十一点游戏 (Blackjack Game)

一个用Go语言编写的完整二十一点游戏，支持终端命令行交互。

## 🚀 快速开始

### 运行游戏
```bash
go run .
```

### 游戏操作
- `h` 或 `hit` - 要牌
- `s` 或 `stand` - 停牌  
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

### 游戏流程
1. 玩家和庄家各发2张牌
2. 庄家有一张暗牌(背面朝上)
3. 玩家选择要牌或停牌
4. 庄家小于17点必须要牌，17点以上必须停牌
5. 比较点数决定胜负

### 特殊情况
- **Blackjack**: 前两张牌就是21点(A+10点牌)，3:2赔率
- **爆牌**: 点数超过21点立即失败
- **平局**: 双方点数相同

## 🏗️ 代码结构

### 文件组织

```
blackjack/
├── main.go      # 程序入口点
├── card.go      # 扑克牌相关定义
├── deck.go      # 牌堆管理
├── hand.go      # 手牌逻辑
├── player.go    # 玩家和庄家结构
├── game.go      # 游戏主逻辑
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

#### 4. `player.go` - 参与者
```go
type Player struct {
    Name string
    Hand Hand
    Bet  int
}

type Dealer struct {
    Hand Hand
}
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
- playerTurn()      // 玩家回合
- dealerTurn()      // 庄家回合
- determineWinner() // 胜负判断
- showMenu()        // 用户界面
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

- [ ] 添加下注系统
- [ ] 支持分牌(Split)功能
- [ ] 支持加倍(Double Down)功能
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