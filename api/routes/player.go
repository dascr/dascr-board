package routes

import (
	"database/sql"

	"github.com/dascr/dascr-board/player"
	"github.com/go-chi/chi"
)

// PlayerRoutes represent the Player API endpoint
func PlayerRoutes(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", player.GetAllPlayer(db))
	r.Get("/{id}", player.GetSpecificPlayer(db))
	r.Post("/", player.AddPlayer(db))
	r.Post("/{id}/image", player.HandlePlayerImage(db))
	r.Patch("/{id}", player.UpdatePlayer(db))
	r.Delete("/{id}", player.DeletePlayer(db))

	return r
}
