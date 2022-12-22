/*

Game logic adapted from https://github.com/nleskiw/blackjack on 22/12/2022
Modifications were made to allow the 'game context' to be played generically using provided callback functions

*/

package main

import (
	"blackjack/context"
	"fmt"
	"strconv"

	"github.com/nleskiw/goplaycards/deck"
)

type cmdLinePlayer struct {
}

func (c cmdLinePlayer) AddToWallet() (amount float64) {
	return float64(getInteger("How much would you like to add to your wallet? "))
}

func (c cmdLinePlayer) HitOrStand(otherHands [][]deck.Card, dealerHand deck.Card) (hitOrStand context.Action) {
	return getPlayerActionAction()
}

func (c cmdLinePlayer) Bet(wallet float64) (bet float64) {
	str := fmt.Sprintf("%f", wallet)
	return float64(getInteger("You currently have " + str + "How much would you like to bet ($5 increments)? "))
}

// getPlayerAction determines what the player will do
func getPlayerActionAction() context.Action {
	input := ""
	for {
		input = getString("[H]it or [S]tand? ")
		if input == "hit" || input == "Hit" || input == "H" || input == "h" {
			return context.Hit
		} else if input == "stand" || input == "Stand" || input == "S" || input == "s" {
			return context.Stand
		} else {
			fmt.Println("Invalid option. H to Hit or S to Stand. ")
		}
	}
}

// getString gets an arbitrary string from the user with a prompt.
func getString(prompt string) string {
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)
	return input
}

// getInteger gets an arbitrary int from the user with a prompt.
// Retries until a valid integer is entered.
func getInteger(prompt string) int {
	valid := true
	input := getString(prompt)
	integer, err := strconv.Atoi(input)
	if err != nil {
		valid = false
	}
	for valid == false {
		fmt.Println("Can't convert your answer into an integer.")
		input = getString(prompt)
		integer, err = strconv.Atoi(input)
		if err == nil {
			valid = true
		}
	}
	return integer
}

func printLicense() {
	fmt.Println("goplaycards Copyright (C) 2017  Nicholas Leskiw")
	fmt.Println("This program comes with ABSOLUTELY NO WARRANTY; for details please visit")
	fmt.Println("<https://www.gnu.org/licenses/gpl-3.0.txt>")
	fmt.Println("This is free software, and you are welcome to redistribute it")
	fmt.Printf("under certain conditions. Please see the GPLv3 license at the URL listed above.\n\n")
}

func main() {
	printLicense()

	c := context.NewGame()
	c.AddPlayer(new(cmdLinePlayer))

	c.PlayRound()
}
