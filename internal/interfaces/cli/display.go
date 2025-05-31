package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/luffy050596/go-blackjack/internal/application/dtos"
	"github.com/luffy050596/go-blackjack/internal/domain/entities"
)

// DisplayService 显示服务
type DisplayService struct{}

// NewDisplayService 创建显示服务
func NewDisplayService() *DisplayService {
	return &DisplayService{}
}

// ShowWelcome 显示欢迎信息
func (d *DisplayService) ShowWelcome() {
	d.clearScreen()
	fmt.Println("🃏 欢迎来到二十一点游戏! 🃏")
	fmt.Println("=======================")
	fmt.Println()
}

// ShowMenu 显示主菜单
func (d *DisplayService) ShowMenu() {
	fmt.Println("请选择:")
	fmt.Print(MenuOptionStart + ". 开始游戏\n")
	fmt.Print(MenuOptionRules + ". 游戏规则\n")
	fmt.Print(MenuOptionExit + ". 退出游戏\n")
	fmt.Println()
}

// ShowGoodbye 显示再见信息
func (d *DisplayService) ShowGoodbye() {
	fmt.Println("感谢游戏！再见！👋")
}

// ShowError 显示错误信息
func (d *DisplayService) ShowError(message string) {
	fmt.Printf("❌ %s\n\n", message)
}

// ShowRoundStart 显示回合开始
func (d *DisplayService) ShowRoundStart(round, chips int) {
	fmt.Printf("🎯 第 %d 轮游戏开始! 💰 当前筹码: %d\n", round, chips)
	fmt.Println(strings.Repeat("=", 40))
	fmt.Println()
}

// ShowBettingSection 显示下注区域
func (d *DisplayService) ShowBettingSection(chips int) {
	fmt.Printf("💰 当前筹码: %d\n", chips)
	fmt.Println("请选择下注金额:")
}

// ShowBetOptions 显示下注选项
func (d *DisplayService) ShowBetOptions(options []int) {
	for i, amount := range options {
		fmt.Printf("%d. %d 筹码\n", i+1, amount)
	}
	fmt.Printf("%d. 退出游戏\n", len(options)+1)
	fmt.Println()
}

// ShowBetSuccess 显示下注成功
func (d *DisplayService) ShowBetSuccess(amount int) {
	fmt.Printf("✅ 下注成功: %d 筹码\n\n", amount)
	time.Sleep(500 * time.Millisecond)
}

// ShowPlayerTurnStart 显示玩家回合开始
func (d *DisplayService) ShowPlayerTurnStart() {
	fmt.Println("🎮 === 玩家回合开始 ===")
}

// ShowDealerTurnStart 显示庄家回合开始
func (d *DisplayService) ShowDealerTurnStart() {
	fmt.Println("\n🤖 === 庄家回合开始 ===")
	time.Sleep(1 * time.Second)
}

type playerPromptOptions struct {
	doubleDown bool
}

type playerPromptOption func(options *playerPromptOptions)

func WithDoubleDown(doubleDown bool) playerPromptOption {
	return func(options *playerPromptOptions) {
		options.doubleDown = doubleDown
	}
}

// buildPlayerPrompt 构建玩家输入提示
func (d *DisplayService) buildPlayerPrompt(options ...playerPromptOption) string {
	opts := playerPromptOptions{}

	for _, option := range options {
		option(&opts)
	}

	prompt := "请选择: (h)要牌 (s)停牌"
	if opts.doubleDown {
		prompt += " (d)加倍"
	}
	prompt += " (q)退出: "
	return prompt
}

