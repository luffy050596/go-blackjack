package services

import (
	"math/rand/v2"
	"slices"
	"time"

	"github.com/luffy050596/go-blackjack/internal/domain/entities"
)

// ProbabilityCalculator 概率计算器
type ProbabilityCalculator struct {
	deck   *entities.Deck
	trials int // 蒙特卡洛模拟次数
	rng    *rand.Rand
}

// NewProbabilityCalculator 创建概率计算器
func NewProbabilityCalculator(deck *entities.Deck) *ProbabilityCalculator {
	return &ProbabilityCalculator{
		deck:   deck,
		trials: 10000,
		rng:    rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano()<<32))),
	}
}

// ProbabilityResult 概率计算结果
type ProbabilityResult struct {
	PlayerWinProbability  float64
	DealerWinProbability  float64
	PushProbability       float64
	PlayerBlackjackProb   float64
	DealerBlackjackProb   float64
	PlayerBustProbability float64
	DealerBustProbability float64
	Player21Probability   float64
	Dealer21Probability   float64

	// 操作胜率分析
	ActionAnalysis *ActionAnalysis
}

// ActionAnalysis 操作胜率分析结果
type ActionAnalysis struct {
	HitWinRate    float64 // 要牌胜率
	StandWinRate  float64 // 停牌胜率
	DoubleWinRate float64 // 加倍胜率
	SplitWinRate  float64 // 分牌胜率（如果可分牌）

	// 操作可用性
	CanHit    bool
	CanStand  bool
	CanDouble bool
	CanSplit  bool

	// 推荐操作
	RecommendedAction string
	ExpectedValue     float64 // 推荐操作的期望值

	// 凯利公式相关
	KellyRecommendation *KellyRecommendation
}

// KellyRecommendation 凯利公式推荐结果
type KellyRecommendation struct {
	// 当前情况下的凯利比例
	StandardKellyFraction  float64 // 普通胜利的凯利比例
	BlackjackKellyFraction float64 // Blackjack胜利的凯利比例
	DoubleKellyFraction    float64 // 加倍的凯利比例

	// 推荐投注金额
	RecommendedBetAmount   int     // 基于凯利公式的推荐投注金额
	RecommendedBetFraction float64 // 推荐投注比例（相对于总筹码）

	// 加倍决策
	ShouldDouble      bool    // 是否推荐加倍
	DoubleExpectedROI float64 // 加倍的期望投资回报率

	// 风险评估
	RiskLevel          string  // 风险等级 (Low/Medium/High)
	ExpectedGrowthRate float64 // 期望资金增长率
}

// CalculateWinProbabilities 计算获胜概率
func (pc *ProbabilityCalculator) CalculateWinProbabilities(
	playerHand *entities.Hand,
	dealerHand *entities.Hand,
	remainingCards []entities.Card,
	currentChips int,
) *ProbabilityResult {
	// 获取当前玩家状态
	currentPlayerValue := playerHand.Value()
	currentPlayerBlackjack := playerHand.IsBlackjack()

	// 初始化概率变量
	var player21Prob float64
	var playerBlackjackProb float64

	// 如果玩家已经有21点，直接返回确定结果
	if currentPlayerValue == 21 {
		player21Prob = 1.0
		if currentPlayerBlackjack {
			playerBlackjackProb = 1.0
		}

		// 对于已经21点的情况，只需要计算庄家的概率
		return pc.calculateProbabilitiesFor21Player(dealerHand, remainingCards, player21Prob, playerBlackjackProb)
	}

	// 计算如果玩家要牌一次的直接概率
	hitAnalysis := pc.calculateHitAnalysis(currentPlayerValue, remainingCards)

	// 进行蒙特卡洛模拟
	playerWins := 0
	dealerWins := 0
	pushes := 0
	playerBlackjacks := 0
	dealerBlackjacks := 0
	playerBusts := 0
	dealerBusts := 0
	player21s := 0
	dealer21s := 0

	// 进行蒙特卡洛模拟
	for i := 0; i < pc.trials; i++ {
		result := pc.simulateGame(playerHand, dealerHand, remainingCards)

		switch result.Winner {
		case "player":
			playerWins++
		case "dealer":
			dealerWins++
		case "push":
			pushes++
		}

		if result.PlayerBlackjack {
			playerBlackjacks++
		}
		if result.DealerBlackjack {
			dealerBlackjacks++
		}
		if result.PlayerBust {
			playerBusts++
		}
		if result.DealerBust {
			dealerBusts++
		}
		if result.PlayerFinalValue == 21 {
			player21s++
		}
		if result.DealerFinalValue == 21 {
			dealer21s++
		}
	}

	trials := float64(pc.trials)

	// 使用直接分析的结果替代模拟中可能不准确的概率
	player21Prob = hitAnalysis.Hit21Probability
	playerBlackjackProb = float64(playerBlackjacks) / trials

	// 计算操作胜率分析
	actionAnalysis := pc.calculateActionAnalysis(playerHand, dealerHand, remainingCards)

	// 创建概率结果
	result := &ProbabilityResult{
		PlayerWinProbability:  float64(playerWins) / trials,
		DealerWinProbability:  float64(dealerWins) / trials,
		PushProbability:       float64(pushes) / trials,
		PlayerBlackjackProb:   playerBlackjackProb,
		DealerBlackjackProb:   float64(dealerBlackjacks) / trials,
		PlayerBustProbability: hitAnalysis.BustProbability,
		DealerBustProbability: float64(dealerBusts) / trials,
		Player21Probability:   player21Prob,
		Dealer21Probability:   float64(dealer21s) / trials,
		ActionAnalysis:        actionAnalysis,
	}

	// 计算凯利公式推荐
	kellyRecommendation := pc.calculateKellyRecommendation(playerHand, dealerHand, remainingCards, currentChips, result)
	actionAnalysis.KellyRecommendation = kellyRecommendation

	return result
}

