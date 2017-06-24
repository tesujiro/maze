package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const Width = 31
const Height = 21

const (
	WALL   = "."
	ROAD   = "#"
	START  = "S"
	FINISH = "F"
)

//type Maze [Width][Height]bool

type Point struct {
	x int
	y int
}

type Maze [][]bool

func NewMaze() Maze {
	m := make([][]bool, Width)
	for i := range m {
		m[i] = make([]bool, Height)
	}
	return Maze(m)
}

func (m Maze) setWall(p Point) {
	m[p.x][p.y] = true
}

func (m Maze) point(p Point) bool {
	return m[p.x][p.y]
}

func (m Maze) print() {
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			if m[x][Height-y-1] == true {
				fmt.Printf("%."+strconv.Itoa(x%2+1)+"s", strings.Repeat(WALL, 3))
			} else {
				fmt.Printf("%."+strconv.Itoa(x%2+1)+"s", strings.Repeat(ROAD, 3))
			}
		}
		fmt.Printf("\n")
	}
}

func (m Maze) drawLine(p1, p2 Point) {
	X, x, Y, y := p1.x, p2.x, p1.y, p2.y
	if x > X {
		X, x = p2.x, p1.x
	}
	if y > Y {
		Y, y = p2.y, p1.y
	}
	if X == x {
		for i := y; i <= Y; i++ {
			m[x][i] = true
		}
	} else if Y == y {
		for i := x; i <= X; i++ {
			m[i][y] = true
		}
	} else {
		fmt.Fprintf(os.Stderr, "drawLine error\n")
	}
}

func (m Maze) drawFrame() {
	m.drawLine(Point{0, 0}, Point{Width - 1, 0})
	m.drawLine(Point{Width - 1, 0}, Point{Width - 1, Height - 1})
	m.drawLine(Point{0, 0}, Point{0, Height - 1})
	m.drawLine(Point{0, Height - 1}, Point{Width - 1, Height - 1})
}

func (m Maze) nextTo(p Point) []Point {
	var plist []Point
	nlist := []Point{
		Point{p.x - 1, p.y - 1},
		Point{p.x - 1, p.y},
		Point{p.x - 1, p.y + 1},
		Point{p.x, p.y + 1},
		Point{p.x + 1, p.y + 1},
		Point{p.x + 1, p.y},
		Point{p.x + 1, p.y - 1},
		Point{p.x, p.y - 1},
	}
	for _, n := range nlist {
		if n.x >= 0 && n.x < Width && n.y >= 0 && n.y < Height {
			plist = append(plist, n)
		}
	}
	return plist
}

func (m Maze) canPlot(p Point) bool {
	count := 0
	for _, n := range m.nextTo(p) {
		if m.point(n) {
			count++
		}
	}
	if count < 4 {
		return true
	} else {
		return false
	}
}

func (m Maze) plotWallAtRandom() {
	//p := Point{rand.Intn(int(Width/2)) * 2, rand.Intn(int(Height/2)) * 2}
	p := Point{rand.Intn(Width), rand.Intn(Height)}
	if m.canPlot(p) {
		m.setWall(p)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	m := NewMaze()
	m.drawFrame()
	for i := 1; i < 500; i++ {
		m.plotWallAtRandom()
	}
	m.print()
}
