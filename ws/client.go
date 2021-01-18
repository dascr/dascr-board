package ws

import (
	"net/http"
	"time"

	"github.com/dascr/dascr-board/logger"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 5) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Connection represents a websocket connection
type Connection struct {
	ws   *websocket.Conn
	send chan []byte
}

func (s Subscription) readPump() {
	c := s.Conn
	defer func() {
		s.Hub.Unregister <- s
		if err := c.ws.Close(); err != nil {
			logger.Errorf("Error closing socket in readPump(): %+v", err)
			return
		}
	}()

	c.ws.SetReadLimit(maxMessageSize)
	if err := c.ws.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		logger.Errorf("Error in websocket readPump(): %+v", err)
	}

	c.ws.SetPongHandler(func(string) error {
		if err := c.ws.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			logger.Errorf("Error in websocket readPump(): %+v", err)
		}
		return nil
	})

	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				logger.Errorf("Websocket closed unexpected: %+v", err)
			}
			break
		}
		m := Message{msg, s.Room}
		logger.Debugf("Message by websocket is: %+v", m)
	}
}

func (c *Connection) write(mt int, payload []byte) error {
	if err := c.ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		logger.Errorf("Error in websocket write(): %+v", err)
	}
	return c.ws.WriteMessage(mt, payload)
}

func (s *Subscription) writePump() {
	c := s.Conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := c.ws.Close(); err != nil {
			// Ignore error by now as it occures when ending game but does not break anything
			// logger.Errorf("Error closing socket in writePump(): %+v", err)
		}
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				if err := c.write(websocket.CloseMessage, []byte{}); err != nil {
					// Ignore error by now as it occures when ending game but does not break anything
					// logger.Errorf("Error in websocket writePump() - closeMessage: %+v", err)
					return
				}
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				// logger.Errorf("Error in websocket writePump() - TextMessage: %+v", err)
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte("ping")); err != nil {
				// logger.Errorf("Error in websocket writePump() - Ping message: %+v", err)
				return
			}
		}
	}
}

// ServeWS will handle the websocket connection
func ServeWS(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Errorf("Unable to upgrade to websocket: %+v", err.Error())
			return
		}

		c := &Connection{
			send: make(chan []byte, 256),
			ws:   ws,
		}

		s := Subscription{
			Hub:  hub,
			Conn: c,
			Room: chi.URLParam(r, "id"),
		}
		s.Hub.Register <- s
		go s.writePump()
		go s.readPump()
	}
}
