package entities

import (
	"fmt"
	"strconv"
)

// Suit 花色枚举
type Suit int

const (
	Hearts Suit = iota
	Diamonds
	Clubs
	Spades
)

func (s Suit) String() string {
	switch s {
	case Hearts:
		return "♥"
	case Diamonds:
		return "♦"
	case Clubs:
		return "♣"
	case Spades:
		return "♠"
	default:
		return "?"
	}
}

// Rank 牌面枚举
type Rank int

const (
	Ace Rank = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

func (r Rank) String() string {
	switch r {
	case Ace:
		return "A"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	default:
		return strconv.Itoa(int(r))
	}
}

// Card 扑克牌结构
type Card struct {
	Suit Suit
	Rank Rank
}

// IsEmpty 检查是否为空牌（零值）
func (c Card) IsEmpty() bool {
	return c.Suit == 0 && c.Rank == 0
}

func (c Card) String() string {
	if c.IsEmpty() {
		return "空牌"
	}
	return fmt.Sprintf("%s%s", c.Rank, c.Suit)
}

// BaseValue 获取牌的基础点数
func (c Card) BaseValue() int {
	switch c.Rank {
	case Ace:
		return 11 // A的默认值是11
	case Jack, Queen, King:
		return 10
	default:
		return int(c.Rank)
	}
}

// Value 获取牌的点数（为了向后兼容保留此方法）
func (c Card) Value() int {
	return c.BaseValue()
}

// IsAce 判断是否为A牌
func (c Card) IsAce() bool {
	return c.Rank == Ace
}
