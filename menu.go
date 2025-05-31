package main

import (
	"fmt"
	"strings"
)

// showMenu 主菜单
func (g *Game) showMenu() {
	for {
		g.display.showWelcome()
		fmt.Println("1. 开始游戏")
		fmt.Println("2. 查看规则")
		fmt.Println("3. 退出游戏")
		fmt.Println()

		choice := g.getInput("请选择 (1-3): ")

		switch choice {
		case MenuOptionStart:
			g.gameLoop()
		case MenuOptionRules:
			g.showRules()
		case MenuOptionExit:
			g.display.showInfo("感谢游戏！再见！")
			return
		default:
			g.display.showError("无效选择，请重新输入")
		}
	}
}

// gameLoop 游戏循环
func (g *Game) gameLoop() {
	for {
		g.clearScreen()

		// 检查玩家是否还有筹码
		if !g.Player.HasChips() {
			g.display.showInfo("你的筹码用完了！")
			choice := g.getInput("是否重新开始游戏？(y/n): ")
			if strings.ToLower(choice) == InputYes || strings.ToLower(choice) == InputYesFull {
				g.resetGame()
				g.display.showInfo("🎉 重新开始！你获得了1000筹码")
				continue
			} else {
				break
			}
		}

		g.playRound()

		// 如果玩家还有筹码，询问是否继续
		if g.Player.HasChips() {
			fmt.Println()
			choice := g.getInput("再来一局？(y/n): ")

			if strings.ToLower(choice) != InputYes && strings.ToLower(choice) != InputYesFull {
				break
			}
		}
	}

	g.display.showGameOver()
}

// resetGame 重置游戏状态
func (g *Game) resetGame() {
	g.Player.Chips = InitialChips
	g.Player.Bet = 0
	g.RoundNumber = 0
}
