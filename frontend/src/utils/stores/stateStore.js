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
            update(async (state) => {
                const res = await api.get(`game/${gameid}/display`);
                state.gameData = await res.json();
                state.players = state.gameData.Player;
                state.activePlayer =
                    state.gameData.Player[state.gameData.ActivePlayer];
                state.message = transformGameMessage(
                    state.gameData,
                    state.activePlayer
                );
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
