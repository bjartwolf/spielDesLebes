package main

import (
    "time"
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
            for vert := -1; vert <= 1; vert++ {
              for horiz := -1; horiz <= 1; horiz++ {
                if vert == 0 && horiz == 0 {
                  continue
                }
                vertical := (i+vert) % height
                horizontal := (j+horiz) % width 
                if vertical < 0 {
                  vertical += height
                }
                if horizontal < 0 {
                  horizontal += width
                }
                cells[vertical][horizontal].Subscribe <- n
              }
            }
            // Starts the game, but all cells are still dead
            go func(i,j int) {
                time.Sleep(time.Second)
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
