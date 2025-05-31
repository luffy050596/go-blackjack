package main

// ResultType 游戏结果类型
type ResultType int

const (
	PlayerBust ResultType = iota
	DealerBust
	BothBlackjack
	PlayerBlackjack
	DealerBlackjack
	PlayerWin
	DealerWin
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
	ActionHit PlayerAction = iota
	ActionStand
	ActionDoubleDown
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

// 显示行动常量
const (
	DisplayActionStand  = "stand"
	DisplayActionDouble = "double"
	DisplayActionHit    = "hit"
	DisplayActionBust   = "bust"
)

// 菜单选项常量
const (
	MenuOptionStart = "1"
	MenuOptionRules = "2"
	MenuOptionExit  = "3"
)
