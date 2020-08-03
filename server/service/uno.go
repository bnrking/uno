package service

import (
	"errors"
	"math/rand"

	"github.com/jak103/uno/db"
	"github.com/jak103/uno/model"
)

////////////////////////////////////////////////////////////
// Utility functions
////////////////////////////////////////////////////////////

// A simple helper function to pull a card from a game and put it in the players hand.
// THis is used in  a lot of places, so this should be  a nice help
func DrawCardHelper(game *model.Game, player *model.Player) {
	lastIndex := len(game.DrawPile) - 1
	card := game.DrawPile[lastIndex]

	player.Cards = append(player.Cards, card)
	game.DrawPile = game.DrawPile[:lastIndex]
}

// A simpler helper function for getting the player with a matching ID to playerID
// from the list of players in the game.
func GetPlayer(game *model.Game, playerID string) *model.Player {
	for _, item := range game.Players {
		if playerID == item.ID {
			return &item
		}
	}
	return nil
}

// given a player and a card look for the card in players hand and return the index
// If it doesn't exists return -1
func CardFromPlayer(player *model.Player, card *model.Card) int {
	// Loop through all cards the player holds
	for index, item := range player.Cards {
		// check if current loop item matches card provided
		if item.Color == card.Color && item.Value == card.Value {
			// If the card matches return the current index
			return index
		}
	}
	// If we get to this point the player does not hold the card so we return nil
	return -1
}

////////////////////////////////////////////////////////////
// These are all of the functions for the game -> essentially public functions
////////////////////////////////////////////////////////////
func GetGameUpdate(gameID string) (*model.Game, error) {
	database, err := db.GetDb()

	if err != nil {
		return nil, err
	}

	gameData, gameErr := database.LookupGameByID(gameID)

	if gameErr != nil {
		return nil, err
	}

	return gameData, nil
}

