package main

import "strings"

// GameService 游戏业务逻辑服务
type GameService struct {
	game    *Game
	display *Display
}

// NewGameService 创建游戏服务
func NewGameService(game *Game) *GameService {
	return &GameService{
		game:    game,
		display: NewDisplay(game),
	}
}

// PrepareBetting 准备下注阶段
func (gs *GameService) PrepareBetting() []int {
	betOptions := []int{10, 25, 50, 100, 200}
	validOptions := []int{}

	for _, amount := range betOptions {
		if amount <= gs.game.Player.Chips {
			validOptions = append(validOptions, amount)
		}
	}

	// 如果玩家筹码很少，添加全押选项
	if gs.game.Player.Chips < betOptions[0] && gs.game.Player.Chips > 0 {
		validOptions = append(validOptions, gs.game.Player.Chips)
	}

	return validOptions
}

// ProcessBet 处理下注
func (gs *GameService) ProcessBet(choice int, validOptions []int) (bool, bool) {
	exitIndex := len(validOptions) + 1

	if choice == exitIndex {
		return false, false // 退出游戏
	}

	if choice < 1 || choice > len(validOptions) {
		return false, true // 无效选择，继续游戏
	}

	betAmount := validOptions[choice-1]
	if gs.game.Player.PlaceBet(betAmount) {
		gs.display.showBetSuccess(betAmount, gs.game.Player.Chips)
		return true, true // 下注成功，继续游戏
	}

	return false, true // 下注失败，继续游戏
}

// ParsePlayerInput 解析玩家输入
func (gs *GameService) ParsePlayerInput(input string) PlayerAction {
	switch strings.ToLower(input) {
	case InputHit, InputHitFull:
		return ActionHit
	case InputStand, InputStandFull:
		return ActionStand
	case InputDouble, InputDoubleFull, InputDoubleDown:
		return ActionDoubleDown
	case InputQuit, InputQuitFull:
		return ActionQuit
	default:
		return -1 // 无效输入
	}
}

// ProcessPlayerAction 处理玩家行动
func (gs *GameService) ProcessPlayerAction(action PlayerAction) ActionResult {
	switch action {
	case ActionHit:
		card := gs.game.Deck.Deal()
		gs.game.Player.Hand.AddCard(card)
		gs.display.showCardDealt(card, true)
		return ActionResult{Action: action, IsValid: true, Continue: true}

	case ActionStand:
		gs.display.showPlayerChoice(DisplayActionStand)
		return ActionResult{Action: action, IsValid: true, Continue: false}

	case ActionDoubleDown:
		if gs.game.Player.CanDoubleDown() {
			if gs.game.Player.DoubleBet() {
				gs.display.showDoubleDownSuccess(gs.game.Player.Bet, gs.game.Player.Chips)
				return ActionResult{Action: action, IsValid: true, Continue: false}
			} else {
				gs.display.showError("筹码不足，无法加倍")
				return ActionResult{Action: action, IsValid: false, Continue: true}
			}
		} else {
			gs.display.showError("现在无法加倍")
			return ActionResult{Action: action, IsValid: false, Continue: true}
		}

	case ActionQuit:
		gs.display.showInfo("感谢游戏！再见！")
		return ActionResult{Action: action, IsValid: true, Continue: false}

	default:
		gs.display.showError("无效选择，请重新输入")
		return ActionResult{Action: action, IsValid: false, Continue: true}
	}
}

// EvaluateGame 评估游戏结果
func (gs *GameService) EvaluateGame() GameResult {
	player := gs.game.Player
	dealer := gs.game.Dealer

	result := GameResult{
		BetAmount: player.Bet,
		IsDoubled: player.DoubledDown,
	}

	playerValue := player.Hand.Value()
	dealerValue := dealer.Hand.Value()
	playerBlackjack := player.Hand.IsBlackjack()
	dealerBlackjack := dealer.Hand.IsBlackjack()

	// 判断游戏结果
	if player.Hand.IsBust() {
		result.ResultType = PlayerBust
		player.LoseBet()
	} else if dealer.Hand.IsBust() {
		result.ResultType = DealerBust
		player.WinBet(1.0)
	} else if playerBlackjack && dealerBlackjack {
		result.ResultType = BothBlackjack
		player.PushBet()
	} else if playerBlackjack {
		result.ResultType = PlayerBlackjack
		if player.DoubledDown {
			player.WinBet(1.0) // 加倍后Blackjack按1:1赔率
		} else {
			player.WinBet(1.5) // 3:2 赔率
		}
	} else if dealerBlackjack {
		result.ResultType = DealerBlackjack
		player.LoseBet()
	} else if playerValue > dealerValue {
		result.ResultType = PlayerWin
		player.WinBet(1.0)
	} else if playerValue < dealerValue {
		result.ResultType = DealerWin
		player.LoseBet()
	} else {
		result.ResultType = Push
		player.PushBet()
	}

	return result
}

// ShouldDealerHit 判断庄家是否应该要牌
func (gs *GameService) ShouldDealerHit() bool {
	return gs.game.Dealer.Hand.Value() < 17
}

// DealerHit 庄家要牌
func (gs *GameService) DealerHit() {
	gs.display.showDealerAction(DisplayActionHit)
	card := gs.game.Deck.Deal()
	gs.game.Dealer.Hand.AddCard(card)
	gs.display.showCardDealt(card, false)
}

// DealerStand 庄家停牌
func (gs *GameService) DealerStand() {
	gs.display.showDealerAction(DisplayActionStand)
}

// CheckDealerBust 检查庄家是否爆牌
func (gs *GameService) CheckDealerBust() bool {
	if gs.game.Dealer.Hand.IsBust() {
		gs.display.showDealerAction(DisplayActionBust)
		return true
	}
	return false
}

// GetPromptText 获取提示文本
func (gs *GameService) GetPromptText() string {
	prompt := "请选择: (h)要牌 (s)停牌"
	if gs.game.Player.CanDoubleDown() {
		prompt += " (d)加倍"
	}
	prompt += " (q)退出游戏: "
	return prompt
}
