package main

import "fmt"

// 基本策略建议
type StrategyAdvice string

const (
	Hit    StrategyAdvice = "要牌"
	Stand  StrategyAdvice = "停牌"
	Double StrategyAdvice = "加倍"
	Split  StrategyAdvice = "分牌"
)

// 基本策略引擎
type BasicStrategy struct{}

// 获取基本策略建议
func (bs *BasicStrategy) GetAdvice(playerHand *Hand, dealerUpCard Card) StrategyAdvice {
	playerValue := playerHand.Value()
	dealerValue := dealerUpCard.BaseValue()

	// 软点数策略（包含被当作11点的A）
	if playerHand.IsSoft() {
		return bs.getSoftHandStrategy(playerValue, dealerValue)
	}

	// 硬点数策略
	return bs.getHardHandStrategy(playerValue, dealerValue)
}

// 软点数策略
func (bs *BasicStrategy) getSoftHandStrategy(playerValue, dealerValue int) StrategyAdvice {
	switch playerValue {
	case 13, 14: // A,2 或 A,3
		if dealerValue >= 5 && dealerValue <= 6 {
			return Double
		}
		return Hit
	case 15, 16: // A,4 或 A,5
		if dealerValue >= 4 && dealerValue <= 6 {
			return Double
		}
		return Hit
	case 17: // A,6
		if dealerValue >= 3 && dealerValue <= 6 {
			return Double
		}
		return Hit
	case 18: // A,7
		if dealerValue >= 3 && dealerValue <= 6 {
			return Double
		}
		if dealerValue >= 9 || dealerValue == 11 { // 9,10,A
			return Hit
		}
		return Stand
	case 19, 20, 21: // A,8 A,9 A,10
		return Stand
	default:
		return Hit
	}
}

// 硬点数策略
func (bs *BasicStrategy) getHardHandStrategy(playerValue, dealerValue int) StrategyAdvice {
	switch {
	case playerValue <= 8:
		return Hit
	case playerValue == 9:
		if dealerValue >= 3 && dealerValue <= 6 {
			return Double
		}
		return Hit
	case playerValue == 10:
		if dealerValue <= 9 {
			return Double
		}
		return Hit
	case playerValue == 11:
		if dealerValue <= 10 {
			return Double
		}
		return Hit
	case playerValue == 12:
		if dealerValue >= 4 && dealerValue <= 6 {
			return Stand
		}
		return Hit
	case playerValue >= 13 && playerValue <= 16:
		if dealerValue <= 6 {
			return Stand
		}
		return Hit
	case playerValue >= 17:
		return Stand
	default:
		return Hit
	}
}

// 显示策略建议
func (bs *BasicStrategy) ShowAdvice(playerHand *Hand, dealerUpCard Card) {
	advice := bs.GetAdvice(playerHand, dealerUpCard)

	fmt.Printf("你的手牌: %s (点数: %d", playerHand.String(), playerHand.Value())

	if playerHand.IsSoft() {
		fmt.Printf(", 软点数")
	} else {
		fmt.Printf(", 硬点数")
	}

	if playerHand.AceCount() > 0 {
		fmt.Printf(", %d张A", playerHand.AceCount())
	}

	fmt.Printf(")\n")
	fmt.Printf("庄家明牌: %s\n", dealerUpCard.String())
	fmt.Printf("建议: %s\n", advice)
	fmt.Println()
}
