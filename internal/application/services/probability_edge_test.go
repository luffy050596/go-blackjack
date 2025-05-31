package services

import (
	"testing"

	"github.com/luffy050596/go-blackjack/internal/domain/entities"
)

// TestEdgeCases 测试边界情况
func TestEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("empty_remaining_cards", func(t *testing.T) {
		t.Parallel()

		deck := entities.NewDeck()
		pc := NewProbabilityCalculator(deck)

		playerHand := entities.NewHand()
		playerHand.AddCard(entities.Card{Suit: entities.Hearts, Rank: entities.Eight})
		playerHand.AddCard(entities.Card{Suit: entities.Spades, Rank: entities.Six})

		dealerHand := entities.NewHand()
		dealerHand.AddCard(entities.Card{Suit: entities.Diamonds, Rank: entities.Ten})

		remainingCards := []entities.Card{} // 空牌堆

		result := pc.CalculateWinProbabilities(playerHand, dealerHand, remainingCards, 1000)

		if result == nil {
			t.Fatal("Expected non-nil result even with empty remaining cards")
		}

		// 验证结果仍然有意义
		if result.ActionAnalysis == nil {
			t.Error("Expected ActionAnalysis to be non-nil")
		}
	})

	t.Run("dealer_empty_hand", func(t *testing.T) {
		t.Parallel()

		deck := entities.NewDeck()
		pc := NewProbabilityCalculator(deck)

		playerHand := entities.NewHand()
		playerHand.AddCard(entities.Card{Suit: entities.Hearts, Rank: entities.Eight})
		playerHand.AddCard(entities.Card{Suit: entities.Spades, Rank: entities.Six})

		dealerHand := entities.NewHand() // 空手牌

		remainingCards := deck.Cards[:10]

		result := pc.CalculateWinProbabilities(playerHand, dealerHand, remainingCards, 1000)

		if result == nil {
			t.Fatal("Expected non-nil result even with empty dealer hand")
		}
	})

	t.Run("player_bust_hand", func(t *testing.T) {
		t.Parallel()

		deck := entities.NewDeck()
		pc := NewProbabilityCalculator(deck)

		playerHand := entities.NewHand()
		playerHand.AddCard(entities.Card{Suit: entities.Hearts, Rank: entities.Ten})
		playerHand.AddCard(entities.Card{Suit: entities.Spades, Rank: entities.Eight})
		playerHand.AddCard(entities.Card{Suit: entities.Diamonds, Rank: entities.Five}) // 总和23，爆牌

		dealerHand := entities.NewHand()
		dealerHand.AddCard(entities.Card{Suit: entities.Clubs, Rank: entities.Ten})

		remainingCards := deck.Cards[:10]

		result := pc.CalculateWinProbabilities(playerHand, dealerHand, remainingCards, 1000)

		if result == nil {
			t.Fatal("Expected non-nil result even with bust hand")
		}

		// 玩家爆牌，胜率应该为0
		if result.PlayerWinProbability > 0.05 { // 允许一些模拟误差
			t.Errorf("Expected very low win probability for bust hand, got %f", result.PlayerWinProbability)
		}
	})

	t.Run("very_low_chips", func(t *testing.T) {
		t.Parallel()

		deck := entities.NewDeck()
		pc := NewProbabilityCalculator(deck)

		result := pc.CalculateBasicKellyFraction(0.48, 0.52, 5) // 很少的筹码

		if result == nil {
			t.Fatal("Expected non-nil result even with very low chips")
		}

		if result.RecommendedBetAmount > 5 {
			t.Error("Recommended bet should not exceed available chips")
		}

		if result.RiskLevel != "High" {
			t.Errorf("Expected High risk level for very low chips, got %s", result.RiskLevel)
		}
	})

	t.Run("zero_chips", func(t *testing.T) {
		t.Parallel()

		deck := entities.NewDeck()
		pc := NewProbabilityCalculator(deck)

		result := pc.CalculateBasicKellyFraction(0.48, 0.52, 0) // 没有筹码

		if result == nil {
			t.Fatal("Expected non-nil result even with zero chips")
		}

		if result.RecommendedBetAmount != 0 {
			t.Error("Recommended bet should be 0 for zero chips")
		}
	})
}

// TestKellyCalculationEdgeCases 测试凯利公式计算的边界情况
func TestKellyCalculationEdgeCases(t *testing.T) {
	t.Parallel()

	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)

	tests := []struct {
		name     string
		winProb  float64
		loseProb float64
		odds     float64
		expect   string
	}{
		{
			name:     "zero_win_probability",
			winProb:  0.0,
			loseProb: 1.0,
			odds:     1.0,
			expect:   "zero_result",
		},
		{
			name:     "zero_odds",
			winProb:  0.5,
			loseProb: 0.5,
			odds:     0.0,
			expect:   "zero_result",
		},
		{
			name:     "negative_odds",
			winProb:  0.5,
			loseProb: 0.5,
			odds:     -1.0,
			expect:   "zero_result",
		},
		{
			name:     "extreme_win_probability",
			winProb:  0.99,
			loseProb: 0.01,
			odds:     1.0,
			expect:   "high_result",
		},
		{
			name:     "very_high_odds",
			winProb:  0.5,
			loseProb: 0.5,
			odds:     10.0,
			expect:   "high_result",
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := pc.calculateKellyFraction(tt.winProb, tt.loseProb, tt.odds)

			// 结果不应该为负数
			if result < 0 {
				t.Error("Kelly fraction should not be negative")
			}

			switch tt.expect {
			case "zero_result":
				if result > 0.01 {
					t.Errorf("Expected near-zero result, got %f", result)
				}
			case "high_result":
				if result < 0.1 {
					t.Errorf("Expected high result, got %f", result)
				}
			}
		})
	}
}

