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
