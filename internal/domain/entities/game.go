package entities

import (
	"errors"

	"github.com/google/uuid"
)

// GameState 游戏状态枚举
type GameState int

const (
	// StateWaitingToBet represents the state when waiting for a bet
	StateWaitingToBet GameState = iota
	// StatePlayerTurn represents the state when it's the player's turn
	StatePlayerTurn
	// StateDealerTurn represents the state when it's the dealer's turn
	StateDealerTurn
	// StateGameOver represents the state when the game is over
	StateGameOver
)

// Game 游戏聚合根
type Game struct {
	ID          string
	Player      *Player
	Dealer      *Dealer
	Deck        *Deck
	State       GameState
	RoundNumber int
	IsActive    bool
}

// NewGame 创建新游戏
func NewGame(playerName string) *Game {
	return &Game{
		ID:          generateGameID(),
		Player:      NewPlayer(playerName, 1000),
		Dealer:      NewDealer(),
		Deck:        NewDeck(),
		State:       StateWaitingToBet,
		RoundNumber: 0,
		IsActive:    true,
	}
}

// StartNewRound 开始新一轮游戏
func (g *Game) StartNewRound() error {
	if g.State != StateWaitingToBet {
		return errors.New("cannot start new round in current state")
	}

	g.RoundNumber++
	g.Player.ResetRound()
	g.Dealer.ResetRound()
	g.ensureDeckSize()

	return nil
}

// PlaceBet 下注
func (g *Game) PlaceBet(amount int) error {
	if g.State != StateWaitingToBet {
		return errors.New("cannot place bet in current state")
	}

	if !g.Player.PlaceBet(amount) {
		return errors.New("cannot place bet")
	}

	g.State = StatePlayerTurn
	return nil
}

// DealInitialCards 发初始牌
func (g *Game) DealInitialCards() error {
	if g.State != StatePlayerTurn {
		return errors.New("cannot deal cards in current state")
	}

	// 发两张牌给玩家和庄家
	for range 2 {
		card, err := g.Deck.Deal()
		if err != nil {
			return err
		}
		g.Player.Hand.AddCard(card)

		card, err = g.Deck.Deal()
		if err != nil {
			return err
		}
		g.Dealer.Hand.AddCard(card)
	}

	return nil
}

// PlayerHit 玩家要牌
func (g *Game) PlayerHit() (Card, error) {
	if g.State != StatePlayerTurn {
		return Card{}, errors.New("not player's turn")
	}

	card, err := g.Deck.Deal()
	if err != nil {
		return Card{}, err
	}

	g.Player.Hand.AddCard(card)

	return card, nil
}

// PlayerStand 玩家停牌
func (g *Game) PlayerStand() {
}

// PlayerDoubleDown 玩家加倍
func (g *Game) PlayerDoubleDown() (Card, error) {
	if g.State != StatePlayerTurn {
		return Card{}, errors.New("not player's turn")
	}

	if !g.Player.CanDoubleDown() {
		return Card{}, errors.New("cannot double down")
	}

	if !g.Player.DoubleBet() {
		return Card{}, errors.New("cannot double down")
	}

	card, err := g.PlayerHit()
	if err != nil {
		return Card{}, err
	}

	return card, nil
}

// DealerTurn 庄家回合
func (g *Game) DealerTurn() error {
	if g.State != StateDealerTurn {
		return errors.New("not dealer's turn")
	}

	// 如果玩家爆牌或任一方有Blackjack，庄家不需要额外要牌
	if g.Player.Hand.IsBust() || g.Player.Hand.IsBlackjack() || g.Dealer.Hand.IsBlackjack() {
		g.State = StateGameOver
		return nil
	}

	// 庄家按规则要牌
	for g.Dealer.Hand.Value() < 17 {
		card, err := g.Deck.Deal()
		if err != nil {
			return err
		}
		g.Dealer.Hand.AddCard(card)
	}

	g.State = StateGameOver
	return nil
}

// EvaluateResult 评估游戏结果
func (g *Game) EvaluateResult() *GameResult {
	if g.State != StateGameOver {
		return nil
	}

	result := &GameResult{
		BetAmount: g.Player.Bet,
		IsDoubled: g.Player.DoubledDown,
	}

	playerValue := g.Player.Hand.Value()
	dealerValue := g.Dealer.Hand.Value()
	playerBlackjack := g.Player.Hand.IsBlackjack()
	dealerBlackjack := g.Dealer.Hand.IsBlackjack()

	// 评估逻辑
	switch {
	case g.Player.Hand.IsBust():
		result.ResultType = PlayerBust
		g.Player.LoseBet()
	case g.Dealer.Hand.IsBust():
		result.ResultType = DealerBust
		g.Player.WinBet(1.0)
	case playerBlackjack && dealerBlackjack:
		result.ResultType = Push
		g.Player.PushBet()
	case playerBlackjack:
		result.ResultType = PlayerBlackjack
		if g.Player.DoubledDown {
			g.Player.WinBet(1.0)
		} else {
			g.Player.WinBet(1.5)
		}
	case dealerBlackjack:
		result.ResultType = DealerBlackjack
		g.Player.LoseBet()
	case playerValue > dealerValue:
		result.ResultType = PlayerWin
		g.Player.WinBet(1.0)
	case playerValue < dealerValue:
		result.ResultType = DealerWin
		g.Player.LoseBet()
	default:
		result.ResultType = Push
		g.Player.PushBet()
	}

	g.State = StateWaitingToBet
	return result
}

// IsGameOver 检查游戏是否结束
func (g *Game) IsGameOver() bool {
	return !g.Player.HasChips()
}

// ensureDeckSize 确保牌堆足够
func (g *Game) ensureDeckSize() {
	if len(g.Deck.Cards) < 10 {
		g.Deck = NewDeck()
	}
}

// generateGameID 生成游戏ID
func generateGameID() string {
	return uuid.New().String()
}

// GetRemainingCards 获取剩余卡牌（用于概率计算）
func (g *Game) GetRemainingCards() []Card {
	return g.Deck.Cards
}

// GetUsedCards 获取已使用的卡牌（玩家和庄家手牌）
func (g *Game) GetUsedCards() []Card {
	used := make([]Card, 0)
	used = append(used, g.Player.Hand.Cards...)
	used = append(used, g.Dealer.Hand.Cards...)
	return used
}
