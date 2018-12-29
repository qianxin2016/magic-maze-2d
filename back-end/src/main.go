// fileUpload project main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sort"
	"encoding/hex"
	"crypto/sha256"
)

type Player struct {
	Account string
	Txid string
	Position uint16
	Steps uint64
	Gems []uint16
}

type Battle struct {
	Id string
	Name string
	MazeHash string
	Bonus uint64
	Creator *Player
	Challenger *Player
	Winner *Player
	maze Maze
}

func (b Battle) String() string {
	return fmt.Sprintf("[Id: %v, Name: %v, MazeHash: %v, Bonus %v, Creator: %v, Challenger: %v, Winner: %v]",
		b.Id,
		b.Name,
		b.MazeHash,
		b.Bonus,
		b.Creator,
		b.Challenger,
		b.Winner)
}

const PAGE_SIZE uint64 = 100
var mazes map[string]Maze = make(map[string]Maze)
var battles map[string]*Battle = make(map[string]*Battle)

func getParam(r *http.Request, key string) string {
	value := r.URL.Query().Get(key)
	if value == "" {
		r.ParseForm()
		value = r.FormValue(key)
	}
	return value
}

func echoJson(w http.ResponseWriter, data map[string]interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func hash(s string) string {
	h := sha256.New()
    h.Write([]byte(s))
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}

func newmaze(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--Received newmaze request--")
	enableCors(&w)
	id := getParam(r, "id")

	if id == "" {
		fmt.Println("Missing params")
		data := map[string]interface{}{
			"status": "Missing params",
		}
		echoJson(w, data)
		return
	}

	maze := generateMaze()
	mazeHash := hash(flatten(maze))
	mazes[id] = maze

	fmt.Printf("maze generated for %v: %v\n", id, mazeHash)
	data := map[string]interface{}{
		"status": "OK",
		"mazeHash": mazeHash,
	}
	echoJson(w, data)
}

func create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--Received create request--")
	enableCors(&w)
	id := getParam(r, "id")
	name := getParam(r, "name")
	txid := getParam(r, "txid")
	bonus := getParam(r, "bonus")
	account := getParam(r, "account")
	mazeHash := getParam(r, "mazeHash")

	if id == "" || name == "" || txid == "" || bonus == "" || account == "" {
		fmt.Println("Missing params")
		data := map[string]interface{}{
			"status": "Missing params",
		}
		echoJson(w, data)
		return
	}

	bonusInt, err := strconv.ParseUint(bonus, 10, 16)
	if err == nil {
		battle := &Battle {
			Id: id,
			Name: name,
			MazeHash: mazeHash,
			Bonus: bonusInt,
			Creator: &Player{
				Account: account,
				Txid: txid,
				Position: 0,
				Steps: 0,
				Gems: []uint16{},
			},
			Challenger: nil,
			Winner: nil,
			maze: mazes[id],
		}

		battles[id] = battle
		fmt.Printf("Battle created: %v\n", battle)
		delete(mazes, id)

		data := map[string]interface{}{
			"status": "OK",
			"mazeHash": battle.MazeHash,
		}
		echoJson(w, data)
	} else {
		fmt.Println("Invalid bonus")
		data := map[string]interface{}{
			"status": "Invalid bonus",
		}
		echoJson(w, data)
	}
}

func join(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--Received join request--")
	enableCors(&w)
	id := getParam(r, "id")
	txid := getParam(r, "txid")
	account := getParam(r, "account")

	if id == "" || txid == "" || account == "" {
		fmt.Println("Missing params")
		data := map[string]interface{}{
			"status": "Missing params",
		}
		echoJson(w, data)
		return
	}

	if battle, ok := battles[id]; ok {
		battle.Challenger = &Player {
			Account: account,
			Txid: txid,
			Position: 0,
			Steps: 0,
			Gems: []uint16{},
		}
	} else {
		fmt.Println("Invalid battle id")
	}

	fmt.Printf("%v has joint room %v\n", account, id)
	data := map[string]interface{}{
		"status": "OK",
	}
	echoJson(w, data)
}

func roomlist(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--Received rootlist request--")
	enableCors(&w)
	page := getParam(r, "page")
	pageInt, err := strconv.ParseUint(page, 10, 16)
	if err == nil {
		start := PAGE_SIZE * pageInt

		sorted_keys := make([]string, 0)
		for k, _ := range battles {
			sorted_keys = append(sorted_keys, k)
		}
		sort.Sort(sort.Reverse(sort.StringSlice(sorted_keys)))

		data := make(map[string]interface{})
		for _, k := range sorted_keys {
			if start > 0 {
				start--
				continue
			}

			data[k] = battles[k]
			fmt.Printf("room: %v\n", data[k])
			if uint64(len(data)) >= PAGE_SIZE {
				break
			}
		}

		if start + uint64(len(data)) < uint64(len(battles)) {
			data["more"] = nil
		}
		echoJson(w, data)
	} else {
		fmt.Println("Invalid page")
	}
}

