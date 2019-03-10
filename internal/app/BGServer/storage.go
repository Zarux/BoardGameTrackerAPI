package BGServer

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type BoardGame struct {
	Id   uint64 `json:"id" validate:"nonzero"`
	Name string `json:"name"`
}

type Player struct {
	Id       uint64    `json:"id,omitempty"`
	JoinTime time.Time `json:"joinTime,omitempty"`
	Name     string    `json:"name"`
	Color    string    `json:"color"`
}

type GamePlayer struct {
	Player Player `json:"player,omitempty"`
	Points int    `json:"points,omitempty"`
}

type Game struct {
	Id          uint64       `json:"id,omitempty"`
	Game        BoardGame    `json:"boardGame,omitempty" validate:"nonzero"`
	GameTime    time.Time    `json:"gameTime,omitempty"`
	GamePlayers []GamePlayer `json:"gamePlayers,omitempty"`
}

type room struct {
	Id         uint64    `json:"id,omitempty"`
	Hash       string    `json:"hash,omitempty"`
	CreateTime time.Time `json:"createTime,omitempty"`
	Players    []Player  `json:"players"`
	Games      []Game    `json:"games"`
}

func (r *room) GetPlayer(playerId uint64) (*Player, error) {

	for _, player := range r.Players {
		if player.Id == playerId {
			return &player, nil
		}
	}

	return nil, errors.New("could not find player")
}

