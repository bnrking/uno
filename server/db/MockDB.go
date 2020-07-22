package db

import (
	"github.com/google/uuid"
	"github.com/jak103/uno/model"
)

// MockDB is an implemenation declaring the unit test db
type mockDB struct{}

// HasGame checks to see if a game with the given ID exists in the database.
func (db *mockDB) HasGameByPassword(password string) bool {
	return password == "12234"
}

func (db *mockDB) HasGameByID(id uuid.UUID) bool {
	return true
}

// CreateGame a game with the given ID. Perhaps this should instead just return an id?
func (db *mockDB) CreateGame(id uuid.UUID) model.Game {
	myGame := model.Game{ID: uuid.Nil, Password: "12234"}
	return myGame
}

// LookupGame looks up an existing game in the database.
func (db *mockDB) LookupGameByID(id uuid.UUID) model.Game {
	myGame := model.Game{ID: uuid.Nil, Password: "12234"}
	return myGame
}

// LookupGame looks up an existing game in the database.
func (db *mockDB) LookupGameByPassword(password string) model.Game {
	myGame := model.Game{ID: uuid.Nil, Password: "12234"}
	return myGame
}

// JoinGame mockDB a player to a game.
func (db *mockDB) JoinGame(id uuid.UUID, username string) {
	return
}

func newMockDB() *mockDB {
	return new(mockDB)
}
