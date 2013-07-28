package main

import (
    "time"
    "runtime"
)

const height int = 20
const width int = 20

func main() {
    runtime.GOMAXPROCS(4)
    world := newWorld()
    cells := world.cells
    for i := range cells {
        for j := range cells[i] {
           // find neighbors and add them to channellist
            n := cells[i][j].neighbors
            if ((i-1) >= 0 && (j-1) >= 0) {
                cells[i-1][j-1].Subscribe(n)
            }
            if ((i-1) >= 0) {
                cells[i-1][j].Subscribe(n)
            }
            if ((j-1) >= 0) {
                cells[i][j-1].Subscribe(n)
            }
            if ((i+1) < 10) {
                cells[i+1][j].Subscribe(n)
            }
            if ((j+1) < 10) {
                cells[i][j+1].Subscribe(n)
            }
            if ((i+1) < 10 && (j+1) < 10) {
                cells[i+1][j+1].Subscribe(n)
            }
            if ((i-1) >= 0 && (j+1) < 10) {
                cells[i-1][j+1].Subscribe(n)
            }
            if ((i+1) < 10 && (j-1) >= 0) {
                cells[i+1][j-1].Subscribe(n)
            }
            go func(i,j int) {
                world.cells[i][j].Spawn()
            }(i,j)
        }
    }
 //  world.InitGleiter()
 //  world.InitBlinker()
    world.InitToad()
    world.Print()
    timer := time.Tick(1000* time.Millisecond)
    for _ = range timer{
        world.Print()
    }
}
