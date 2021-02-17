package api

import (
	"database/sql"
	"embed"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dascr/dascr-board/api/routes"
	"github.com/dascr/dascr-board/logger"
	"github.com/dascr/dascr-board/ws"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

// Embed static images
//go:embed images
var static embed.FS

func getDir(dir string) string {
	cwd, err := os.Getwd()
	if err != nil {
		logger.Fatal("Unable to get current working directory")
	}
	folder := filepath.Join(cwd, dir)

	return folder
}

func allowOriginFunc(_ *http.Request, _ string) bool {
	return true
}

// SetRoutes will instantiate the REST routes
func SetRoutes(db *sql.DB, h *ws.Hub) *chi.Mux {
	// Router setup
	r := chi.NewRouter()

	// Middleware
	r.Use(
		render.SetContentType(render.ContentTypeJSON),
		cors.Handler(cors.Options{
			AllowOriginFunc:  allowOriginFunc,
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}),
		middleware.Recoverer,
		middleware.RedirectSlashes,
	)
	if os.Getenv("DEBUG") == "TRUE" {
		r.Use(middleware.Logger)
	}

	/* API Endpoints */
	r.Route("/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("ok")); err != nil {
				logger.Errorf("Error writing response back to browser: %+v", err)
			}
		})

		// Game routes
		r.Mount("/game", routes.GameRoutes(db, h))

		// Player routes
		r.Mount("/player", routes.PlayerRoutes(db))

		// Debug routes
		if os.Getenv("DEBUG") == "TRUE" {
			r.Mount("/debug", routes.DebugRoutes(h))
		}
	})

	// Static content
	r.Route("/images", func(r chi.Router) {
		r.Mount("/", http.FileServer(http.FS(static)))
	})
	r.Mount("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(getDir("uploads")))))

	// Socket route
	r.Mount("/ws", routes.SocketRoutes(h))

	// Debug output for registered route
	logger.Debug("All routes are")
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		logger.Debugf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		logger.Errorf("Logging err: %s\n", err.Error())
	}

	return r
}
