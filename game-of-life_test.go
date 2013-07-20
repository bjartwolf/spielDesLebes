package main

import "testing"

func Test_initialize_world(t *testing.T) {
    world := initWorld(10)
    if (len(world) != 10) {
        t.Fatal("World should have 10 columns")
    }
    if (len (world[5]) != 10) {
        t.Fatal("World should have 10 rows")
    }
}
