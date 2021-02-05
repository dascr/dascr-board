package player

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dascr/dascr-board/logger"
	"github.com/dascr/dascr-board/score"
	"github.com/dascr/dascr-board/throw"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
)

// Player holds all the players information
type Player struct {
	UID             string
	Name            string
	Nickname        string
	Image           string
	ThrowRounds     []throw.Round
	TotalThrowCount int
	Score           score.BaseScore
	LastThrows      []throw.Throw
	ThrowSum        int
	Average         float64
}

// Image holds the b64 string from a uploaded avatar
type Image struct {
	B64 string `json:"b64"`
}

// GetAllPlayer will return all player
func GetAllPlayer(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var players []Player

		rows, err := db.Query("SELECT * from player")
		if err != nil {
			logger.Errorf("Reading players from database: %+v", err)
			return
		}

		for rows.Next() {
			var player Player
			err := rows.Scan(&player.UID, &player.Name, &player.Nickname, &player.Image)
			if err != nil {
				logger.Errorf("Reading player row in table player")
				continue
			}
			players = append(players, player)
		}

		if err := rows.Close(); err != nil {
			logger.Errorf("Error closing the database row: %+v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(players); err != nil {
			logger.Errorf("Error writing response back to browser: %+v", err)
		}
	}
}

// GetSpecificPlayer will return a specific player by id
func GetSpecificPlayer(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		playerID := chi.URLParam(r, "id")
		var player Player

		stmt, err := db.Prepare("SELECT * FROM player WHERE id = ? LIMIT 1")
		if err != nil {
			logger.Errorf("Preparing the player query statement: %+v", err)
			return
		}

		row, err := stmt.Query(playerID)
		if err != nil {
			logger.Errorf("Fetching the player in database: %+v", err)
			return
		}

		for !row.Next() {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		err = row.Scan(&player.UID, &player.Name, &player.Nickname, &player.Image)
		if err != nil {
			logger.Errorf("Reading player row in table player: %+v", err)
			return
		}

		if err := row.Close(); err != nil {
			logger.Errorf("Error closing the database row: %+v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(player); err != nil {
			logger.Errorf("Error writing response back to browser: %+v", err)
		}
	}
}

// GetPlayerList will return a list of player using the id
func GetPlayerList(db *sql.DB, ids []int) []Player {
	var players []Player
	var questionMarks []string
	var params []interface{}

	// IDs need to be 1,2,3,4
	// from slice of ints
	for _, id := range ids {
		questionMarks = append(questionMarks, "?")
		params = append(params, id)
	}
	query := fmt.Sprintf("SELECT * FROM player WHERE id IN (%s)", strings.Join(questionMarks, ","))

	stmt, err := db.Prepare(query)
	if err != nil {
		logger.Errorf("Preparing the player list query statement: %+v", err)
		return nil
	}

	rows, err := stmt.Query(params...)
	if err != nil {
		logger.Errorf("Reading players from database: %+v", err)
		return nil
	}

	for rows.Next() {
		var player Player
		err := rows.Scan(&player.UID, &player.Name, &player.Nickname, &player.Image)
		if err != nil {
			logger.Errorf("Reading player row in table player")
			continue
		}
		players = append(players, player)
	}

	// Sort to initial order
	var returnPlayer []Player

	for len(ids) > 0 {
		for i, player := range players {
			if players[i].UID == fmt.Sprintf("%d", ids[0]) {
				returnPlayer = append(returnPlayer, player)
				ids = ids[1:]
				if len(ids) == 0 {
					break
				}
			}
		}
	}

	if err := rows.Close(); err != nil {
		logger.Errorf("Error closing the database row: %+v", err)
		return nil
	}

	return returnPlayer
}

// GetPlayer will return a player as struct, not as json
func GetPlayer(db *sql.DB, id string) Player {
	var player Player

	stmt, err := db.Prepare("SELECT * FROM player WHERE id = ? LIMIT 1")
	if err != nil {
		logger.Errorf("Preparing the player query statement: %+v", err)
		return player
	}

	row, err := stmt.Query(id)
	if err != nil {
		logger.Errorf("Fetching the player in database: %+v", err)
		return player
	}

	for !row.Next() {
		logger.Error("Player was not found")
		return player
	}

	err = row.Scan(&player.UID, &player.Name, &player.Nickname, &player.Image)
	if err != nil {
		logger.Errorf("Reading player row in table player: %+v", err)
		return player
	}

	if err := row.Close(); err != nil {
		logger.Errorf("Error closing the database row: %+v", err)
		return player
	}

	return player
}

// AddPlayer will add a player
func AddPlayer(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var player *Player
		err := json.NewDecoder(r.Body).Decode(&player)
		if err != nil {
			logger.Errorf("Unable to get request body when creating player: %+v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create statement
		create := `
		INSERT INTO player (name, nickname, image)
		VALUES(?,?,?)
		`
		stmt, err := db.Prepare(create)
		if err != nil {
			logger.Errorf("Preparing the player creation statement: %+v", err)
			return
		}

		stmtRes, err := stmt.Exec(player.Name, player.Nickname, player.Image)
		if err != nil {
			logger.Errorf("Creating the player in database: %+v", err)
			return
		}

		newID, err := stmtRes.LastInsertId()
		if err != nil {
			logger.Errorf("Creating the player in database: %+v", err)
			return
		}
		player.UID = strconv.FormatInt(newID, 10)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(player); err != nil {
			logger.Errorf("Error writing response back to browser: %+v", err)
		}
	}
}

// UpdatePlayer will update player details
func UpdatePlayer(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		playerID := chi.URLParam(r, "id")
		var player Player
		var updatedPlayer Player

		stmt, err := db.Prepare("SELECT * FROM player WHERE id = ? LIMIT 1")
		if err != nil {
			logger.Errorf("Preparing the player query statement: %+v", err)
			return
		}

		row, err := stmt.Query(playerID)
		if err != nil {
			logger.Errorf("Fetching the player in database: %+v", err)
			return
		}

		for !row.Next() {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		err = row.Scan(&player.UID, &player.Name, &player.Nickname, &player.Image)
		if err != nil {
			logger.Errorf("Reading player row in table player: %+v", err)
			return
		}

		if err := row.Close(); err != nil {
			logger.Errorf("Error closing the database row: %+v", err)
			return
		}

		err = json.NewDecoder(r.Body).Decode(&updatedPlayer)
		if err != nil {
			logger.Errorf("Unable to get request body when updating player: %+v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		player.Name = updatedPlayer.Name
		player.Nickname = updatedPlayer.Nickname
		player.Image = updatedPlayer.Image

		update := `
		UPDATE player
		SET name = ?,
		nickname = ?,
		image = ?
		WHERE id = ?
		`

		stmt, err = db.Prepare(update)
		if err != nil {
			logger.Errorf("Preparing the player query statement: %+v", err)
			return
		}

		_, err = stmt.Exec(player.Name, player.Nickname, player.Image, player.UID)
		if err != nil {
			logger.Errorf("Updating the player in database: %+v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(player); err != nil {
			logger.Errorf("Error writing response back to browser: %+v", err)
		}
	}
}

// DeletePlayer will delete player
func DeletePlayer(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		playerID := chi.URLParam(r, "id")
		var player Player

		stmt, err := db.Prepare("SELECT * FROM player WHERE id = ? LIMIT 1")
		if err != nil {
			logger.Errorf("Preparing the player query statement: %+v", err)
			return
		}

		row, err := stmt.Query(playerID)
		if err != nil {
			logger.Errorf("Fetching the player in database: %+v", err)
			return
		}

		for !row.Next() {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		err = row.Scan(&player.UID, &player.Name, &player.Nickname, &player.Image)
		if err != nil {
			logger.Errorf("Reading player row in table player: %+v", err)
			return
		}

		if err := row.Close(); err != nil {
			logger.Errorf("Error closing the database row: %+v", err)
			return
		}

		stmt, err = db.Prepare("DELETE FROM player WHERE id = ?")
		if err != nil {
			logger.Errorf("Preparing the player deletion statement: %+v", err)
			return
		}

		res, err := stmt.Exec(playerID)
		if err != nil {
			logger.Errorf("Deleting the player from the database: %+v", err)
			return
		}

		affected, err := res.RowsAffected()
		if err != nil {
			logger.Errorf("Reading the affected rows")
		}

		if affected == 0 {
			w.WriteHeader(http.StatusNotFound)
			if _, err := w.Write([]byte("Player was not found and therefore not deleted")); err != nil {
				logger.Errorf("Error writing response back to browser: %+v", err)
			}
			return
		}

		// Delete image from disk
		if err := handlePlayerImageDeletion(player.Image); err != nil {
			logger.Warnf("Unable to delete player image from disk: %+v. Image will still be there.", err)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// HandlePlayerImage will handle the upload of an player Image
func HandlePlayerImage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		playerID := chi.URLParam(r, "id")
		var err error
		var imagePath string

		updateImage := `
		UPDATE player
		SET image = ?
		WHERE id = ?
		`

		var b64image *Image

		if err := json.NewDecoder(r.Body).Decode(&b64image); err != nil {
			logger.Errorf("Unable to get request body when creating player image: %+v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if b64image.B64 != "" {
			// have image
			idx := strings.Index(b64image.B64, ";base64,")
			if idx < 0 {
				logger.Errorf("Unable to decode player image from base64. Invalid Image: %+v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			decodedImg, err := base64.StdEncoding.DecodeString(b64image.B64[idx+8:])
			if err != nil {
				logger.Errorf("Unable to decode player image from base64: %+v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if filetype.IsImage(decodedImg) {
				kind, _ := filetype.Match(decodedImg)
				// UUID for filename
				id := uuid.New()
				savepath := fmt.Sprintf("./uploads/user_%s_%s.%s", playerID, id, kind.Extension)
				_, err := os.Create(savepath)
				if err != nil {
					logger.Errorf("Creating the file on disk: %+v", err)
					w.WriteHeader(http.StatusBadRequest)
					if _, err := w.Write([]byte("notdone")); err != nil {
						logger.Errorf("Error writing response back to browser: %+v", err)
					}
					return
				}
				if err := ioutil.WriteFile(savepath, decodedImg, 0600); err != nil {
					logger.Errorf("Writing image content to file: %+v", err)
					w.WriteHeader(http.StatusBadRequest)
					if _, err := w.Write([]byte("notdone")); err != nil {
						logger.Errorf("Error writing response back to browser: %+v", err)
					}
					return
				}
				imageName := fmt.Sprintf("user_%s_%s.%s", playerID, id, kind.Extension)
				imagePath = fmt.Sprintf("uploads/%s", imageName)
			} else {
				logger.Error("Provided file is not an image")
				w.WriteHeader(http.StatusBadRequest)
				if _, err := w.Write([]byte("notdone")); err != nil {
					logger.Errorf("Error writing response back to browser: %+v", err)
				}
				return
			}
		} else {
			// fake image
			// Choose from 0 - 8 for getting random image name
			rand.Seed(time.Now().UnixNano())
			randomNumber := rand.Intn(7-0+1) + 1
			imageName := fmt.Sprintf("static_%v.png", randomNumber)
			imagePath = fmt.Sprintf("images/%s", imageName)
		}

		// Update path to image in db (player)
		stmt, err := db.Prepare(updateImage)
		if err != nil {
			logger.Errorf("Preparing the player image update statement: %+v", err)
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte("notdone")); err != nil {
				logger.Errorf("Error writing response back to browser: %+v", err)
			}
			return
		}

		_, err = stmt.Exec(imagePath, playerID)
		if err != nil {
			logger.Errorf("Updating the players image path in database: %+v", err)
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte("notdone")); err != nil {
				logger.Errorf("Error writing response back to browser: %+v", err)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("done")); err != nil {
			logger.Errorf("Error writing response back to browser: %+v", err)
		}
	}
}

// handlePlayerImageDeletion will detele the players image on disk when DELETE endpoint was hit
func handlePlayerImageDeletion(path string) error {
	// Delete image in path
	if !strings.Contains(path, "static") {
		if err := os.Remove("./" + path); err != nil {
			logger.Errorf("Removing players image from disk: %+v", err)
			return err
		}
	}
	return nil
}
