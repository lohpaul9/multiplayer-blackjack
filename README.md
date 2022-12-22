Blackjack
=========

Golang implementation of Blackjack (21) with a generic interface to play the game
logic using any extendable logic.

Installation instructions 
------------
get module using
```
go get github.com/lohpaul9/multiplayer-blackjack
```

Basic Usage
------------
```
c := context.NewGame() // Create a new game context
c.AddPlayer(new(PlayerActionImplementation)) // Add players as desired
c.PlayRound() // Play a single round of the game
```

The game context will, at appropriate moments, call the callback functions provided by the any implementation of the ```PlayerAction``` interface. 

Ideas for usages of the interface would be to link up the interface to a CLI, web app game etc. 

* Example code is given in main.go, where a basic CLI implementation of the interface demonstrates the reusability of the client. 

Table rules:
------------
Table rules was refactored from: https://github.com/nleskiw/blackjack

* Basic play only (Hit / Stand) 
* Dealer stands on soft 17.
* Single deck
* Reshuffle if less than 17 cards in deck.
* Bet in $5 increments only
* Player starts with $100

AFAIK 17 is the most cards you'd need with one deck:
Player (A A A A 2 2 2 2 3 3 3) Dealer (3 4 4 4 4)

Requires https://github.com/nleskiw/goplaycards

Report bugs and/or style suggestions via Github issues.