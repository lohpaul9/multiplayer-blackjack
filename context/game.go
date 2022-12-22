package context

import (
	"errors"
	"fmt"

	"github.com/nleskiw/goplaycards/deck"
)

// The main interface to implement for a single player
// The game context will call these methods to interact with the player
// Example usages could be a CLI player, or a web player
type Playeractor interface {
	HitOrStand(ownHand []deck.Card, otherHands [][]deck.Card, dealerHand deck.Card) (hitOrStand Action)
	Bet(wallet float64) (bid float64)
	AddToWallet() (amount float64)
}

// Action is the action a player can take (Hit or Stand)
type Action uint

const (
	Hit Action = iota
	Stand
)

// Context is the main game context, it holds the deck, players, and dealer
// Create a context instance, add players and run the game
type Context struct {
	deck    *deck.Deck
	players []*player
	dealer  *player
}

// NewGame creates a new game context
func NewGame() *Context {
	context := &Context{}
	context.deck = &deck.Deck{}
	return context
}

// AddPlayer adds a player to the game (max 3)
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

// PlayRound starts a new round of blackjack
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

// getBids gets the bids from all players (and asks if they want to add to wallet)
func (c *Context) getBids() {
	for i, player := range c.players {
		fmt.Println("Getting bid for player", i+1)
		player.addToWallet()
		player.bid()
	}
}

// Initializes the dealer
func (c *Context) addDealer() {
	c.dealer = &player{}
}

// draws two cards for a player
func (c *Context) drawTwo(p *player) {
	cards, err := c.deck.Draw(2)
	if err != nil {
		panic(err)
	}
	p.addCards(cards)
}

// calculates and distributes winnings to all players wallets at the end of the turn
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

// draws the initial cards for all players and the dealer
func (c *Context) drawInitial() {
	for _, player := range c.players {
		player.clearHand()
		c.drawTwo(player)

	}
	c.drawTwo(c.dealer)
}

// a single winner's round to hit and stand till they either bust or stand
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

// draws cards for dealer at end of round
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

// helper function to print the winnings of all players
func (c *Context) printWinnings() {
	for i, player := range c.players {
		fmt.Println("Player ", i+1, " wallet: ", player.wallet)
	}
}
