import BaseService from "./baseService";

export default {
    setToken(token) {
        BaseService.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    },

    async login(username) {
        return BaseService.post(`/login/`, {username});
    },

    async newGame() {
        return BaseService.post(`/games/`);
    },

    async newGameWithName(name) {
        return BaseService.post(`/games/`, {name})
    },

    update(game) {
        return BaseService.get(`/games/${game}`);
    },

    startGame(game) {
        return BaseService.post(`/games/${game}/start`);
    },

    playCard(game, cardNumber, cardColor) {
        return BaseService.post(`/games/${game}/play`, {number: cardNumber, color: cardColor});
    },

    drawCard(game) {
        return BaseService.post(`/games/${game}/draw`);
    },

    joinGame(game, password) {
        return BaseService.post(`/games/${game}/join`, {password});
    }
}