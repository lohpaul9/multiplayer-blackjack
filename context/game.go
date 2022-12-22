package context

//to do - change winnings to floats

import (
	"errors"

	"github.com/nleskiw/goplaycards/deck"
)

type Playeractor interface {
	HitOrStand(otherHands [][]deck.Card, dealerHand deck.Card) (hitOrStand Action)
	Bet(wallet float64) (bid float64)
	AddToWallet() (amount float64)
}

type Action uint

const (
	Hit Action = iota
	Stand
)

type Context struct {
	deck    *deck.Deck
	players []*player
	dealer  *player
}

func NewGame() *Context {
	context := &Context{}
	context.deck = &deck.Deck{}
	return context
}

func (c *Context) AddPlayer(actor Playeractor) (err error) {
	if len(c.players) >= 3 {
		return errors.New("too many players")
	}
	player := &player{actor: actor}
	player.addToWalletNonZero()
	c.players = append(c.players, player)
	return
}

func (c *Context) PlayRound() {

	if len(c.players) == 0 {
		return
	}
	c.addDealer()
	c.deck.Initialize()
	c.deck.Shuffle()

	for _, player := range c.players {
		player.addToWallet()
		player.bid()
	}

	c.drawInitial()

	if c.dealer.isBlackjack() {
		for _, player := range c.players {
			if player.isBlackjack() {
				player.wallet += player.bet
			}
		}
		return
	}

	// Dealer does not have blackjack so everyone plays
	for i, player := range c.players {
		if !player.isBlackjack() {
			c.playerRound(i)
		}
	}

	// Dealer plays
	c.dealerRound()
	c.dealer.printHand()

	c.distributeWinnings()
}

func (c *Context) addDealer() {
	c.dealer = &player{}
}

func (c *Context) drawTwo(p *player) {
	cards, err := c.deck.Draw(2)
	if err != nil {
		panic(err)
	}
	p.addCards(cards)
}

func (c *Context) distributeWinnings() {
	// Check who won and distribute winnings
	for _, player := range c.players {
		if player.isBlackjack() {
			// Player is blackjack (always beats dealer)
			player.wallet += player.bet * 2.5
		} else if c.dealer.isBust() {
			// Player not blackjack AND dealer bust
			if !player.isBust() {
				// Player not bust so wins
				player.wallet += player.bet * 2
			} else {
				// Player bust and dealer bust so push
				player.wallet += player.bet
			}
		} else {
			// Player not blackjack and dealer not bust
			if !player.isBust() {
				if player.handTotal() > c.dealer.handTotal() {
					player.wallet += player.bet * 2
				} else if player.handTotal() == c.dealer.handTotal() {
					player.wallet += player.bet
				}
			}
		}
	}
}

func (c *Context) drawInitial() {
	for _, player := range c.players {
		c.drawTwo(player)
	}
	c.drawTwo(c.dealer)
}

func (c *Context) playerRound(playerIndex int) {
	player := c.players[playerIndex]
	otherPlayers := append(c.players[0:playerIndex], c.players[playerIndex+1:]...)
	otherHands := make([][]deck.Card, len(otherPlayers))
	for i, otherPlayer := range otherPlayers {
		otherHands[i] = otherPlayer.hand
	}
	playerDone := false
	for !playerDone {
		action := player.hitOrStand(otherHands, c.dealer.hand[0])
		switch action {
		case Hit:
			player.printHand()
			cards, err := c.deck.Draw(1)
			if err != nil {
				panic(err)
			}
			player.addCards(cards)
			if player.isBust() {
				playerDone = true
			}
		case Stand:
			playerDone = true
		}
	}
}

func (c *Context) dealerRound() {
	dealerDone := false
	for !dealerDone {
		if c.dealer.handTotal() >= 17 {
			dealerDone = true
		} else {
			cards, err := c.deck.Draw(1)
			if err != nil {
				panic(err)
			}
			c.dealer.addCards(cards)
		}
	}
	c.dealer.printHand()
}