// HitAnalysis 要牌分析结果
type HitAnalysis struct {
	BustProbability    float64
	Hit21Probability   float64
	SafeHitProbability float64
}

// calculateHitAnalysis 计算玩家要牌一次的直接概率分析
func (pc *ProbabilityCalculator) calculateHitAnalysis(currentValue int, remainingCards []entities.Card) *HitAnalysis {
	if len(remainingCards) == 0 {
		return &HitAnalysis{
			BustProbability:    0.0,
			Hit21Probability:   0.0,
			SafeHitProbability: 0.0,
		}
	}

	totalCards := len(remainingCards)
	bustCards := 0
	hit21Cards := 0
	safeCards := 0

	// 分析每张剩余牌的效果
	for _, card := range remainingCards {
		cardValue := card.BaseValue()

		// 处理A牌的特殊情况
		if card.IsAce() {
			// A牌可以是1或11，选择最优值
			if currentValue+1 <= 21 {
				cardValue = 1 // 使用1避免爆牌
			} else {
				cardValue = 11 // 如果+1还是会爆，那+11也会爆
			}
		}

		newValue := currentValue + cardValue

		switch {
		case newValue > 21:
			bustCards++
		case newValue == 21:
			hit21Cards++
		default:
			safeCards++
		}
	}

	return &HitAnalysis{
		BustProbability:    float64(bustCards) / float64(totalCards),
		Hit21Probability:   float64(hit21Cards) / float64(totalCards),
		SafeHitProbability: float64(safeCards) / float64(totalCards),
	}
}

