package main

// Player 玩家结构
type Player struct {
	Name string
	Hand Hand
	Bet  int
}

// Dealer 庄家结构
type Dealer struct {
	Hand Hand
}
