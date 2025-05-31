# 🃏 Intelligent Blackjack Game

**🌐 Language / 语言选择**: [English](README.md) | [中文](README_CN.md)

A fully-featured blackjack game written in Go, integrated with **Monte Carlo probability analysis**, **Kelly Criterion bankroll management**, and **intelligent decision recommendation** systems.

## ✨ Core Features

### 🎯 **Intelligent Probability Analysis**
- **Real-time probability calculation**: Based on Monte Carlo simulation (10,000 trials) to calculate winning probabilities
- **Action win rate comparison**: Analyzes expected win rates for hit, stand, double down, and other actions
- **Optimal strategy recommendation**: Automatically recommends the best decision based on basic strategy

### 💰 **Kelly Criterion Bankroll Management**
- **Intelligent betting suggestions**: Provides scientific betting amount recommendations based on bankroll status
- **Risk assessment**: Real-time evaluation of current bankroll status and risk level
- **Double down decision analysis**: Evaluates expected ROI and risk-reward ratio for doubling down

### 🧠 **Decision Support System**
- **Basic strategy integration**: Built-in professional blackjack basic strategy
- **Real-time data analysis**: Displays key indicators like bust probability, 21-point probability, etc.
- **Entertainment cost estimation**: Helps players understand expected entertainment costs

## 🚀 Quick Start

### Requirements
- Go 1.22+ 
- UTF-8 compatible terminal

### Install Dependencies
```bash
go mod tidy
```

### Run the Game
```bash
# Method 1: Run directly
go run ./cmd

# Method 2: Build and run
go build -o blackjack ./cmd
./blackjack
```

## 🎮 Game Controls

### Basic Actions
- `h` / `hit` - Hit (take a card)
- `s` / `stand` - Stand
- `d` / `double` / `doubledown` - Double down
- `q` / `quit` - Quit game
- `y` / `yes` - Continue game
- `n` / `no` - End game

### Menu Options
- `1` - Start game
- `2` - View game rules
- `3` - Exit program

## 🃏 Game Rules

### 🎯 Game Objective
Get your hand as close to 21 as possible without going over, while beating the dealer's hand.

### 🎴 Card Values
- **Number cards (2-10)**: Face value
- **Face cards (J,Q,K)**: 10 points each
- **Ace**: Intelligently calculated as 1 or 11 (whichever is more favorable)

### 💰 Betting System
- **Starting chips**: 1000
- **Betting options**: 10, 25, 50, 100, 200 chips
- **Smart betting**: Automatic all-in option when chips are insufficient
- **Payout rules**:
  - 🏆 Regular win: 1:1
  - 🌟 Blackjack win: 3:2 (non-double situations)
  - 🤝 Push: Return original amount

### ⚡ Double Down Feature
- **Trigger condition**: Available on first two cards
- **Chip requirement**: Current chips ≥ current bet amount
- **Game rule**: Can only take one more card after doubling
- **Payout adjustment**: Blackjack pays 1:1 after doubling

### 🎲 Game Flow
1. **Betting phase** → Choose bet amount + bankroll management advice
2. **Dealing phase** → Player and dealer each get 2 cards (dealer has 1 hidden card)
3. **Player turn** → Choose hit/stand/double + probability analysis
4. **Dealer turn** → Dealer follows automatic rules
5. **Settlement phase** → Compare points and settle chips

## 📊 Intelligent Analysis System

### 🎯 **Real-time Probability Analysis**
The game displays detailed probability analysis each turn:

```
📊 Current Win Probability Analysis
────────────────────────────────────────
🟢 Player Win Probability: 45.6%
🔴 Dealer Win Probability: 16.7%  
🟡 Push Probability:       37.6%

📈 Detailed Analysis:
   💥 Player Bust Probability: 85.1%
   💥 Dealer Bust Probability: 22.7%
   🎯 Player 21 Probability: 8.5%
   🎯 Dealer 21 Probability: 6.3%

🎯 Action Win Rate Comparison:
   ✋ Stand: 45.8% ⭐ (Recommended)
   👆 Hit: 41.9%
   ⚡ Double: 36.3%
```

### 💰 **Kelly Criterion Bankroll Management**

#### **Betting Phase Advice**
```
💰 Bankroll Management Advice:
📊 Recommended Bet: 15 chips (1.5% of bankroll)
💡 Your bankroll status is good, moderate betting recommended
🟢 Risk Status: Sufficient funds, controllable risk
🎮 Expected Entertainment Cost: 0.06% per hand
```

