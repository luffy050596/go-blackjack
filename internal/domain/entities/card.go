// Package entities contains domain entities for the blackjack game.
package entities

import (
	"fmt"
	"strconv"
)

// Suit 花色枚举
type Suit int

const (
	// Hearts represents the hearts suit
	Hearts Suit = iota
	// Diamonds represents the diamonds suit
	Diamonds
	// Clubs represents the clubs suit
	Clubs
	// Spades represents the spades suit
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
	// Ace represents the ace rank (value 1)
	Ace Rank = iota + 1
	// Two represents the two rank (value 2)
	Two
	// Three represents the three rank (value 3)
	Three
	// Four represents the four rank (value 4)
	Four
	// Five represents the five rank (value 5)
	Five
	// Six represents the six rank (value 6)
	Six
	// Seven represents the seven rank (value 7)
	Seven
	// Eight represents the eight rank (value 8)
	Eight
	// Nine represents the nine rank (value 9)
	Nine
	// Ten represents the ten rank (value 10)
	Ten
	// Jack represents the jack rank (value 11)
	Jack
	// Queen represents the queen rank (value 12)
	Queen
	// King represents the king rank (value 13)
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
