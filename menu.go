package main

import (
	"fmt"
	"strings"
	"time"
)

// showMenu 主菜单
func (g *Game) showMenu() {
	for {
		g.clearScreen()
		fmt.Println("🃏 欢迎来到二十一点游戏！ 🃏")
		fmt.Println()
		fmt.Println("1. 开始游戏")
		fmt.Println("2. 查看规则")
		fmt.Println("3. 退出游戏")
		fmt.Println()

		choice := g.getInput("请选择 (1-3): ")

		switch choice {
		case "1":
			g.gameLoop()
		case "2":
			g.showRules()
		case "3":
			fmt.Println("感谢游戏！再见！")
			return
		default:
			fmt.Println("无效选择，请重新输入")
			time.Sleep(1 * time.Second)
		}
	}
}

// getInput 获取用户输入
func (g *Game) getInput(prompt string) string {
	fmt.Print(prompt)
	g.Scanner.Scan()
	return strings.TrimSpace(g.Scanner.Text())
}

// gameLoop 游戏循环
func (g *Game) gameLoop() {
	for {
		g.clearScreen()

		// 检查玩家是否还有筹码
		if !g.Player.HasChips() {
			fmt.Println("💸 你的筹码用完了！")
			fmt.Println()
			choice := g.getInput("是否重新开始游戏？(y/n): ")
			if strings.ToLower(choice) == "y" || strings.ToLower(choice) == "yes" {
				// 重置玩家筹码
				g.Player.Chips = 1000
				g.Player.Bet = 0
				g.RoundNumber = 0
				fmt.Println("🎉 重新开始！你获得了1000筹码")
				time.Sleep(1 * time.Second)
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

			if strings.ToLower(choice) != "y" && strings.ToLower(choice) != "yes" {
				break
			}
		}
	}
}
