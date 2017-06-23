package main

import (
	"fmt"
	"strconv"
)

const Height = 20
const Width = 30

//type Maze [Width][Height]bool

type Maze [][]bool

func NewMaze() Maze {
	m := make([][]bool, Width)
	for i := range m {
		m[i] = make([]bool, Height)
	}
	return Maze(m)
}

func (m Maze) print() {
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			//fmt.Printf("(%v,%v)", x, y)
			//if m[x][y] == true {
			if m[x][Height-y-1] == true {
				fmt.Printf("%."+strconv.Itoa(x%2+1)+"s", "###")
			} else {
				fmt.Printf("%."+strconv.Itoa(x%2+1)+"s", "...")
			}
		}
		fmt.Printf("\n")
	}
}

func main() {
	m := NewMaze()
	m.print()
}
