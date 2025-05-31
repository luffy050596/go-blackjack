package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	InitialChips = 1000
)

// Game 游戏结构
type Game struct {
	Deck        *Deck
	Player      *Player
	Dealer      *Dealer
	Scanner     *bufio.Scanner
	RoundNumber int
	service     *GameService
	display     *Display
}

// NewGame 创建新游戏
func NewGame() *Game {
	game := &Game{
		Deck:        NewDeck(),
		Player:      NewPlayer("玩家", InitialChips),
		Dealer:      &Dealer{},
		Scanner:     bufio.NewScanner(os.Stdin),
		RoundNumber: 0,
	}

	game.service = NewGameService(game)
	game.display = NewDisplay(game)

	return game
}

// clearScreen 清屏
func (g *Game) clearScreen() {
	g.display.clearScreen()
}

// dealInitialCards 发初始牌
func (g *Game) dealInitialCards() {
	for range 2 {
		g.Player.Hand.AddCard(g.Deck.Deal())
		g.Dealer.Hand.AddCard(g.Deck.Deal())
	}
}

// playRound 游戏主循环
func (g *Game) playRound() {
	g.RoundNumber++
	g.display.showRoundStart(g.RoundNumber, g.Player.Chips)

	if !g.Player.HasChips() {
		g.display.showInfo("你的筹码用完了！游戏结束！")
		return
	}

	g.Player.ResetRound()

	if !g.placeBet() {
		return
	}

	g.ensureDeckSize()
	g.dealInitialCards()
	g.playerTurn()

	if !g.Player.Hand.IsBust() && !g.Player.Hand.IsBlackjack() {
		g.dealerTurn()
	}

	result := g.service.EvaluateGame()
	g.display.showGameState(false)
	g.display.showRoundResult(result)
}

// ensureDeckSize 确保牌堆足够
func (g *Game) ensureDeckSize() {
	if len(g.Deck.Cards) < 10 {
		g.display.showInfo("牌不够了，重新洗牌...")
		g.Deck = NewDeck()
		time.Sleep(1 * time.Second)
	}
}

// placeBet 下注阶段
func (g *Game) placeBet() bool {
	g.display.showBetSection(g.Player.Chips)
	validOptions := g.service.PrepareBetting()
	g.display.showBetOptions(validOptions, g.Player.Chips)

	for {
		input := g.getInput("请选择 (输入选项编号): ")
		var choice int

		if _, err := fmt.Sscanf(input, "%d", &choice); err != nil {
			g.display.showError("请输入有效的选项编号")
			continue
		}

		success, continueGame := g.service.ProcessBet(choice, validOptions)
		if !continueGame {
			return false
		}
		if success {
			return true
		}

		g.display.showError("请重试")
	}
}

// playerTurn 玩家回合
func (g *Game) playerTurn() {
	g.display.showPlayerTurnStart()

	for {
		g.display.showGameState(true)

		if g.Player.Hand.IsBlackjack() {
			g.display.showBlackjack()
			return
		}

		if g.Player.Hand.IsBust() {
			g.display.showPlayerBust()
			return
		}

		// 如果已经加倍，只能拿一张牌后必须停牌
		if g.Player.DoubledDown {
			g.display.showPlayerChoice(DisplayActionDouble)
			card := g.Deck.Deal()
			g.Player.Hand.AddCard(card)
			g.display.showCardDealt(card, true)
			return
		}

		prompt := g.service.GetPromptText()
		choice := g.getInput(prompt)

		action := g.service.ParsePlayerInput(choice)
		result := g.service.ProcessPlayerAction(action)

		if action == ActionQuit {
			os.Exit(0)
		}

		if !result.Continue {
			return
		}
	}
}

// dealerTurn 庄家回合
func (g *Game) dealerTurn() {
	g.display.showDealerTurnStart()

	for {
		g.display.showGameState(false)

		if g.service.ShouldDealerHit() {
			g.service.DealerHit()
		} else {
			g.service.DealerStand()
			break
		}

		if g.service.CheckDealerBust() {
			break
		}
	}
}

// getInput 获取用户输入
func (g *Game) getInput(prompt string) string {
	fmt.Print(prompt)
	g.Scanner.Scan()
	return strings.TrimSpace(g.Scanner.Text())
}
