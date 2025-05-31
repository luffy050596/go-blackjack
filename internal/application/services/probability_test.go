package services

import (
	"math"
	"testing"

	"github.com/luffy050596/go-blackjack/internal/domain/entities"
)

// TestNewProbabilityCalculator 测试创建概率计算器
func TestNewProbabilityCalculator(t *testing.T) {
	t.Parallel()

	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)

	if pc == nil {
		t.Fatal("Expected non-nil ProbabilityCalculator")
	}

	if pc.deck != deck {
		t.Error("Expected deck to be set correctly")
	}

	if pc.trials != 10000 {
		t.Errorf("Expected trials to be 10000, got %d", pc.trials)
	}

	if pc.rng == nil {
		t.Error("Expected rng to be initialized")
	}
}

// TestCalculateWinProbabilities 测试概率计算
func TestCalculateWinProbabilities(t *testing.T) {
	t.Parallel()

	t.Run("player_blackjack", func(t *testing.T) {
		t.Parallel()
		testPlayerBlackjackScenario(t)
	})

	t.Run("normal_hand", func(t *testing.T) {
		t.Parallel()
		testNormalHandScenario(t)
	})

	t.Run("player_21_not_blackjack", func(t *testing.T) {
		t.Parallel()
		testPlayer21NotBlackjackScenario(t)
	})
}

// testPlayerBlackjackScenario 测试玩家Blackjack场景
func testPlayerBlackjackScenario(t *testing.T) {
	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)

	playerCards := []entities.Card{
		{Suit: entities.Hearts, Rank: entities.Ace},
		{Suit: entities.Spades, Rank: entities.King},
	}
	dealerCards := []entities.Card{
		{Suit: entities.Diamonds, Rank: entities.Seven},
		{Suit: entities.Clubs, Rank: entities.Five},
	}

	result := runProbabilityTest(t, pc, playerCards, dealerCards, 1000)

	// 对于Blackjack情况的特殊验证
	if result.Player21Probability != 1.0 {
		t.Errorf("Expected Player21Probability to be 1.0 for blackjack, got %f", result.Player21Probability)
	}
}

// testNormalHandScenario 测试普通手牌场景
func testNormalHandScenario(t *testing.T) {
	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)

	playerCards := []entities.Card{
		{Suit: entities.Hearts, Rank: entities.Eight},
		{Suit: entities.Spades, Rank: entities.Six},
	}
	dealerCards := []entities.Card{
		{Suit: entities.Diamonds, Rank: entities.Ten},
		{Suit: entities.Clubs, Rank: entities.Five},
	}

	runProbabilityTest(t, pc, playerCards, dealerCards, 500)
}

// testPlayer21NotBlackjackScenario 测试21点非Blackjack场景
func testPlayer21NotBlackjackScenario(t *testing.T) {
	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)

	playerCards := []entities.Card{
		{Suit: entities.Hearts, Rank: entities.Seven},
		{Suit: entities.Spades, Rank: entities.Seven},
		{Suit: entities.Diamonds, Rank: entities.Seven},
	}
	dealerCards := []entities.Card{
		{Suit: entities.Diamonds, Rank: entities.Ten},
		{Suit: entities.Clubs, Rank: entities.Five},
	}

	runProbabilityTest(t, pc, playerCards, dealerCards, 200)
}

// runProbabilityTest 运行概率测试的通用函数
func runProbabilityTest(t *testing.T, pc *ProbabilityCalculator, playerCards, dealerCards []entities.Card, currentChips int) *ProbabilityResult {
	// 创建手牌
	playerHand := entities.NewHand()
	for _, card := range playerCards {
		playerHand.AddCard(card)
	}

	dealerHand := entities.NewHand()
	for _, card := range dealerCards {
		dealerHand.AddCard(card)
	}

	// 创建剩余牌堆
	remainingCards := createRemainingCards(playerCards, dealerCards)

	result := pc.CalculateWinProbabilities(playerHand, dealerHand, remainingCards, currentChips)

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	validateProbabilityResult(t, result)
	return result
}

