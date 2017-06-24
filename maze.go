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
	WALL   = "#"
	ROAD   = " "
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
	//m.print()
}

func (m Maze) setRoad(p Point) {
	m[p.x][p.y] = false
	//m.print()
}

func (m Maze) point(p Point) bool {
	return m[p.x][p.y]
}

func (m Maze) print() {
	fmt.Print("\x1b[H\x1b[2J") // Clear Screen
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			if m.point(Point{x, Height - y - 1}) == true {
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
			m.setWall(Point{x, i})
		}
	} else if Y == y {
		for i := x; i <= X; i++ {
			m.setWall(Point{i, y})
		}
	} else {
		fmt.Fprintf(os.Stderr, "drawLine error\n")
	}
}

func (m Maze) drawFrame(p1, p2 Point) {
	m.drawLine(Point{p1.x, p1.y}, Point{p2.x, p1.y})
	m.drawLine(Point{p2.x, p1.y}, Point{p2.x, p2.y})
	m.drawLine(Point{p2.x, p2.y}, Point{p1.x, p2.y})
	m.drawLine(Point{p1.x, p2.y}, Point{p1.x, p1.y})
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
	if p.x%2 == 1 && p.y%2 == 1 {
		return false
	}
	count := 0
	for _, n := range m.nextTo(p) {
		if m.point(n) {
			count++
		}
	}
	switch {
	case count <= 3:
		return true
	default:
		return false
	}
}

func (m Maze) randomPoint() Point {
	//p := Point{rand.Intn(int(Width/2)) * 2, rand.Intn(int(Height/2)) * 2}
	//return Point{rand.Intn(Width), rand.Intn(Height)}
	return Point{rand.Intn(int(Width/2)) * 2, rand.Intn(int(Height/2)) * 2}
}

func (m Maze) plotWallAtRandom() {
	p := m.randomPoint()
	if m.canPlot(p) {
		m.setWall(p)
	}
}

func getPointAtRandom(in []Point) []Point {
	if len(in) == 1 {
		return []Point{in[0]}
	} else {
		i := rand.Intn(len(in))
		newList := []Point{}
		for k, p := range in {
			if k != i {
				newList = append(newList, p)
			}
		}
		return append(getPointAtRandom(newList), in[i])
	}
}

func (m Maze) getWallCandidate(p Point) []Point {
	var result []Point
	var list []Point = []Point{Point{p.x - 1, p.y}, Point{p.x, p.y + 1}, Point{p.x + 1, p.y}, Point{p.x, p.y - 1}}
	for _, p := range list {
		if p.x >= 0 && p.x < Width && p.y >= 0 && p.y < Height && m.point(p) == false {
			result = append(result, p)
		}
	}
	return result
}

func (m Maze) extendWall(p Point) {
	fmt.Printf("p.x=%v\tp.y=%v\n", p.x, p.y)
	for _, wc := range getPointAtRandom(m.getWallCandidate(p)) {
		fmt.Printf("\twc.x=%v\twc.y=%v\n", wc.x, wc.y)
		if m.canPlot(wc) {
			m.setWall(wc)
			m.extendWall(wc)
		}
	}
}

func (m Maze) makeMaze() {
	m.drawFrame(Point{0, 0}, Point{Width - 1, Height - 1})
	//var list [(int(Width/2) + int(Height/2)) * 2 - 4 ]Point
	var list []Point
	for i := 2; i < Width-1; i += 2 {
		list = append(list, Point{i, 0})
		list = append(list, Point{i, Height - 1})
	}
	for i := 2; i < Height-1; i += 2 {
		list = append(list, Point{0, i})
		list = append(list, Point{Width - 1, i})
	}
	//fmt.Printf("len(list)=%v\n", len(list))
	for _, p := range getPointAtRandom(list) {
		m.extendWall(p)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	m := NewMaze()
	/*
		m.drawFrame(m.randomPoint(), m.randomPoint())
		for i := 1; i < 10; i++ {
			m.plotWallAtRandom()
		}
	*/
	m.makeMaze()
	m.print()
}