// TestActionAnalysisEdgeCases 测试操作分析的边界情况
func TestActionAnalysisEdgeCases(t *testing.T) {
	t.Parallel()

	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)

	t.Run("single_card_remaining", func(t *testing.T) {
		t.Parallel()

		playerHand := entities.NewHand()
		playerHand.AddCard(entities.Card{Suit: entities.Hearts, Rank: entities.Eight})
		playerHand.AddCard(entities.Card{Suit: entities.Spades, Rank: entities.Six})

		dealerHand := entities.NewHand()
		dealerHand.AddCard(entities.Card{Suit: entities.Diamonds, Rank: entities.Ten})

		remainingCards := []entities.Card{
			{Suit: entities.Clubs, Rank: entities.Five}, // 只有一张牌
		}

		analysis := pc.calculateActionAnalysis(playerHand, dealerHand, remainingCards)

		if analysis == nil {
			t.Fatal("Expected non-nil ActionAnalysis")
		}

		// 验证基本逻辑仍然有效
		if !analysis.CanStand {
			t.Error("Should always be able to stand")
		}
	})

	t.Run("ace_heavy_remaining_cards", func(t *testing.T) {
		t.Parallel()

		playerHand := entities.NewHand()
		playerHand.AddCard(entities.Card{Suit: entities.Hearts, Rank: entities.Ten})
		playerHand.AddCard(entities.Card{Suit: entities.Spades, Rank: entities.Five})

		dealerHand := entities.NewHand()
		dealerHand.AddCard(entities.Card{Suit: entities.Diamonds, Rank: entities.Six})

		// 剩余都是A牌
		remainingCards := []entities.Card{
			{Suit: entities.Hearts, Rank: entities.Ace},
			{Suit: entities.Spades, Rank: entities.Ace},
			{Suit: entities.Diamonds, Rank: entities.Ace},
			{Suit: entities.Clubs, Rank: entities.Ace},
		}

		analysis := pc.calculateActionAnalysis(playerHand, dealerHand, remainingCards)

		if analysis == nil {
			t.Fatal("Expected non-nil ActionAnalysis")
		}

		// A牌对玩家有利（15 + 1 = 16 或 15 + 11 = 26）
		// 所以要牌胜率应该合理
		if analysis.HitWinRate < 0 || analysis.HitWinRate > 1 {
			t.Errorf("HitWinRate out of range: %f", analysis.HitWinRate)
		}
	})

	t.Run("ten_heavy_remaining_cards", func(t *testing.T) {
		t.Parallel()

		playerHand := entities.NewHand()
		playerHand.AddCard(entities.Card{Suit: entities.Hearts, Rank: entities.Eight})
		playerHand.AddCard(entities.Card{Suit: entities.Spades, Rank: entities.Eight})

		dealerHand := entities.NewHand()
		dealerHand.AddCard(entities.Card{Suit: entities.Diamonds, Rank: entities.Six})

		// 剩余都是10点牌
		remainingCards := []entities.Card{
			{Suit: entities.Hearts, Rank: entities.Ten},
			{Suit: entities.Spades, Rank: entities.Jack},
			{Suit: entities.Diamonds, Rank: entities.Queen},
			{Suit: entities.Clubs, Rank: entities.King},
		}

		analysis := pc.calculateActionAnalysis(playerHand, dealerHand, remainingCards)

		if analysis == nil {
			t.Fatal("Expected non-nil ActionAnalysis")
		}

		// 16点要10点牌会爆牌，所以要牌胜率应该很低
		if analysis.HitWinRate > 0.1 {
			t.Errorf("Expected low HitWinRate with ten-heavy deck, got %f", analysis.HitWinRate)
		}
	})
}

