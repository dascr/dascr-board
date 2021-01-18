package routes

import (
	"database/sql"

	"github.com/dascr/dascr-board/game"
	"github.com/dascr/dascr-board/ws"
	"github.com/go-chi/chi"
)

// GameRoutes represent the Game API endpoint.
func GameRoutes(db *sql.DB, h *ws.Hub) *chi.Mux {
	r := chi.NewRouter()

	// Unprotected Group
	r.Group(func(r chi.Router) {
		r.Get("/", game.GetAllGames())
		r.Get("/{id}", game.GetSpecificGame())
		r.Get("/{id}/display", game.GetSpecificGameDisplay())
		r.Post("/{id}", game.CreateGame(db, h))
		r.Post("/{id}/nextPlayer", game.NextPlayer(h))
		r.Post("/{id}/throw/{number}/{modifier}", game.InsertThrow(h))
		r.Post("/{id}/undo", game.Undo(h))
		r.Post("/{id}/rematch", game.Rematch(h))
		r.Delete("/{id}", game.DeleteGame(h))
	})

	return r
}
