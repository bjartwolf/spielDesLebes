package main

import (
    "time"
    "fmt"
    "runtime"
)

const height int = 10
const width int = 10

func main() {
    runtime.GOMAXPROCS(4)
    world := newWorld()
    cells := world.cells
    for i := range cells {
        for j := range cells[i] {
           // find neighbors and add them to channellist
            n := cells[i][j].neighbors
            if ((i-1) >= 0 && (j-1) >= 0) {
                cells[i-1][j-1].Subscribe <- n
            }
            if ((i-1) >= 0) {
                cells[i-1][j].Subscribe <- n
            }
            if ((j-1) >= 0) {
                cells[i][j-1].Subscribe <- n
            }
            if ((i+1) < height) {
                cells[i+1][j].Subscribe <- n
            }
            if ((j+1) < width) {
                cells[i][j+1].Subscribe <- n
            }
            if ((i+1) < height && (j+1) < width) {
                cells[i+1][j+1].Subscribe <- n
            }
            if ((i-1) >= 0 && (j+1) < width) {
                cells[i-1][j+1].Subscribe <- n
            }
            if ((i+1) < height && (j-1) >= 0) {
                cells[i+1][j-1].Subscribe <- n
            }
//            fmt.Printf("Cell %d, %d has %d neighbors\n", 
            // Starts the game, but all cells are still dead
            go func(i,j int) {
                time.Sleep(time.Second)
                fmt.Println(len(world.cells[i][j].subscribers))
                world.cells[i][j].Spawn()
            }(i,j)
        }
    }
   world.InitGleiter()
 //  world.InitBlinker()
 //   world.InitToad()
    world.Print()
    timer := time.Tick(1000* time.Millisecond)
    for _ = range timer{
        world.Print()
    }
}