#### **Double Down Decision Analysis**
```
💰 Kelly Criterion Double Analysis:
⚡ Recommend Double (Expected ROI: 18.0%)
🔴 Double Risk Level: High (Kelly fraction: 0.180)
```

### 🧮 **Tiered Bankroll Management Strategy**

| Bankroll Level | Recommended % | Risk Level | Strategy Description               |
| -------------- | ------------- | ---------- | ---------------------------------- |
| ≥ 1000 chips   | 1.5%          | 🟢 Low      | Sufficient funds, moderate betting |
| ≥ 500 chips    | 1.0%          | 🟢 Low      | Conservative betting, risk control |
| ≥ 200 chips    | 0.5%          | 🟡 Medium   | More conservative, minimum betting |
| < 200 chips    | Minimum       | 🔴 High     | Recommend caution or leave game    |

## 🧠 Probability Calculation Principles

### 🎲 **Monte Carlo Simulation Method**
The game uses Monte Carlo simulation to calculate various probabilities:

1. **Simulation count**: 10,000 simulations per analysis
2. **Strategy simulation**: Player uses basic strategy, dealer follows fixed rules
3. **Remaining deck**: Based on current known cards and remaining deck
4. **Statistical results**: Statistics on frequency of various outcomes

### 📈 **Advantages Over Direct Calculation**

| Aspect                      | Monte Carlo Simulation                 | Direct Calculation                     |
| --------------------------- | -------------------------------------- | -------------------------------------- |
| **Complexity Management**   | ✅ Handles complex strategies simply    | ❌ Combinatorial explosion              |
| **Strategy Integration**    | ✅ Naturally integrates basic strategy  | ❌ Difficult to handle strategy changes |
| **Extensibility**           | ✅ Easy to add new rules                | ❌ Need to rewrite calculation logic    |
| **Multi-variable Handling** | ✅ Naturally handles multiple factors   | ❌ Complexity explodes with dimensions  |
| **Dynamic Adaptation**      | ✅ Automatically adapts to deck changes | ❌ Need to re-derive formulas           |

### 💡 **Kelly Criterion Application**

#### **Basic Formula**
```
f* = (bp - q) / b
```
- `f*`: Optimal betting fraction
- `b`: Odds (1:1 or 3:2)
- `p`: Win probability
- `q`: Loss probability

#### **Practical Application Strategy**
1. **Conservative factor**: Use 25% of Kelly fraction to reduce risk
2. **Tiered management**: Different strategies based on bankroll status
3. **Entertainment-oriented**: Focus on bankroll management rather than strict expected returns

## 🏗️ Architecture Design

### 📁 Project Structure
```
go-blackjack/
├── cmd/                         # 🚀 Program entry
│   └── main.go                   
├── internal/                    # 🔒 Internal modules
│   ├── domain/                  # 🎯 Domain layer - Core business logic
│   │   └── entities/             
│   │       ├── game.go          # Game aggregate root
│   │       ├── player.go        # Player entity
│   │       ├── dealer.go        # Dealer entity
│   │       ├── card.go          # Card entity
│   │       ├── deck.go          # Deck entity
│   │       ├── hand.go          # Hand entity
│   │       └── types.go         # Type definitions
│   ├── application/             # 🔄 Application layer - Use case orchestration
│   │   ├── services/           
│   │   │   ├── game.go          # Game application service
│   │   │   └── probability.go   # Probability calculation service
│   │   └── dtos/               
│   │       └── game.go          # Data transfer objects
│   └── interfaces/              # 🖥️ Interface layer - User interaction
│       └── cli/                
│           ├── game.go          # CLI handler
│           └── display.go       # Display service
├── go.mod                       # 📦 Dependency management
├── go.sum                       
└── README.md                    # 📖 Project documentation
```

### 🎯 Architecture Layers

#### **Domain Layer**
- **Responsibility**: Core business logic and rules
- **Characteristics**: No external dependencies, pure business code
- **Contains**: Game entities, business rules, state management

```go
// Game aggregate root - Unified game state management
type Game struct {
    ID          string
    Player      *Player  
    Dealer      *Dealer
    Deck        *Deck
    State       GameState // State machine management
    RoundNumber int
}
```

#### **Application Layer**
- **Responsibility**: Use case orchestration, domain object coordination
- **Characteristics**: Handles business processes, data transformation
- **Contains**: Application services, probability calculation, DTO objects

```go
// Game application service - Orchestrates game use cases
type GameApplicationService struct {
    game            *entities.Game
    probabilityCalc *ProbabilityCalculator
}

// Probability calculator - Monte Carlo simulation
type ProbabilityCalculator struct {
    trials int // Number of simulations
    rng    *rand.Rand
}
```

