package entities

import (
	"errors"

	"github.com/google/uuid"
)

// GameState 游戏状态枚举
type GameState int

const (
	StateWaitingToBet GameState = iota
	StatePlayerTurn
	StateDealerTurn
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

	// 检查是否爆牌或需要转到庄家回合
	if g.Player.Hand.IsBust() || g.Player.Hand.IsBlackjack() {
		g.State = StateDealerTurn
	}

	return card, nil
}

// PlayerStand 玩家停牌
func (g *Game) PlayerStand() {
	if g.State == StatePlayerTurn {
		g.State = StateDealerTurn
	}
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

	// 加倍后只能拿一张牌
	g.State = StateDealerTurn
	return card, nil
}

// DealerTurn 庄家回合
func (g *Game) DealerTurn() error {
	if g.State != StateDealerTurn {
		return errors.New("not dealer's turn")
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
