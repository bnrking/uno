package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jak103/uno/db"
	"github.com/jak103/uno/model"
	"github.com/labstack/echo/v4"
)

var sim bool = true

type Response struct {
	ValidGame bool                   `json:"valid"` // Valid game id/game id is in JWT
	Payload   map[string]interface{} `json:"payload"`
}

type request = map[string]interface{}

func setupRoutes(e *echo.Echo) {
	e.POST("/games", newGame)
	e.GET("/games/:game", update)
	e.POST("/games/:game/start", start)
	e.POST("/login", login)
	e.POST("/games/:game/join", join)
	e.POST("/games/:game/play", play)
	e.POST("/games/:game/draw", draw)
}

func newGame(c echo.Context) error {
	authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
	player, validPlayer, err := getPlayerFromHeader(authHeader)

	if err != nil {
		return err
	}

	if !validPlayer {
		return c.JSONPretty(http.StatusUnauthorized, &Response{false, nil}, " ")
	}

	// Get request body.
	body := make(request)
	err = getBody(c, &body)
	if err != nil {
		return err
	}
	// Get password or default it to empty string.
	var password string
	rawPassword, ok := body["password"]
	if !ok {
		password = ""
	} else {
		password = rawPassword.(string)
	}
	// Get name or default it to "NAME's game"
	var name string
	rawName, ok := body["name"]
	if !ok {
		name = player.Name + "'s Game"
	} else {
		name = rawName.(string)
	}

	game, gameErr := createNewGame(player, password, name)

	if gameErr != nil {
		return gameErr
	}
	c.Response().Header().Set(echo.HeaderLocation, "/games/"+game.ID)
	return c.JSONPretty(http.StatusCreated, &Response{true, newPayload(game)}, "  ")
}

func login(c echo.Context) error {
	rawUsername := getPostParam(c, "username")
	if rawUsername == nil {
		return errors.New("Empty username")
	}
	username := rawUsername.(string)

	database, err := db.GetDb()
	if err != nil {
		return err
	}

	player, playerErr := database.CreatePlayer(username)

	if playerErr != nil {
		return playerErr
	}
	fmt.Println(player.ID)

	token, err := newJWT(username, player.ID)

	if err != nil {
		return err
	}

	return c.JSONPretty(http.StatusOK, &Response{true, makeJWTPayload(token)}, "  ")
}

func join(c echo.Context) error {
	authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
	player, validPlayer, err := getPlayerFromHeader(authHeader)

	if err != nil {
		return err
	}

	if !validPlayer {
		return c.JSONPretty(http.StatusUnauthorized, &Response{false, nil}, " ")
	}

	// Get request body.
	body := make(request)
	err = getBody(c, &body)
	if err != nil {
		return err
	}
	// Get password or default it to empty string.
	var password string
	rawPassword, ok := body["password"]
	if !ok {
		password = ""
	} else {
		password = rawPassword.(string)
	}

	game, err := joinGame(c.Param("game"), player, password)

	if err != nil {
		var passwordError *InvalidPasswordError
		if errors.As(err, &passwordError) {
			return c.JSONPretty(http.StatusBadRequest, map[string]interface{}{
				"error": "Invalid password.",
			}, " ")
		}
		return err
	}

	return c.JSONPretty(http.StatusOK, &Response{true, newPayload(game)}, "  ")
}

func start(c echo.Context) error {
	authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
	player, validPlayer, err := getPlayerFromHeader(authHeader)

	if err != nil {
		return err
	}

	if !validPlayer {
		return c.JSONPretty(http.StatusUnauthorized, &Response{false, nil}, " ")
	}

	game, err := startGame(c.Param("game"), player)
	if err != nil {
		if _, ok := err.(*InvalidPlayerError); ok {
			return c.JSONPretty(http.StatusForbidden, map[string]interface{}{
				"error": "You are not a member of this game.",
			}, " ")
		}
		return err
	}
	return c.JSONPretty(http.StatusOK, &Response{true, newPayload(game)}, "  ")
}

func update(c echo.Context) error {
	authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
	player, validPlayer, err := getPlayerFromHeader(authHeader)

	if err != nil {
		return err
	}

	if !validPlayer {
		return c.JSONPretty(http.StatusUnauthorized, &Response{false, nil}, " ")
	}

	game, err := startGame(c.Param("game"), player)
	if err != nil {
		if _, ok := err.(*InvalidPlayerError); ok {
			return c.JSONPretty(http.StatusForbidden, map[string]interface{}{
				"error": "You are not a member of this game.",
			}, " ")
		}
		return err
	}
	return c.JSONPretty(http.StatusOK, &Response{true, newPayload(game)}, "  ")
}

func play(c echo.Context) error {
	authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
	player, validPlayer, err := getPlayerFromHeader(authHeader)

	if err != nil {
		return err
	}

	if !validPlayer {
		return c.JSONPretty(http.StatusUnauthorized, &Response{false, nil}, " ")
	}

	body := make(request)
	err = getBody(c, &body)
	if err != nil {
		return err
	}
	card := model.Card{Value: body["value"].(string), Color: body["color"].(string)}

	game, err := playCard(c.Param("game"), player, card)

	if err != nil {
		return err
	}

	return c.JSONPretty(http.StatusOK, &Response{true, newPayload(game)}, "  ")
}

func draw(c echo.Context) error {
	authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
	player, validPlayer, err := getPlayerFromHeader(authHeader)

	if err != nil {
		return err
	}

	if !validPlayer {
		return c.JSONPretty(http.StatusUnauthorized, &Response{false, nil}, " ")
	}

	game, err := drawCard(c.Param("game"), player)

	if err != nil {
		return err
	}

	return c.JSONPretty(http.StatusOK, &Response{true, newPayload(game)}, "  ")

}

func newPayload(game *model.Game) map[string]interface{} {
	payload := make(map[string]interface{})

	// Update known variables
	payload["direction"] = game.Direction
	payload["current_player"] = game.CurrentPlayer
	payload["all_players"] = game.Players
	payload["draw_pile"] = game.DrawPile
	payload["discard_pile"] = game.DiscardPile
	payload["game_id"] = game.ID
	payload["game_over"] = game.Status == "Finished"
	payload["host"] = game.Host
	payload["name"] = game.Name

	return payload
}

func getBody(c echo.Context, output *request) error {
	if err := c.Bind(output); err != nil {
		return err
	}
	return nil
}

func getPostParam(c echo.Context, key string) interface{} {
	req := make(request)
	err := getBody(c, &req)
	if err != nil {
		fmt.Println(err)
	}

	res, ok := req[key]
	if ok {
		return res
	}
	return nil
}
