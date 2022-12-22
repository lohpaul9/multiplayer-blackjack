package context

// Stretch features: Implement allowing player to see
// other players hands.

import (
	"fmt"

	"github.com/nleskiw/goplaycards/deck"
)

type player struct {
	hand   []deck.Card
	bet    float64
	wallet float64
	actor  Playeractor
}

func (p *player) addToWallet() {
	p.wallet += p.actor.AddToWallet()
}

func (p *player) addToWalletNonZero() {
	toAdd := p.actor.AddToWallet()
	for toAdd <= 0 {
		fmt.Println("Give a positive number amount to add to your wallet.")
		toAdd = p.actor.AddToWallet()
	}
	p.wallet += toAdd
}

func (p *player) bid() {
	bid := p.actor.Bet(p.wallet)
	for bid > p.wallet {
		fmt.Println("You can't bet more than you have in your wallet.")
		bid = p.actor.Bet(p.wallet)
	}
	p.bet = bid
	p.wallet -= bid
}

func (p *player) hitOrStand(otherHands [][]deck.Card, dealerHand deck.Card) Action {
	return p.actor.HitOrStand(p.hand, otherHands, dealerHand)
}

func (p *player) addCards(cards []deck.Card) {
	p.hand = append(p.hand, cards...)
}

// handTotal returns the numerical value of a Blackjack hand
func (p *player) handTotal() int {
	hand := p.hand
	total := 0
	numberOfAces := 0
	for _, card := range hand {
		if card.Value.Name == "Ace" {
			numberOfAces = numberOfAces + 1
		} else {
			if card.Facecard() {
				total = total + 10
			} else {
				total = total + card.Value.Value
			}
		}
	}

	// If there's at least one Ace, deal with it.
	// In multi-shoe decks, there could be many Aces (more than 4) in a hand.
	if numberOfAces > 0 {
		// All but the last Ace must be a one, because 11 + 11 = 22 (bust)
		// This loop shouldn't run if there's only one Ace
		for numberOfAces > 1 {
			total = total + 1
			numberOfAces = numberOfAces - 1
		}
		// There should now only be one Ace
		// if the last Ace being 11 doesn't cause a bust, make it an 11
		if total+11 > 21 {
			total = total + 1
		} else {
			// If 11 doesn't cause a bust, make it worth 11
			total = total + 11
		}
	}
	return total
}

// Returns true if a hand is bust / over 21
func (p *player) isBust() bool {
	return p.handTotal() > 21
}

// Returns true if a hand is a Blacjack (Ace + [10 | K | Q | J])
func (p *player) isBlackjack() bool {
	hand := p.hand
	// A Blackjack is exactly one Ace and Exactly one 10, K, Q, or A
	if len(hand) != 2 {
		return false
	}
	// In the goplaycards library, the values enumerate from 2 to 14.
	// J = 11, Q = 12, K = 13, Ace = 14
	if hand[0].Value.Name == "Ace" {
		if hand[1].Value.Value >= 10 && hand[1].Value.Value <= 13 {
			return true
		}
	}
	if hand[1].Value.Name == "Ace" {
		if hand[0].Value.Value >= 10 && hand[0].Value.Value <= 13 {
			return true
		}
	}
	return false
}

func (p *player) printHand() {
	hand := p.hand
	fmt.Printf("Player Hand: ")
	for _, card := range hand {
		fmt.Printf("%s  ", card.ToStr())
	}
	fmt.Printf(" Total: %d\n", p.handTotal())
}

func (p *player) printDealerHandHidden() {
	hand := p.hand
	fmt.Printf("Dealer Hand: ")
	fmt.Printf("XX  %s  \n", hand[1].ToStr())
}

func (p *player) printDealerHand() {
	hand := p.hand
	fmt.Printf("Dealer Hand: ")
	for _, card := range hand {
		fmt.Printf("%s  ", card.ToStr())
	}
	fmt.Printf(" Total: %d\n", p.handTotal())
}

func (p *player) clearHand() {
	p.hand = nil
}
