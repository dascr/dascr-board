import { dc, tc } from './checkouts';

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
