// Package dtos contains data transfer objects for the blackjack game application.
package dtos

import "github.com/luffy050596/go-blackjack/internal/domain/entities"

// GameStateDTO 游戏状态数据传输对象
type GameStateDTO struct {
	RoundNumber int                `json:"round_number"`
	PlayerChips int                `json:"player_chips"`
	PlayerBet   int                `json:"player_bet"`
	PlayerHand  *HandDTO           `json:"player_hand"`
	DealerHand  *HandDTO           `json:"dealer_hand"`
	State       entities.GameState `json:"state"`
	IsGameOver  bool               `json:"is_game_over"`
}

// HandDTO 手牌数据传输对象
type HandDTO struct {
	Cards []*CardDTO `json:"cards"`
	Value int        `json:"value"`
}

// CardDTO 卡牌数据传输对象
type CardDTO struct {
	Suit  string `json:"suit"`
	Rank  string `json:"rank"`
	Value int    `json:"value"`
}

// ActionResultDTO 行动结果数据传输对象
type ActionResultDTO struct {
	Action   entities.PlayerAction `json:"action"`
	Success  bool                  `json:"success"`
	Continue bool                  `json:"continue"`
	Card     *CardDTO              `json:"card,omitempty"`
	Message  string                `json:"message,omitempty"`
}

// GameResultDTO 游戏结果数据传输对象
type GameResultDTO struct {
	Type        entities.ResultType `json:"type"`
	BetAmount   int                 `json:"bet_amount"`
	IsDoubled   bool                `json:"is_doubled"`
	PlayerChips int                 `json:"player_chips"`
}

// BetOptionDTO 下注选项数据传输对象
type BetOptionDTO struct {
	Amount  int    `json:"amount"`
	Display string `json:"display"`
}

// MenuOptionDTO 菜单选项数据传输对象
type MenuOptionDTO struct {
	Key     string `json:"key"`
	Display string `json:"display"`
}

// ProbabilityResultDTO 概率计算结果数据传输对象
type ProbabilityResultDTO struct {
	PlayerWinProbability  float64 `json:"player_win_probability"`
	DealerWinProbability  float64 `json:"dealer_win_probability"`
	PushProbability       float64 `json:"push_probability"`
	PlayerBlackjackProb   float64 `json:"player_blackjack_probability"`
	DealerBlackjackProb   float64 `json:"dealer_blackjack_probability"`
	PlayerBustProbability float64 `json:"player_bust_probability"`
	DealerBustProbability float64 `json:"dealer_bust_probability"`
	Player21Probability   float64 `json:"player_21_probability"`
	Dealer21Probability   float64 `json:"dealer_21_probability"`

	// 操作胜率分析
	ActionAnalysis *ActionAnalysisDTO `json:"action_analysis,omitempty"`
}

// ActionAnalysisDTO 操作胜率分析数据传输对象
type ActionAnalysisDTO struct {
	HitWinRate    float64 `json:"hit_win_rate"`    // 要牌胜率
	StandWinRate  float64 `json:"stand_win_rate"`  // 停牌胜率
	DoubleWinRate float64 `json:"double_win_rate"` // 加倍胜率
	SplitWinRate  float64 `json:"split_win_rate"`  // 分牌胜率（如果可分牌）

	// 操作可用性
	CanHit    bool `json:"can_hit"`
	CanStand  bool `json:"can_stand"`
	CanDouble bool `json:"can_double"`
	CanSplit  bool `json:"can_split"`

	// 推荐操作
	RecommendedAction string  `json:"recommended_action"`
	ExpectedValue     float64 `json:"expected_value"` // 推荐操作的期望值

	// 凯利公式相关
	KellyRecommendation *KellyRecommendationDTO `json:"kelly_recommendation,omitempty"`
}

// KellyRecommendationDTO 凯利公式推荐数据传输对象
type KellyRecommendationDTO struct {
	// 当前情况下的凯利比例
	StandardKellyFraction  float64 `json:"standard_kelly_fraction"`  // 普通胜利的凯利比例
	BlackjackKellyFraction float64 `json:"blackjack_kelly_fraction"` // Blackjack胜利的凯利比例
	DoubleKellyFraction    float64 `json:"double_kelly_fraction"`    // 加倍的凯利比例

	// 推荐投注金额
	RecommendedBetAmount   int     `json:"recommended_bet_amount"`   // 基于凯利公式的推荐投注金额
	RecommendedBetFraction float64 `json:"recommended_bet_fraction"` // 推荐投注比例（相对于总筹码）

	// 加倍决策
	ShouldDouble      bool    `json:"should_double"`       // 是否推荐加倍
	DoubleExpectedROI float64 `json:"double_expected_roi"` // 加倍的期望投资回报率

	// 风险评估
	RiskLevel          string  `json:"risk_level"`           // 风险等级 (Low/Medium/High)
	ExpectedGrowthRate float64 `json:"expected_growth_rate"` // 期望资金增长率
}
