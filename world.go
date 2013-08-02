package main
import "fmt"

type World struct {
    cells [][]*Cell
}

func newWorld() World {
    world := World { make([][]*Cell, height)}
    cells := world.cells
    for i := range cells {
        cells[i] = make([]*Cell, width)
        for j := range(cells[i]) {
            cells[i][j] = NewCell(i,j)
        }
    }
    return world
}


func (w *World) InitBlinker() {
    world := *w
    world.cells[4][5].våkne()
    world.cells[5][5].våkne()
    world.cells[6][5].våkne()
}

func (w *World) InitGleiter() {
    world := *w
    world.cells[0][7].våkne()
    world.cells[1][7].våkne()
    world.cells[2][7].våkne()
    world.cells[2][8].våkne()
    world.cells[1][9].våkne()
}

func (w *World) InitToad() {
    world := *w
    world.cells[4][4].våkne()
    world.cells[4][5].våkne()
    world.cells[4][6].våkne()
    world.cells[5][5].våkne()
    world.cells[5][6].våkne()
    world.cells[5][7].våkne()
}


func (w *World) Print() {
   for i := range w.cells {
        for j := range w.cells[i] {
            if (w.cells[i][j].levande) {
                fmt.Printf("*")
            } else {
                fmt.Printf("X")
            }
        }
        fmt.Printf("\n")
   }
   fmt.Printf("\n")
}
