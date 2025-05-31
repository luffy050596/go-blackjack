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
	Player      *Player // 玩家
	Dealer      *Dealer // 庄家
	Scanner     *bufio.Scanner
	RoundNumber int // 当前轮数
}

// NewGame 创建新游戏
func NewGame() *Game {
	return &Game{
		Deck:        NewDeck(),
		Player:      NewPlayer("玩家", InitialChips),
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

	fmt.Printf("=== 第%d轮开始 ===\n", g.RoundNumber)
	fmt.Printf("💰 当前筹码: %d\n", g.Player.Chips)
	fmt.Println()

	if !g.Player.HasChips() {
		fmt.Println("💸 你的筹码用完了！游戏结束！")
		return
	}

	g.Player.ResetRound()

	// 下注阶段
	if !g.placeBet() {
		return
	}

	// 检查牌堆是否足够
	if len(g.Deck.Cards) < 10 {
		fmt.Println("牌不够了，重新洗牌...")
		g.Deck = NewDeck()
		time.Sleep(1 * time.Second)
	}

	g.dealInitialCards()

	g.playerTurn()

	if !g.Player.Hand.IsBust() && !g.Player.Hand.IsBlackjack() {
		g.dealerTurn()
	}

	g.determineWinner()
}

// 预设的下注选项
var betOptions = []int{10, 25, 50, 100, 200}

// placeBet 下注阶段
func (g *Game) placeBet() bool {
	fmt.Println("=== 下注阶段 ===")
	fmt.Printf("💰 当前筹码: %d\n", g.Player.Chips)
	fmt.Println()

	// 显示可用的下注选项
	fmt.Println("请选择下注金额:")
	validOptions := []int{}
	optionIndex := 1

	for _, amount := range betOptions {
		if amount <= g.Player.Chips {
			fmt.Printf("%d. %d 筹码\n", optionIndex, amount)
			validOptions = append(validOptions, amount)
			optionIndex++
		}
	}

	// 如果玩家筹码很少，添加全押选项
	if g.Player.Chips < betOptions[0] && g.Player.Chips > 0 {
		fmt.Printf("%d. %d 筹码 (全押)\n", optionIndex, g.Player.Chips)
		validOptions = append(validOptions, g.Player.Chips)
		optionIndex++
	}

	fmt.Printf("%d. 退出游戏\n", optionIndex)
	fmt.Println()

	for {
		input := g.getInput("请选择 (输入选项编号): ")

		var choice int
		if _, err := fmt.Sscanf(input, "%d", &choice); err != nil {
			fmt.Println("❌ 请输入有效的选项编号")
			continue
		}

		if choice == optionIndex {
			return false
		}

		if choice < 1 || choice > len(validOptions) {
			fmt.Printf("❌ 请输入 1-%d 之间的选项编号\n", optionIndex)
			continue
		}

		betAmount := validOptions[choice-1]

		if g.Player.PlaceBet(betAmount) {
			fmt.Printf("✅ 下注成功！下注金额: %d\n", betAmount)
			fmt.Printf("💰 剩余筹码: %d\n", g.Player.Chips)
			fmt.Println()
			time.Sleep(1 * time.Second)
			return true
		} else {
			fmt.Println("❌ 下注失败，请重试")
		}
	}
}

// playerTurn 玩家回合
func (g *Game) playerTurn() {
	fmt.Println("\n=== 玩家回合 ===")

	for {
		g.displayGame(true)

		if g.Player.Hand.IsBlackjack() {
			fmt.Println("🎉 恭喜！你拿到了Blackjack！")
			return
		}

		if g.Player.Hand.IsBust() {
			fmt.Println("💥 爆牌了！你输了！")
			return
		}

		// 如果已经加倍，只能拿一张牌后必须停牌
		if g.Player.DoubledDown {
			fmt.Println("你已经加倍，必须拿一张牌后停牌...")
			card := g.Deck.Deal()
			g.Player.Hand.AddCard(card)
			fmt.Printf("你拿到了: %s\n", card)
			time.Sleep(2 * time.Second)
			return
		}

		// 构建选项提示
		prompt := "请选择: (h)要牌 (s)停牌"
		if g.Player.CanDoubleDown() {
			prompt += " (d)加倍"
		}
		prompt += " (q)退出游戏: "

		choice := g.getInput(prompt)

		switch strings.ToLower(choice) {
		case "h", "hit":
			g.Player.Hand.AddCard(g.Deck.Deal())
			fmt.Printf("你拿到了: %s\n", g.Player.Hand.Cards[len(g.Player.Hand.Cards)-1])
			time.Sleep(1 * time.Second)
		case "s", "stand":
			fmt.Println("你选择停牌")
			time.Sleep(1 * time.Second)
			return
		case "d", "double", "doubledown":
			if g.Player.CanDoubleDown() {
				if g.Player.DoubleBet() {
					fmt.Printf("✅ 加倍成功！新的下注金额: %d\n", g.Player.Bet)
					fmt.Printf("💰 剩余筹码: %d\n", g.Player.Chips)
					fmt.Println("现在你必须拿一张牌后停牌...")
					time.Sleep(2 * time.Second)
				} else {
					fmt.Println("❌ 筹码不足，无法加倍")
					time.Sleep(1 * time.Second)
				}
			} else {
				fmt.Println("❌ 现在无法加倍")
				time.Sleep(1 * time.Second)
			}
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
	betAmount := g.Player.Bet // 保存下注金额用于显示
	isDoubled := g.Player.DoubledDown

	fmt.Printf("=== 第%d轮游戏结果 ===\n", g.RoundNumber)
	fmt.Printf("💰 下注金额: %d", betAmount)
	if isDoubled {
		fmt.Printf(" (加倍)")
	}
	fmt.Println()

	if g.Player.Hand.IsBust() {
		fmt.Println("💀 你爆牌了，庄家获胜！")
		g.Player.LoseBet()
		fmt.Printf("💸 损失: %d 筹码", betAmount)
		if isDoubled {
			fmt.Printf(" (加倍后)")
		}
		fmt.Println()
		g.showChipsStatus()
		return
	}

	if g.Dealer.Hand.IsBust() {
		fmt.Println("🎉 庄家爆牌，你获胜！")
		g.Player.WinBet(1.0) // 1:1 赔率
		winnings := betAmount
		fmt.Printf("💰 获得: %d 筹码 (1:1)", winnings)
		if isDoubled {
			fmt.Printf(" (加倍后)")
		}
		fmt.Println()
		g.showChipsStatus()
		return
	}

	if playerBlackjack && dealerBlackjack {
		fmt.Println("🤝 双方都是Blackjack，平局！")
		g.Player.PushBet()
		fmt.Println("💰 返还下注金额")
		g.showChipsStatus()
		return
	}

	if playerBlackjack {
		// 注意：加倍后的Blackjack通常只按1:1赔率，不是3:2
		if isDoubled {
			fmt.Println("🎉 你拿到Blackjack，获胜！(加倍后按1:1赔率)")
			g.Player.WinBet(1.0) // 加倍后Blackjack按1:1赔率
			winnings := betAmount
			fmt.Printf("💰 获得: %d 筹码 (1:1 加倍后)\n", winnings)
		} else {
			fmt.Println("🎉 你拿到Blackjack，获胜！(3:2赔率)")
			winnings := int(float64(betAmount) * 1.5)
			g.Player.WinBet(1.5) // 3:2 赔率
			fmt.Printf("💰 获得: %d 筹码 (3:2)\n", winnings)
		}
		g.showChipsStatus()
		return
	}

	if dealerBlackjack {
		fmt.Println("💀 庄家拿到Blackjack，你输了！")
		g.Player.LoseBet()
		fmt.Printf("💸 损失: %d 筹码", betAmount)
		if isDoubled {
			fmt.Printf(" (加倍后)")
		}
		fmt.Println()
		g.showChipsStatus()
		return
	}

	if playerValue > dealerValue {
		fmt.Println("🎉 你的点数更高，获胜！")
		g.Player.WinBet(1.0) // 1:1 赔率
		winnings := betAmount
		fmt.Printf("💰 获得: %d 筹码 (1:1)", winnings)
		if isDoubled {
			fmt.Printf(" (加倍后)")
		}
		fmt.Println()
		g.showChipsStatus()
	} else if playerValue < dealerValue {
		fmt.Println("💀 庄家点数更高，你输了！")
		g.Player.LoseBet()
		fmt.Printf("💸 损失: %d 筹码", betAmount)
		if isDoubled {
			fmt.Printf(" (加倍后)")
		}
		fmt.Println()
		g.showChipsStatus()
	} else {
		fmt.Println("🤝 点数相同，平局！")
		g.Player.PushBet()
		fmt.Println("💰 返还下注金额")
		g.showChipsStatus()
	}
}

// showChipsStatus 显示筹码状态
func (g *Game) showChipsStatus() {
	fmt.Printf("💰 当前筹码: %d\n", g.Player.Chips)
	if g.Player.Chips <= 0 {
		fmt.Println("💸 你的筹码用完了！")
	}
	fmt.Println()
}

// displayGame 显示游戏状态
func (g *Game) displayGame(hideDealer bool) {
	// 显示筹码和下注信息
	fmt.Printf("💰 筹码: %d | 💸 下注: %d\n", g.Player.Chips, g.Player.Bet)
	fmt.Println()

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

// 游戏结算
func (g *Game) displayGameOver() {
	fmt.Println()
	fmt.Println("🎉 游戏结算：")
	fmt.Printf("💰 共进行%d轮游戏。最终筹码: %d ", g.RoundNumber, g.Player.Chips)
	if g.Player.Chips > g.Player.InitialChips {
		fmt.Printf("盈利: %d\n", g.Player.Chips-g.Player.InitialChips)
	} else {
		fmt.Printf("亏损: %d\n", g.Player.InitialChips-g.Player.Chips)
	}
	fmt.Println("🎉 游戏结束！欢迎下次再来！")
	fmt.Println()
}
