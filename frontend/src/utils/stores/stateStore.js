import { writable } from 'svelte/store';
import api from '../api';
import { transformGameMessage } from '../methods';

const store = () => {
    const state = {
        double: false,
        triple: false,
        gameData: {},
        players: [],
        activePlayer: {},
        message: '',
        numbers: [],
        revealed: [],
        allRevealed: false,
    };

    const { subscribe, set, update } = writable(state);

    const methods = {
        toggleDouble() {
            update((state) => {
                state.triple = false;
                state.double = !state.double;
                return state;
            });
        },
        toggleTriple() {
            update((state) => {
                state.double = false;
                state.triple = !state.triple;
                return state;
            });
        },
        async updateState(gameid) {
            const res = await api.get(`game/${gameid}/display`);
            const data = await res.json();
            const players = data.Player;
            const activePlayer = data.Player[data.ActivePlayer];
            if (data.CricketController !== null) {
                const numbers = data.CricketController.Numbers;
                const revealed = data.CricketController.NumberRevealed;
                const allRevealed = revealed.every((n) => n === true);
                update((state) => {
                    state.numbers = numbers;
                    state.revealed = revealed;
                    state.allRevealed = allRevealed;
                    return state;
                });
            }
            const transformMessage = transformGameMessage(data, activePlayer);
            update((state) => {
                state.gameData = data;
                state.players = players;
                state.activePlayer = activePlayer;
                state.message = transformMessage;
                return state;
            });
        },
    };

    return {
        subscribe,
        set,
        update,
        ...methods,
    };
};

export default store();
