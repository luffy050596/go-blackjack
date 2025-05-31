package main

import (
	"fmt"
	"time"
)

// Display 显示相关的结构
type Display struct {
	game *Game
}

// NewDisplay 创建显示管理器
func NewDisplay(game *Game) *Display {
	return &Display{game: game}
}

// showWelcome 显示欢迎信息
func (d *Display) showWelcome() {
	fmt.Println("🃏 欢迎来到二十一点游戏！ 🃏")
	fmt.Println()
}

// showRoundStart 显示回合开始信息
func (d *Display) showRoundStart(roundNumber int, chips int) {
	fmt.Printf("=== 第%d轮开始 ===\n", roundNumber)
	fmt.Printf("💰 当前筹码: %d\n", chips)
	fmt.Println()
}

// showBetSection 显示下注区域标题
func (d *Display) showBetSection(chips int) {
	fmt.Println("=== 下注阶段 ===")
	fmt.Printf("💰 当前筹码: %d\n", chips)
	fmt.Println()
}

// showBetOptions 显示下注选项
func (d *Display) showBetOptions(validOptions []int, chips int) {
	fmt.Println("请选择下注金额:")

	for i, amount := range validOptions {
		if amount == chips && amount < 50 { // 如果是全押且金额较小
			fmt.Printf("%d. %d 筹码 (全押)\n", i+1, amount)
		} else {
			fmt.Printf("%d. %d 筹码\n", i+1, amount)
		}
	}

	fmt.Printf("%d. 退出游戏\n", len(validOptions)+1)
	fmt.Println()
}

// showBetSuccess 显示下注成功信息
func (d *Display) showBetSuccess(betAmount, remainingChips int) {
	fmt.Printf("✅ 下注成功！下注金额: %d\n", betAmount)
	fmt.Printf("💰 剩余筹码: %d\n", remainingChips)
	fmt.Println()
	time.Sleep(1 * time.Second)
}

// showPlayerTurnStart 显示玩家回合开始
func (d *Display) showPlayerTurnStart() {
	fmt.Println("\n=== 玩家回合 ===")
}

// showDealerTurnStart 显示庄家回合开始
func (d *Display) showDealerTurnStart() {
	fmt.Println("\n=== 庄家回合 ===")
}

// showBlackjack 显示Blackjack信息
func (d *Display) showBlackjack() {
	fmt.Println("🎉 恭喜！你拿到了Blackjack！")
}

// showPlayerBust 显示玩家爆牌
func (d *Display) showPlayerBust() {
	fmt.Println("💥 爆牌了！你输了！")
}

// showCardDealt 显示发牌信息
func (d *Display) showCardDealt(card Card, isPlayer bool) {
	if isPlayer {
		fmt.Printf("你拿到了: %s\n", card)
	} else {
		fmt.Printf("庄家拿到: %s\n", card)
	}
	time.Sleep(1 * time.Second)
}

// showPlayerChoice 显示玩家选择信息
func (d *Display) showPlayerChoice(choice string) {
	switch choice {
	case DisplayActionStand:
		fmt.Println("你选择停牌")
	case DisplayActionDouble:
		fmt.Println("现在你必须拿一张牌后停牌...")
	}
	time.Sleep(1 * time.Second)
}

// showDoubleDownSuccess 显示加倍成功信息
func (d *Display) showDoubleDownSuccess(newBet, remainingChips int) {
	fmt.Printf("✅ 加倍成功！新的下注金额: %d\n", newBet)
	fmt.Printf("💰 剩余筹码: %d\n", remainingChips)
	fmt.Println("现在你必须拿一张牌后停牌...")
	time.Sleep(2 * time.Second)
}

// showDealerAction 显示庄家行动
func (d *Display) showDealerAction(action string) {
	switch action {
	case DisplayActionHit:
		fmt.Println("庄家点数小于17，必须要牌...")
		time.Sleep(2 * time.Second)
	case DisplayActionStand:
		fmt.Println("庄家停牌")
		time.Sleep(1 * time.Second)
	case DisplayActionBust:
		fmt.Println("💥 庄家爆牌！")
		time.Sleep(1 * time.Second)
	}
}

