package main

import (
	"github.com/luffy050596/go-blackjack/internal/interfaces/cli"
)

func main() {
	// 创建命令行游戏处理器
	gameHandler := cli.NewGameHandler()

	// 运行游戏
	gameHandler.Run()
}