// TestSpecialHandScenarios 测试特殊手牌场景
func TestSpecialHandScenarios(t *testing.T) {
	t.Parallel()

	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)

	t.Run("soft_ace_scenarios", func(t *testing.T) {
		t.Parallel()

		// 软17（A-6）
		playerHand := entities.NewHand()
		playerHand.AddCard(entities.Card{Suit: entities.Hearts, Rank: entities.Ace})
		playerHand.AddCard(entities.Card{Suit: entities.Spades, Rank: entities.Six})

		dealerHand := entities.NewHand()
		dealerHand.AddCard(entities.Card{Suit: entities.Diamonds, Rank: entities.Ten})

		remainingCards := deck.Cards[:20]

		analysis := pc.calculateActionAnalysis(playerHand, dealerHand, remainingCards)

		if analysis == nil {
			t.Fatal("Expected non-nil ActionAnalysis")
		}

		// 软17通常应该要牌
		if analysis.RecommendedAction == "stand" {
			t.Error("Soft 17 should usually recommend hit")
		}
	})

	t.Run("hard_17_scenarios", func(t *testing.T) {
		t.Parallel()

		// 硬17（10-7）
		playerHand := entities.NewHand()
		playerHand.AddCard(entities.Card{Suit: entities.Hearts, Rank: entities.Ten})
		playerHand.AddCard(entities.Card{Suit: entities.Spades, Rank: entities.Seven})

		dealerHand := entities.NewHand()
		dealerHand.AddCard(entities.Card{Suit: entities.Diamonds, Rank: entities.Six})

		remainingCards := deck.Cards[:20]

		analysis := pc.calculateActionAnalysis(playerHand, dealerHand, remainingCards)

		if analysis == nil {
			t.Fatal("Expected non-nil ActionAnalysis")
		}

		// 硬17应该停牌
		if analysis.RecommendedAction != "stand" {
			t.Error("Hard 17 should recommend stand")
		}
	})

	t.Run("pair_scenarios", func(t *testing.T) {
		t.Parallel()

		// 对子8-8
		playerHand := entities.NewHand()
		playerHand.AddCard(entities.Card{Suit: entities.Hearts, Rank: entities.Eight})
		playerHand.AddCard(entities.Card{Suit: entities.Spades, Rank: entities.Eight})

		dealerHand := entities.NewHand()
		dealerHand.AddCard(entities.Card{Suit: entities.Diamonds, Rank: entities.Six})

		remainingCards := deck.Cards[:20]

		analysis := pc.calculateActionAnalysis(playerHand, dealerHand, remainingCards)

		if analysis == nil {
			t.Fatal("Expected non-nil ActionAnalysis")
		}

		// 应该能够分牌
		if !analysis.CanSplit {
			t.Error("Should be able to split pair of 8s")
		}

		// 分牌胜率应该被计算
		if analysis.SplitWinRate <= 0 {
			t.Error("Split win rate should be calculated for valid split")
		}
	})
}

// TestSimulationConsistency 测试模拟一致性
func TestSimulationConsistency(t *testing.T) {
	t.Parallel()

	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)
	pc.trials = 1000 // 使用较少试验次数以加快测试

	playerHand := entities.NewHand()
	playerHand.AddCard(entities.Card{Suit: entities.Hearts, Rank: entities.Ten})
	playerHand.AddCard(entities.Card{Suit: entities.Spades, Rank: entities.Six})

	dealerHand := entities.NewHand()
	dealerHand.AddCard(entities.Card{Suit: entities.Diamonds, Rank: entities.Seven})
	dealerHand.AddCard(entities.Card{Suit: entities.Clubs, Rank: entities.Five})

	remainingCards := createRemainingCards(
		[]entities.Card{
			{Suit: entities.Hearts, Rank: entities.Ten},
			{Suit: entities.Spades, Rank: entities.Six},
		},
		[]entities.Card{
			{Suit: entities.Diamonds, Rank: entities.Seven},
			{Suit: entities.Clubs, Rank: entities.Five},
		},
	)

	// 运行多次，检查结果的一致性
	var results []*ProbabilityResult
	for i := 0; i < 5; i++ {
		result := pc.CalculateWinProbabilities(playerHand, dealerHand, remainingCards, 1000)
		results = append(results, result)
	}

	// 检查结果是否在合理范围内变化
	for i := 1; i < len(results); i++ {
		diff := abs(results[0].PlayerWinProbability - results[i].PlayerWinProbability)
		if diff > 0.1 { // 允许10%的变化
			t.Errorf("Results vary too much between runs: %f vs %f",
				results[0].PlayerWinProbability, results[i].PlayerWinProbability)
		}
	}
}

// TestMemoryUsage 测试内存使用
func TestMemoryUsage(t *testing.T) {
	t.Parallel()

	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)
	pc.trials = 100 // 减少试验次数

	playerHand := entities.NewHand()
	playerHand.AddCard(entities.Card{Suit: entities.Hearts, Rank: entities.Eight})
	playerHand.AddCard(entities.Card{Suit: entities.Spades, Rank: entities.Six})

	dealerHand := entities.NewHand()
	dealerHand.AddCard(entities.Card{Suit: entities.Diamonds, Rank: entities.Ten})

	remainingCards := deck.Cards[:30]

	// 运行多次计算，确保没有内存泄漏
	for i := 0; i < 100; i++ {
		result := pc.CalculateWinProbabilities(playerHand, dealerHand, remainingCards, 500)
		if result == nil {
			t.Fatalf("Iteration %d returned nil result", i)
		}
		// 不需要显式设置result = nil，Go的垃圾回收器会处理
	}
}

// Helper function
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
