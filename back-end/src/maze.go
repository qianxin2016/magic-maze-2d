package main

import (
    "fmt"
    "time"
    "math/rand"
)

type Stack []int

func (s Stack) Push(v int) Stack {
    return append(s, v)
}

func (s Stack) Pop() (Stack, int) {
    l := len(s)
    if l == 0 {
        return s, 0
    }
    return s[:l - 1], s[l - 1]
}

func ( s Stack) Peek() int {
    l := len(s)
    if l == 0 {
        return 0
    }

    return s[l - 1]
}

func (s Stack) Length() int {
    return len(s)
}

type Cell struct {
    visited bool
    hasGem bool
    left bool
    right bool
    top bool
    bottom bool
}

type Maze []*Cell

const DEBUG = false
const ROW_SIZE = 7
const MAZE_SIZE = ROW_SIZE * ROW_SIZE
const LEFT = -2
const RIGHT = 2
const TOP = -14
const BOTTOM = 14
const ENTRY_NUM = 2

var gemSpots []int = []int{
    8, 10, 12, 22, 24, 26, 36, 38, 40,
}

var stack Stack

func printMaze(maze Maze) {
    start := 0
    for i := 0; i < 7; i++ {
        start = i * 7

        for j := start; j < start+7; j++ {
            if maze[j].top {
                fmt.Print("\u3000\u2501")
            } else {
                fmt.Print("\u3000\u3000")
            }
        }
        fmt.Print("\n")

        for j := start; j < start+7; j++ {
            if maze[j].left {
                if isGemSpot(j) {
                    fmt.Print("\u2503\u2573")
                } else {
                    fmt.Print("\u2503\u3000")
                }
            } else {
                if isGemSpot(j) {
                    fmt.Print("\u3000\u2573")
                } else {
                    fmt.Print("\u3000\u3000")
                }
            }
        }
        fmt.Print("\u2503\n")
    }

    for i := 0; i < 7; i++ {
        fmt.Print("\u3000\u2501")
    }
    fmt.Print("\n")
}

func isGemSpot(spot int) bool {
    return (spot < MAZE_SIZE) && (spot % 2 == 0) && (spot % 14 >= 8)
}

func getNeighbours(maze Maze, footHold int) (neighbours []int) {
    if isGemSpot(footHold + LEFT) && !maze[footHold + LEFT].visited {
        neighbours = append(neighbours, footHold + LEFT)
    }

    if isGemSpot(footHold + RIGHT) && !maze[footHold + RIGHT].visited {
        neighbours = append(neighbours, footHold + RIGHT)
    }

    if isGemSpot(footHold + TOP) && !maze[footHold + TOP].visited {
        neighbours = append(neighbours, footHold + TOP)
    }

    if isGemSpot(footHold + BOTTOM) && !maze[footHold + BOTTOM].visited {
        neighbours = append(neighbours, footHold + BOTTOM)
    }

    return
}

func initCells() Maze {
    var maze Maze = make(Maze, MAZE_SIZE)
    for i := 0; i < MAZE_SIZE; i++ {
        maze[i] = &Cell {
            visited: false,
            hasGem: false,
            left: true,
            right: true,
            top: true,
            bottom: true,
        }

        if i % 7 == 0 || i % 7 == 6 {
            maze[i].top = false
            maze[i].bottom = false
        }

        if (i > 0 && i < 6) || (i > 42 && i < 48) {
            maze[i].left = false
            maze[i].right = false
        }
    }

    maze[0] = &Cell {false, false, true, false, true, false}
    maze[6] = &Cell {false, false, false, true, true, false}
    maze[42] = &Cell {false, false, true, false, false, true}
    maze[48] = &Cell {false, false, false, true, false, true}

    for _, spot := range gemSpots {
        maze[spot].hasGem = true
    }

    return maze
}

func connectGems(maze Maze) {
    rand.Seed(time.Now().UnixNano())
    footHold := gemSpots[rand.Intn(len(gemSpots))]
    maze[footHold].visited = true
    stack = stack.Push(footHold)
    if (DEBUG) {
        fmt.Println("start: ", footHold)
    }

    for {
        if stack.Length() == 0 {
            break
        }

        neighbours := getNeighbours(maze, footHold)
        if len(neighbours) > 0 {
            next := neighbours[rand.Intn(len(neighbours))]

            switch next {
            case footHold + LEFT:
                maze[footHold].left = false
                maze[next].right = false
                maze[footHold + LEFT/2].left = false
                maze[footHold + LEFT/2].right = false
                if (DEBUG) {
                    fmt.Println("move left: ", next)
                }
            case footHold + RIGHT:
                maze[footHold].right = false
                maze[next].left = false
                maze[footHold + RIGHT/2].left = false
                maze[footHold + RIGHT/2].right = false
                if (DEBUG) {
                    fmt.Println("move right: ", next)
                }
            case footHold + TOP:
                maze[footHold].top = false
                maze[next].bottom = false
                maze[footHold + TOP/2].top = false
                maze[footHold + TOP/2].bottom = false
                if (DEBUG) {
                    fmt.Println("move up: ", next)
                }
            case footHold + BOTTOM:
                maze[footHold].bottom = false
                maze[next].top = false
                maze[footHold + BOTTOM/2].top = false
                maze[footHold + BOTTOM/2].bottom = false
                if (DEBUG) {
                    fmt.Println("move down: ", next)
                }
            }

            footHold = next
            maze[footHold].visited = true
            stack = stack.Push(footHold)
        } else {
            stack, _ = stack.Pop()
            footHold = stack.Peek()
            if (DEBUG) {
                fmt.Println("backtracking: ", footHold)
            }
        }
    }
}

func digEntryPoint(maze Maze) {
    scope := ROW_SIZE - 2
    for i := 0; i < ENTRY_NUM; i++ {
        random := rand.Intn(scope)
        maze[random + 1].bottom = false
        maze[random + 1 + ROW_SIZE].top = false
    }

    for i := 0; i < ENTRY_NUM; i++ {
        random := rand.Intn(scope)
        maze[random + ROW_SIZE * (ROW_SIZE - 2)].bottom = false
        maze[random + ROW_SIZE * (ROW_SIZE - 1)].top = false
    }

    for i := 0; i < ENTRY_NUM; i++ {
        random := rand.Intn(scope)
        maze[(random + 1) * ROW_SIZE].right = false
        maze[(random + 1) * ROW_SIZE + 1].left = false
    }

    for i := 0; i < ENTRY_NUM; i++ {
        random := rand.Intn(scope)
        maze[(random + 1) * ROW_SIZE + ROW_SIZE - 2].right = false
        maze[(random + 1) * ROW_SIZE + ROW_SIZE - 1].left = false
    }
}

func generateMaze() Maze {
    maze := initCells()
    connectGems(maze)
    digEntryPoint(maze)
    printMaze(maze)
    return maze
}

func flatten(maze Maze) string {
    str := ""
    for i := 0; i < MAZE_SIZE; i++ {
        if maze[i].left {
            str += "1"
        } else {
            str += "0"
        }

        if maze[i].right {
            str += "1"
        } else {
            str += "0"
        }

        if maze[i].top {
            str += "1"
        } else {
            str += "0"
        }

        if maze[i].bottom {
            str += "1"
        } else {
            str += "0"
        }
    }
    return str
}