// ShowGameState 显示游戏状态
func (d *DisplayService) ShowGameState(gameState *dtos.GameStateDTO, hideFirstDealerCard bool) {
	fmt.Print("\n👨 庄家手牌")
	if hideFirstDealerCard && len(gameState.DealerHand.Cards) > 1 {
		fmt.Println(" (第一张牌隐藏):")
		d.showHand(gameState.DealerHand, true)
	} else {
		fmt.Printf(" (点数: %d):\n", gameState.DealerHand.Value)
		d.showHand(gameState.DealerHand, false)
	}

	fmt.Printf("\n👨 玩家手牌 (点数: %d):\n", gameState.PlayerHand.Value)
	d.showHand(gameState.PlayerHand, false)

	fmt.Println()
}

// showHand 显示手牌
func (d *DisplayService) showHand(hand *dtos.HandDTO, hideFirst bool) {
	for i, card := range hand.Cards {
		if hideFirst && i == 0 {
			fmt.Print("🂠 ")
		} else {
			fmt.Printf("%s%s ", d.getSuitSymbol(card.Suit), card.Rank)
		}
	}
	fmt.Println()
}

// getSuitSymbol 获取花色符号
func (d *DisplayService) getSuitSymbol(suit string) string {
	switch suit {
	case "Hearts":
		return "♥️"
	case "Diamonds":
		return "♦️"
	case "Clubs":
		return "♣️"
	case "Spades":
		return "♠️"
	default:
		return "🃏"
	}
}

// ShowBlackjack 显示21点
func (d *DisplayService) ShowBlackjack() {
	fmt.Println("🎉 21点! 🎉")
}

// ShowPlayerBust 显示玩家爆牌
func (d *DisplayService) ShowPlayerBust() {
	fmt.Println("💥 爆牌了! 💥")
}

// ShowActionResult 显示行动结果
func (d *DisplayService) ShowActionResult(result *dtos.ActionResultDTO) {
	if !result.Success {
		d.ShowError(result.Message)
		return
	}

	switch result.Action {
	case entities.ActionHit:
		if result.Card != nil {
			fmt.Printf("🃏 获得一张牌: %s%s\n",
				d.getSuitSymbol(result.Card.Suit), result.Card.Rank)
		}
	case entities.ActionStand:
		fmt.Println("✋ 停牌")
	case entities.ActionDoubleDown:
		fmt.Println("🎯 加倍下注!自动要牌")
		if result.Card != nil {
			fmt.Printf("🃏 获得一张牌: %s%s\n",
				d.getSuitSymbol(result.Card.Suit), result.Card.Rank)
		}
	}

	time.Sleep(500 * time.Millisecond)
}

// ShowGameResult 显示游戏结果
func (d *DisplayService) ShowGameResult(result *dtos.GameResultDTO) {
	fmt.Println("\n" + strings.Repeat("=", 40))
	fmt.Println("🎯 游戏结果")
	fmt.Println(strings.Repeat("=", 40))
	fmt.Printf("结果: %s\n", GetResultMessage(result.Type))
	fmt.Printf("本轮下注: %d 筹码", result.BetAmount)
	if result.IsDoubled {
		fmt.Print(" (已加倍)")
	}
	fmt.Printf("\n当前筹码: %d\n", result.PlayerChips)
	fmt.Println(strings.Repeat("=", 40))
	fmt.Println()
}

// ShowGameOver 显示游戏结束
func (d *DisplayService) ShowGameOver() {
	fmt.Println("💸 筹码用完了！游戏结束！")
	fmt.Println("感谢游戏！")
}

