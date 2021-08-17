package game

import (
	"database/sql"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/dascr/dascr-board/logger"
	"github.com/dascr/dascr-board/player"
	"github.com/dascr/dascr-board/settings"
	"github.com/dascr/dascr-board/ws"
	"github.com/go-chi/chi"
)

// Games hold the running games
var Games = make(map[string]Data)

// Data stuct will hold the generic game data
// handed over by the frontend. It might be used in
// the different game types
type Data struct {
	UID           string `json:"uid"`
	Player        []int  `json:"player"`
	Game          string `json:"game"`
	Variant       string `json:"variant"`
	In            string `json:"in"`
	Out           string `json:"out"`
	Punisher      bool   `json:"punisher"`
	Sound         bool   `json:"sound"`
	Podium        bool   `json:"podium"`
	AutoSwitch    bool   `json:"autoswitch"`
	CricketRandom bool   `json:"cricketrandom"`
	CricketGhost  bool   `json:"cricketghost"`
	GameObject    Game
}

// GetAllGames lists all in memory games
func GetAllGames() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var gamesSlice []Data
		for _, element := range Games {
			gamesSlice = append(gamesSlice, element)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(gamesSlice); err != nil {
			logger.Errorf("Error when encoding json: %+v", err)
		}
	}
}

// GetSpecificGame lists a specific in memory game
func GetSpecificGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := chi.URLParam(r, "id")
		if _, ok := Games[uid]; ok {
			game := Games[uid]
			status := game.GameObject.GetStatus()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(status); err != nil {
				logger.Errorf("Error when encoding json: %+v", err)
			}
		} else {
			http.Error(w, "There is no game with this id", http.StatusNotFound)
		}
	}
}

// GetSpecificGameDisplay lists a specific in memory game
// and only provides the data necessary to display the
// Scoreboard information (small data size)
func GetSpecificGameDisplay() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := chi.URLParam(r, "id")
		if _, ok := Games[uid]; ok {
			game := Games[uid]
			status := game.GameObject.GetStatusDisplay()

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(status); err != nil {
				logger.Errorf("Error when encoding json: %+v", err)
			}
		} else {
			http.Error(w, "There is no game with this id", http.StatusNotFound)
		}

	}
}

// CreateGame will dispatch the creation of a game
// to the corresponding game creation handler after
// adding it to the in memory games map
func CreateGame(db *sql.DB, h *ws.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data Data
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte("Error parsing json data on backend")); err != nil {
				logger.Errorf("Error when writing response back to browser: %+v", err)
			}
			return
		}

		// Init random number generator
		s := rand.NewSource(time.Now().Unix())
		rg := rand.New(s)
		activePlayer := rg.Intn(len(data.Player))

		// Get players from database
		var players = player.GetPlayerList(db, data.Player)

		// Dispatch to handler
		switch game := data.Game; game {
		case "x01":
			data.GameObject = &X01Game{
				Base: BaseGame{
					UID:          data.UID,
					Game:         data.Game,
					Player:       players,
					Variant:      data.Variant,
					In:           data.In,
					Out:          data.Out,
					Punisher:     data.Punisher,
					ActivePlayer: activePlayer,
					ThrowRound:   1,
					GameState:    "THROW",
					Message:      "",
					Settings: &settings.Settings{
						Podium:     data.Podium,
						Sound:      data.Sound,
						AutoSwitch: data.AutoSwitch,
					},
				},
			}

		case "cricket":
			data.GameObject = &CricketGame{
				Base: BaseGame{
					UID:          data.UID,
					Game:         data.Game,
					Player:       players,
					Variant:      data.Variant,
					ActivePlayer: activePlayer,
					ThrowRound:   1,
					GameState:    "THROW",
					Message:      "",
					Settings: &settings.Settings{
						Podium:     data.Podium,
						Sound:      data.Sound,
						AutoSwitch: data.AutoSwitch,
					},
					CricketController: &CricketGameController{
						Random: data.CricketRandom,
						Ghost:  data.CricketGhost,
					},
				},
			}

		case "atc":
			data.GameObject = &ATCGame{
				Base: BaseGame{
					UID:          data.UID,
					Game:         data.Game,
					Player:       players,
					Variant:      data.Variant,
					ActivePlayer: activePlayer,
					ThrowRound:   1,
					GameState:    "THROW",
					Message:      "",
					Settings: &settings.Settings{
						Podium:     data.Podium,
						Sound:      data.Sound,
						AutoSwitch: data.AutoSwitch,
					},
				},
			}
		case "split":
			data.GameObject = &SplitGame{
				Base: BaseGame{
					UID:          data.UID,
					Game:         data.Game,
					Player:       players,
					Variant:      data.Variant,
					ActivePlayer: activePlayer,
					ThrowRound:   1,
					GameState:    "THROW",
					Message:      "",
					Settings: &settings.Settings{
						Podium:     data.Podium,
						Sound:      data.Sound,
						AutoSwitch: data.AutoSwitch,
					},
				},
			}
		}

		if err := data.GameObject.StartGame(); err != nil {
			logger.Errorf("Game cannot be started: %+v", err)
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte("Error occurred starting a Game")); err != nil {
				logger.Errorf("Error writing response back to browser: %+v", err)
			}
		}

		Games[data.UID] = data

		// trigger redirect of scanpage and game setup page
		message := ws.Message{
			Room: chi.URLParam(r, "id"),
			Data: []byte("redirect"),
		}
		h.Broadcast <- message

		// json encode game as response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(data); err != nil {
			logger.Errorf("Error when encoding json: %+v", err)
		}

	}
}

