package main

import "fmt"
type Cell struct {
    x int
    y int
    alive int
}

func initWorld(size int) [][]Cell {
    world := make([][]Cell, size)
    for i := range world {
        world[i] = make([]Cell, size)
    }
    return world
}
func main() {
    fmt.Printf("Hello, world\n")
}

