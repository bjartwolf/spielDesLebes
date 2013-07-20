package main

import "fmt"
import "time"

type World struct {
    cells [][]Cell
}

type Cell struct {
    x int
    y int
    alive int
    //inbox chan int // this is where you recieve messages if neighbors are dead or alive
    outbox chan int // this is where you send messages if neighbors are dead or alive
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
    return world }

func (w *World) NextRound() {
    world := *w
    cells := world.cells
    // All cells should listen for messages from neighbors and then do calculations when messages
    // have been recievd
    for i := range cells {
        for j := range cells {
            cell := cells[i][j]
            go func(cell Cell) {
                aliveNeighbors := 0
                i := cell.x
                j := cell.y
                    if ((i-1) >= 0 && (j-1) >= 0) {
                        var k int  = <-cells[i-1][j-1].outbox
                        aliveNeighbors += k
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
                if ((aliveNeighbors == 2 || aliveNeighbors ==3)&& cell.alive == 1) {
                    fmt.Printf("cell keeps on living with 2 or 3 neighbors\n")
                } else if (cell.alive == 0 && aliveNeighbors == 3) {
                    fmt.Printf("Cell spawned in position %d, %d\n", cell.x, cell.y)
                    cell.alive = 1
                } else if (aliveNeighbors > 3) {
                    cell.alive = 0
                    fmt.Printf("Killed a cell in position %d, %d\n", cell.x, cell.y)
                } else if (cell.alive == 1 && aliveNeighbors < 2) {
                    fmt.Printf("cell dies because fewer than 2 neighbors alive in position %d, %d\n", cell.x, cell.y)
                    cell.alive = 0
                }
            }(cell)
        }
    }
    // Then send all messages
    for i := range cells {
        for j := range cells[i] {
            cell := cells[i][j]
            go func(cell Cell) {
                for messages := 0; messages < cell.nrOfNeighbors(); messages++ {
                    if (cell.alive == 1) {
                        cell.outbox <- 1
                    } else {
                        cell.outbox <- 0
                    }
                }
            }(cell)
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

func (w *World) Print() {
   world := *w
   cells := world.cells
   fmt.Printf("\n")
   for i := range cells {
        for j := range cells {
            fmt.Printf("%d ", cells[i][j].alive)
        }
        fmt.Printf("\n")
   }
}

func main() {
    world := newWorld(10)
    world.InitBlinker()
    world.Print()
    world.NextRound()
    //time.Sleep(100000000)
    time.Sleep(1)
    world.Print()
}

