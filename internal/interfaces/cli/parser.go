package cli

import (
	"strings"

	"github.com/luffy050596/go-blackjack/internal/domain/entities"
)

// ParsePlayerInput 解析玩家输入
func ParsePlayerInput(input string) entities.PlayerAction {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case entities.InputHit, entities.InputHitFull:
		return entities.ActionHit
	case entities.InputStand, entities.InputStandFull:
		return entities.ActionStand
	case entities.InputDouble, entities.InputDoubleFull, entities.InputDoubleDown:
		return entities.ActionDoubleDown
	case entities.InputQuit, entities.InputQuitFull:
		return entities.ActionQuit
	default:
		return entities.ActionInvalid
	}
}
