let socket;
let wsURL = import.meta.env.VITE_WS_URL;

const eventRedirect = new Event('redirect');
const eventUpdate = new Event('update');

const init = (gameid, page) => {
  socket = new WebSocket(`${wsURL}/${gameid}`);

  socket.onopen = () => {
    console.log(`${page} connected via websockets`);
  };

  socket.onclose = () => {
    console.log(`${page} disconnected from websocket`);
  };

  socket.onmessage = (event) => {
    switch (event.data) {
      case 'redirect':
        socket.close(1000, 'redirect');
        socket.dispatchEvent(eventRedirect);
        break;
      case 'update':
        socket.dispatchEvent(eventUpdate);
        break;
    }
  };

  socket.onerror = (error) => {
    console.error(error);
  };

  return socket;
};

const ws = {
  init,
};

export default ws;