// calculateProbabilitiesFor21Player 为已经有21点的玩家计算概率
func (pc *ProbabilityCalculator) calculateProbabilitiesFor21Player(
	dealerHand *entities.Hand,
	remainingCards []entities.Card,
	player21Prob float64,
	playerBlackjackProb float64,
) *ProbabilityResult {
	// 对于已经21点的玩家，只模拟庄家的行为
	dealerWins := 0
	pushes := 0
	dealerBlackjacks := 0
	dealerBusts := 0
	dealer21s := 0

	for i := 0; i < pc.trials; i++ {
		// 创建庄家模拟手牌
		simDealerHand := entities.NewHand()
		if len(dealerHand.Cards) > 0 {
			simDealerHand.AddCard(dealerHand.Cards[0])
		}

		// 创建剩余牌的副本并洗牌
		simDeck := pc.createShuffledDeckWithHiddenCard(remainingCards, dealerHand)
		deckIndex := 0

		// 如果庄家有隐藏牌，先为庄家发隐藏牌
		if len(dealerHand.Cards) > 1 && deckIndex < len(simDeck) {
			simDealerHand.AddCard(simDeck[deckIndex])
			deckIndex++
		}

		// 庄家按规则要牌
		for simDealerHand.Value() < 17 && deckIndex < len(simDeck) {
			simDealerHand.AddCard(simDeck[deckIndex])
			deckIndex++
		}

		// 评估结果（玩家固定21点）
		dealerValue := simDealerHand.Value()
		dealerBlackjack := simDealerHand.IsBlackjack()
		dealerBust := simDealerHand.IsBust()

		if dealerBlackjack {
			dealerBlackjacks++
		}
		if dealerBust {
			dealerBusts++
		}
		if dealerValue == 21 {
			dealer21s++
		}

		// 判断胜负（玩家固定21点）
		switch {
		case dealerBust:
			// 玩家胜利（已计入playerWins）
		case dealerBlackjack && playerBlackjackProb < 1.0:
			// 庄家Blackjack而玩家不是Blackjack
			dealerWins++
		case dealerValue == 21:
			// 平局
			pushes++
		default:
			// 玩家21点胜过庄家非21点（已计入playerWins）
		}
	}

	trials := float64(pc.trials)
	playerWins := pc.trials - dealerWins - pushes

	// 为已经21点的玩家创建操作分析
	actionAnalysis := &ActionAnalysis{
		HitWinRate:        0.0,                          // 已经21点，不能再要牌
		StandWinRate:      float64(playerWins) / trials, // 停牌胜率即为当前胜率
		DoubleWinRate:     0.0,                          // 不能加倍
		SplitWinRate:      0.0,                          // 不能分牌
		CanHit:            false,
		CanStand:          true,
		CanDouble:         false,
		CanSplit:          false,
		RecommendedAction: "stand",
		ExpectedValue:     float64(playerWins) / trials,
	}

	return &ProbabilityResult{
		PlayerWinProbability:  float64(playerWins) / trials,
		DealerWinProbability:  float64(dealerWins) / trials,
		PushProbability:       float64(pushes) / trials,
		PlayerBlackjackProb:   playerBlackjackProb,
		DealerBlackjackProb:   float64(dealerBlackjacks) / trials,
		PlayerBustProbability: 0.0, // 已经21点，不会爆牌
		DealerBustProbability: float64(dealerBusts) / trials,
		Player21Probability:   player21Prob,
		Dealer21Probability:   float64(dealer21s) / trials,
		ActionAnalysis:        actionAnalysis,
	}
}

// SimulationResult 模拟结果
type SimulationResult struct {
	Winner           string
	PlayerFinalValue int
	DealerFinalValue int
	PlayerBlackjack  bool
	DealerBlackjack  bool
	PlayerBust       bool
	DealerBust       bool
}

// simulateGame 模拟一局游戏
func (pc *ProbabilityCalculator) simulateGame(
	playerHand *entities.Hand,
	dealerHand *entities.Hand,
	remainingCards []entities.Card,
) *SimulationResult {
	// 复制玩家手牌
	simPlayerHand := pc.copyHand(playerHand)

	// 创建庄家模拟手牌 - 只包含明牌
	simDealerHand := entities.NewHand()
	if len(dealerHand.Cards) > 0 {
		// 只添加第一张牌（明牌）
		simDealerHand.AddCard(dealerHand.Cards[0])
	}

	// 创建剩余牌的副本并洗牌，包含庄家的隐藏牌
	simDeck := pc.createShuffledDeckWithHiddenCard(remainingCards, dealerHand)
	deckIndex := 0

	// 先为庄家发隐藏牌
	if len(dealerHand.Cards) > 1 && deckIndex < len(simDeck) {
		simDealerHand.AddCard(simDeck[deckIndex])
		deckIndex++
	}

	// 玩家决策（使用基本策略）
	for !simPlayerHand.IsBust() && simPlayerHand.Value() < 21 {
		action := pc.getBasicStrategyAction(simPlayerHand, simDealerHand)
		if action == "stand" {
			break
		}
		if action == "hit" && deckIndex < len(simDeck) {
			simPlayerHand.AddCard(simDeck[deckIndex])
			deckIndex++
		} else {
			break
		}
	}

	// 庄家按规则要牌
	for simDealerHand.Value() < 17 && deckIndex < len(simDeck) {
		simDealerHand.AddCard(simDeck[deckIndex])
		deckIndex++
	}

	// 评估结果
	return pc.evaluateResult(simPlayerHand, simDealerHand)
}

