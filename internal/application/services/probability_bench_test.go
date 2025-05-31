package services

import (
	"fmt"
	"testing"

	"github.com/luffy050596/go-blackjack/internal/domain/entities"
)

// BenchmarkNewProbabilityCalculator 基准测试创建概率计算器
func BenchmarkNewProbabilityCalculator(b *testing.B) {
	deck := entities.NewDeck()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewProbabilityCalculator(deck)
	}
}

// BenchmarkCalculateWinProbabilities 基准测试概率计算
func BenchmarkCalculateWinProbabilities(b *testing.B) {
	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)

	// 设置较少的试验次数以加快基准测试
	pc.trials = 1000

	// 创建测试手牌
	playerHand := entities.NewHand()
	playerHand.AddCard(entities.Card{Suit: entities.Hearts, Rank: entities.Eight})
	playerHand.AddCard(entities.Card{Suit: entities.Spades, Rank: entities.Six})

	dealerHand := entities.NewHand()
	dealerHand.AddCard(entities.Card{Suit: entities.Diamonds, Rank: entities.Ten})
	dealerHand.AddCard(entities.Card{Suit: entities.Clubs, Rank: entities.Five})

	remainingCards := createRemainingCards(
		[]entities.Card{
			{Suit: entities.Hearts, Rank: entities.Eight},
			{Suit: entities.Spades, Rank: entities.Six},
		},
		[]entities.Card{
			{Suit: entities.Diamonds, Rank: entities.Ten},
			{Suit: entities.Clubs, Rank: entities.Five},
		},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pc.CalculateWinProbabilities(playerHand, dealerHand, remainingCards, 1000)
	}
}

// BenchmarkCalculateWinProbabilities_Blackjack 基准测试Blackjack场景概率计算
func BenchmarkCalculateWinProbabilities_Blackjack(b *testing.B) {
	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)
	pc.trials = 1000

	// 创建Blackjack手牌
	playerHand := entities.NewHand()
	playerHand.AddCard(entities.Card{Suit: entities.Hearts, Rank: entities.Ace})
	playerHand.AddCard(entities.Card{Suit: entities.Spades, Rank: entities.King})

	dealerHand := entities.NewHand()
	dealerHand.AddCard(entities.Card{Suit: entities.Diamonds, Rank: entities.Seven})
	dealerHand.AddCard(entities.Card{Suit: entities.Clubs, Rank: entities.Five})

	remainingCards := createRemainingCards(
		[]entities.Card{
			{Suit: entities.Hearts, Rank: entities.Ace},
			{Suit: entities.Spades, Rank: entities.King},
		},
		[]entities.Card{
			{Suit: entities.Diamonds, Rank: entities.Seven},
			{Suit: entities.Clubs, Rank: entities.Five},
		},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pc.CalculateWinProbabilities(playerHand, dealerHand, remainingCards, 1000)
	}
}

// BenchmarkCalculateBasicKellyFraction 基准测试基础凯利公式计算
func BenchmarkCalculateBasicKellyFraction(b *testing.B) {
	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pc.CalculateBasicKellyFraction(0.48, 0.52, 1000)
	}
}

// BenchmarkCalculateHitAnalysis 基准测试要牌分析
func BenchmarkCalculateHitAnalysis(b *testing.B) {
	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)

	remainingCards := deck.Cards[:40] // 使用前40张牌作为剩余牌堆

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pc.calculateHitAnalysis(15, remainingCards)
	}
}

// BenchmarkCalculateActionAnalysis 基准测试操作分析
func BenchmarkCalculateActionAnalysis(b *testing.B) {
	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)
	pc.trials = 500 // 减少试验次数以提高基准测试速度

	// 创建测试手牌
	playerHand := entities.NewHand()
	playerHand.AddCard(entities.Card{Suit: entities.Hearts, Rank: entities.Eight})
	playerHand.AddCard(entities.Card{Suit: entities.Spades, Rank: entities.Six})

	dealerHand := entities.NewHand()
	dealerHand.AddCard(entities.Card{Suit: entities.Diamonds, Rank: entities.Ten})

	remainingCards := createRemainingCards(
		[]entities.Card{
			{Suit: entities.Hearts, Rank: entities.Eight},
			{Suit: entities.Spades, Rank: entities.Six},
		},
		[]entities.Card{
			{Suit: entities.Diamonds, Rank: entities.Ten},
		},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pc.calculateActionAnalysis(playerHand, dealerHand, remainingCards)
	}
}

