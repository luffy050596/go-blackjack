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
		g.playRound()

		fmt.Println()
		choice := g.getInput("再来一局？(y/n): ")

		if strings.ToLower(choice) != "y" && strings.ToLower(choice) != "yes" {
			break
		}
	}
}
