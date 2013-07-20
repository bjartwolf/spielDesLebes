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

func (w *World) Print() {
   world := *w
   cells := world.cells
   for i := range cells {
        for j := range cells {
            fmt.Printf("%d ", cells[i][j].alive)
        }
        fmt.Printf("\n")
   }
}

func main() {
    world := initWorld(10)
    world.InitBlinker()
    world.Print()
}

