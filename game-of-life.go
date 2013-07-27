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
    subscribers []chan bool // This is everyone the cell should notify about turns
}

func (w *World) Subscribe (subscriber chan bool) {
    w.subscribers = append(w.subscribers,subscriber)
}

func (w *World) proceed(msg bool) {
    for _, s := range w.subscribers {
        s<-msg
    }
}
type Cell struct {
    x int
    y int
    alive bool
    neighbors chan bool // This is the channel that subscribes to other neighbors
    subscribers []chan bool // This is everyone the cell should notify
    done chan bool
 }


func (c *Cell) Subscribe(subscriber chan bool) { // could return dispose method to unsubscribe like rx?  
    c.subscribers = append(c.subscribers,subscriber)
}

func (c *Cell) notify() {
    if (c.alive) {
        for _, s := range c.subscribers {
            s <-true
        }
    }
}

func newWorld() World {
    var subscribers []chan bool
    world := World { make([][]Cell, height), subscribers}
    cells := world.cells
    for i := range cells {
        cells[i] = make([]Cell, width)
        for j := range(cells[i]) {
            cells[i][j] = Cell{i,j,false, make(chan bool, 8), nil, make(chan bool)}
            world.Subscribe(cells[i][j].done)
        }
    }
    return world
}

func (c *Cell) StartPlaying() {
         nrOfAliveNeighbors := 0
         go func() { // Delays notification to make sure all cells are alive and recieving
                time.Sleep(500*time.Millisecond)
                c.notify()
         }()

         for {
             select {
                case <-c.neighbors:
                   nrOfAliveNeighbors++
                case <- time.Tick(time.Second):
                    if (!c.alive && nrOfAliveNeighbors== 3) {
                        c.alive = true
                    } else if (c.alive && nrOfAliveNeighbors> 3) {
                        c.alive = false
                    } else if (c.alive && nrOfAliveNeighbors< 2) {
                        c.alive = false
                    }
                    nrOfAliveNeighbors = 0
                    go func() {
                        time.Sleep(500*time.Millisecond)
                        c.notify()
                    }()
                }
          }
}

func (c *Cell) nrOfNeighbors() int{
    nrOfNeighbors := 8
    cell := *c
    x := cell.x
    y := cell.y
    if (x == 0 || x == width-1 || y == 0 || y == height-1) { // sidewall
        nrOfNeighbors -= 3
    }
    if ( (x== 0 && y == 0) ||
         (x== width-1 && y == 0) ||
         (x== 0 && y == height-1) ||
         (x== width-1 && y == height-1)) {
         nrOfNeighbors = 3
    }
    return nrOfNeighbors
}

func (w *World) InitBlinker() {
    world := *w
    world.cells[4][5].alive = true
    world.cells[5][5].alive = true
    world.cells[6][5].alive = true
}

func (w *World) InitGleiter() {
    world := *w
    world.cells[0][7].alive = true
    world.cells[1][7].alive = true
    world.cells[2][7].alive = true
    world.cells[2][8].alive = true
    world.cells[1][9].alive = true
}

func (w *World) InitToad() {
    world := *w
    world.cells[4][4].alive = true
    world.cells[4][5].alive = true
    world.cells[4][6].alive = true
    world.cells[5][5].alive = true
    world.cells[5][6].alive = true
    world.cells[5][7].alive = true
}


var generations = 0
func (w *World) Print() {
   fmt.Println(generations)
   generations++
   fmt.Printf("\n")
   for i := range w.cells {
        for j := range w.cells[i] {
            if (w.cells[i][j].alive) { fmt.Printf("*") } else {
                fmt.Printf("X")
            }
        }
        fmt.Printf("\n")
   }
}

func main() {
    runtime.GOMAXPROCS(4)
   world := newWorld()
 //   world.InitGleiter()
 //  world.InitBlinker()
   world.InitToad()
    world.Print()
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
    time.Sleep(100*time.Millisecond)
    timer := time.Tick(200* time.Millisecond)
    for _ = range timer{
        world.Print()
    }
}