// validateProbabilityResult 验证概率结果的有效性
func validateProbabilityResult(t *testing.T, result *ProbabilityResult) {
	// 验证概率范围 [0, 1]
	if result.PlayerWinProbability < 0 || result.PlayerWinProbability > 1 {
		t.Errorf("PlayerWinProbability out of range [0,1]: %f", result.PlayerWinProbability)
	}

	if result.DealerWinProbability < 0 || result.DealerWinProbability > 1 {
		t.Errorf("DealerWinProbability out of range [0,1]: %f", result.DealerWinProbability)
	}

	if result.PushProbability < 0 || result.PushProbability > 1 {
		t.Errorf("PushProbability out of range [0,1]: %f", result.PushProbability)
	}

	// 验证概率总和接近1 (允许一定误差)
	total := result.PlayerWinProbability + result.DealerWinProbability + result.PushProbability
	if math.Abs(total-1.0) > 0.05 {
		t.Errorf("Total probability should be close to 1.0, got %f", total)
	}

	// 验证操作分析
	if result.ActionAnalysis == nil {
		t.Error("Expected ActionAnalysis to be non-nil")
	}
}

// TestCalculateBasicKellyFraction 测试基础凯利公式计算
func TestCalculateBasicKellyFraction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		winRate      float64
		loseRate     float64
		currentChips int
		expectLow    bool
	}{
		{
			name:         "high_chips_conservative",
			winRate:      0.45,
			loseRate:     0.55,
			currentChips: 2000,
			expectLow:    true,
		},
		{
			name:         "medium_chips",
			winRate:      0.48,
			loseRate:     0.52,
			currentChips: 750,
			expectLow:    true,
		},
		{
			name:         "low_chips",
			winRate:      0.50,
			loseRate:     0.50,
			currentChips: 150,
			expectLow:    false,
		},
		{
			name:         "very_low_chips",
			winRate:      0.40,
			loseRate:     0.60,
			currentChips: 50,
			expectLow:    false,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deck := entities.NewDeck()
			pc := NewProbabilityCalculator(deck)

			result := pc.CalculateBasicKellyFraction(tt.winRate, tt.loseRate, tt.currentChips)

			validateKellyResult(t, result, tt.currentChips, tt.expectLow)
		})
	}
}

// validateKellyResult 验证凯利公式结果
func validateKellyResult(t *testing.T, result *KellyRecommendation, currentChips int, expectLow bool) {
	if result == nil {
		t.Fatal("Expected non-nil KellyRecommendation")
	}

	// 验证推荐投注金额
	if result.RecommendedBetAmount < 0 {
		t.Error("RecommendedBetAmount should not be negative")
	}

	if result.RecommendedBetAmount > currentChips {
		t.Errorf("RecommendedBetAmount (%d) should not exceed currentChips (%d)",
			result.RecommendedBetAmount, currentChips)
	}

	// 验证推荐投注比例
	if result.RecommendedBetFraction < 0 || result.RecommendedBetFraction > 1 {
		t.Errorf("RecommendedBetFraction out of range [0,1]: %f", result.RecommendedBetFraction)
	}

	// 验证风险等级
	validRiskLevels := []string{"Low", "Medium", "High"}
	found := false
	for _, level := range validRiskLevels {
		if result.RiskLevel == level {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Invalid RiskLevel: %s", result.RiskLevel)
	}

	// 验证保守策略
	if expectLow && result.RecommendedBetFraction > 0.02 {
		t.Errorf("Expected conservative recommendation for high chips, got %f", result.RecommendedBetFraction)
	}
}

// TestHitAnalysis 测试要牌分析
func TestHitAnalysis(t *testing.T) {
	t.Parallel()

	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)

	tests := []struct {
		name           string
		currentValue   int
		remainingCards []entities.Card
		expectBust     bool
	}{
		{
			name:         "safe_hit",
			currentValue: 10,
			remainingCards: []entities.Card{
				{Suit: entities.Hearts, Rank: entities.Two},
				{Suit: entities.Spades, Rank: entities.Three},
			},
			expectBust: false,
		},
		{
			name:         "risky_hit",
			currentValue: 20,
			remainingCards: []entities.Card{
				{Suit: entities.Hearts, Rank: entities.King},
				{Suit: entities.Spades, Rank: entities.Queen},
			},
			expectBust: true,
		},
		{
			name:           "empty_deck",
			currentValue:   15,
			remainingCards: []entities.Card{},
			expectBust:     false,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			analysis := pc.calculateHitAnalysis(tt.currentValue, tt.remainingCards)
			validateHitAnalysis(t, analysis, tt.remainingCards, tt.expectBust)
		})
	}
}

