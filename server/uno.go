package main

import (
	"errors"
	"fmt"

	"github.com/jak103/uno/db"
	"github.com/jak103/uno/model"
)

type InvalidPasswordError struct {
	Password string
	Err      error
}

func (e InvalidPasswordError) Error() string {
	return e.Err.Error() + ": " + e.Password
}

type InvalidPlayerError struct {
	Err error
}

func (e InvalidPlayerError) Error() string {
	return e.Err.Error()
}

////////////////////////////////////////////////////////////
// Utility functions used in place of firebase
////////////////////////////////////////////////////////////
func randColor(i int) string {
	switch i {
	case 0:
		return "red"
	case 1:
		return "blue"
	case 2:
		return "green"
	case 3:
		return "yellow"
	}
	return ""
}

////////////////////////////////////////////////////////////
// Utility functions
////////////////////////////////////////////////////////////

// TODO: make sure this reflects on the front end
// func checkForWinner(game *model.Game) string {
// 	for k := range game.Players {
// 		if len(allCards[players[k]]) == 0 {
// 			return players[k]
// 		}
// 	}
// 	return ""
// }

////////////////////////////////////////////////////////////
// These are all of the functions for the game -> essentially public functions
////////////////////////////////////////////////////////////
func updateGame(game string, reqPlayer *model.Player) (*model.Game, error) {
	database, err := db.GetDb()

	if err != nil {
		return nil, err
	}

	gameData, gameErr := database.LookupGameByID(game)

	if gameErr != nil {
		return nil, err
	}

	found := false
	for i := 0; i < len(gameData.Players); i++ {
		loopPlayer := gameData.Players[i]
		if loopPlayer.ID == reqPlayer.ID {
			found = true
			break
		}
	}

	if !found {
		e := errors.New("Player not in game, cannot start")
		return nil, InvalidPlayerError{Err: e}
	}

	return gameData, nil
}

func startGame(game string, player *model.Player) (*model.Game, error) {
	gameData, err := updateGame(game, player)
	if err != nil {
		return nil, err
	}

	database, err := db.GetDb()

	if err != nil {
		return nil, err
	}

	if gameData.Status != "Playing" {
		gameData.Status = "Playing"
	}

	err = database.SaveGame(*gameData)

	if err != nil {
		return nil, err
	}

	return gameData, nil
}

func createNewGame(player *model.Player, password string, name string) (*model.Game, error) {
	database, err := db.GetDb()

	if err != nil {
		return nil, err
	}

	if player == nil {
		panic("Player should not be nil")
	}
	fmt.Println(player)

	game, err := database.CreateGame()
	game.Host = player.ID
	game.Name = name
	game.Password = password
	fmt.Println("Setting game host to " + player.ID)
	database.SaveGame(*game)
	database.JoinGame(game.ID, player.ID)

	if err != nil {
		return nil, err
	}

	return game, nil
}

func joinGame(game string, player *model.Player, password string) (*model.Game, error) {
	database, err := db.GetDb()
	if err != nil {
		return nil, err
	}
	gameData, err := database.LookupGameByID(game)
	if err != nil {
		return nil, err
	}

	if password != gameData.Password {
		return nil, InvalidPasswordError{Password: password, Err: errors.New("Invalid password")}
	}

	gameData, gameErr := database.JoinGame(game, player.ID)

	if gameErr != nil {
		return nil, gameErr
	}

	return gameData, nil
}

func playCard(game string, player *model.Player, card model.Card) (*model.Game, error) {
	database, err := db.GetDb()

	if err != nil {
		return nil, err
	}

	gameData, gameErr := database.LookupGameByID(game)

	if gameErr != nil {
		return nil, err
	}

	// if gameData.CurrentPlayer == username {
	// 	cards := allCards[username]
	// 	if card.Color == currCard[0].Color || card.Value == currCard[0].Value {
	// 		// Valid card can be played
	// 		playerIndex = (playerIndex + 1) % len(players)
	// 		currPlayer = players[playerIndex]
	// 		currCard[0] = card

	// 		for index, item := range cards {
	// 			if item == currCard[0] {
	// 				allCards[username] = append(cards[:index], cards[index+1:]...)
	// 				break
	// 			}
	// 		}
	// 	}
	// 	return true
	// }
	return gameData, nil
}

// TODO: Keep track of current card that is top of the deck
func drawCard(game string, player *model.Player) (*model.Game, error) {
	database, err := db.GetDb()

	if err != nil {
		return nil, err
	}

	gameData, gameErr := database.LookupGameByID(game)

	if gameErr != nil {
		return nil, err
	}

	// if checkID(game) && username == currPlayer {
	// 	playerIndex = (playerIndex + 1) % len(players)
	// 	currPlayer = players[playerIndex]
	// 	// TODO: Use deck utils instead
	// 	//allCards[username] = append(allCards[username], newRandomCard()[0])
	// 	return true
	// }
	return gameData, nil
}

// TODO: need to deal the actual cards, not just random numbers
func dealCards(game string, player *model.Player) (*model.Game, error) {
	database, err := db.GetDb()

	if err != nil {
		return nil, err
	}

	gameData, gameErr := database.LookupGameByID(game)

	if gameErr != nil {
		return nil, err
	}

	// The game has started, no more players are joining
	// loop through players, set their cards
	// gameStarted = true
	// currPlayer = players[rand.Intn(len(players))]
	// deck := generateShuffledDeck()

	// for k := range players {
	// 	cards := []model.Card{}
	// 	for i := 0; i < 7; i++ {

	// 		drawnCard := deck[len(deck)-1]
	// 		deck = deck[:len(deck)-1]
	// 		cards = append(cards, drawnCard)
	// 		//cards = append(cards, model.Card{rand.Intn(10), randColor(rand.Intn(4))})
	// 	}
	// 	allCards[players[k]] = cards
	// }

	// currCard = deck
	//currCard = newRandomCard()

	return gameData, nil
}