// BenchmarkSimulateGame 基准测试游戏模拟
func BenchmarkSimulateGame(b *testing.B) {
	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)

	// 创建测试手牌
	playerHand := entities.NewHand()
	playerHand.AddCard(entities.Card{Suit: entities.Hearts, Rank: entities.Eight})
	playerHand.AddCard(entities.Card{Suit: entities.Spades, Rank: entities.Six})

	dealerHand := entities.NewHand()
	dealerHand.AddCard(entities.Card{Suit: entities.Diamonds, Rank: entities.Ten})

	remainingCards := []entities.Card{
		{Suit: entities.Clubs, Rank: entities.Five},
		{Suit: entities.Hearts, Rank: entities.Seven},
		{Suit: entities.Spades, Rank: entities.Nine},
		{Suit: entities.Diamonds, Rank: entities.Two},
		{Suit: entities.Hearts, Rank: entities.Three},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pc.simulateGame(playerHand, dealerHand, remainingCards)
	}
}

// BenchmarkPlayerWins 基准测试玩家获胜判断
func BenchmarkPlayerWins(b *testing.B) {
	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)

	// 创建测试手牌
	playerHand := entities.NewHand()
	playerHand.AddCard(entities.Card{Suit: entities.Hearts, Rank: entities.Ten})
	playerHand.AddCard(entities.Card{Suit: entities.Spades, Rank: entities.Nine})

	dealerHand := entities.NewHand()
	dealerHand.AddCard(entities.Card{Suit: entities.Diamonds, Rank: entities.Ten})
	dealerHand.AddCard(entities.Card{Suit: entities.Clubs, Rank: entities.Seven})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pc.playerWins(playerHand, dealerHand)
	}
}

// BenchmarkCopyHand 基准测试手牌复制
func BenchmarkCopyHand(b *testing.B) {
	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)

	// 创建测试手牌
	hand := entities.NewHand()
	hand.AddCard(entities.Card{Suit: entities.Hearts, Rank: entities.Ace})
	hand.AddCard(entities.Card{Suit: entities.Spades, Rank: entities.King})
	hand.AddCard(entities.Card{Suit: entities.Diamonds, Rank: entities.Queen})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pc.copyHand(hand)
	}
}

// BenchmarkRemoveCard 基准测试移除牌
func BenchmarkRemoveCard(b *testing.B) {
	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)

	cards := deck.Cards[:20]  // 使用前20张牌
	cardToRemove := cards[10] // 移除中间的牌

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pc.removeCard(cards, cardToRemove)
	}
}

// BenchmarkCalculateKellyFraction 基准测试凯利公式计算
func BenchmarkCalculateKellyFraction(b *testing.B) {
	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pc.calculateKellyFraction(0.55, 0.45, 1.0)
	}
}

// BenchmarkAssessRiskLevel 基准测试风险等级评估
func BenchmarkAssessRiskLevel(b *testing.B) {
	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pc.assessRiskLevel(0.03)
	}
}

// BenchmarkCalculateStandWinRate 基准测试停牌胜率计算
func BenchmarkCalculateStandWinRate(b *testing.B) {
	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)
	pc.trials = 500 // 减少试验次数

	// 创建测试手牌
	playerHand := entities.NewHand()
	playerHand.AddCard(entities.Card{Suit: entities.Hearts, Rank: entities.Ten})
	playerHand.AddCard(entities.Card{Suit: entities.Spades, Rank: entities.Seven})

	dealerHand := entities.NewHand()
	dealerHand.AddCard(entities.Card{Suit: entities.Diamonds, Rank: entities.Six})

	remainingCards := createRemainingCards(
		[]entities.Card{
			{Suit: entities.Hearts, Rank: entities.Ten},
			{Suit: entities.Spades, Rank: entities.Seven},
		},
		[]entities.Card{
			{Suit: entities.Diamonds, Rank: entities.Six},
		},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pc.calculateStandWinRate(playerHand, dealerHand, remainingCards)
	}
}