// validateHitAnalysis 验证要牌分析结果
func validateHitAnalysis(t *testing.T, analysis *HitAnalysis, remainingCards []entities.Card, expectBust bool) {
	if analysis == nil {
		t.Fatal("Expected non-nil HitAnalysis")
	}

	// 验证概率范围
	if analysis.BustProbability < 0 || analysis.BustProbability > 1 {
		t.Errorf("BustProbability out of range [0,1]: %f", analysis.BustProbability)
	}

	if analysis.Hit21Probability < 0 || analysis.Hit21Probability > 1 {
		t.Errorf("Hit21Probability out of range [0,1]: %f", analysis.Hit21Probability)
	}

	if analysis.SafeHitProbability < 0 || analysis.SafeHitProbability > 1 {
		t.Errorf("SafeHitProbability out of range [0,1]: %f", analysis.SafeHitProbability)
	}

	// 验证概率总和
	total := analysis.BustProbability + analysis.Hit21Probability + analysis.SafeHitProbability
	if len(remainingCards) > 0 && math.Abs(total-1.0) > 0.01 {
		t.Errorf("Total probability should be close to 1.0, got %f", total)
	}

	// 验证爆牌概率期望
	if expectBust && analysis.BustProbability == 0 {
		t.Error("Expected non-zero bust probability for risky hit")
	}
}

// TestActionAnalysis 测试操作分析
func TestActionAnalysis(t *testing.T) {
	t.Parallel()

	t.Run("first_turn_different_cards", func(t *testing.T) {
		t.Parallel()
		testActionAnalysisFirstTurnDifferent(t)
	})

	t.Run("first_turn_pair", func(t *testing.T) {
		t.Parallel()
		testActionAnalysisFirstTurnPair(t)
	})

	t.Run("after_hit", func(t *testing.T) {
		t.Parallel()
		testActionAnalysisAfterHit(t)
	})

	t.Run("player_21", func(t *testing.T) {
		t.Parallel()
		testActionAnalysisPlayer21(t)
	})
}

// testActionAnalysisFirstTurnDifferent 测试第一轮不同牌的操作分析
func testActionAnalysisFirstTurnDifferent(t *testing.T) {
	playerCards := []entities.Card{
		{Suit: entities.Hearts, Rank: entities.Eight},
		{Suit: entities.Spades, Rank: entities.Six},
	}
	dealerCards := []entities.Card{
		{Suit: entities.Diamonds, Rank: entities.Ten},
	}

	analysis := runActionAnalysisTest(t, playerCards, dealerCards)

	validateActionCapabilities(t, analysis, true, true, false)
}

// testActionAnalysisFirstTurnPair 测试第一轮对子的操作分析
func testActionAnalysisFirstTurnPair(t *testing.T) {
	playerCards := []entities.Card{
		{Suit: entities.Hearts, Rank: entities.Eight},
		{Suit: entities.Spades, Rank: entities.Eight},
	}
	dealerCards := []entities.Card{
		{Suit: entities.Diamonds, Rank: entities.Six},
	}

	analysis := runActionAnalysisTest(t, playerCards, dealerCards)

	validateActionCapabilities(t, analysis, true, true, true)
}

// testActionAnalysisAfterHit 测试要牌后的操作分析
func testActionAnalysisAfterHit(t *testing.T) {
	playerCards := []entities.Card{
		{Suit: entities.Hearts, Rank: entities.Eight},
		{Suit: entities.Spades, Rank: entities.Six},
		{Suit: entities.Diamonds, Rank: entities.Three},
	}
	dealerCards := []entities.Card{
		{Suit: entities.Clubs, Rank: entities.Seven},
	}

	analysis := runActionAnalysisTest(t, playerCards, dealerCards)

	validateActionCapabilities(t, analysis, true, false, false)
}

// testActionAnalysisPlayer21 测试玩家21点的操作分析
func testActionAnalysisPlayer21(t *testing.T) {
	playerCards := []entities.Card{
		{Suit: entities.Hearts, Rank: entities.Ten},
		{Suit: entities.Spades, Rank: entities.Ace},
	}
	dealerCards := []entities.Card{
		{Suit: entities.Diamonds, Rank: entities.Seven},
	}

	analysis := runActionAnalysisTest(t, playerCards, dealerCards)

	validateActionCapabilities(t, analysis, false, false, false)

	if analysis.RecommendedAction != "stand" {
		t.Errorf("Expected RecommendedAction stand, got %s", analysis.RecommendedAction)
	}
}

