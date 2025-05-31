package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// Game 游戏结构
type Game struct {
	Deck        *Deck
	Player      *Player // 玩家
	Dealer      *Dealer // 庄家
	Scanner     *bufio.Scanner
	RoundNumber int // 当前轮数
}

// NewGame 创建新游戏
func NewGame() *Game {
	return &Game{
		Deck:        NewDeck(),
		Player:      &Player{Name: "玩家"},
		Dealer:      &Dealer{},
		Scanner:     bufio.NewScanner(os.Stdin),
		RoundNumber: 0, // 初始轮数为0
	}
}

// clearScreen 清屏
func (g *Game) clearScreen() {
	fmt.Print("\033[2J\033[H")
}

// dealInitialCards 发初始牌
func (g *Game) dealInitialCards() {
	// 重置手牌
	g.Player.Hand = Hand{}
	g.Dealer.Hand = Hand{}

	// 发两张牌给玩家和庄家
	for range 2 {
		g.Player.Hand.AddCard(g.Deck.Deal())
		g.Dealer.Hand.AddCard(g.Deck.Deal())
	}
}

// playRound 游戏主循环
func (g *Game) playRound() {
	// 增加轮数
	g.RoundNumber++

	// 显示轮数
	fmt.Printf("=== 第%d轮开始 ===\n", g.RoundNumber)
	fmt.Println()

	// 检查牌堆是否足够
	if len(g.Deck.Cards) < 10 {
		fmt.Println("牌不够了，重新洗牌...")
		g.Deck = NewDeck()
		time.Sleep(1 * time.Second)
	}

	// 发初始牌
	g.dealInitialCards()

	// 玩家回合
	g.playerTurn()

	// 如果玩家没有爆牌且没有Blackjack，进行庄家回合
	if !g.Player.Hand.IsBust() && !g.Player.Hand.IsBlackjack() {
		g.dealerTurn()
	}

	// 判断结果
	g.determineWinner()
}

// playerTurn 玩家回合
func (g *Game) playerTurn() {
	fmt.Println("\n=== 玩家回合 ===")

	for {
		g.displayGame(true)

		// 检查是否为Blackjack
		if g.Player.Hand.IsBlackjack() {
			fmt.Println("🎉 恭喜！你拿到了Blackjack！")
			return
		}

		// 检查是否爆牌
		if g.Player.Hand.IsBust() {
			fmt.Println("💥 爆牌了！你输了！")
			return
		}

		// 获取玩家选择
		choice := g.getInput("请选择: (h)要牌 (s)停牌 (q)退出游戏: ")

		switch strings.ToLower(choice) {
		case "h", "hit":
			g.Player.Hand.AddCard(g.Deck.Deal())
			fmt.Printf("你拿到了: %s\n", g.Player.Hand.Cards[len(g.Player.Hand.Cards)-1])
			time.Sleep(1 * time.Second)
		case "s", "stand":
			fmt.Println("你选择停牌")
			time.Sleep(1 * time.Second)
			return
		case "q", "quit":
			fmt.Println("感谢游戏！再见！")
			os.Exit(0)
		default:
			fmt.Println("无效选择，请重新输入")
			time.Sleep(1 * time.Second)
		}
	}
}

// dealerTurn 庄家回合
func (g *Game) dealerTurn() {
	fmt.Println("\n=== 庄家回合 ===")

	for {
		g.displayGame(false)

		// 庄家规则：小于17必须要牌
		if g.Dealer.Hand.Value() < 17 {
			fmt.Println("庄家点数小于17，必须要牌...")
			time.Sleep(2 * time.Second)
			card := g.Deck.Deal()
			g.Dealer.Hand.AddCard(card)
			fmt.Printf("庄家拿到: %s\n", card)
			time.Sleep(1 * time.Second)
		} else {
			fmt.Println("庄家停牌")
			time.Sleep(1 * time.Second)
			break
		}

		// 检查庄家是否爆牌
		if g.Dealer.Hand.IsBust() {
			fmt.Println("💥 庄家爆牌！")
			time.Sleep(1 * time.Second)
			break
		}
	}
}

// determineWinner 判断游戏结果
func (g *Game) determineWinner() {
	g.displayGame(false)

	playerValue := g.Player.Hand.Value()
	dealerValue := g.Dealer.Hand.Value()
	playerBlackjack := g.Player.Hand.IsBlackjack()
	dealerBlackjack := g.Dealer.Hand.IsBlackjack()

	fmt.Printf("=== 第%d轮游戏结果 ===\n", g.RoundNumber)

	// 玩家爆牌
	if g.Player.Hand.IsBust() {
		fmt.Println("💀 你爆牌了，庄家获胜！")
		return
	}

	// 庄家爆牌
	if g.Dealer.Hand.IsBust() {
		fmt.Println("🎉 庄家爆牌，你获胜！")
		return
	}

	// 双方都有Blackjack
	if playerBlackjack && dealerBlackjack {
		fmt.Println("🤝 双方都是Blackjack，平局！")
		return
	}

	// 只有玩家有Blackjack
	if playerBlackjack {
		fmt.Println("🎉 你拿到Blackjack，获胜！(3:2赔率)")
		return
	}

	// 只有庄家有Blackjack
	if dealerBlackjack {
		fmt.Println("💀 庄家拿到Blackjack，你输了！")
		return
	}

	// 比较点数
	if playerValue > dealerValue {
		fmt.Println("🎉 你的点数更高，获胜！")
	} else if playerValue < dealerValue {
		fmt.Println("💀 庄家点数更高，你输了！")
	} else {
		fmt.Println("🤝 点数相同，平局！")
	}
}

// displayGame 显示游戏状态
func (g *Game) displayGame(hideDealer bool) {
	// 显示庄家手牌
	if hideDealer && len(g.Dealer.Hand.Cards) > 0 {
		fmt.Printf("庄家: %s [?] (点数: ?)\n", g.Dealer.Hand.Cards[0])
	} else {
		fmt.Printf("庄家: %s (点数: %d)\n", g.Dealer.Hand.String(), g.Dealer.Hand.Value())
	}

	// 显示玩家手牌
	fmt.Printf("%s: %s (点数: %d)\n", g.Player.Name, g.Player.Hand.String(), g.Player.Hand.Value())
	fmt.Println()
}
