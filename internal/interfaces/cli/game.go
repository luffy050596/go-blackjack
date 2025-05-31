package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/luffy050596/go-blackjack/internal/application/services"
	"github.com/luffy050596/go-blackjack/internal/domain/entities"
)

// GameHandler 游戏命令行处理器
type GameHandler struct {
	gameService *services.GameApplicationService
	scanner     *bufio.Scanner
	display     *DisplayService
}

// NewGameHandler 创建游戏处理器
func NewGameHandler() *GameHandler {
	handler := &GameHandler{
		gameService: services.NewGameApplicationService("玩家"),
		scanner:     bufio.NewScanner(os.Stdin),
		display:     NewDisplayService(),
	}
	return handler
}

// Run 运行游戏
func (h *GameHandler) Run() {
	h.display.ShowWelcome()

	for {
		h.display.ShowMenu()
		choice := h.getInput("请选择选项: ")

		switch choice {
		case MenuOptionStart:
			h.playGame()
		case MenuOptionRules:
			h.display.ShowRules()
		case MenuOptionExit:
			h.display.ShowGoodbye()
			return
		default:
			h.display.ShowError("无效的选择，请重试")
		}
	}
}

// playGame 游戏主循环
func (h *GameHandler) playGame() {
	for !h.gameService.IsGameOver() {
		if err := h.playRound(); err != nil {
			h.display.ShowError(fmt.Sprintf("游戏错误: %v", err))
			return
		}

		if h.gameService.IsGameOver() {
			h.display.ShowGameOver()
			return
		}

		if !h.askPlayAgain() {
			return
		}
	}
}

// playRound 单轮游戏
func (h *GameHandler) playRound() error {
	// 开始新一轮
	if err := h.gameService.StartNewRound(); err != nil {
		return err
	}

	gameState := h.gameService.GetGameState()
	h.display.ShowRoundStart(gameState.RoundNumber, gameState.PlayerChips)

	// 下注阶段
	if !h.handleBetting() {
		return nil // 玩家选择退出
	}

	// 发初始牌
	if err := h.gameService.DealInitialCards(); err != nil {
		return err
	}

	// 玩家回合
	if err := h.handlePlayerTurn(); err != nil {
		return err
	}

	// 庄家回合（如果需要）
	gameState = h.gameService.GetGameState()
	if gameState.State == entities.StateDealerTurn {
		if err := h.handleDealerTurn(); err != nil {
			return err
		}
	}

	// 显示结果
	result := h.gameService.EvaluateGame()
	if result != nil {
		h.display.ShowGameResult(result)
	}

	return nil
}

// handleBetting 处理下注阶段
func (h *GameHandler) handleBetting() bool {
	gameState := h.gameService.GetGameState()
	h.display.ShowBettingSection(gameState.PlayerChips)

	betOptions := h.gameService.GetBetOptions()
	h.display.ShowBetOptions(betOptions)

	for {
		input := h.getInput("请选择下注金额 (输入选项编号或 'q' 退出): ")

		if strings.ToLower(input) == entities.InputQuit {
			return false
		}

		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(betOptions) {
			h.display.ShowError("请输入有效的选项编号")
			continue
		}

		betAmount := betOptions[choice-1]
		if err := h.gameService.PlaceBet(betAmount); err != nil {
			h.display.ShowError(fmt.Sprintf("下注失败: %v", err))
			continue
		}

		h.display.ShowBetSuccess(betAmount)
		return true
	}
}

// handlePlayerTurn 处理玩家回合
func (h *GameHandler) handlePlayerTurn() error {
	h.display.ShowPlayerTurnStart()

	defer func() {
		h.gameService.StartDealerTurn()
	}()

	for {
		gameState := h.gameService.GetGameState()
		h.display.ShowGameState(gameState, true)

		// 显示获胜概率
		probabilities := h.gameService.CalculateWinProbabilities()
		h.display.ShowProbabilities(probabilities)

		// 检查21点
		if gameState.PlayerHand.Value == 21 && len(gameState.PlayerHand.Cards) == 2 {
			h.display.ShowBlackjack()
			break
		}

		// 检查爆牌
		if gameState.PlayerHand.Value > 21 {
			h.display.ShowPlayerBust()
			break
		}

		// 获取玩家输入
		prompt := h.display.buildPlayerPrompt(WithDoubleDown(h.gameService.CanPlayerDoubleDown()))
		input := h.getInput(prompt)

		// 处理玩家行动
		action := ParsePlayerInput(input)
		if action == entities.ActionInvalid {
			h.display.ShowError("无效的输入，请重试")
			continue
		}

		result, err := h.gameService.ProcessPlayerAction(action)
		if err != nil {
			return err
		}

		// 显示行动结果
		h.display.ShowActionResult(result)

		if result.Action == entities.ActionQuit {
			os.Exit(0)
		}

		if !result.Continue {
			break
		}
	}

	return nil
}

// handleDealerTurn 处理庄家回合
func (h *GameHandler) handleDealerTurn() error {
	h.display.ShowDealerTurnStart()

	if err := h.gameService.ProcessDealerTurn(); err != nil {
		return err
	}

	gameState := h.gameService.GetGameState()
	h.display.ShowGameState(gameState, false)

	return nil
}

// getInput 获取用户输入
func (h *GameHandler) getInput(prompt string) string {
	fmt.Print(prompt)
	h.scanner.Scan()
	return strings.TrimSpace(h.scanner.Text())
}

// askPlayAgain 询问是否继续游戏
func (h *GameHandler) askPlayAgain() bool {
	input := h.getInput("是否继续游戏? (y/n): ")
	return strings.ToLower(strings.TrimSpace(input)) == "y" ||
		strings.ToLower(strings.TrimSpace(input)) == "yes"
}
