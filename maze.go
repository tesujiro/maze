package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	WALL   = " "
	ROAD   = "#"
	START  = "S"
	FINISH = "F"
)

type Point struct {
	x int
	y int
}

func (p1 Point) is(p2 Point) bool {
	return p1.x == p2.x && p1.y == p2.y
}

func (p2 Point) opposite(p1 Point) Point {
	return Point{2*p1.x - p2.x, 2*p1.y - p2.y}
}

func getPointAtRandom(in []Point) []Point {
	if len(in) <= 1 {
		return in
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

type Edge struct {
	start    Point
	end      Point
	p        []Point
	distance int
	//next     []*Edge
}

func (e *Edge) print(indent int) {
	fmt.Printf("%v%v->%v (distance:%v)(p:%v)\n", strings.Repeat("  ", indent), e.start, e.end, e.distance, e.p)
	/*
		for _, e := range e.next {
			e.print(indent + 1)
		}
	*/
}

func (e *Edge) hasPoint(p Point) bool {
	for _, ep := range e.p {
		if p.is(ep) {
			return true
		}
	}
	return false
}

func (e *Edge) nextTo(p Point) bool {
	nextTo := len(e.p) < 2 || (!p.is(e.p[1]) && !p.is(e.p[len(e.p)-2]))
	return nextTo
}

type Graph struct {
	edgelist   []Edge
	vertexlist []Point
}

func (g *Graph) print() {
	for _, e := range g.edgelist {
		e.print(0)
	}
}

func (g *Graph) addVertex(p Point) {
	g.vertexlist = append(g.vertexlist, p)
	//fmt.Printf("addVertext :%v\n", p)
}

func (g *Graph) addEdge(e Edge) {
	g.edgelist = append(g.edgelist, e)
}

func (g *Graph) hasVertex(p Point) bool {
	for _, v := range g.vertexlist {
		if v.is(p) {
			return true
		}
	}
	return false
}

func (g *Graph) hasEdge(e Edge) bool {
	for _, e2 := range g.edgelist {
		if e2.start.is(e.start) && e2.end.is(e.end) || e2.start.is(e.end) && e2.end.is(e.start) {
			if e2.hasPoint(e.p[1]) {
				return true
			}
		}
	}
	return false
}

type Maze struct {
	width  int
	height int
	data   [][]bool
	start  Point
	finish Point
}

func NewMaze(w, h int) *Maze {
	m := make([][]bool, w)
	for i := range m {
		m[i] = make([]bool, h)
	}
	return &Maze{
		width:  w,
		height: h,
		data:   m,
	}
}

func (m *Maze) isInside(p Point) bool {
	return p.x >= 0 && p.x < m.width && p.y >= 0 && p.y < m.height
}

func (m *Maze) setRoad(p Point) {
	m.data[p.x][p.y] = true
	m.printPoint(p)
}

func (m *Maze) isRoad(p Point) bool {
	return m.data[p.x][p.y]
}

func (m *Maze) setWall(p Point) {
	m.data[p.x][p.y] = false
	m.printPoint(p)
}

func (m *Maze) point(p Point) bool {
	if p.x < 0 || p.x >= m.width || p.y < 0 || p.y >= m.height {
		return false
	} else {
		return m.data[p.x][p.y]
	}
}

func (m *Maze) printRoad(x, y int) {
	sym := ROAD
	fmt.Printf("\x1b[%v;%vH%."+strconv.Itoa(x%2+1)+"s", y, x+1, strings.Repeat(sym, 3))
}

func (m *Maze) printInit() {
	// Clear Screen
	fmt.Print("\x1b[H\x1b[2J")

}

func (m *Maze) printFinish() {
	fmt.Printf("\x1b[%v;%vH", m.height+3, 1)
}

func (m *Maze) print() {
	// Clear Screen
	fmt.Print("\x1b[H\x1b[2J")

	// Print Outer Wall
	for x := 0; x < (m.width+1)/2*3; x++ {
		m.printRoad(x, 0)
		m.printRoad(x, m.height+2)
	}
	//time.Sleep(1 * time.Second)
	for y := 0; y < m.height+2; y++ {
		m.printRoad(0, y)
		m.printRoad((m.width+1)/2*3-1, y)
	}
	//time.Sleep(1 * time.Second)

	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			m.printPoint(Point{x, m.height - y - 1})
		}
	}
	//time.Sleep(1 * time.Second)
}