// runActionAnalysisTest 运行操作分析测试的通用函数
func runActionAnalysisTest(t *testing.T, playerCards, dealerCards []entities.Card) *ActionAnalysis {
	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)

	// 创建手牌
	playerHand := entities.NewHand()
	for _, card := range playerCards {
		playerHand.AddCard(card)
	}

	dealerHand := entities.NewHand()
	for _, card := range dealerCards {
		dealerHand.AddCard(card)
	}

	// 创建剩余牌堆
	remainingCards := createRemainingCards(playerCards, dealerCards)

	analysis := pc.calculateActionAnalysis(playerHand, dealerHand, remainingCards)

	if analysis == nil {
		t.Fatal("Expected non-nil ActionAnalysis")
	}

	validateActionAnalysisResult(t, analysis)
	return analysis
}

// validateActionCapabilities 验证操作可用性
func validateActionCapabilities(t *testing.T, analysis *ActionAnalysis, canHit, canDouble, canSplit bool) {
	if analysis.CanHit != canHit {
		t.Errorf("Expected CanHit %v, got %v", canHit, analysis.CanHit)
	}

	if analysis.CanStand != true {
		t.Errorf("Expected CanStand %v, got %v", true, analysis.CanStand)
	}

	if analysis.CanDouble != canDouble {
		t.Errorf("Expected CanDouble %v, got %v", canDouble, analysis.CanDouble)
	}

	if analysis.CanSplit != canSplit {
		t.Errorf("Expected CanSplit %v, got %v", canSplit, analysis.CanSplit)
	}
}

// validateActionAnalysisResult 验证操作分析结果
func validateActionAnalysisResult(t *testing.T, analysis *ActionAnalysis) {
	// 验证胜率范围
	if analysis.StandWinRate < 0 || analysis.StandWinRate > 1 {
		t.Errorf("StandWinRate out of range [0,1]: %f", analysis.StandWinRate)
	}

	if analysis.CanHit && (analysis.HitWinRate < 0 || analysis.HitWinRate > 1) {
		t.Errorf("HitWinRate out of range [0,1]: %f", analysis.HitWinRate)
	}

	// 验证期望值
	if analysis.ExpectedValue < 0 || analysis.ExpectedValue > 1 {
		t.Errorf("ExpectedValue out of range [0,1]: %f", analysis.ExpectedValue)
	}
}

// TestKellyFractionCalculation 测试凯利公式计算
func TestKellyFractionCalculation(t *testing.T) {
	t.Parallel()

	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)

	tests := []struct {
		name     string
		winProb  float64
		loseProb float64
		odds     float64
		expect   float64
	}{
		{
			name:     "favorable_odds",
			winProb:  0.6,
			loseProb: 0.4,
			odds:     1.0,
			expect:   0.2, // (1.0*0.6 - 0.4) / 1.0 = 0.2
		},
		{
			name:     "unfavorable_odds",
			winProb:  0.4,
			loseProb: 0.6,
			odds:     1.0,
			expect:   0.0, // Should be 0 due to negative result
		},
		{
			name:     "even_odds",
			winProb:  0.5,
			loseProb: 0.5,
			odds:     1.0,
			expect:   0.0, // (1.0*0.5 - 0.5) / 1.0 = 0.0
		},
		{
			name:     "blackjack_odds",
			winProb:  0.5,
			loseProb: 0.5,
			odds:     1.5,
			expect:   0.25, // (1.5*0.5 - 0.5) / 1.5 = 0.25/1.5 ≈ 0.167
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := pc.calculateKellyFraction(tt.winProb, tt.loseProb, tt.odds)

			if result < 0 {
				t.Error("Kelly fraction should not be negative")
			}

			// 对于有利赔率，验证结果接近期望值
			if tt.expect > 0 && math.Abs(result-tt.expect) > 0.1 {
				t.Errorf("Expected Kelly fraction close to %f, got %f", tt.expect, result)
			}

			// 对于不利赔率，验证结果为0或接近0
			if tt.expect == 0 && result > 0.01 {
				t.Errorf("Expected Kelly fraction to be 0 for unfavorable odds, got %f", result)
			}
		})
	}
}

