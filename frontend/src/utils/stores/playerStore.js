import { writable } from 'svelte/store';

const store = () => {
  const state = {
    name: '',
    nickname: '',
    image: null,
    cropfactor: {
      x: 0,
      y: 0,
      width: 0,
      height: 0,
    },
  };

  const { subscribe, set, update } = writable(state);

  const methods = {
    setImage(file) {
      update((state) => {
        state.image = URL.createObjectURL(file);
        return state;
      });
    },
    setCropFactor(factor) {
      update((state) => {
        state.cropfactor = factor;
        return state;
      });
    },
    resetAll() {
      update((state) => {
        state = {
          name: '',
          nickname: '',
          image: null,
          cropfactor: {
            x: 0,
            y: 0,
            width: 0,
            height: 0,
          },
        };
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
