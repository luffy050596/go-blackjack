package entities

// Dealer 庄家结构
type Dealer struct {
	Hand *Hand
}

func NewDealer() *Dealer {
	return &Dealer{
		Hand: NewHand(),
	}
}

func (d *Dealer) ResetRound() {
	d.Hand = NewHand()
}
