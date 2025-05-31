package entities

// ResultType 游戏结果类型
type ResultType int

const (
	// PlayerBust represents the result when player busts
	PlayerBust ResultType = iota
	// DealerBust represents the result when dealer busts
	DealerBust
	// BothBlackjack represents the result when both have blackjack
	BothBlackjack
	// PlayerBlackjack represents the result when player has blackjack
	PlayerBlackjack
	// DealerBlackjack represents the result when dealer has blackjack
	DealerBlackjack
	// PlayerWin represents the result when player wins
	PlayerWin
	// DealerWin represents the result when dealer wins
	DealerWin
	// Push represents the result when it's a tie
	Push
)

// GameResult 游戏结果结构
type GameResult struct {
	ResultType ResultType
	BetAmount  int
	IsDoubled  bool
}

// PlayerAction 玩家行动类型
type PlayerAction int

const (
	// ActionInvalid represents an invalid action
	ActionInvalid PlayerAction = iota
	// ActionHit represents the hit action
	ActionHit
	// ActionStand represents the stand action
	ActionStand
	// ActionDoubleDown represents the double down action
	ActionDoubleDown
	// ActionQuit represents the quit action
	ActionQuit
)

// ActionResult 行动结果
type ActionResult struct {
	Action   PlayerAction
	IsValid  bool
	Continue bool
}

// 玩家输入常量
const (
	InputHit        = "h"
	InputHitFull    = "hit"
	InputStand      = "s"
	InputStandFull  = "stand"
	InputDouble     = "d"
	InputDoubleFull = "double"
	InputDoubleDown = "doubledown"
	InputQuit       = "q"
	InputQuitFull   = "quit"
	InputYes        = "y"
	InputYesFull    = "yes"
	InputNo         = "n"
	InputNoFull     = "no"
)
