package routes

import (
	"net/http"

	"github.com/dascr/dascr-board/ws"
	"github.com/go-chi/chi"
)

func websocketSendUpdate(h *ws.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := chi.URLParam(r, "id")
		message := ws.Message{
			Room: uid,
			Data: []byte("update"),
		}
		h.Broadcast <- message
		w.WriteHeader(http.StatusOK)
	}
}

func websocketSendRedirect(h *ws.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := chi.URLParam(r, "id")
		message := ws.Message{
			Room: uid,
			Data: []byte("redirect"),
		}
		h.Broadcast <- message
		w.WriteHeader(http.StatusOK)
	}
}

// DebugRoutes represent the Game API endpoint.
func DebugRoutes(h *ws.Hub) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/{id}/update", websocketSendUpdate(h))
	r.Get("/{id}/redirect", websocketSendRedirect(h))

	return r
}
