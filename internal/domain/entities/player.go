package entities

// Player 玩家结构
type Player struct {
	Name         string
	Hand         *Hand
	InitialChips int  // 初始筹码
	Chips        int  // 玩家筹码总数
	Bet          int  // 当前下注金额
	DoubledDown  bool // 是否已经加倍
}

// NewPlayer 创建新玩家
func NewPlayer(name string, initialChips int) *Player {
	return &Player{
		Name:         name,
		Hand:         NewHand(),
		InitialChips: initialChips,
		Chips:        initialChips,
		Bet:          0,
	}
}

// CanBet 检查玩家是否有足够筹码下注
func (p *Player) CanBet(amount int) bool {
	return amount > 0 && p.Chips >= amount
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

// DoubleBet 加倍下注
func (p *Player) DoubleBet() bool {
	if !p.CanBet(p.Bet) {
		return false
	}
	p.Chips -= p.Bet // 扣除额外的下注金额
	p.Bet *= 2       // 下注金额翻倍
	p.DoubledDown = true
	return true
}

// CanDoubleDown 检查是否可以加倍
func (p *Player) CanDoubleDown() bool {
	// 只有在前两张牌且有足够筹码时才能加倍
	return len(p.Hand.Cards) == 2 && !p.DoubledDown && p.CanBet(p.Bet)
}

// ResetRound 重置回合状态
func (p *Player) ResetRound() {
	p.Hand = NewHand()
	p.Bet = 0
	p.DoubledDown = false
}