// showGameState 显示游戏状态
func (d *Display) showGameState(hideDealer bool) {
	player := d.game.Player
	dealer := d.game.Dealer

	// 显示筹码和下注信息
	fmt.Printf("💰 筹码: %d | 💸 下注: %d\n", player.Chips, player.Bet)
	fmt.Println()

	// 显示庄家手牌
	if hideDealer && len(dealer.Hand.Cards) > 0 {
		fmt.Printf("庄家: %s [?] (点数: ?)\n", dealer.Hand.Cards[0])
	} else {
		fmt.Printf("庄家: %s (点数: %d)\n", dealer.Hand.String(), dealer.Hand.Value())
	}

	// 显示玩家手牌
	fmt.Printf("%s: %s (点数: %d)\n", player.Name, player.Hand.String(), player.Hand.Value())
	fmt.Println()
}

// showRoundResult 显示回合结果
func (d *Display) showRoundResult(result GameResult) {
	fmt.Printf("=== 第%d轮游戏结果 ===\n", d.game.RoundNumber)
	fmt.Printf("💰 下注金额: %d", result.BetAmount)
	if result.IsDoubled {
		fmt.Printf(" (加倍)")
	}
	fmt.Println()

	switch result.ResultType {
	case PlayerBust:
		fmt.Println("💀 你爆牌了，庄家获胜！")
		fmt.Printf("💸 损失: %d 筹码", result.BetAmount)
	case DealerBust:
		fmt.Println("🎉 庄家爆牌，你获胜！")
		fmt.Printf("💰 获得: %d 筹码 (1:1)", result.BetAmount)
	case BothBlackjack:
		fmt.Println("🤝 双方都是Blackjack，平局！")
		fmt.Println("💰 返还下注金额")
	case PlayerBlackjack:
		if result.IsDoubled {
			fmt.Println("🎉 你拿到Blackjack，获胜！(加倍后按1:1赔率)")
			fmt.Printf("💰 获得: %d 筹码 (1:1 加倍后)", result.BetAmount)
		} else {
			fmt.Println("🎉 你拿到Blackjack，获胜！(3:2赔率)")
			fmt.Printf("💰 获得: %d 筹码 (3:2)", int(float64(result.BetAmount)*1.5))
		}
	case DealerBlackjack:
		fmt.Println("💀 庄家拿到Blackjack，你输了！")
		fmt.Printf("💸 损失: %d 筹码", result.BetAmount)
	case PlayerWin:
		fmt.Println("🎉 你的点数更高，获胜！")
		fmt.Printf("💰 获得: %d 筹码 (1:1)", result.BetAmount)
	case DealerWin:
		fmt.Println("💀 庄家点数更高，你输了！")
		fmt.Printf("💸 损失: %d 筹码", result.BetAmount)
	case Push:
		fmt.Println("🤝 点数相同，平局！")
		fmt.Println("💰 返还下注金额")
	}

	if result.IsDoubled && result.ResultType != BothBlackjack && result.ResultType != Push {
		fmt.Printf(" (加倍后)")
	}
	fmt.Println()

	d.showChipsStatus()
}

// showChipsStatus 显示筹码状态
func (d *Display) showChipsStatus() {
	fmt.Printf("💰 当前筹码: %d\n", d.game.Player.Chips)
	if d.game.Player.Chips <= 0 {
		fmt.Println("💸 你的筹码用完了！")
	}
	fmt.Println()
}

// showGameOver 显示游戏结束信息
func (d *Display) showGameOver() {
	player := d.game.Player
	fmt.Println()
	fmt.Println("🎉 游戏结算：")
	fmt.Printf("💰 共进行%d轮游戏。最终筹码: %d ", d.game.RoundNumber, player.Chips)
	if player.Chips > player.InitialChips {
		fmt.Printf("盈利: %d\n", player.Chips-player.InitialChips)
	} else {
		fmt.Printf("亏损: %d\n", player.InitialChips-player.Chips)
	}
	fmt.Println("🎉 游戏结束！欢迎下次再来！")
	fmt.Println()
}

// showError 显示错误信息
func (d *Display) showError(msg string) {
	fmt.Printf("❌ %s\n", msg)
	time.Sleep(1 * time.Second)
}

// showInfo 显示一般信息
func (d *Display) showInfo(msg string) {
	fmt.Println(msg)
	time.Sleep(1 * time.Second)
}

// clearScreen 清屏
func (d *Display) clearScreen() {
	fmt.Print("\033[2J\033[H")
}