// DeleteGame will delete game from in memory map
func DeleteGame(h *ws.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		delete(Games, chi.URLParam(r, "id"))

		// redirect scoreboard and controller page
		message := ws.Message{
			Room: chi.URLParam(r, "id"),
			Data: []byte("redirect"),
		}
		h.Broadcast <- message

		w.WriteHeader(http.StatusNoContent)
	}
}

// NextPlayer will switch to next player in game
func NextPlayer(h *ws.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := chi.URLParam(r, "id")
		if _, ok := Games[uid]; ok {
			game := Games[uid]
			game.GameObject.NextPlayer(h)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(game.GameObject); err != nil {
				logger.Errorf("Error when encoding json: %+v", err)
			}
		} else {
			http.Error(w, "There is no game with this id", http.StatusNotFound)
		}
	}
}

// InsertThrow will handle to add a throw to a game
func InsertThrow(h *ws.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var iNum, iMod int
		var err error
		uid := chi.URLParam(r, "id")
		if _, ok := Games[uid]; ok {
			number := chi.URLParam(r, "number")
			modifier := chi.URLParam(r, "modifier")
			if iNum, err = strconv.Atoi(number); err != nil {
				logger.Errorf("When adding throw. Invalid number: %+v", err)
				w.WriteHeader(http.StatusBadRequest)
				if _, err := w.Write([]byte("Throw was not added")); err != nil {
					logger.Errorf("Error writing response back to browser: %+v", err)
				}
				return
			}

			if iMod, err = strconv.Atoi(modifier); err != nil {
				logger.Errorf("When adding throw. Invalid modifier: %+v", err)
				w.WriteHeader(http.StatusBadRequest)
				if _, err := w.Write([]byte("Throw was not added")); err != nil {
					logger.Errorf("Error writing response back to browser: %+v", err)
				}
				return
			}
			game := Games[uid]
			if err := game.GameObject.RequestThrow(iNum, iMod, h); err != nil {
				logger.Errorf("When adding throw: %+v", err)
				w.WriteHeader(http.StatusBadRequest)
				if _, err := w.Write([]byte("Throw was not added")); err != nil {
					logger.Errorf("Error writing response back to browser: %+v", err)
				}
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(game.GameObject); err != nil {
				logger.Errorf("Error when encoding json: %+v", err)
			}
		} else {
			http.Error(w, "There is no game with this id", http.StatusNotFound)
		}
	}
}

// Undo will handle the undo actions per game
func Undo(h *ws.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if data, ok := Games[chi.URLParam(r, "id")]; ok {
			if err := data.GameObject.Undo(h); err != nil {
				logger.Errorf("There was an error undoing last turn: %+v", err)
				http.Error(w, "There was an error undoing last turn", 500)
				return
			}
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("Turn was undone")); err != nil {
				logger.Errorf("Error writing response back to browser: %+v", err)
			}

		} else {
			http.Error(w, "There is no game with this id", http.StatusNotFound)
		}
	}
}

// Rematch will handle the rematch action per game
func Rematch(h *ws.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if data, ok := Games[chi.URLParam(r, "id")]; ok {
			if err := data.GameObject.Rematch(h); err != nil {
				logger.Errorf("There was an error triggering rematch: %+v", err)
				http.Error(w, "There was an error triggering", 500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(data.GameObject.GetStatusDisplay()); err != nil {
				logger.Errorf("Error when encoding json: %+v", err)
			}

		} else {
			http.Error(w, "There is no game with this id", http.StatusNotFound)
		}
	}
}
