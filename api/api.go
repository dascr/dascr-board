package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/dascr/dascr-board/config"
	"github.com/dascr/dascr-board/logger"
	"github.com/dascr/dascr-board/ws"
)

// API holds the relevant information of the backend API
type API struct {
	APIConfig  *config.APIConfig
	DB         *sql.DB
	HTTPServer *http.Server
	Hub        *ws.Hub
}

// SetupAPI will instantiate an API and return it
func SetupAPI(db *sql.DB) *API {
	// Setup Websocket hub
	hub := ws.Hub{
		Broadcast:  make(chan ws.Message),
		Register:   make(chan ws.Subscription),
		Unregister: make(chan ws.Subscription),
		Rooms:      make(map[string]map[*ws.Connection]bool),
	}
	go hub.Run()

	return &API{
		APIConfig: &config.APIConfig{
			IP:   config.MustGet("API_IP"),
			Port: config.MustGet("API_PORT"),
		},
		DB:  db,
		Hub: &hub,
	}
}

// Start will start the backend API
func (a *API) Start() {
	// Build Address
	addr := fmt.Sprintf("%s:%s", a.APIConfig.IP, a.APIConfig.Port)
	// Create router
	r := SetRoutes(a.DB, a.Hub)
	// Create server
	a.HTTPServer = &http.Server{
		Addr:    addr,
		Handler: r,
	}

	logger.Infof("Starting API at: %+v", addr)

	// Start http listener
	logger.Panic(a.HTTPServer.ListenAndServe())
}
