package main

import "fmt"
import "time"
import "os"
import "runtime"
import "os/exec"

type World struct {
    cells [][]Cell
}

type Cell struct {
    x int
    y int
    alive int
    outbox chan int 
}

func newWorld(size int) World {
    world := World { make([][]Cell, size)}
    cells := world.cells
    for i := range cells {
        cells[i] = make([]Cell, size)
        for j := range(cells[i]) {
            cells[i][j] = Cell{i,j,0, make(chan int, 8)}
        }
    }
    return world
}

func (w *World) NextRound() {
    world := *w
    cells := world.cells
    // All cells should listen for messages from neighbors and then do calculations when messages
    // have been recievd
    for i := range cells {
        for j := range cells[i] {
            go func(i, j int) {
                cell := cells[i][j]
                aliveNeighbors := 0
                    if ((i-1) >= 0 && (j-1) >= 0) {
                        aliveNeighbors += <-cells[i-1][j-1].outbox
                    }
                    if ((i-1) >= 0) {
                        aliveNeighbors += <-cells[i-1][j].outbox
                    }
                    if ((j-1) >= 0) {
                        aliveNeighbors += <-cells[i][j-1].outbox
                    }
                    if ((i+1) <= 9) {
                        aliveNeighbors += <-cells[i+1][j].outbox
                    }
                    if ((j+1) <= 9) {
                        aliveNeighbors += <-cells[i][j+1].outbox
                    }
                    if ((i+1) <= 9 && (j+1) <= 9) {
                        aliveNeighbors += <-cells[i+1][j+1].outbox
                    }
                    if ((i-1) >= 0 && (j+1) <= 9) {
                        aliveNeighbors += <-cells[i-1][j+1].outbox
                    }
                    if ((i+1) <= 9 && (j-1) >= 0) {
                        aliveNeighbors += <-cells[i+1][j-1].outbox
                    }
                // Had to go through world to make it work...
                if ((aliveNeighbors == 2 || aliveNeighbors ==3)&& cell.alive == 1) {
                } else if (cell.alive == 0 && aliveNeighbors == 3) {
                    world.cells[i][j].alive = 1
                } else if (aliveNeighbors > 3) {
                    world.cells[i][j].alive = 0
                } else if (cell.alive == 1 && aliveNeighbors < 2) {
                    world.cells[i][j].alive = 0
                }
            }(i, j)
        }
    }
    // Then send all messages
    for i := range cells {
        for j := range cells[i] {
            cell := cells[i][j]
            for messages := 0; messages < cell.nrOfNeighbors(); messages++ {
                if (cell.alive == 1) {
                    cell.outbox <- 1
                } else {
                    cell.outbox <- 0
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
    world.cells[4][5].alive = 1
    world.cells[5][5].alive = 1
    world.cells[6][5].alive = 1
}

func (w *World) InitGleiter() {
    world := *w
    world.cells[0][7].alive = 1
    world.cells[1][7].alive = 1
    world.cells[2][7].alive = 1
    world.cells[2][8].alive = 1
    world.cells[1][9].alive = 1
}

func (w *World) Print() {
   c := exec.Command("clear")
   c.Stdout = os.Stdout
   c.Run()
   world := *w
   cells := world.cells
   fmt.Printf("\n")
   for i := range cells {
        for j := range cells[i] {
            fmt.Printf("%d ", cells[i][j].alive)
        }
        fmt.Printf("\n")
   }
}

func main() {
    runtime.GOMAXPROCS(2)
    world := newWorld(10)
    world.InitGleiter()
    //world.InitBlinker()
    world.Print()
    worlds := 0
    c := time.Tick(10 * time.Second)
    go func() {for now := range c {
        fmt.Printf("Worlds: %d\n", worlds, now)
    }}()
    for {
        //time.Sleep(100000000)
        worlds++
        world.NextRound()
//        world.Print()
    }
}