func getPlayer(battle *Battle, account string) *Player {
	if battle.Creator.Account == account {
		return battle.Creator

	} else if (battle.Challenger != nil && battle.Challenger.Account == account) {
		return battle.Challenger
	} else {
		return nil
	}
}

func tryCollectGem(player *Player, maze Maze, index int64) {
	if maze[index].hasGem {
		player.Gems = append(player.Gems, uint16(index))
		maze[index].hasGem = false
		fmt.Printf("%v collect %v gems!\n", player.Account, len(player.Gems))
	}
}

func step(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--Received step request--")
	enableCors(&w)
	id := getParam(r, "id")
	f := getParam(r, "from")
	t := getParam(r, "to")
	account := getParam(r, "account")

	if id == "" || f == "" || t == "" || account == "" {
		fmt.Println("Invalid params")
		data := map[string]interface{}{
			"status": "Invalid params",
		}
		echoJson(w, data)
		return
	}

	data := make(map[string]interface{})
	if battle, ok := battles[id]; ok {
		player := getPlayer(battle, account)
		if player == nil {
			data["status"] = "Not a player"
			fmt.Println("Not a player")
			echoJson(w, data)
			return
		}

		from, err1 := strconv.ParseInt(f, 10, 16)
		to, err2 := strconv.ParseInt(t, 10, 16)

		if err1 == nil && err2 == nil && from < MAZE_SIZE && to < MAZE_SIZE {
			data["status"] = "OK"
			maze := battle.maze
			switch (from - to) {
			case 1:
				if maze[from].left || maze[to].right {
					data["hitWall"] = true
				} else {
					data["hitWall"] = false
					player.Position = uint16(to)
					tryCollectGem(player, maze, to)
				}
			case -1:
				if maze[from].right || maze[to].left {
					data["hitWall"] = true
				} else {
					data["hitWall"] = false
					player.Position = uint16(to)
					tryCollectGem(player, maze, to)
				}
			case ROW_SIZE:
				if maze[from].top || maze[to].bottom {
					data["hitWall"] = true
				} else {
					data["hitWall"] = false
					player.Position = uint16(to)
					tryCollectGem(player, maze, to)
				}
			case -ROW_SIZE:
				if maze[from].bottom || maze[to].top {
					data["hitWall"] = true
				} else {
					data["hitWall"] = false
					player.Position = uint16(to)
					tryCollectGem(player, maze, to)
				}
			default:
				data["status"] = "Can not jump"
				fmt.Println("Can not jump")
				echoJson(w, data)
				return
			}

			fmt.Printf("step from %v to %v, hitWall: %v\n", from, to, data["hitWall"])
			player.Steps++
		} else {
			data["status"] = "Invalid position"
			fmt.Println("Invalid position")
		}
	} else {
		data["status"] = "Invalid maze id"
		fmt.Println("Invalid maze id")
	}

	echoJson(w, data)
}

func claim(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--Received claim request--")
	enableCors(&w)
	id := getParam(r, "id")
	account := getParam(r, "account")

	if id == "" || account == "" {
		fmt.Println("Missing params")
		data := map[string]interface{}{
			"status": "Missing params",
		}
		echoJson(w, data)
		return
	}

	fmt.Printf("%v is trying to claim winner of %v\n", account, id)

	data := make(map[string]interface{})
	if battle, ok := battles[id]; ok {
		if battle.Creator.Account == account {
			if len(battle.Creator.Gems) >= 5 {
				battle.Winner = battle.Creator
				data["status"] = "OK"
				data["mazeInfo"] = flatten(battle.maze)
				fmt.Println("Congratulations!")
			} else {
				data["status"] = "Insufficient gems"
				fmt.Println("Insufficient gems")
			}
		} else if battle.Challenger != nil && battle.Challenger.Account == account {
			if len(battle.Challenger.Gems) >= 5 {
				battle.Winner = battle.Challenger
				data["status"] = "OK"
				data["mazeInfo"] = flatten(battle.maze)
				fmt.Println("Congratulations!")
			} else {
				data["status"] = "Insufficient gems"
				fmt.Println("Insufficient gems")
			}
		} else {
			data["status"] = "Not a player"
			fmt.Println("Not a player")
		}
	} else {
		data["status"] = "Invalid battle id"
		fmt.Println("Invalid battle id")
	}

	echoJson(w, data)
}

func main() {
	http.HandleFunc("/newmaze", newmaze)
	http.HandleFunc("/create", create)
	http.HandleFunc("/join", join)
	http.HandleFunc("/roomlist", roomlist)
	http.HandleFunc("/step", step)
	http.HandleFunc("/claim", claim)

	err := http.ListenAndServe(":6688", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
