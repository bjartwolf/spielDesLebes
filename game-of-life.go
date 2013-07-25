package main

import (
    "fmt"
    "time"
    "runtime"
)

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
            go func(s chan bool) { s <-true}(s)
        }
    }
}

func newWorld(size int) World {
    var subscribers []chan bool
    world := World { make([][]Cell, size), subscribers}
    cells := world.cells
    for i := range cells {
        cells[i] = make([]Cell, size)
        for j := range(cells[i]) {
            cells[i][j] = Cell{i,j,false, make(chan bool, 8), nil, make(chan bool)}
            world.Subscribe(cells[i][j].done)
        }
    }
    return world
}

func (c *Cell) StartPlaying() {
     for {
         nrOfAliveNeighbors := 0
         for {
             select {
                case <-c.neighbors:
                   nrOfAliveNeighbors++
                case msg := <-c.done:
                    switch msg {
//                    fmt.Printf("I am cell %d, %d and I got the done signal\n", c.x, c.y)
                        case true:
                            // Ignore rule about keep on living
                                if (!c.alive && nrOfAliveNeighbors== 3) {
                                    c.alive = true
                                } else if (c.alive && nrOfAliveNeighbors> 3) {
                                    c.alive = false
                                } else if (c.alive && nrOfAliveNeighbors< 2) {
                                    c.alive = false
                                }
                        case false:
                            nrOfAliveNeighbors = 0
                            c.notify()
                    }
                }
             }
         }
}

func (c *Cell) nrOfNeighbors() int{
    nrOfNeighbors := 8
    cell := *c
    x := cell.x
    y := cell.y
    if (x == 0 || x == 9 || y == 0 || y == 9) { // sidewall
        nrOfNeighbors -= 3
    }
    if ( (x== 0 && y == 0) ||
         (x== 9 && y == 0) ||
         (x== 0 && y == 9) ||
         (x== 9 && y == 9)) {
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
    world := newWorld(10)
 //   world.InitGleiter()
   world.InitBlinker()
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
    timer := time.Tick(40* time.Millisecond)
    i:= 0
    for _ = range timer{
        i++
        if (i % 11 == 0) {
            world.Print()
        }
        world.proceed(false)
        time.Sleep(30*time.Millisecond)
        world.proceed(true)
    }
}