func (r *room) EditPlayer(player *Player) error {
	if db == nil {
		connect()
	}
	_, err := db.Exec(`UPDATE Player SET name=?, color=? 
							WHERE room_id = ? AND player_id = ?`, player.Name, player.Color, r.Id, player.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *room) AddPlayer(player *Player) error {
	if db == nil {
		connect()
	}
	res, err := db.Exec("INSERT INTO Player(room_id, name, color) VALUES (?, ?, ?)", r.Id, player.Name, player.Color)
	if err != nil {
		return err
	} else {
		playerId, err := res.LastInsertId()
		if err != nil {
			return err
		}
		player.Id = uint64(playerId)
		player.JoinTime = time.Now()
	}
	return nil
}

func (r *room) EditGame(game *Game) error {
	return nil
}

func (r *room) AddGame(game *Game) error {
	var gameId int64
	if db == nil {
		connect()
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec("INSERT INTO Game(room_id, boardgame_id) VALUES (?, ?)", r.Id, game.Game.Id)
	if err != nil {
		_ = tx.Rollback()
		return err
	} else {
		gameId, err = res.LastInsertId()
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	game.Id = uint64(gameId)
	game.GameTime = time.Now()
	if len(game.GamePlayers) > 0 {
		query := "INSERT INTO GamePlayers(game_id, player_id, points) VALUES"
		var values []interface{}
		for _, gamePlayer := range game.GamePlayers {
			query += "(?, ?, ?),"
			values = append(values, gameId, gamePlayer.Player.Id, gamePlayer.Points)
		}
		query = query[0 : len(query)-1]
		query += " ON DUPLICATE KEY UPDATE points=VALUES(points)"

		_, err = tx.Exec(query, values...)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func addRoom(hash string) (*room, error) {
	var id int64
	if db == nil {
		connect()
	}
	res, err := db.Exec("INSERT INTO Room(room_hash) VALUES(?)", hash)
	if err != nil {
		return nil, err
	} else {
		id, err = res.LastInsertId()
		if err != nil {
			return nil, err
		}
	}
	newRoom := &room{Id: uint64(id), Hash: hash}
	return newRoom, nil
}

func GetRoom(name string, create bool) (*room, error) {
	var (
		count      int
		roomId     uint64
		roomHash   string
		createTime time.Time
	)
	if db == nil {
		connect()
	}
	hash := md5.New()
	hash.Write([]byte(name))
	hexHash := hex.EncodeToString(hash.Sum(nil))

	err := db.QueryRow("SELECT COUNT(*) FROM Room WHERE room_hash = ?", hexHash).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		if create {
			return addRoom(hexHash)
		} else {
			return nil, errors.New("could not find room")
		}
	}

	err = db.QueryRow("SELECT room_id, room_hash, create_time FROM Room WHERE room_hash = ?", hexHash).
		Scan(&roomId, &roomHash, &createTime)
	if err != nil {
		return nil, err
	}

	// For json purposes these slices are initialised
	players := []Player{}
	games := []Game{}
	room := &room{Id: roomId, Hash: roomHash, CreateTime: createTime, Players: players, Games: games}

	// Get players
	var (
		playerId   uint64
		playerName string
		color      string
		joinTime   time.Time
	)
	rows, err := db.Query(`SELECT player_id, name, color, join_time FROM Player WHERE room_id = ?`, room.Id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&playerId, &playerName, &color, &joinTime)
		if err != nil {
			return nil, err
		}
		player := Player{Id: playerId, Name: playerName, Color: color, JoinTime: joinTime}
		room.Players = append(room.Players, player)
	}

	// Get games
	var (
		gameId      uint64
		boardgameId uint64
		gameTime    time.Time
		points      int
	)
	rows, err = db.Query(`SELECT game_id, boardgame_id, game_time FROM Game WHERE room_id = ?`, room.Id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&gameId, &boardgameId, &gameTime)
		if err != nil {
			return nil, err
		}
		boardGame, err := getBoardGame(boardgameId)
		if err != nil {
			return nil, err
		}

		gpRows, err := db.Query(`SELECT player_id, points FROM GamePlayers WHERE game_id = ?`, gameId)
		if err != nil {
			return nil, err
		}
		// For json purposes this is initialised
		gamePlayers := []GamePlayer{}
		for gpRows.Next() {
			if err := gpRows.Scan(&playerId, &points); err != nil {
				return nil, err
			}
			player, err := room.GetPlayer(playerId)
			if err != nil {
				return nil, err
			}
			gamePlayer := GamePlayer{Player: *player, Points: points}
			gamePlayers = append(gamePlayers, gamePlayer)
		}
		gpRows.Close()

		game := Game{Id: gameId, GameTime: gameTime, Game: *boardGame, GamePlayers: gamePlayers}
		room.Games = append(room.Games, game)
	}

	return room, nil
}

func getBoardGame(id uint64) (*BoardGame, error) {
	var (
		err  error
		name string
	)
	if db == nil {
		connect()
	}
	err = db.QueryRow("SELECT name FROM BoardGame WHERE boardgame_id = ?", id).Scan(&name)
	if err != nil {
		return nil, err
	}
	bg := &BoardGame{Id: id, Name: name}
	return bg, nil
}

func SearchBoardGames(search string, limit int, roomId int) ([]BoardGame, error) {
	var (
		id         uint64
		name       string
		boardGames []BoardGame
		query      string
		params     []interface{}
	)
	if db == nil {
		connect()
	}
	if roomId != 0 {
		query = `SELECT DISTINCT boardgame_id, name FROM (
  					SELECT boardgame_id, name, 1 AS sort_score 
					FROM Game 
					LEFT JOIN BoardGame BG USING(boardgame_id) 
					WHERE room_id = 1 AND name LIKE CONCAT('%', ?, '%') 
					UNION 
					SELECT boardgame_id, name, 0 AS sort_score 
					FROM BoardGame 
					WHERE name LIKE CONCAT('%', ?, '%') 
					ORDER BY sort_score DESC, LENGTH(name) ASC, boardgame_id ASC 
					LIMIT ?
				) X`
		params = append(params, search, search, limit)
	} else {
		query = `SELECT boardgame_id, name 
				FROM BoardGame 
				WHERE name LIKE CONCAT('%', ?, '%')
				ORDER BY LENGTH(name) ASC, boardgame_id ASC
				LIMIT ?`
		params = append(params, search, limit)
	}
	rows, err := db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		game := BoardGame{Id: id, Name: name}
		boardGames = append(boardGames, game)
	}
	return boardGames, nil
}