func (m *Maze) printPoint(p Point) {
	sym := WALL
	if !m.point(p) {
		sym = ROAD
	}
	fmt.Printf("\x1b[%v;%vH%."+strconv.Itoa(p.x%2+1)+"s", m.height-p.y+1, p.x/2*3+p.x%2+2, strings.Repeat(sym, 3))
}

func (m *Maze) drawLine(p1, p2 Point) {
	X, x, Y, y := p1.x, p2.x, p1.y, p2.y
	if x > X {
		X, x = p2.x, p1.x
	}
	if y > Y {
		Y, y = p2.y, p1.y
	}
	if X == x {
		for i := y; i <= Y; i++ {
			m.setRoad(Point{x, i})
		}
	} else if Y == y {
		for i := x; i <= X; i++ {
			m.setRoad(Point{i, y})
		}
	} else {
		fmt.Fprintf(os.Stderr, "drawLine error\n")
	}
}

func (m *Maze) drawFrame(p1, p2 Point) {
	m.drawLine(Point{p1.x, p1.y}, Point{p2.x, p1.y})
	m.drawLine(Point{p2.x, p1.y}, Point{p2.x, p2.y})
	m.drawLine(Point{p2.x, p2.y}, Point{p1.x, p2.y})
	m.drawLine(Point{p1.x, p2.y}, Point{p1.x, p1.y})
}

func (m *Maze) nextTo(p Point) []Point {
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
		if n.x >= 0 && n.x < m.width && n.y >= 0 && n.y < m.height {
			plist = append(plist, n)
		}
	}
	return plist
}

func (m *Maze) canPlot(p Point) bool {
	if p.x%2 == 1 && p.y%2 == 1 {
		return false
	}
	//fmt.Printf("p.x=%v\tp.y=%v\n", p.x, p.y)
	for _, n := range m.nextTo(p) {
		if m.point(n) {
			if m.point(n.opposite(p)) {
				//fmt.Printf("\topposite false p=%v\tn=%v\topposite(p,n)=%v\n", p, n, opposite(p, n))
				return false
			}
		}
	}
	return true
}

func (m *Maze) randomPoint() Point {
	//p := Point{rand.Intn(int(m.width/2)) * 2, rand.Intn(int(m.height/2)) * 2}
	//return Point{rand.Intn(m.width), rand.Intn(m.height)}
	return Point{rand.Intn(int(m.width/2)) * 2, rand.Intn(int(m.height/2)) * 2}
}

func (m *Maze) getRoadCandidate(p Point) []Point {
	var result []Point
	var list []Point = []Point{Point{p.x - 1, p.y}, Point{p.x, p.y + 1}, Point{p.x + 1, p.y}, Point{p.x, p.y - 1}}
	for _, p := range list {
		//if p.x >= 0 && p.x < m.width && p.y >= 0 && p.y < m.height && m.point(p) == false {
		if m.isInside(p) && !m.isRoad(p) {
			result = append(result, p)
		}
	}
	return result
}

func (m *Maze) extendRoad(p Point) {
	//fmt.Printf("p.x=%v\tp.y=%v\n", p.x, p.y)
	for _, wc := range getPointAtRandom(m.getRoadCandidate(p)) {
		//fmt.Printf("\twc.x=%v\twc.y=%v\n", wc.x, wc.y)
		if m.canPlot(wc) {
			m.setRoad(wc)
			m.extendRoad(wc)
		}
	}
}

func (m *Maze) makeMaze() {
	var list []Point
	list = []Point{m.randomPoint()}
	for _, p := range getPointAtRandom(list) {
		m.extendRoad(p)
	}
}

