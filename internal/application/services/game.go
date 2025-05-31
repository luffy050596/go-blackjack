package services

import (
	"errors"

	"github.com/luffy050596/go-blackjack/internal/application/dtos"
	"github.com/luffy050596/go-blackjack/internal/domain/entities"
)

// GameApplicationService 游戏应用服务
type GameApplicationService struct {
	game            *entities.Game
	probabilityCalc *ProbabilityCalculator
}

// NewGameApplicationService 创建游戏应用服务
func NewGameApplicationService(playerName string) *GameApplicationService {
	game := entities.NewGame(playerName)
	return &GameApplicationService{
		game:            game,
		probabilityCalc: NewProbabilityCalculator(game.Deck),
	}
}

// StartNewRound 开始新一轮
func (s *GameApplicationService) StartNewRound() error {
	return s.game.StartNewRound()
}

func (s *GameApplicationService) StartDealerTurn() bool {
	if s.game.State != entities.StatePlayerTurn {
		return false
	}
	s.game.State = entities.StateDealerTurn
	return true
}

// GetGameState 获取游戏状态
func (s *GameApplicationService) GetGameState() *dtos.GameStateDTO {
	return &dtos.GameStateDTO{
		RoundNumber: s.game.RoundNumber,
		PlayerChips: s.game.Player.Chips,
		PlayerBet:   s.game.Player.Bet,
		PlayerHand:  convertHandToDTO(s.game.Player.Hand),
		DealerHand:  convertHandToDTO(s.game.Dealer.Hand),
		State:       s.game.State,
		IsGameOver:  s.game.IsGameOver(),
	}
}

// PlaceBet 下注
func (s *GameApplicationService) PlaceBet(amount int) error {
	return s.game.PlaceBet(amount)
}

// DealInitialCards 发初始牌
func (s *GameApplicationService) DealInitialCards() error {
	return s.game.DealInitialCards()
}

// ProcessPlayerAction 处理玩家行动
func (s *GameApplicationService) ProcessPlayerAction(action entities.PlayerAction) (*dtos.ActionResultDTO, error) {
	switch action {
	case entities.ActionHit:
		card, err := s.game.PlayerHit()
		if err != nil {
			return nil, err
		}
		return &dtos.ActionResultDTO{
			Action:   entities.ActionHit,
			Success:  true,
			Continue: !s.game.Player.Hand.IsBust() && !s.game.Player.Hand.IsBlackjack(),
			Card:     convertCardToDTO(card),
		}, nil

	case entities.ActionStand:
		s.game.PlayerStand()
		return &dtos.ActionResultDTO{
			Action:   entities.ActionStand,
			Success:  true,
			Continue: false,
		}, nil

	case entities.ActionDoubleDown:
		card, err := s.game.PlayerDoubleDown()
		if err != nil {
			return nil, err
		}
		return &dtos.ActionResultDTO{
			Action:   entities.ActionDoubleDown,
			Success:  true,
			Continue: false,
			Card:     convertCardToDTO(card),
		}, nil

	case entities.ActionQuit:
		return &dtos.ActionResultDTO{
			Action:   entities.ActionQuit,
			Success:  true,
			Continue: false,
		}, nil

	default:
		return &dtos.ActionResultDTO{
			Action:  entities.ActionInvalid,
			Success: false,
			Message: "无效的输入",
		}, errors.New("invalid input")
	}
}

// ProcessDealerTurn 处理庄家回合
func (s *GameApplicationService) ProcessDealerTurn() error {
	return s.game.DealerTurn()
}

// EvaluateGame 评估游戏结果
func (s *GameApplicationService) EvaluateGame() *dtos.GameResultDTO {
	result := s.game.EvaluateResult()
	if result == nil {
		return nil
	}

	return &dtos.GameResultDTO{
		Type:        result.ResultType,
		BetAmount:   result.BetAmount,
		IsDoubled:   result.IsDoubled,
		PlayerChips: s.game.Player.Chips,
	}
}

// GetBetOptions 获取下注选项
func (s *GameApplicationService) GetBetOptions() []int {
	chips := s.game.Player.Chips
	betOptions := []int{10, 25, 50, 100, 200}
	validOptions := []int{}

	for _, amount := range betOptions {
		if amount <= chips {
			validOptions = append(validOptions, amount)
		}
	}

	// 如果玩家筹码很少，添加全押选项
	if chips < betOptions[0] && chips > 0 {
		validOptions = append(validOptions, chips)
	}

	return validOptions
}