#### **Interface Layer**
- **Responsibility**: User interaction, input/output handling
- **Characteristics**: Extensible for multiple interface implementations
- **Contains**: CLI handler, display service

```go
// CLI game handler
type GameHandler struct {
    gameService *services.GameApplicationService
    display     *DisplayService
}
```

### 🔄 Data Flow
```
User Input → CLI Handler → Application Service → Probability Calculator
                    ↓                              ↓
Domain Entities ← Application Service ← Monte Carlo Simulation
                    ↓
CLI Display ← Kelly Formula Analysis
```

## 🧪 Core Features

### 🎮 **State Machine Management**
```go
type GameState int
const (
    StateWaitingToBet GameState = iota
    StatePlayerTurn  
    StateDealerTurn
    StateGameOver
)
```

### 🃏 **Intelligent Card Value Calculation**
Automatically calculates optimal Ace value (1 or 11) to ensure best possible hand value.

### 📊 **Probability Analysis Engine**
```go
type ProbabilityResult struct {
    PlayerWinProbability  float64
    DealerWinProbability  float64
    PlayerBustProbability float64
    ActionAnalysis        *ActionAnalysis
}
```

### 💰 **Kelly Criterion Decision Making**
```go
type KellyRecommendation struct {
    RecommendedBetAmount   int
    RecommendedBetFraction float64
    ShouldDouble           bool
    RiskLevel              string
}
```

### 💡 **Business Rule Validation**
Strict business rule validation in the domain layer ensures game logic correctness.

### 🔧 **Extensible Design**
- **Multi-interface support**: Easy to add web, mobile interfaces
- **Feature extension**: Easy to add split, surrender features
- **Algorithm optimization**: Supports different probability calculation methods

## 🎯 Game Screenshots Examples

### Betting Phase - Bankroll Management Advice
```
💰 Current Chips: 1000
Please select bet amount:
1. 10 chips
2. 25 chips
3. 50 chips
4. 100 chips
5. 200 chips
6. Exit game

💰 Bankroll Management Advice:
📊 Recommended Bet: 15 chips (1.5% of bankroll)
💡 Your bankroll status is good, moderate betting recommended
🟢 Risk Status: Sufficient funds, controllable risk
🎮 Expected Entertainment Cost: 0.06% per hand
```

### Game Round - Intelligent Analysis
```
👨 Dealer Hand (first card hidden):
🂠 🃏A 

👨 Player Hand (Value: 9):
🃏5 🃏4 

────────────────────────────────────────
📊 Current Win Probability Analysis
────────────────────────────────────────
🟢 Player Win Probability: 59.1%
🔴 Dealer Win Probability: 34.0%
🟡 Push Probability:       6.8%

📈 Detailed Analysis:
   💥 Player Bust Probability: 0.0%
   💥 Dealer Bust Probability: 44.4%
   🎯 Player 21 Probability: 0.0%
   🎯 Dealer 21 Probability: 8.7%

🎯 Action Win Rate Comparison:
   ✋ Stand: 44.2%
   👆 Hit: 56.4%
   ⚡ Double: 59.0% ⭐ (Recommended)

🏆 Optimal Strategy Expected Win Rate: 59.0%

💰 Kelly Criterion Double Analysis:
   ⚡ Recommend Double (Expected ROI: 18.0%)
   🔴 Double Risk Level: High (Kelly fraction: 0.180)
────────────────────────────────────────
```

## 🎯 Future Plans

### Short-term Goals
- [ ] 🧪 Complete unit test coverage
- [ ] 📈 Add historical statistics feature
- [ ] 🎮 Implement split functionality
- [ ] 🔄 Add surrender option

### Medium-term Goals  
- [ ] 🌐 Develop web interface version
- [ ] 💾 Implement game data persistence
- [ ] 📱 Mobile adaptation
- [ ] 🤖 AI opponent mode

### Long-term Goals
- [ ] 🏆 Multiplayer functionality
- [ ] 📊 Detailed statistical reports
- [ ] 🎯 Custom rule settings
- [ ] 🌍 Internationalization support

## 🤝 Contributing

Issues and Pull Requests are welcome!

### Development Environment Setup
```bash
git clone https://github.com/yourusername/go-blackjack.git
cd go-blackjack
go mod tidy
go test ./...
```

### Code Style
- Follow Go official coding standards
- Use meaningful variable and function names
- Add appropriate comments and documentation

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

---

⭐ If this project helps you, please give us a star! 