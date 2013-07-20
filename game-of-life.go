package main

import "fmt"

type World struct {
    cells [][]Cell
}

type Cell struct {
    x int
    y int
    alive int
}

func initWorld(size int) World {
    world := World { make([][]Cell, size)}
    cells := world.cells
    for i := range cells {
        cells[i] = make([]Cell, size)
        for j := range(cells[i]) {
            cells[i][j] = Cell{i,j,0}
        }
    }
    return world
}

func (w *World) InitBlinker() {
    world := *w
    world.cells[4][5].alive = 1
    world.cells[5][5].alive = 1
    world.cells[6][5].alive = 1
    world.cells[5][4].alive = 1
    world.cells[5][6].alive = 1
}

func main() {
    fmt.Printf("Hello, world\n")
}