func (m *Maze) getNextRoad(p Point) []Point {
	var result []Point
	var list []Point = []Point{Point{p.x - 1, p.y}, Point{p.x, p.y + 1}, Point{p.x + 1, p.y}, Point{p.x, p.y - 1}}
	for _, p := range list {
		if m.isInside(p) && m.isRoad(p) {
			result = append(result, p)
		}
	}
	return result
}

func (m *Maze) extendGraph(cpoint Point, cedge Edge, g *Graph) *Graph {

	// Loop
	if g.hasVertex(cpoint) {
		if !g.hasEdge(cedge) {
			//g.edgelist = append(g.edgelist, cedge)
			g.addEdge(cedge)
		}
		return g
	}

	var nextRoads []Point
	for _, road := range m.getNextRoad(cpoint) {
		if cedge.nextTo(road) {
			nextRoads = append(nextRoads, road)
		}
	}

	if len(nextRoads) == 1 { //One way to go
		// Extend Current Edge
		for _, next := range nextRoads {
			cedge = Edge{start: cedge.start, end: next, p: append(cedge.p, next), distance: cedge.distance + 1}
			g = m.extendGraph(next, cedge, g)
		}
	} else {
		// Save Last Edge
		if cedge.distance > 0 {
			g.addVertex(cpoint)
			//g.edgelist = append(g.edgelist, cedge)
			g.addEdge(cedge)
		}

		// New Edge
		for _, next := range nextRoads {
			ne := Edge{start: cpoint, end: next, distance: 1, p: []Point{cpoint, next}} //Next Edge List
			g = m.extendGraph(next, ne, g)
		}
	}
	return g
}

func (m *Maze) makeGraph() *Graph {
	start := Point{x: 0, y: 0}
	e := Edge{start: start, end: start, distance: 0, p: []Point{start}} //Next Edge List
	g := &Graph{}

	g = m.extendGraph(start, e, g)
	g.addVertex(start)
	return g
}

const (
	White = iota
	Gray
	Black
)

//func (g *Graph) dearch(black, gray []Point) {
//}

//func (g *Graph) depthFirstSearch(s Point)
func (g *Graph) dijkstraSearch(s Point) {
	//PQ := []Point{}
	dist := make(map[Point]int)
	HIGH := int(math.Pow(2, 31) - 1)

	for _, v := range g.vertexlist {
		dist[v] = HIGH
	}

	//for _, v := range g.vertexlist {
	//PQ = append(PQ, v)
	//}

	BlackList := []Point{}
	dist[s] = 0
	GrayList := []Point{s}

	for len(GrayList) > 0 {
		// Search Shortest
		shortest := GrayList[0]
		for _, v := range GrayList {
			if dist[v] < dist[shortest] {
				shortest = v
			}
		}

		// Add Shortest V to BlackList
		BlackList = append(BlackList, shortest)
		fmt.Printf("dist[%v] => %v\n", shortest, dist[shortest])

		// Remov Shortest V from GrayList
		newGrayList := []Point{}
		for _, v := range GrayList {
			if !v.is(shortest) {
				newGrayList = append(newGrayList, v)
			}
		}
		GrayList = newGrayList

		//
		for _, e := range g.edgelist {
			if e.start.is(shortest) {
				GrayList = append(GrayList, e.end)
				if dist[e.end] < dist[e.start]+e.distance {
					dist[e.end] = dist[e.start] + e.distance
				}
			}
		}
	}

	//for _, v := range g.vertexlist {
	//fmt.Printf("dist[%v] => %v\n", v, dist[v])
	//}
}

//func (g *Graph) makeRoute(start Point, end Point) {
//for _, e := range g.searchEdge(start) {
//}
//}

func main() {
	var width *int = flag.Int("width", 30, "Width of the maze.")
	var height *int = flag.Int("height", 20, "Height of the maze.")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	m := NewMaze(*width*2+1, *height*2+1)
	m.print() // print Init Maze
	m.makeMaze()
	m.printFinish()
	g := m.makeGraph()
	g.print()
	start := Point{x: 0, y: 0}
	//end := Point{x: m.width - 1, y: m.height - 1}
	//g.makeRoute(start, end)
	g.dijkstraSearch(start)
}
