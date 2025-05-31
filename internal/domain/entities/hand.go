package entities

import "strings"

// Hand 手牌结构
type Hand struct {
	Cards []Card
}

func NewHand() *Hand {
	return &Hand{
		Cards: make([]Card, 0, 8),
	}
}

// AddCard 添加牌到手牌
func (h *Hand) AddCard(card Card) {
	h.Cards = append(h.Cards, card)
}

// Value 计算手牌点数（动态调整A牌点数）
func (h *Hand) Value() int {
	value := 0
	aceCount := 0

	for _, card := range h.Cards {
		if card.IsAce() {
			aceCount++
		}
		value += card.BaseValue()
	}

	// 当总点数超过21时，将A从11调整为1
	for aceCount > 0 && value > 21 {
		value -= 10
		aceCount--
	}

	return value
}

// AceCount 获取A牌数量
func (h *Hand) AceCount() int {
	count := 0
	for _, card := range h.Cards {
		if card.IsAce() {
			count++
		}
	}
	return count
}

// IsSoft 获取软点数信息（是否包含被当作11点的A）
func (h *Hand) IsSoft() bool {
	if h.AceCount() == 0 {
		return false
	}

	value := 0
	aceCount := 0

	for _, card := range h.Cards {
		if card.IsAce() {
			aceCount++
		}
		value += card.BaseValue()
	}

	// 模拟调整过程，看最终是否还有A被当作11点
	adjustments := 0
	for aceCount > adjustments && value-adjustments*10 > 21 {
		adjustments++
	}

	// 如果调整次数少于A的总数，说明还有A被当作11点
	return adjustments < aceCount
}

// IsBlackjack 是否为Blackjack
func (h *Hand) IsBlackjack() bool {
	return len(h.Cards) == 2 && h.Value() == 21
}

// IsBust 是否爆牌
func (h *Hand) IsBust() bool {
	return h.Value() > 21
}

// String 显示手牌
func (h *Hand) String() string {
	var cards []string
	for _, card := range h.Cards {
		cards = append(cards, card.String())
	}
	return strings.Join(cards, " ")
}
