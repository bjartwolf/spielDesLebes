package main

import "fmt"
import "time"
import "os"
import "log"
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

func (w *World) proceed(msg bool) {
    for _, s := range w.subscribers {
        s<-msg
    }
}

func (c *Cell) Subscribe(subscriber chan bool) { // could return dispose method to unsubscribe like rx?  
    c.subscribers = append(c.subscribers,subscriber)
}

func (c *Cell) notify() {
    if (c.alive) {
        for _, s := range c.subscribers {
            go func(s chan bool) { s <-true}(s)
        }
        log.Println("*Notified")
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
     // Notify other neigbors that if you are alive
     // Go notify?
    // RACE CONDITION!!! Rewrite to use messages or something
    // Or is it that i am blocking...
    // Send message to inc-function?
    // Sending and recieving works, same nr (from 14 to 22 something depending on cores)
    // might be sending on different something... something is not sync'ed
//     if (c.alive) { go log.Println("Notified")}
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
                            if (nrOfAliveNeighbors > 0) {
            //                    fmt.Printf("I am cell %d, %d and I have %d alive neighbors\n", c.x, c.y, nrOfAliveNeighbors)
                            }
                            // Move on and kill yourself or spawn 
                            if ((nrOfAliveNeighbors== 2 || nrOfAliveNeighbors==3)&& c.alive) {
            //                        fmt.Println("I keep on living")
                                    //c.alive = true // live ons
                                } else if (!c.alive && nrOfAliveNeighbors== 3) {
            //                        fmt.Println("I spawn")
                                    c.alive = true
                                } else if (c.alive && nrOfAliveNeighbors> 3) {
            //                        fmt.Println("Starvation")
                                    c.alive = false
                                } else if (c.alive && nrOfAliveNeighbors< 2) {
            //                        fmt.Println("killing myself since there are less than two")
                                    c.alive = false
                                }
                        case false:
                         nrOfAliveNeighbors = 0
                         c.notify()
                    // Send done signal to confirm, for now just use timer in main loop
                    // Need to wait until everybody is ready to do this... Sort of post-done
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
   world := *w // gues better to dereference when not mutating , or make function that takes world instead of method
   c := exec.Command("clear")
   c.Stdout = os.Stdout
   c.Run()
   cells := world.cells
   fmt.Println(generations)
   generations++
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
    runtime.GOMAXPROCS(4)
    world := newWorld(10)
//    world.InitGleiter()
    world.InitBlinker()
    world.Print()
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
            if ((i+1) < 10) {
                cells[i+1][j].Subscribe(cells[i][j].neighbors)
            }
            if ((j+1) < 10) {
                cells[i][j+1].Subscribe(cells[i][j].neighbors)
            }
            if ((i+1) < 10 && (j+1) < 10) {
                cells[i+1][j+1].Subscribe(cells[i][j].neighbors)
            }
            if ((i-1) >= 0 && (j+1) < 10) {
                cells[i-1][j+1].Subscribe(cells[i][j].neighbors)
            }
            if ((i+1) < 10 && (j-1) >= 0) {
                cells[i+1][j-1].Subscribe(cells[i][j].neighbors)
            }
            go func(i,j int) {
                world.cells[i][j].StartPlaying()
            }(i,j)
        }
    }
    time.Sleep(100*time.Millisecond)
    timer := time.Tick(1000* time.Millisecond)
    for _ = range timer{
        world.Print()
        world.proceed(false)
        time.Sleep(500*time.Millisecond)
        world.proceed(true)
    }
}