// BenchmarkCalculateHitWinRate 基准测试要牌胜率计算
func BenchmarkCalculateHitWinRate(b *testing.B) {
	deck := entities.NewDeck()
	pc := NewProbabilityCalculator(deck)
	pc.trials = 100 // 大幅减少试验次数以提高速度

	// 创建测试手牌
	playerHand := entities.NewHand()
	playerHand.AddCard(entities.Card{Suit: entities.Hearts, Rank: entities.Eight})
	playerHand.AddCard(entities.Card{Suit: entities.Spades, Rank: entities.Five})

	dealerHand := entities.NewHand()
	dealerHand.AddCard(entities.Card{Suit: entities.Diamonds, Rank: entities.Ten})

	remainingCards := createRemainingCards(
		[]entities.Card{
			{Suit: entities.Hearts, Rank: entities.Eight},
			{Suit: entities.Spades, Rank: entities.Five},
		},
		[]entities.Card{
			{Suit: entities.Diamonds, Rank: entities.Ten},
		},
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pc.calculateHitWinRate(playerHand, dealerHand, remainingCards)
	}
}

// BenchmarkVariousTrials 比较不同试验次数的性能
func BenchmarkVariousTrials(b *testing.B) {
	trials := []int{100, 500, 1000, 5000, 10000}

	for _, trial := range trials {
		b.Run(fmt.Sprintf("trials_%d", trial), func(b *testing.B) {
			deck := entities.NewDeck()
			pc := NewProbabilityCalculator(deck)
			pc.trials = trial

			// 创建测试手牌
			playerHand := entities.NewHand()
			playerHand.AddCard(entities.Card{Suit: entities.Hearts, Rank: entities.Eight})
			playerHand.AddCard(entities.Card{Suit: entities.Spades, Rank: entities.Six})

			dealerHand := entities.NewHand()
			dealerHand.AddCard(entities.Card{Suit: entities.Diamonds, Rank: entities.Ten})
			dealerHand.AddCard(entities.Card{Suit: entities.Clubs, Rank: entities.Five})

			remainingCards := createRemainingCards(
				[]entities.Card{
					{Suit: entities.Hearts, Rank: entities.Eight},
					{Suit: entities.Spades, Rank: entities.Six},
				},
				[]entities.Card{
					{Suit: entities.Diamonds, Rank: entities.Ten},
					{Suit: entities.Clubs, Rank: entities.Five},
				},
			)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = pc.CalculateWinProbabilities(playerHand, dealerHand, remainingCards, 1000)
			}
		})
	}
}

// BenchmarkMemoryAllocation 测试内存分配
func BenchmarkMemoryAllocation(b *testing.B) {
	b.ReportAllocs()

	deck := entities.NewDeck()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pc := NewProbabilityCalculator(deck)

		// 创建测试手牌
		playerHand := entities.NewHand()
		playerHand.AddCard(entities.Card{Suit: entities.Hearts, Rank: entities.Eight})
		playerHand.AddCard(entities.Card{Suit: entities.Spades, Rank: entities.Six})

		dealerHand := entities.NewHand()
		dealerHand.AddCard(entities.Card{Suit: entities.Diamonds, Rank: entities.Ten})

		remainingCards := []entities.Card{
			{Suit: entities.Clubs, Rank: entities.Five},
			{Suit: entities.Hearts, Rank: entities.Seven},
		}

		pc.trials = 100 // 减少试验次数关注内存分配
		_ = pc.CalculateWinProbabilities(playerHand, dealerHand, remainingCards, 1000)
	}
}
