package main

import (
    "fmt"
    "time"
    "runtime"
)

const height int = 10
const width int = 10

type World struct {
    cells [][]Cell
}

type Cell struct {
    x int
    y int
    alive bool
    neighbors chan string// This is the channel that subscribes to other neighbors
    subscribers []chan string// This is everyone the cell should notify done chan bool
 }

// ctrl+v u nnnn
func (c *Cell) Subscribe(subscriber chan string) { // could return dispose method to unsubscribe like rx?  
    c.subscribers = append(c.subscribers,subscriber)
}

func (c *Cell) Die() {
    c.alive = false
    c.notify("eg đøyr")
}

func (c *Cell) Spawn() {
    c.alive = true
    c.notify("eg łevar")
}

func (c *Cell) notify(msg string) {
    go func(c *Cell, msg string) {
            time.Sleep(500*time.Millisecond)
            for _, s := range c.subscribers {
                s <- msg
            }
    }(c, msg)
}

func newWorld() World {
    world := World { make([][]Cell, height)}
    cells := world.cells
    for i := range cells {
        cells[i] = make([]Cell, width)
        for j := range(cells[i]) {
            cells[i][j] = Cell{i,j,false, make(chan string, 8), nil}
        }
    }
    return world
}

func (c *Cell) StartPlaying() {
         nrOfAliveNeighbors := 0
         for {
             select {
                 case msg := <-c.neighbors:
                    switch msg {
                        case "eg đøyr":
                            nrOfAliveNeighbors--
                        case "eg łevar":
                            nrOfAliveNeighbors++
                    }
                 case <- time.Tick(time.Second):
                    if (nrOfAliveNeighbors > 0) {
//                        fmt.Printf("Nr of: %d for cell %d, %d\n", nrOfAliveNeighbors, c.x, c.y)
                    }
                    if (!c.alive && nrOfAliveNeighbors== 3) {
                        c.Spawn()
                    } else if (c.alive && nrOfAliveNeighbors> 3) {
                        c.Die()
                    } else if (c.alive && nrOfAliveNeighbors< 2) {
                        c.Die()
                    }
                }
          }
}

func (w *World) InitBlinker() {
    world := *w
    world.cells[4][5].Spawn()
    world.cells[5][5].Spawn()
    world.cells[6][5].Spawn()
}

func (w *World) InitGleiter() {
    world := *w
    world.cells[0][7].Spawn()
    world.cells[1][7].Spawn()
    world.cells[2][7].Spawn()
    world.cells[2][8].Spawn()
    world.cells[1][9].Spawn()
}

func (w *World) InitToad() {
    world := *w
    world.cells[4][4].Spawn()
    world.cells[4][5].Spawn()
    world.cells[4][6].Spawn()
    world.cells[5][5].Spawn()
    world.cells[5][6].Spawn()
    world.cells[5][7].Spawn()
}


var generations = 0
func (w *World) Print() {
   fmt.Println(generations)
   generations++
   fmt.Printf("\n")
   for i := range w.cells {
        for j := range w.cells[i] {
            if (w.cells[i][j].alive) {
                fmt.Printf("*")
            } else {
                fmt.Printf("X")
            }
        }
        fmt.Printf("\n")
   }
}

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
                world.cells[i][j].StartPlaying()
            }(i,j)
        }
    }
 //   world.InitGleiter()
 //  world.InitBlinker()
   // need to init with goroutine
    time.Sleep(50 * time.Millisecond)
    world.InitToad()
    time.Sleep(50 * time.Millisecond)
    world.Print()
    timer := time.Tick(500* time.Millisecond)
    for _ = range timer{
        world.Print()
    }
}
