package context

import (
	"errors"
	"fmt"

	"github.com/nleskiw/goplaycards/deck"
)

type Playeractor interface {
	HitOrStand(ownHand []deck.Card, otherHands [][]deck.Card, dealerHand deck.Card) (hitOrStand Action)
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
	fmt.Println("New player added!")
	player.addToWalletNonZero()
	c.players = append(c.players, player)
	return
}

func (c *Context) PlayRound() {
	fmt.Println("Starting new round...")

	if len(c.players) == 0 {
		return
	}
	c.addDealer()
	c.deck.Initialize()
	c.deck.Shuffle()

	c.getBids()
	c.drawInitial()

	if c.dealer.isBlackjack() {
		fmt.Println("Dealer has blackjack!")
		for _, player := range c.players {
			if player.isBlackjack() {
				player.wallet += player.bet
			}
		}
	} else {
		// Dealer does not have blackjack so everyone plays
		for i, player := range c.players {
			if !player.isBlackjack() {
				c.playerTurn(i)
			} else {
				fmt.Println("Player", i+1, "has blackjack!")
			}
		}

		// Dealer plays
		c.dealerRound()

		c.distributeWinnings()
	}
	c.printWinnings()
}

func (c *Context) getBids() {
	for i, player := range c.players {
		fmt.Println("Getting bid for player", i+1)
		player.addToWallet()
		player.bid()
	}
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
		player.clearHand()
		c.drawTwo(player)

	}
	c.drawTwo(c.dealer)
}

func (c *Context) playerTurn(playerIndex int) {

	player := c.players[playerIndex]
	otherPlayers := append(c.players[0:playerIndex], c.players[playerIndex+1:]...)
	otherHands := make([][]deck.Card, len(otherPlayers))
	for i, otherPlayer := range otherPlayers {
		otherHands[i] = otherPlayer.hand
	}
	playerDone := false
	for !playerDone {
		fmt.Println("PLAYER TURN:", playerIndex+1)
		action := player.hitOrStand(otherHands, c.dealer.hand[1])
		switch action {
		case Hit:
			cards, err := c.deck.Draw(1)
			if err != nil {
				panic(err)
			}
			player.addCards(cards)
			if player.isBust() {
				fmt.Println("Player ", playerIndex+1, " busts!")
				player.printHand()
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
	fmt.Println("Dealer drawing:")
	c.dealer.printDealerHand()
}

func (c *Context) printWinnings() {
	for i, player := range c.players {
		fmt.Println("Player ", i+1, " wallet: ", player.wallet)
	}
}
