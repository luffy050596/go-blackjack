package main

import (
	"math/rand/v2"
	"time"
)

// Deck 牌堆结构
type Deck struct {
	Cards []Card
}

// NewDeck 创建新牌堆
func NewDeck() *Deck {
	deck := &Deck{}

	// 创建52张牌
	for suit := Hearts; suit <= Spades; suit++ {
		for rank := Ace; rank <= King; rank++ {
			deck.Cards = append(deck.Cards, Card{Suit: suit, Rank: rank})
		}
	}

	// 洗牌
	deck.Shuffle()
	return deck
}

// Shuffle 洗牌
func (d *Deck) Shuffle() {
	seed := uint64(time.Now().UnixNano())
	rd := rand.New(rand.NewPCG(uint64(seed), uint64(seed>>32)))

	for i := len(d.Cards) - 1; i > 0; i-- {
		j := rd.IntN(i + 1)
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	}
}

// Deal 发牌
func (d *Deck) Deal() Card {
	if len(d.Cards) == 0 {
		panic("牌堆已空！")
	}
	card := d.Cards[0]
	d.Cards = d.Cards[1:]
	return card
}