// TestPlayerWins 测试玩家获胜判断
func TestPlayerWins(t *testing.T) {
	t.Parallel()

	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)

	tests := []struct {
		name            string
		playerCards     []entities.Card
		dealerCards     []entities.Card
		expectPlayerWin bool
	}{
		{
			name: "player_blackjack_vs_dealer_20",
			playerCards: []entities.Card{
				{Suit: entities.Hearts, Rank: entities.Ace},
				{Suit: entities.Spades, Rank: entities.King},
			},
			dealerCards: []entities.Card{
				{Suit: entities.Diamonds, Rank: entities.Ten},
				{Suit: entities.Clubs, Rank: entities.King},
			},
			expectPlayerWin: true,
		},
		{
			name: "both_blackjack",
			playerCards: []entities.Card{
				{Suit: entities.Hearts, Rank: entities.Ace},
				{Suit: entities.Spades, Rank: entities.King},
			},
			dealerCards: []entities.Card{
				{Suit: entities.Diamonds, Rank: entities.Ace},
				{Suit: entities.Clubs, Rank: entities.Queen},
			},
			expectPlayerWin: false, // Push, not a win
		},
		{
			name: "player_bust",
			playerCards: []entities.Card{
				{Suit: entities.Hearts, Rank: entities.Ten},
				{Suit: entities.Spades, Rank: entities.Eight},
				{Suit: entities.Diamonds, Rank: entities.Five},
			},
			dealerCards: []entities.Card{
				{Suit: entities.Clubs, Rank: entities.Ten},
				{Suit: entities.Hearts, Rank: entities.Seven},
			},
			expectPlayerWin: false,
		},
		{
			name: "dealer_bust",
			playerCards: []entities.Card{
				{Suit: entities.Hearts, Rank: entities.Ten},
				{Suit: entities.Spades, Rank: entities.Eight},
			},
			dealerCards: []entities.Card{
				{Suit: entities.Clubs, Rank: entities.Ten},
				{Suit: entities.Hearts, Rank: entities.Seven},
				{Suit: entities.Diamonds, Rank: entities.Five},
			},
			expectPlayerWin: true,
		},
		{
			name: "player_higher_value",
			playerCards: []entities.Card{
				{Suit: entities.Hearts, Rank: entities.Ten},
				{Suit: entities.Spades, Rank: entities.Nine},
			},
			dealerCards: []entities.Card{
				{Suit: entities.Clubs, Rank: entities.Ten},
				{Suit: entities.Hearts, Rank: entities.Seven},
			},
			expectPlayerWin: true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// 创建手牌
			playerHand := entities.NewHand()
			for _, card := range tt.playerCards {
				playerHand.AddCard(card)
			}

			dealerHand := entities.NewHand()
			for _, card := range tt.dealerCards {
				dealerHand.AddCard(card)
			}

			result := pc.playerWins(playerHand, dealerHand)

			if result != tt.expectPlayerWin {
				t.Errorf("Expected playerWins %v, got %v", tt.expectPlayerWin, result)
			}
		})
	}
}

// TestCopyHand 测试手牌复制
func TestCopyHand(t *testing.T) {
	t.Parallel()

	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)

	// 创建原始手牌
	originalHand := entities.NewHand()
	originalHand.AddCard(entities.Card{Suit: entities.Hearts, Rank: entities.Ace})
	originalHand.AddCard(entities.Card{Suit: entities.Spades, Rank: entities.King})

	// 复制手牌
	copiedHand := pc.copyHand(originalHand)

	// 验证复制结果
	if copiedHand == nil {
		t.Fatal("Expected non-nil copied hand")
	}

	if len(copiedHand.Cards) != len(originalHand.Cards) {
		t.Errorf("Expected same number of cards, original: %d, copied: %d",
			len(originalHand.Cards), len(copiedHand.Cards))
	}

	// 验证牌是否相同
	for i, card := range originalHand.Cards {
		if copiedHand.Cards[i] != card {
			t.Errorf("Card %d mismatch: original %v, copied %v", i, card, copiedHand.Cards[i])
		}
	}

	// 验证是否为深拷贝（修改复制的手牌不影响原手牌）
	copiedHand.AddCard(entities.Card{Suit: entities.Diamonds, Rank: entities.Two})

	if len(originalHand.Cards) == len(copiedHand.Cards) {
		t.Error("Expected deep copy - modifying copied hand should not affect original")
	}
}

