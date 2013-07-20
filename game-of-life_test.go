package main

import "testing"

// In the future they should not really be in a grid
// but it is much easier for now...
func Test_initialize_world_correct_size(t *testing.T) {
    world := initWorld(10)
    if (len(world.cells) != 10) {
        t.Fatal("World should have 10 columns")
    }
    if (len (world.cells[5]) != 10) {
        t.Fatal("World should have 10 rows")
    }
}

func Test_initialize_world_correct_positions(t *testing.T) {
    world := initWorld(10)
    if (world.cells[5][5].x != 5) {
        t.Fatal("should be in position 5")
    }
}

func Test_initialize_blinker(t *testing.T) {
    world := initWorld(10)
    world.InitBlinker()
    if (world.cells[5][7].alive != 0 ||
        world.cells[6][7].alive != 0 ||
        world.cells[4][8].alive != 0) {
        t.Fatal("All these cells should be dead")
    }
    if (world.cells[5][6].alive != 1 ||
        world.cells[6][5].alive != 1 ||
        world.cells[4][5].alive != 1) {
        t.Fatal("All these cells should be alive")
    }
}
