package routes

import (
	"github.com/dascr/dascr-board/ws"
	"github.com/go-chi/chi"
)

// SocketRoutes provide websocket routes
func SocketRoutes(h *ws.Hub) chi.Router {
	r := chi.NewRouter()
	r.Get("/{id}", ws.ServeWS(h))
	return r
}
