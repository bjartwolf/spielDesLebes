package main

import "fmt"
import "time"
import "os"
import "runtime"
import "os/exec"

type World struct {
    cells [][]Cell
    subscribers []chan bool // This is everyone the cell should notify
}

type Cell struct {
    x int
    y int
    alive bool
    neighbors chan bool // This is the channel that subscribes to other neighbors
    subscribers []chan bool // This is everyone the cell should notify
    done chan bool
 }

func (w *World) Subscribe (subscriber chan bool) {
    w.subscribers = append(w.subscribers,subscriber)
}

func (w *World) notify() {
    for _, subscriber := range w.subscribers {
        go func(s chan bool) {
            s<-true
        }(subscriber)
    }
}

func (c *Cell) Subscribe(subscriber chan bool) { // could return dispose method to unsubscribe like rx?  
    c.subscribers = append(c.subscribers,subscriber)
}

func (c *Cell) notify() {
    for _, subscriber := range c.subscribers {
        if (c.alive) {
            go func(s chan bool) {s<-true}(subscriber)
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
        }
    }
    return world
}

func (c *Cell) StartPlaying() {
    round:= 1
         // Notify other neigbors that if you are alive
         // Go notify?
         nrOfAliveNeighbors := 0
         c.notify()
         for {
             select {
                case <-c.neighbors:
                    nrOfAliveNeighbors++
                case <-c.done:
//                    fmt.Printf("I am cell %d, %d and I got the done signal\n", c.x, c.y)
                    if (nrOfAliveNeighbors > 0) {
                        fmt.Printf("I am cell %d, %d and I have %d alive neighbors\n", c.x, c.y, nrOfAliveNeighbors)
                    }
                    round++
//                    fmt.Println(round)
                    // Move on and kill yourself or spawn 
                    if ((nrOfAliveNeighbors== 2 || nrOfAliveNeighbors==3)&& c.alive) {
                            fmt.Println("I keep on living")
                            //c.alive = true // live ons
                        } else if (!c.alive && nrOfAliveNeighbors== 3) {
                            fmt.Println("I spawn")
                            c.alive = true
                        } else if (c.alive && nrOfAliveNeighbors> 3) {
                            fmt.Println("killing myself")
                            c.alive = false
                        } else if (c.alive && nrOfAliveNeighbors< 2) {
                            fmt.Println("killing myself since there are less than two")
                            c.alive = false
                        }
                    // Send done signal to confirm, for now just use timer in main loop
                    nrOfAliveNeighbors = 0
                    c.notify()
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
   world := *w // gues better to dereference when not mutating , or make function that takes world instead of method
   c := exec.Command("clear")
   c.Stdout = os.Stdout
   c.Run()
   cells := world.cells
   generations++
   fmt.Println(generations)
   fmt.Printf("\n")
   for i := range cells {
        for j := range cells[i] {
            if (cells[i][j].alive) {
                fmt.Printf("*")
            } else {
                fmt.Printf("X")
            }
        }
        fmt.Printf("\n")
   }
}

func main() {
    runtime.GOMAXPROCS(1)
    world := newWorld(10)
    world.InitGleiter()
   //world.InitBlinker()
    world.Print()
    timer := time.Tick(1000 * time.Millisecond)
    cells := world.cells
    for i := range cells {
        for j := range cells[i] {
           // find neighbors and add them to channellist
            if ((i-1) >= 0 && (j-1) >= 0) {
                cells[i-1][j-1].Subscribe(cells[i][j].neighbors)
            }
            if ((i-1) >= 0) {
                cells[i-1][j].Subscribe(cells[i][j].neighbors)
            }
            if ((j-1) >= 0) {
                cells[i][j-1].Subscribe(cells[i][j].neighbors)
            }
            if ((i+1) <= 9) {
                cells[i+1][j].Subscribe(cells[i][j].neighbors)
            }
            if ((j+1) <= 9) {
                cells[i][j+1].Subscribe(cells[i][j].neighbors)
            }
            if ((i+1) <= 9 && (j+1) <= 9) {
                cells[i+1][j+1].Subscribe(cells[i][j].neighbors)
            }
            if ((i-1) >= 0 && (j+1) <= 9) {
                cells[i-1][j+1].Subscribe(cells[i][j].neighbors)
            }
            if ((i+1) <= 9 && (j-1) >= 0) {
                cells[i+1][j-1].Subscribe(cells[i][j].neighbors)
            }
            world.Subscribe(cells[i][j].done)
            go func(i,j int) {
                world.cells[i][j].StartPlaying()
            }(i,j)
        }
    }
    time.Sleep(1*time.Second)
    for _ = range timer{
        world.Print()
        go world.notify()
    }
}

