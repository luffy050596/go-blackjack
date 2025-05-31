package main

import (
	"testing"
)

func TestHandValue(t *testing.T) {
	tests := []struct {
		name     string
		cards    []Card
		expected int
	}{
		{
			name:     "普通牌组合",
			cards:    []Card{{Hearts, Five}, {Spades, Seven}},
			expected: 12,
		},
		{
			name:     "单张A（软11）",
			cards:    []Card{{Hearts, Ace}, {Spades, Five}},
			expected: 16,
		},
		{
			name:     "A被调整为1",
			cards:    []Card{{Hearts, Ace}, {Spades, King}, {Clubs, Five}},
			expected: 16, // A=1, K=10, 5=5
		},
		{
			name:     "Blackjack",
			cards:    []Card{{Hearts, Ace}, {Spades, King}},
			expected: 21,
		},
		{
			name:     "双A（一个11，一个1）",
			cards:    []Card{{Hearts, Ace}, {Spades, Ace}, {Clubs, Nine}},
			expected: 21, // A=11, A=1, 9=9
		},
		{
			name:     "双A都调整为1",
			cards:    []Card{{Hearts, Ace}, {Spades, Ace}, {Clubs, King}},
			expected: 12, // A=1, A=1, K=10
		},
		{
			name:     "三A",
			cards:    []Card{{Hearts, Ace}, {Spades, Ace}, {Clubs, Ace}, {Diamonds, Eight}},
			expected: 21, // A=11, A=1, A=1, 8=8
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand := Hand{}
			for _, card := range tt.cards {
				hand.AddCard(card)
			}

			result := hand.Value()
			if result != tt.expected {
				t.Errorf("Hand.Value() = %d, expected %d for cards %v", result, tt.expected, tt.cards)
			}
		})
	}
}

func TestHandIsBlackjack(t *testing.T) {
	tests := []struct {
		name     string
		cards    []Card
		expected bool
	}{
		{
			name:     "A + K = Blackjack",
			cards:    []Card{{Hearts, Ace}, {Spades, King}},
			expected: true,
		},
		{
			name:     "A + Q = Blackjack",
			cards:    []Card{{Hearts, Ace}, {Spades, Queen}},
			expected: true,
		},
		{
			name:     "A + J = Blackjack",
			cards:    []Card{{Hearts, Ace}, {Spades, Jack}},
			expected: true,
		},
		{
			name:     "A + 10 = Blackjack",
			cards:    []Card{{Hearts, Ace}, {Spades, Ten}},
			expected: true,
		},
		{
			name:     "三张牌21点不是Blackjack",
			cards:    []Card{{Hearts, Seven}, {Spades, Seven}, {Clubs, Seven}},
			expected: false,
		},
		{
			name:     "A + 9 不是Blackjack",
			cards:    []Card{{Hearts, Ace}, {Spades, Nine}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand := Hand{}
			for _, card := range tt.cards {
				hand.AddCard(card)
			}

			result := hand.IsBlackjack()
			if result != tt.expected {
				t.Errorf("Hand.IsBlackjack() = %t, expected %t for cards %v", result, tt.expected, tt.cards)
			}
		})
	}
}

func TestHandIsSoft(t *testing.T) {
	tests := []struct {
		name     string
		cards    []Card
		expected bool
	}{
		{
			name:     "A + 5 = 软16",
			cards:    []Card{{Hearts, Ace}, {Spades, Five}},
			expected: true,
		},
		{
			name:     "A + K + 5 = 硬16",
			cards:    []Card{{Hearts, Ace}, {Spades, King}, {Clubs, Five}},
			expected: false,
		},
		{
			name:     "无A = 硬点数",
			cards:    []Card{{Hearts, Ten}, {Spades, Five}},
			expected: false,
		},
		{
			name:     "Blackjack = 软21",
			cards:    []Card{{Hearts, Ace}, {Spades, King}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand := Hand{}
			for _, card := range tt.cards {
				hand.AddCard(card)
			}

			result := hand.IsSoft()
			if result != tt.expected {
				t.Errorf("Hand.IsSoft() = %t, expected %t for cards %v", result, tt.expected, tt.cards)
			}
		})
	}
}