// getBasicStrategyAction 获取基本策略行动 - 只基于庄家明牌
func (pc *ProbabilityCalculator) getBasicStrategyAction(playerHand *entities.Hand, dealerHand *entities.Hand) string {
	playerValue := playerHand.Value()

	// 如果庄家手牌为空，无法获取明牌
	if len(dealerHand.Cards) < 1 {
		// 默认策略：小于17点继续要牌
		if playerValue < 17 {
			return "hit"
		}
		return "stand"
	}

	// 只使用庄家的第一张牌（明牌）进行决策
	dealerUpCard := dealerHand.Cards[0].Value()

	// 简化的基本策略
	if playerHand.IsSoft() {
		// 软牌策略
		if playerValue <= 17 {
			return "hit"
		}
		if playerValue == 18 && (dealerUpCard >= 9 || dealerUpCard == 1) {
			return "hit"
		}
		return "stand"
	}

	// 硬牌策略
	if playerValue <= 11 {
		return "hit"
	}
	if playerValue == 12 && (dealerUpCard >= 4 && dealerUpCard <= 6) {
		return "stand"
	}
	if playerValue >= 13 && playerValue <= 16 && dealerUpCard <= 6 {
		return "stand"
	}
	if playerValue >= 17 {
		return "stand"
	}
	return "hit"
}

// copyHand 复制手牌
func (pc *ProbabilityCalculator) copyHand(hand *entities.Hand) *entities.Hand {
	newHand := entities.NewHand()
	newHand.Cards = slices.Clone(hand.Cards)
	return newHand
}

