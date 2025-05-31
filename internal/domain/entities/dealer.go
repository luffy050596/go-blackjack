package entities

// Dealer 庄家结构
type Dealer struct {
	Hand *Hand
}

// NewDealer creates a new dealer with an empty hand
func NewDealer() *Dealer {
	return &Dealer{
		Hand: NewHand(),
	}
}

// ResetRound resets the dealer's hand for a new round
func (d *Dealer) ResetRound() {
	d.Hand = NewHand()
}
