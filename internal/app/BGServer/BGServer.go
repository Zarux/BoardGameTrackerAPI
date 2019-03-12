package BGServer

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type retError struct {
	Code   int    `json:"code,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func handlerHandler (router http.Handler) http.Handler{
	router = handlers.LoggingHandler(os.Stdout, router)
	return handlers.RecoveryHandler()(router)
}

func Run(port string, configFile string) {
	if err := initConfig(configFile); err != nil {
		log.Fatal(err)
	}

	keepAliveDB := time.NewTicker(10 * time.Minute)
	go func() {
		for range keepAliveDB.C {
			if db != nil {
				err := db.Ping()
				log.Println("Pinging database")
				if err != nil{
					panic("database not responding")
				}
			}
		}
	}()

	router := mux.NewRouter()
	router.HandleFunc("/boardGames", getBoardGames).Methods("GET")
	router.HandleFunc("/room/{room}", getRoom).Methods("GET")

	router.HandleFunc("/room/{room}/player", addPlayer).Methods("POST")
	router.HandleFunc("/room/{room}/player/{id}", editPlayer).Methods("PUT")
	router.HandleFunc("/room/{room}/player/{id}", deletePlayer).Methods("DELETE")

	router.HandleFunc("/room/{room}/game", addGame).Methods("POST")
	router.HandleFunc("/room/{room}/game/{id}", editGame).Methods("PUT")
	router.HandleFunc("/room/{room}/game/{id}", deleteGame).Methods("DELETE")

	log.Println("Running on " + port)
	log.Println(http.ListenAndServe("127.0.0.1:"+port, handlerHandler(router)))
}

/*GET*/
func getBoardGames(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		retErrors []retError
	)
	params := r.URL.Query()
	query := params.Get("search")
	limitStr := params.Get("limit")
	roomIdStr := params.Get("room")
	limit := int64(100)
	roomId := int64(0)
	w.Header().Set("Content-Type", "application/json")
	if limitStr != "" {
		limit, err = strconv.ParseInt(limitStr, 10, 0)
		if err != nil {
			retErrors = append(retErrors, retError{Code: http.StatusBadRequest, Detail: "Invalid limit"})
		}
	}
	if query == "" {
		retErrors = append(retErrors, retError{Code: http.StatusBadRequest, Detail: "Missing query"})
	}
	if roomIdStr != "" {
		roomId, err = strconv.ParseInt(roomIdStr, 10, 0)
		if err != nil {
			retErrors = append(retErrors, retError{Code: http.StatusBadRequest, Detail: "Invalid room"})
		}
	}

	if len(retErrors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(retErrors); err != nil {
			log.Println(err)
		}
	} else {
		boardGames, err := SearchBoardGames(query, int(limit), int(roomId))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			if err = json.NewEncoder(w).Encode(boardGames); err != nil {
				log.Println(err)
			}
		}
	}
}

func getRoom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var (
		err       error
		retErrors []retError
	)
	w.Header().Set("Content-Type", "application/json")
	room := params["room"]
	if room == "" {
		retErrors = append(retErrors, retError{Code: http.StatusBadRequest, Detail: "Missing room"})
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(retErrors); err != nil {
			log.Println(err)
		}
	} else {
		bgRoom, err := GetRoom(room, true)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			err = json.NewEncoder(w).Encode(bgRoom)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

/*POST*/
func addGame(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		newGame   Game
		retErrors []retError
	)
	params := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)
	roomName := params["room"]
	w.Header().Set("Content-Type", "application/json")
	if err = decoder.Decode(&newGame); err != nil {
		log.Println(err)
	}

	if err = validator.Validate(newGame); err != nil {
		retErrors = append(retErrors, retError{Code: http.StatusBadRequest, Detail: err.Error()})
	}
	room, err := GetRoom(roomName, false)
	if err != nil {
		retErrors = append(retErrors, retError{Code: http.StatusBadRequest, Detail: err.Error()})
	}

	boardGame, err := getBoardGame(newGame.Game.Id)
	if err != nil && err == sql.ErrNoRows {
		retErrors = append(retErrors, retError{Code: http.StatusBadRequest, Detail: "Invalid boardGame.id"})
	} else if err != nil {
		retErrors = append(retErrors, retError{Code: http.StatusBadRequest, Detail: err.Error()})
	}

	if len(retErrors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(retErrors); err != nil {
			log.Println(err)
		}
	} else {
		newGame.Game = *boardGame
		if err = room.AddGame(&newGame); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusCreated)
			err = json.NewEncoder(w).Encode(newGame)
			if err != nil {
				log.Println(err)
			}
		}
	}

}

func addPlayer(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		newPlayer Player
		retErrors []retError
	)
	params := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)
	roomName := params["room"]
	w.Header().Set("Content-Type", "application/json")
	if err = decoder.Decode(&newPlayer); err != nil {
		log.Println(err)
	}
	if newPlayer.Name == "" {
		retErrors = append(retErrors, retError{Code: http.StatusBadRequest, Detail: "Missing name"})
	}

	if newPlayer.Color == "" {
		retErrors = append(retErrors, retError{Code: http.StatusBadRequest, Detail: "Missing color"})
	}

	if err = validator.Validate(newPlayer); err != nil {
		retErrors = append(retErrors, retError{Code: http.StatusBadRequest, Detail: err.Error()})
	}
	room, err := GetRoom(roomName, false)
	if err != nil {
		retErrors = append(retErrors, retError{Code: http.StatusBadRequest, Detail: err.Error()})
	}

	if len(retErrors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(retErrors); err != nil {
			log.Println(err)
		}
	} else {
		if err = room.AddPlayer(&newPlayer); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusCreated)
			err = json.NewEncoder(w).Encode(newPlayer)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

/*PUT*/
func editGame(w http.ResponseWriter, r *http.Request) {

}

func editPlayer(w http.ResponseWriter, r *http.Request) {

}

/*DELETE*/
func deleteGame(w http.ResponseWriter, r *http.Request) {

}

func deletePlayer(w http.ResponseWriter, r *http.Request) {

}