// createShuffledDeckWithHiddenCard 创建包含庄家隐藏牌的洗牌牌组
func (pc *ProbabilityCalculator) createShuffledDeckWithHiddenCard(remainingCards []entities.Card, dealerHand *entities.Hand) []entities.Card {
	// 复制剩余卡牌
	deck := make([]entities.Card, 0, len(remainingCards)+1)
	deck = append(deck, remainingCards...)

	// 如果庄家有隐藏牌（第二张牌），将其添加到牌堆中
	if len(dealerHand.Cards) > 1 {
		deck = append(deck, dealerHand.Cards[1])
	}

	// 使用Fisher-Yates洗牌算法
	for i := len(deck) - 1; i > 0; i-- {
		j := pc.rng.IntN(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}

	return deck
}

// evaluateResult 评估游戏结果
func (pc *ProbabilityCalculator) evaluateResult(playerHand *entities.Hand, dealerHand *entities.Hand) *SimulationResult {
	result := &SimulationResult{
		PlayerFinalValue: playerHand.Value(),
		DealerFinalValue: dealerHand.Value(),
		PlayerBlackjack:  playerHand.IsBlackjack(),
		DealerBlackjack:  dealerHand.IsBlackjack(),
		PlayerBust:       playerHand.IsBust(),
		DealerBust:       dealerHand.IsBust(),
	}

	// 判断胜负
	switch {
	case result.PlayerBust:
		result.Winner = "dealer"
	case result.DealerBust:
		result.Winner = "player"
	case result.PlayerBlackjack && result.DealerBlackjack:
		result.Winner = "push"
	case result.PlayerBlackjack:
		result.Winner = "player"
	case result.DealerBlackjack:
		result.Winner = "dealer"
	case result.PlayerFinalValue > result.DealerFinalValue:
		result.Winner = "player"
	case result.PlayerFinalValue < result.DealerFinalValue:
		result.Winner = "dealer"
	default:
		result.Winner = "push"
	}

	return result
}

// calculateActionAnalysis 计算操作胜率分析
func (pc *ProbabilityCalculator) calculateActionAnalysis(playerHand *entities.Hand, dealerHand *entities.Hand, remainingCards []entities.Card) *ActionAnalysis {
	currentValue := playerHand.Value()
	isFirstTurn := len(playerHand.Cards) == 2 // 是否为第一轮（可以加倍/分牌）

	// 检查操作可用性
	canHit := currentValue < 21 && !playerHand.IsBust()
	canStand := true
	canDouble := isFirstTurn && canHit
	canSplit := isFirstTurn && len(playerHand.Cards) == 2 &&
		playerHand.Cards[0].Rank == playerHand.Cards[1].Rank

	actionAnalysis := &ActionAnalysis{
		CanHit:    canHit,
		CanStand:  canStand,
		CanDouble: canDouble,
		CanSplit:  canSplit,
	}

	// 如果玩家已经有21点或爆牌，只能停牌
	if currentValue >= 21 {
		actionAnalysis.StandWinRate = pc.calculateStandWinRate(playerHand, dealerHand, remainingCards)
		actionAnalysis.RecommendedAction = "stand"
		actionAnalysis.ExpectedValue = actionAnalysis.StandWinRate
		return actionAnalysis
	}

	// 计算停牌胜率
	standWinRate := pc.calculateStandWinRate(playerHand, dealerHand, remainingCards)
	actionAnalysis.StandWinRate = standWinRate

	// 计算要牌胜率
	hitWinRate := 0.0
	if canHit {
		hitWinRate = pc.calculateHitWinRate(playerHand, dealerHand, remainingCards)
		actionAnalysis.HitWinRate = hitWinRate
	}

	// 计算加倍胜率
	doubleWinRate := 0.0
	if canDouble {
		doubleWinRate = pc.calculateDoubleWinRate(playerHand, dealerHand, remainingCards)
		actionAnalysis.DoubleWinRate = doubleWinRate
	}

	// 计算分牌胜率
	splitWinRate := 0.0
	if canSplit {
		splitWinRate = pc.calculateSplitWinRate(playerHand, dealerHand, remainingCards)
		actionAnalysis.SplitWinRate = splitWinRate
	}

	// 确定推荐操作
	bestAction := "stand"
	bestValue := standWinRate

	if canHit && hitWinRate > bestValue {
		bestAction = "hit"
		bestValue = hitWinRate
	}

	if canDouble && doubleWinRate > bestValue {
		bestAction = "double"
		bestValue = doubleWinRate
	}

	if canSplit && splitWinRate > bestValue {
		bestAction = "split"
		bestValue = splitWinRate
	}

	actionAnalysis.RecommendedAction = bestAction
	actionAnalysis.ExpectedValue = bestValue

	return actionAnalysis
}

// calculateStandWinRate 计算停牌胜率
func (pc *ProbabilityCalculator) calculateStandWinRate(playerHand *entities.Hand, dealerHand *entities.Hand, remainingCards []entities.Card) float64 {
	playerWins := 0
	trials := pc.trials

	for i := 0; i < trials; i++ {
		// 模拟庄家完成手牌
		finalDealerHand := pc.simulateDealerPlay(dealerHand, remainingCards)

		// 判断结果
		if pc.playerWins(playerHand, finalDealerHand) {
			playerWins++
		}
	}

	return float64(playerWins) / float64(trials)
}

// calculateHitWinRate 计算要牌胜率
func (pc *ProbabilityCalculator) calculateHitWinRate(playerHand *entities.Hand, dealerHand *entities.Hand, remainingCards []entities.Card) float64 {
	if len(remainingCards) == 0 {
		return 0.0
	}

	totalWinRate := 0.0
	totalCards := len(remainingCards)
	trials := 100 // 减少每张牌的模拟次数

	// 对每张可能的牌计算期望胜率
	for _, card := range remainingCards {
		// 复制玩家手牌并添加这张牌
		newPlayerHand := pc.copyHand(playerHand)
		newPlayerHand.AddCard(card)

		// 创建不包含这张牌的剩余牌堆
		newRemainingCards := pc.removeCard(remainingCards, card)

		var winRate float64
		if newPlayerHand.IsBust() {
			winRate = 0.0 // 爆牌必败
		} else {
			// 无论是21点还是其他情况，都使用相同的计算方法
			winRate = pc.calculateStandWinRateSimple(newPlayerHand, dealerHand, newRemainingCards, trials)
		}

		totalWinRate += winRate
	}

	return totalWinRate / float64(totalCards)
}

// calculateStandWinRateSimple 简化版停牌胜率计算
func (pc *ProbabilityCalculator) calculateStandWinRateSimple(playerHand *entities.Hand, dealerHand *entities.Hand, remainingCards []entities.Card, trials int) float64 {
	playerWins := 0

	for i := 0; i < trials; i++ {
		// 模拟庄家完成手牌
		finalDealerHand := pc.simulateDealerPlay(dealerHand, remainingCards)

		// 判断结果
		if pc.playerWins(playerHand, finalDealerHand) {
			playerWins++
		}
	}

	return float64(playerWins) / float64(trials)
}

// calculateDoubleWinRate 计算加倍胜率（只要一张牌然后停牌）
func (pc *ProbabilityCalculator) calculateDoubleWinRate(playerHand *entities.Hand, dealerHand *entities.Hand, remainingCards []entities.Card) float64 {
	if len(remainingCards) == 0 {
		return 0.0
	}

	totalWinRate := 0.0
	totalCards := len(remainingCards)
	trials := 100 // 减少模拟次数

	// 对每张可能的牌计算胜率
	for _, card := range remainingCards {
		// 复制玩家手牌并添加这张牌
		newPlayerHand := pc.copyHand(playerHand)
		newPlayerHand.AddCard(card)

		// 创建不包含这张牌的剩余牌堆
		newRemainingCards := pc.removeCard(remainingCards, card)

		var winRate float64
		if newPlayerHand.IsBust() {
			winRate = 0.0 // 爆牌必败
		} else {
			// 加倍后必须停牌，使用简化计算
			winRate = pc.calculateStandWinRateSimple(newPlayerHand, dealerHand, newRemainingCards, trials)
		}

		totalWinRate += winRate
	}

	return totalWinRate / float64(totalCards)
}

// calculateSplitWinRate 计算分牌胜率
func (pc *ProbabilityCalculator) calculateSplitWinRate(playerHand *entities.Hand, dealerHand *entities.Hand, remainingCards []entities.Card) float64 {
	if len(playerHand.Cards) != 2 || playerHand.Cards[0].Rank != playerHand.Cards[1].Rank {
		return 0.0
	}

	// 简化处理：计算单张牌开始的期望胜率的平均值
	// 实际实现可能需要更复杂的逻辑来处理分牌后的游戏
	card := playerHand.Cards[0]

	// 创建两个新手牌，每个都从一张相同的牌开始
	hand1 := entities.NewHand()
	hand1.AddCard(card)

	hand2 := entities.NewHand()
	hand2.AddCard(card)

	// 计算两手牌的最优胜率平均值
	winRate1 := pc.calculateOptimalWinRate(hand1, dealerHand, remainingCards)
	winRate2 := pc.calculateOptimalWinRate(hand2, dealerHand, remainingCards)

	return (winRate1 + winRate2) / 2.0
}

// calculateOptimalWinRate 计算当前手牌的最优策略胜率
func (pc *ProbabilityCalculator) calculateOptimalWinRate(playerHand *entities.Hand, dealerHand *entities.Hand, remainingCards []entities.Card) float64 {
	// 计算停牌和要牌的胜率，选择更高的
	standRate := pc.calculateStandWinRate(playerHand, dealerHand, remainingCards)

	if playerHand.Value() >= 21 {
		return standRate
	}

	hitRate := pc.calculateHitWinRate(playerHand, dealerHand, remainingCards)

	if hitRate > standRate {
		return hitRate
	}
	return standRate
}

// simulateDealerPlay 模拟庄家完成手牌
func (pc *ProbabilityCalculator) simulateDealerPlay(dealerHand *entities.Hand, remainingCards []entities.Card) *entities.Hand {
	// 创建庄家模拟手牌
	simDealerHand := entities.NewHand()
	if len(dealerHand.Cards) > 0 {
		simDealerHand.AddCard(dealerHand.Cards[0])
	}

	// 创建剩余牌的副本并洗牌
	simDeck := pc.createShuffledDeckWithHiddenCard(remainingCards, dealerHand)
	deckIndex := 0

	// 如果庄家有隐藏牌，先为庄家发隐藏牌
	if len(dealerHand.Cards) > 1 && deckIndex < len(simDeck) {
		simDealerHand.AddCard(simDeck[deckIndex])
		deckIndex++
	}

	// 庄家按规则要牌
	for simDealerHand.Value() < 17 && deckIndex < len(simDeck) {
		simDealerHand.AddCard(simDeck[deckIndex])
		deckIndex++
	}

	return simDealerHand
}

// playerWins 判断玩家是否获胜
func (pc *ProbabilityCalculator) playerWins(playerHand *entities.Hand, dealerHand *entities.Hand) bool {
	if playerHand.IsBust() {
		return false
	}
	if dealerHand.IsBust() {
		return true
	}

	playerValue := playerHand.Value()
	dealerValue := dealerHand.Value()

	// 处理Blackjack优先级
	playerBlackjack := playerHand.IsBlackjack()
	dealerBlackjack := dealerHand.IsBlackjack()

	if playerBlackjack && !dealerBlackjack {
		return true
	}
	if !playerBlackjack && dealerBlackjack {
		return false
	}

	return playerValue > dealerValue
}

// removeCard 从牌堆中移除指定卡牌
func (pc *ProbabilityCalculator) removeCard(cards []entities.Card, cardToRemove entities.Card) []entities.Card {
	result := make([]entities.Card, 0, len(cards)-1)
	removed := false

	for _, card := range cards {
		if !removed && card.Suit == cardToRemove.Suit && card.Rank == cardToRemove.Rank {
			removed = true
			continue
		}
		result = append(result, card)
	}

	return result
}

// calculateKellyRecommendation 计算凯利公式推荐
func (pc *ProbabilityCalculator) calculateKellyRecommendation(
	_ *entities.Hand,
	_ *entities.Hand,
	_ []entities.Card,
	currentChips int,
	probResult *ProbabilityResult,
) *KellyRecommendation {
	// 基础概率数据
	winProb := probResult.PlayerWinProbability
	blackjackProb := probResult.PlayerBlackjackProb
	loseProb := probResult.DealerWinProbability

	// 调整概率：排除平局的影响
	totalNonPushProb := winProb + loseProb
	if totalNonPushProb > 0 {
		winProb /= totalNonPushProb
		loseProb /= totalNonPushProb
	}

	// 标准凯利计算（1:1赔率）
	standardKellyFraction := pc.calculateKellyFraction(winProb, loseProb, 1.0)

	// Blackjack凯利计算（3:2赔率）
	blackjackKellyFraction := pc.calculateKellyFraction(blackjackProb, 1.0-blackjackProb, 1.5)

	// 加倍决策凯利计算
	doubleWinProb := 0.0
	if probResult.ActionAnalysis != nil {
		doubleWinProb = probResult.ActionAnalysis.DoubleWinRate
	}
	doubleLoseProb := 1.0
	if doubleWinProb > 0 {
		doubleLoseProb = 1.0 - doubleWinProb
	}
	doubleKellyFraction := pc.calculateKellyFraction(doubleWinProb, doubleLoseProb, 1.0)

	// 计算推荐投注金额
	recommendedFraction := standardKellyFraction

	// 选择最优的凯利比例
	if blackjackKellyFraction > recommendedFraction {
		recommendedFraction = blackjackKellyFraction
	}

	// 安全限制：不超过总筹码的10%
	recommendedFraction = minFloat64(recommendedFraction, 0.10)

	recommendedAmount := int(float64(currentChips) * recommendedFraction)

	// 最小投注限制
	if recommendedAmount < 10 && currentChips >= 10 {
		recommendedAmount = 10
	}

	// 风险评估
	riskLevel := pc.assessRiskLevel(recommendedFraction)

	// 加倍决策
	shouldDouble := doubleKellyFraction > 0.02 && probResult.ActionAnalysis != nil && probResult.ActionAnalysis.CanDouble
	doubleROI := doubleWinProb*2.0 - 1.0

	// 计算期望增长率
	expectedGrowthRate := pc.calculateExpectedGrowthRate(winProb, loseProb, recommendedFraction)

	return &KellyRecommendation{
		StandardKellyFraction:  standardKellyFraction,
		BlackjackKellyFraction: blackjackKellyFraction,
		DoubleKellyFraction:    doubleKellyFraction,
		RecommendedBetAmount:   recommendedAmount,
		RecommendedBetFraction: recommendedFraction,
		ShouldDouble:           shouldDouble,
		DoubleExpectedROI:      doubleROI,
		RiskLevel:              riskLevel,
		ExpectedGrowthRate:     expectedGrowthRate,
	}
}

// calculateKellyFraction 计算凯利比例
// winProb: 获胜概率, loseProb: 失败概率, odds: 赔率
func (pc *ProbabilityCalculator) calculateKellyFraction(winProb, loseProb, odds float64) float64 {
	if odds <= 0 || winProb <= 0 {
		return 0
	}

	// 凯利公式: f* = (bp - q) / b
	// 其中 b = odds, p = winProb, q = loseProb
	kelly := (odds*winProb - loseProb) / odds

	// 确保非负值
	return maxFloat64(0, kelly)
}

// assessRiskLevel 评估风险等级
func (pc *ProbabilityCalculator) assessRiskLevel(kellyFraction float64) string {
	switch {
	case kellyFraction <= 0.02:
		return "Low"
	case kellyFraction <= 0.05:
		return "Medium"
	default:
		return "High"
	}
}

// calculateExpectedGrowthRate 计算期望资金增长率
func (pc *ProbabilityCalculator) calculateExpectedGrowthRate(winProb, loseProb, kellyFraction float64) float64 {
	if kellyFraction <= 0 {
		return 0
	}

	// 使用Kelly公式的对数增长率计算
	// G = p * log(1 + f * b) + q * log(1 - f)
	// 这里简化为线性近似
	expectedReturn := winProb*(1+kellyFraction) + loseProb*(1-kellyFraction) - 1
	return expectedReturn
}

// maxFloat64 returns the maximum of two float64 values
func maxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}

	return b
}