func CreatePlayer(name string) (*model.Player, error) {
	database, err := db.GetDb()
	if err != nil {
		return nil, err
	}

	player, err := database.CreatePlayer(name)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func CreateNewGame(gameName string, creatorName string) (*model.Game, *model.Player, error) {
	database, err := db.GetDb()
	if err != nil {
		return nil, nil, err
	}

	creator, err := database.CreatePlayer(creatorName)
	if err != nil {
		return nil, nil, err
	}

	game, err := database.CreateGame(gameName, creator.ID)
	if err != nil {
		return nil, nil, err
	}

	err = database.SaveGame(*game)
	if err != nil {
		return nil, nil, err
	}

	return game, creator, nil
}

func JoinGame(game string, player *model.Player) (*model.Game, error) {
	database, err := db.GetDb()
	if err != nil {
		return nil, err
	}

	gameData, gameErr := database.JoinGame(game, player.ID)

	if gameErr != nil {
		return nil, gameErr
	}

	return gameData, nil
}

func PlayCard(game string, playerID string, card model.Card) (*model.Game, error) {
	database, err := db.GetDb()

	if err != nil {
		return nil, err
	}

	gameData, gameErr := database.LookupGameByID(game)

	if gameErr != nil {
		return nil, err
	}

	if gameData.Players[gameData.CurrentPlayer].ID == playerID {
		hand := gameData.Players[gameData.CurrentPlayer].Cards
		if CheckForCardInHand(card, hand) && (card.Color == gameData.DiscardPile[len(gameData.DiscardPile)-1].Color || card.Value == gameData.DiscardPile[len(gameData.DiscardPile)-1].Value || card.Value == "W4" || card.Value == "W") {
			// Valid card can be played

			gameData.DiscardPile = append(gameData.DiscardPile, card)

			for index, item := range hand {
				if item == card || (item.Value == "W" && card.Value == "W") || (item.Value == "W4" && card.Value == "W4") {
					gameData.Players[gameData.CurrentPlayer].Cards = append(hand[:index], hand[index+1:]...)
					break
				}
			}

			if card.Value == "S" {
				gameData = GoToNextPlayer(gameData)
			}

			if card.Value == "D2" {
				gameData = GoToNextPlayer(gameData)
				gameData = DrawNCards(gameData, 2)
			}

			if card.Value == "W4" {
				gameData = GoToNextPlayer(gameData)
				gameData = DrawNCards(gameData, 4)
			}

			if card.Value == "R" {
				gameData.Direction = !gameData.Direction
			}

			gameData = GoToNextPlayer(gameData)

		}
	}

	err = database.SaveGame(*gameData)

	if err != nil {
		return nil, err
	}

	return gameData, nil
}

func CheckForCardInHand(card model.Card, hand []model.Card) bool {
	for _, c := range hand {
		// the wild cards, W4 and W, don't need to match in color; not for the previous card, and not with the hand. The card itself can become any color.
		if c.Value == card.Value && (c.Color == card.Color || card.Value == "W4" || card.Value == "W") {
			return true
		}
	}
	return false
}

// TODO: Keep track of current card that is top of the deck
func DrawCard(gameID, playerID string) (*model.Game, error) {
	database, err := db.GetDb()

	if err != nil {
		return nil, err
	}

	gameData, gameErr := database.LookupGameByID(gameID)

	if gameErr != nil {
		return nil, err
	}

	if gameData.Players[gameData.CurrentPlayer].ID != playerID {
		return nil, errors.New("Wrong player")
	}

	var drawnCard model.Card
	gameData, drawnCard = DrawTopCard(gameData)
	gameData.Players[gameData.CurrentPlayer].Cards = append(gameData.Players[gameData.CurrentPlayer].Cards, drawnCard)

	gameData = GoToNextPlayer(gameData)

	database.SaveGame(*gameData)

	return gameData, nil
}

func GoToNextPlayer(gameData *model.Game) *model.Game {
	if gameData.Direction {
		gameData.CurrentPlayer++
		gameData.CurrentPlayer %= len(gameData.Players)
	} else {
		gameData.CurrentPlayer--
		if gameData.CurrentPlayer < 0 {
			gameData.CurrentPlayer = len(gameData.Players) - 1
		}
	}

	return gameData
}

func DrawNCards(gameData *model.Game, nCards uint) *model.Game {
	for i := uint(0); i < nCards; i++ {
		var drawnCard model.Card
		gameData, drawnCard = DrawTopCard(gameData)
		gameData.Players[gameData.CurrentPlayer].Cards = append(gameData.Players[gameData.CurrentPlayer].Cards, drawnCard)
	}
	return gameData
}

// TODO: need to deal the actual cards, not just random numbers
func DealCards(game *model.Game) (*model.Game, error) {

	// pick a starting player
	game.CurrentPlayer = rand.Intn(len(game.Players))

	// get a deck
	game.DrawPile = generateShuffledDeck(len(game.Players))

	// give everyone a hand of seven cards
	for k := range game.Players {
		cards := []model.Card{}
		for i := 0; i < 7; i++ {
			var drawnCard model.Card
			game, drawnCard = DrawTopCard(game)
			cards = append(cards, drawnCard)
		}
		game.Players[k].Cards = cards
	}

	// draw a card for the discard
	var drawnCard model.Card
	game, drawnCard = DrawTopCard(game)
	game.DiscardPile = append(game.DiscardPile, drawnCard)

	game.Status = "Playing"

	database, err := db.GetDb()

	if err != nil {
		return nil, err
	}

	// save the new game status
	err = database.SaveGame(*game)

	return game, err
}

func DrawTopCard(game *model.Game) (*model.Game, model.Card) {
	drawnCard := game.DrawPile[len(game.DrawPile)-1]
	game.DrawPile = game.DrawPile[:len(game.DrawPile)-1]
	return game, drawnCard
}

func CheckGameExists(gameID string) (bool, error) {
	database, err := db.GetDb()

	if err != nil {
		return false, err
	}

	_, gameErr := database.LookupGameByID(gameID)

	if gameErr != nil {
		return false, gameErr
	}

	return true, nil
}