// ShowProbabilities 显示获胜概率
func (d *DisplayService) ShowProbabilities(probabilities *dtos.ProbabilityResultDTO) {
	if probabilities == nil {
		return
	}

	fmt.Println(strings.Repeat("─", 40))
	fmt.Println("📊 当前获胜概率分析")
	fmt.Println(strings.Repeat("─", 40))

	// 主要概率
	fmt.Printf("🟢 玩家获胜概率: %.1f%%\n", probabilities.PlayerWinProbability*100)
	fmt.Printf("🔴 庄家获胜概率: %.1f%%\n", probabilities.DealerWinProbability*100)
	fmt.Printf("🟡 平局概率:     %.1f%%\n", probabilities.PushProbability*100)

	fmt.Println()

	// 详细概率
	fmt.Println("📈 详细分析:")
	fmt.Printf("   💥 玩家爆牌概率: %.1f%%\n", probabilities.PlayerBustProbability*100)
	fmt.Printf("   💥 庄家爆牌概率: %.1f%%\n", probabilities.DealerBustProbability*100)
	fmt.Printf("   🎯 玩家21点概率: %.1f%%\n", probabilities.Player21Probability*100)
	fmt.Printf("   🎯 庄家21点概率: %.1f%%\n", probabilities.Dealer21Probability*100)

	// 如果有自然21点（Blackjack），也显示出来
	if probabilities.PlayerBlackjackProb > 0 {
		fmt.Printf("   🌟 玩家Blackjack概率: %.1f%%\n", probabilities.PlayerBlackjackProb*100)
	}
	if probabilities.DealerBlackjackProb > 0 {
		fmt.Printf("   🌟 庄家Blackjack概率: %.1f%%\n", probabilities.DealerBlackjackProb*100)
	}

	// 操作胜率分析
	if probabilities.ActionAnalysis != nil {
		d.showActionAnalysis(probabilities.ActionAnalysis)
	}

	fmt.Println(strings.Repeat("─", 40))
	fmt.Println()
}

// showActionAnalysis 显示操作胜率分析
func (d *DisplayService) showActionAnalysis(analysis *dtos.ActionAnalysisDTO) {
	fmt.Println()
	fmt.Println("🎯 操作胜率对比:")

	actions := []struct {
		name    string
		winRate float64
		canUse  bool
		symbol  string
	}{
		{"停牌", analysis.StandWinRate, analysis.CanStand, "✋"},
		{"要牌", analysis.HitWinRate, analysis.CanHit, "👆"},
		{"加倍", analysis.DoubleWinRate, analysis.CanDouble, "⚡"},
		{"分牌", analysis.SplitWinRate, analysis.CanSplit, "✂️"},
	}

	// 显示可用操作的胜率
	for _, action := range actions {
		if action.canUse {
			// 如果是推荐操作，添加特殊标记
			if analysis.RecommendedAction == getActionKey(action.name) {
				fmt.Printf("   %s %s: %.1f%% ⭐ (推荐)\n", action.symbol, action.name, action.winRate*100)
			} else {
				fmt.Printf("   %s %s: %.1f%%\n", action.symbol, action.name, action.winRate*100)
			}
		}
	}

	// 显示最优期望值
	if analysis.ExpectedValue > 0 {
		fmt.Printf("\n🏆 最优策略期望胜率: %.1f%%\n", analysis.ExpectedValue*100)
	}
}

// getActionKey 将操作名称转换为操作键
func getActionKey(actionName string) string {
	switch actionName {
	case "停牌":
		return "stand"
	case "要牌":
		return "hit"
	case "加倍":
		return "double"
	case "分牌":
		return "split"
	default:
		return ""
	}
}

// clearScreen 清屏
func (d *DisplayService) clearScreen() {
	fmt.Print("\033[2J\033[H")
}

// GetResultMessage 获取结果消息
func GetResultMessage(resultType entities.ResultType) string {
	switch resultType {
	case entities.PlayerBust:
		return "玩家爆牌，庄家获胜！"
	case entities.DealerBust:
		return "庄家爆牌，玩家获胜！"
	case entities.BothBlackjack:
		return "双方都是21点，平局！"
	case entities.PlayerBlackjack:
		return "玩家21点，获胜！"
	case entities.DealerBlackjack:
		return "庄家21点，玩家失败！"
	case entities.PlayerWin:
		return "玩家获胜！"
	case entities.DealerWin:
		return "庄家获胜！"
	case entities.Push:
		return "平局！"
	default:
		return "未知结果"
	}
}