// minFloat64 returns the minimum of two float64 values
func minFloat64(a, b float64) float64 {
	if a < b {
		return a
	}

	return b
}

// CalculateBasicKellyFraction 计算基础凯利公式推荐（用于下注阶段）
func (pc *ProbabilityCalculator) CalculateBasicKellyFraction(winRate, loseRate float64, currentChips int) *KellyRecommendation {
	// 在娱乐性游戏中，重点是资金管理而非严格的期望收益
	// 提供保守的资金管理建议

	// 计算理论凯利比例（仅供参考）
	standardKelly := pc.calculateKellyFraction(winRate, loseRate, 1.0)
	blackjackRate := 0.048 // 约4.8%的概率获得blackjack
	blackjackKelly := pc.calculateKellyFraction(blackjackRate, loseRate, 1.5)

	// 实用的资金管理建议
	// 1. 单次下注不超过总资金的1-2%（保守策略）
	// 2. 根据筹码数量调整建议
	var recommendedFraction float64
	var recommendedAmount int
	var riskLevel string

	switch {
	case currentChips >= 1000:
		// 资金充足：建议保守下注（1-2%）
		recommendedFraction = 0.015 // 1.5%
		calculatedAmount := int(float64(currentChips) * recommendedFraction)
		if calculatedAmount > 10 {
			recommendedAmount = calculatedAmount
		} else {
			recommendedAmount = 10
		}
		riskLevel = "Low"
	case currentChips >= 500:
		// 资金中等：稍微保守（1%）
		recommendedFraction = 0.01 // 1%
		calculatedAmount := int(float64(currentChips) * recommendedFraction)
		if calculatedAmount > 10 {
			recommendedAmount = calculatedAmount
		} else {
			recommendedAmount = 10
		}
		riskLevel = "Low"
	case currentChips >= 200:
		// 资金较少：更加保守（0.5%）
		recommendedFraction = 0.005 // 0.5%
		recommendedAmount = 10      // 最小下注
		riskLevel = "Medium"
	default:
		// 资金紧张：最小下注或考虑离开
		recommendedFraction = 0.0
		recommendedAmount = 10 // 最小下注
		riskLevel = "High"
	}

	// 确保推荐金额不超过可用筹码
	if recommendedAmount > currentChips {
		recommendedAmount = currentChips
	}

	// 期望增长率（娱乐成本的角度）
	expectedLoss := (loseRate - winRate) * recommendedFraction // 预期每次下注的损失比例

	return &KellyRecommendation{
		StandardKellyFraction:  standardKelly,
		BlackjackKellyFraction: blackjackKelly,
		DoubleKellyFraction:    0,
		RecommendedBetAmount:   recommendedAmount,
		RecommendedBetFraction: recommendedFraction,
		ShouldDouble:           false,
		DoubleExpectedROI:      0,
		RiskLevel:              riskLevel,
		ExpectedGrowthRate:     -expectedLoss, // 显示为预期损失率
	}
}
