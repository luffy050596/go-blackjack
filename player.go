package main

// Player 玩家结构
type Player struct {
	Name  string
	Hand  Hand
	Bet   int // 当前下注金额
	Chips int // 玩家筹码总数
}

// Dealer 庄家结构
type Dealer struct {
	Hand Hand
}

// NewPlayer 创建新玩家
func NewPlayer(name string, initialChips int) *Player {
	return &Player{
		Name:  name,
		Chips: initialChips,
		Bet:   0,
	}
}

// CanBet 检查玩家是否有足够筹码下注
func (p *Player) CanBet(amount int) bool {
	return p.Chips >= amount && amount > 0
}

// PlaceBet 下注
func (p *Player) PlaceBet(amount int) bool {
	if !p.CanBet(amount) {
		return false
	}
	p.Bet = amount
	p.Chips -= amount
	return true
}

// WinBet 赢得下注（包括本金）
func (p *Player) WinBet(multiplier float64) {
	winnings := int(float64(p.Bet) * multiplier)
	p.Chips += p.Bet + winnings // 返还本金 + 奖金
	p.Bet = 0
}

// LoseBet 输掉下注
func (p *Player) LoseBet() {
	// 下注金额已经在PlaceBet时扣除，这里只需要清零当前下注
	p.Bet = 0
}

// PushBet 平局，返还下注
func (p *Player) PushBet() {
	p.Chips += p.Bet // 返还本金
	p.Bet = 0
}

// HasChips 检查玩家是否还有筹码
func (p *Player) HasChips() bool {
	return p.Chips > 0
}