// GetKellyBettingRecommendation 获取凯利公式下注建议
func (s *GameApplicationService) GetKellyBettingRecommendation() *dtos.KellyRecommendationDTO {
	// 基于历史概率或基本策略估算胜率
	// 在没有具体手牌的情况下，使用保守的估算值
	estimatedWinRate := 0.48 // blackjack的基本胜率约为48%
	estimatedLoseRate := 0.52

	// 计算基础凯利比例
	kelly := s.probabilityCalc.CalculateBasicKellyFraction(estimatedWinRate, estimatedLoseRate, s.game.Player.Chips)

	// 转换为DTO
	return &dtos.KellyRecommendationDTO{
		StandardKellyFraction:  kelly.StandardKellyFraction,
		BlackjackKellyFraction: kelly.BlackjackKellyFraction,
		DoubleKellyFraction:    kelly.DoubleKellyFraction,
		RecommendedBetAmount:   kelly.RecommendedBetAmount,
		RecommendedBetFraction: kelly.RecommendedBetFraction,
		ShouldDouble:           kelly.ShouldDouble,
		DoubleExpectedROI:      kelly.DoubleExpectedROI,
		RiskLevel:              kelly.RiskLevel,
		ExpectedGrowthRate:     kelly.ExpectedGrowthRate,
	}
}

// CanPlayerDoubleDown 检查玩家是否可以加倍
func (s *GameApplicationService) CanPlayerDoubleDown() bool {
	return s.game.Player.CanDoubleDown()
}

// IsGameOver 检查游戏是否结束
func (s *GameApplicationService) IsGameOver() bool {
	return s.game.IsGameOver()
}

// CalculateWinProbabilities 计算当前状态下的获胜概率
func (s *GameApplicationService) CalculateWinProbabilities() *dtos.ProbabilityResultDTO {
	// 只在玩家回合计算概率
	if s.game.State != entities.StatePlayerTurn {
		return nil
	}

	// 获取剩余卡牌
	remainingCards := s.game.GetRemainingCards()

	// 计算概率（传递当前筹码）
	result := s.probabilityCalc.CalculateWinProbabilities(
		s.game.Player.Hand,
		s.game.Dealer.Hand,
		remainingCards,
		s.game.Player.Chips,
	)

	// 转换操作分析为DTO
	var actionAnalysisDTO *dtos.ActionAnalysisDTO
	if result.ActionAnalysis != nil {
		// 只在可以加倍时显示凯利公式推荐
		var kellyRecommendationDTO *dtos.KellyRecommendationDTO
		if result.ActionAnalysis.CanDouble && result.ActionAnalysis.KellyRecommendation != nil {
			kelly := result.ActionAnalysis.KellyRecommendation
			kellyRecommendationDTO = &dtos.KellyRecommendationDTO{
				StandardKellyFraction:  kelly.StandardKellyFraction,
				BlackjackKellyFraction: kelly.BlackjackKellyFraction,
				DoubleKellyFraction:    kelly.DoubleKellyFraction,
				RecommendedBetAmount:   kelly.RecommendedBetAmount,
				RecommendedBetFraction: kelly.RecommendedBetFraction,
				ShouldDouble:           kelly.ShouldDouble,
				DoubleExpectedROI:      kelly.DoubleExpectedROI,
				RiskLevel:              kelly.RiskLevel,
				ExpectedGrowthRate:     kelly.ExpectedGrowthRate,
			}
		}

		actionAnalysisDTO = &dtos.ActionAnalysisDTO{
			HitWinRate:          result.ActionAnalysis.HitWinRate,
			StandWinRate:        result.ActionAnalysis.StandWinRate,
			DoubleWinRate:       result.ActionAnalysis.DoubleWinRate,
			SplitWinRate:        result.ActionAnalysis.SplitWinRate,
			CanHit:              result.ActionAnalysis.CanHit,
			CanStand:            result.ActionAnalysis.CanStand,
			CanDouble:           result.ActionAnalysis.CanDouble,
			CanSplit:            result.ActionAnalysis.CanSplit,
			RecommendedAction:   result.ActionAnalysis.RecommendedAction,
			ExpectedValue:       result.ActionAnalysis.ExpectedValue,
			KellyRecommendation: kellyRecommendationDTO,
		}
	}

	// 转换为DTO
	return &dtos.ProbabilityResultDTO{
		PlayerWinProbability:  result.PlayerWinProbability,
		DealerWinProbability:  result.DealerWinProbability,
		PushProbability:       result.PushProbability,
		PlayerBlackjackProb:   result.PlayerBlackjackProb,
		DealerBlackjackProb:   result.DealerBlackjackProb,
		PlayerBustProbability: result.PlayerBustProbability,
		DealerBustProbability: result.DealerBustProbability,
		Player21Probability:   result.Player21Probability,
		Dealer21Probability:   result.Dealer21Probability,
		ActionAnalysis:        actionAnalysisDTO,
	}
}

// 辅助函数：转换Hand到DTO
func convertHandToDTO(hand *entities.Hand) *dtos.HandDTO {
	cards := make([]*dtos.CardDTO, len(hand.Cards))
	for i, card := range hand.Cards {
		cards[i] = convertCardToDTO(card)
	}

	return &dtos.HandDTO{
		Cards: cards,
		Value: hand.Value(),
	}
}

// 辅助函数：转换Card到DTO
func convertCardToDTO(card entities.Card) *dtos.CardDTO {
	return &dtos.CardDTO{
		Suit:  card.Suit.String(),
		Rank:  card.Rank.String(),
		Value: card.Value(),
	}
}
