import { dc, tc } from './checkouts';
import api from './api';

export const getCheckout = (mode, score) => {
    let returnString = '-';
    if (mode === 'double') {
        if (dc[score]) {
            returnString = dc[score];
        }
    } else if (mode === 'triple') {
        if (tc[score]) {
            returnString = tc[score];
        }
    }

    return returnString;
};

export const transformGameMessage = (data, active) => {
    // Check if double or master out and checkout possible
    if (
        (data.Out === 'double' && active.Score.Score <= 170) ||
        (data.Out === 'master' && active.Score.Score <= 180)
    ) {
        if (
            data.GameState !== 'NEXTPLAYER' &&
            !data.GameState.includes('BUST') &&
            data.GameState !== 'NEXTPLAYERWON' &&
            data.GameState !== 'WON'
        ) {
            data.Message = getCheckout(data.Out, active.Score.Score);
        }
    }

    // If above case does not fit check if message is "-" or empty
    // and then substitute with nicer message
    if (data.Message === '-' || data.Message === '') {
        // Construct players name X's oder X'
        let playerMessageName;
        active.Name.slice(-1) === 's'
            ? (playerMessageName = active.Name + "'")
            : (playerMessageName = active.Name + "'s");
        return `Round ${data.ThrowRound} - ${playerMessageName} turn`;
    }

    // Otherwise just return the message
    return data.Message;
};

export const scoreOrPodium = (player, gameData) => {
    // Choose what to display
    // If player is not on Podium display score
    // Else display podium Number
    var onPodium = gameData.Podium.find((p) => p.Player.UID === player.UID);
    if (typeof onPodium === 'undefined') {
        return player.Score.Score;
    } else {
        return 'Place ' + onPodium.Number;
    }
};

export const scoreOrCurrentNumber = (player, gameData) => {
    // Choose what to display
    // If player is not on Podium display score
    // Else display podium Number
    var onPodium = gameData.Podium.find((p) => p.Player.UID === player.UID);
    if (typeof onPodium === 'undefined') {
        return player.Score.CurrentNumber;
    } else {
        return 'Place ' + onPodium.Number;
    }
};

export const scoreOrHitorder = (player, gameData, hitorder) => {
    // Choose what to display
    // If player is not on Podium display score
    // Else display podium Number
    var onPodium = gameData.Podium.find((p) => p.Player.UID === player.UID);
    if (typeof onPodium === 'undefined') {
        return hitorder;
    } else {
        return 'Place ' + onPodium.Number;
    }
};

export const insertThrow = (gameid, number, double, triple) => {
    let modifier = 1;
    if (double) {
        modifier = 2;
    }
    if (triple) {
        modifier = 3;
    }
    api.post(`game/${gameid}/throw/${number}/${modifier}`);

    navigator.vibrate(200);
};

export const insertThrowNumberMod = (gameid, cn, mod) => {
    api.post(`game/${gameid}/throw/${cn}/${mod}`);
    navigator.vibrate(200);
};

export const rematch = (gameid) => {
    api.post(`game/${gameid}/rematch`);
};

export const nextPlayer = (gameid) => {
    api.post(`game/${gameid}/nextPlayer`);
};

export const undo = (gameid) => {
    api.post(`game/${gameid}/undo`);
};

export const endGame = (gameid) => {
    if (confirm('Really end game?')) {
        api.delete(`game/${gameid}`);
    }
};

export const miss = (gameid) => {
    api.post(`game/${gameid}/throw/0/1`);
    navigator.vibrate(200);
};
