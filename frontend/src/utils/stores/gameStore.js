import { writable } from 'svelte/store';

const store = () => {
  const state = {
    player: '',
    game: 'x01',
    variant: '501',
    in: 'straight',
    out: 'double',
    elimination: false,
    settings: {
      sound: true,
      podium: false,
      autoswitch: false,
      punisher: false,
    },
    cricket: {
      random: false,
      ghost: false,
    },
  };

  const { subscribe, set, update } = writable(state);

  const methods = {};

  return {
    subscribe,
    set,
    update,
    ...methods,
  };
};

export default store();