// TestRemoveCard 测试移除牌
func TestRemoveCard(t *testing.T) {
	t.Parallel()

	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)

	// 创建牌堆
	cards := []entities.Card{
		{Suit: entities.Hearts, Rank: entities.Ace},
		{Suit: entities.Spades, Rank: entities.King},
		{Suit: entities.Diamonds, Rank: entities.Queen},
	}

	cardToRemove := entities.Card{Suit: entities.Spades, Rank: entities.King}

	result := pc.removeCard(cards, cardToRemove)

	// 验证移除结果
	if len(result) != len(cards)-1 {
		t.Errorf("Expected %d cards after removal, got %d", len(cards)-1, len(result))
	}

	// 验证指定牌已被移除
	for _, card := range result {
		if card.Suit == cardToRemove.Suit && card.Rank == cardToRemove.Rank {
			t.Error("Card should have been removed")
		}
	}

	// 测试移除不存在的牌
	nonExistentCard := entities.Card{Suit: entities.Clubs, Rank: entities.Two}
	result2 := pc.removeCard(cards, nonExistentCard)

	if len(result2) != len(cards) {
		t.Error("Removing non-existent card should not change the slice length")
	}
}

// TestAssessRiskLevel 测试风险等级评估
func TestAssessRiskLevel(t *testing.T) {
	t.Parallel()

	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)

	tests := []struct {
		name          string
		kellyFraction float64
		expectedRisk  string
	}{
		{
			name:          "low_risk",
			kellyFraction: 0.01,
			expectedRisk:  "Low",
		},
		{
			name:          "medium_risk",
			kellyFraction: 0.03,
			expectedRisk:  "Medium",
		},
		{
			name:          "high_risk",
			kellyFraction: 0.08,
			expectedRisk:  "High",
		},
		{
			name:          "zero_risk",
			kellyFraction: 0.0,
			expectedRisk:  "Low",
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := pc.assessRiskLevel(tt.kellyFraction)

			if result != tt.expectedRisk {
				t.Errorf("Expected risk level %s, got %s", tt.expectedRisk, result)
			}
		})
	}
}

// TestSimulateGame 测试游戏模拟
func TestSimulateGame(t *testing.T) {
	t.Parallel()

	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)

	// 创建测试手牌
	playerHand := entities.NewHand()
	playerHand.AddCard(entities.Card{Suit: entities.Hearts, Rank: entities.Eight})
	playerHand.AddCard(entities.Card{Suit: entities.Spades, Rank: entities.Six})

	dealerHand := entities.NewHand()
	dealerHand.AddCard(entities.Card{Suit: entities.Diamonds, Rank: entities.Ten})

	// 创建剩余牌堆
	remainingCards := []entities.Card{
		{Suit: entities.Clubs, Rank: entities.Five},
		{Suit: entities.Hearts, Rank: entities.Seven},
		{Suit: entities.Spades, Rank: entities.Nine},
	}

	result := pc.simulateGame(playerHand, dealerHand, remainingCards)

	if result == nil {
		t.Fatal("Expected non-nil simulation result")
	}

	// 验证结果字段
	validWinners := []string{"player", "dealer", "push"}
	found := false
	for _, winner := range validWinners {
		if result.Winner == winner {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Invalid winner: %s", result.Winner)
	}

	if result.PlayerFinalValue < 0 || result.PlayerFinalValue > 30 {
		t.Errorf("Invalid PlayerFinalValue: %d", result.PlayerFinalValue)
	}

	if result.DealerFinalValue < 0 || result.DealerFinalValue > 30 {
		t.Errorf("Invalid DealerFinalValue: %d", result.DealerFinalValue)
	}
}

// createRemainingCards 创建剩余牌堆（排除已使用的牌）
func createRemainingCards(playerCards, dealerCards []entities.Card) []entities.Card {
	deck := entities.NewDeck()
	used := make(map[string]bool)

	// 标记已使用的牌
	for _, card := range playerCards {
		key := cardKey(card)
		used[key] = true
	}
	for _, card := range dealerCards {
		key := cardKey(card)
		used[key] = true
	}

	// 创建剩余牌堆
	var remaining []entities.Card
	for _, card := range deck.Cards {
		key := cardKey(card)
		if !used[key] {
			remaining = append(remaining, card)
		}
	}

	return remaining
}

// cardKey 创建牌的唯一标识
func cardKey(card entities.Card) string {
	return string(rune(card.Suit)) + string(rune(card.Rank))
}
