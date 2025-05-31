package main

import "fmt"

// showRules 显示游戏规则
func (g *Game) showRules() {
	g.clearScreen()
	fmt.Println("=== 二十一点游戏规则 ===")
	fmt.Println()
	fmt.Println("🎯 游戏目标:")
	fmt.Println("   让手中牌的点数尽可能接近21点，但不能超过21点")
	fmt.Println("   点数比庄家更接近21点就获胜")
	fmt.Println()
	fmt.Println("🃏 牌面点数:")
	fmt.Println("   • 数字牌(2-10): 按牌面数字计算")
	fmt.Println("   • 花牌(J,Q,K): 每张都是10点")
	fmt.Println("   • A: 可以是1点或11点(自动选择最优)")
	fmt.Println()
	fmt.Println("🎮 游戏流程:")
	fmt.Println("   1. 玩家和庄家各发2张牌")
	fmt.Println("   2. 玩家选择要牌(h)或停牌(s)")
	fmt.Println("   3. 庄家小于17点必须要牌，17点以上必须停牌")
	fmt.Println("   4. 比较点数决定胜负")
	fmt.Println()
	fmt.Println("🏆 特殊情况:")
	fmt.Println("   • Blackjack: 前两张牌就是21点(A+10点牌)")
	fmt.Println("   • 爆牌: 点数超过21点立即失败")
	fmt.Println("   • 平局: 双方点数相同")
	fmt.Println()
	fmt.Print("按回车键继续...")
	g.Scanner.Scan()
}
