package conway

// Any live cell with fewer than two live neighbors dies, as if by underpopulation. -> false
// Any live cell with more than three live neighbors dies, as if by overpopulation. -> false

// Any live cell with two or three live neighbors lives on to the next generation. -> true
// Any dead cell with exactly three live neighbors becomes a live cell, as if by reproduction. -> true

import (
	"fmt"
	"math/rand"
	"sync"
)

type Pair struct {
	x int
	y int
}

type System struct {
	mu    sync.Mutex
	Alive []Pair
	Size  int
}

func (s *System) contains(pair Pair) bool {
	for _, p := range s.Alive {
		if p.x == pair.x && p.y == pair.y {
			return true
		}
	}

	return false
}

func (s *System) append(pair Pair) {
	s.Alive = append(s.Alive, pair)
}

func Wrap(val int, min int, max int) int {
	if val < min {
		return max
	}

	if val > max {
		return min
	}

	return val
}

type Neighbors struct {
	mu    sync.Mutex
	count int
}

func CountNeighbors(p Pair, system *System) int {
	xRange := [3]int{
		p.x,
		Wrap(p.x-1, 0, system.Size-1),
		Wrap(p.x+1, 0, system.Size-1),
	}

	yRange := [3]int{
		p.y,
		Wrap(p.y-1, 0, system.Size-1),
		Wrap(p.y+1, 0, system.Size-1),
	}

	neighbors := new(Neighbors)
	neighbors.count = 0

	var wg sync.WaitGroup
	wg.Add(9)
	for _, j := range yRange {
		for _, i := range xRange {
			go func(x int, y int) {
				defer wg.Done()
				if p.x == x && p.y == y {
					return
				}
				if system.contains(Pair{x: x, y: y}) {
					neighbors.mu.Lock()
					neighbors.count++
					neighbors.mu.Unlock()
				}
			}(i, j)
		}
	}

	wg.Wait()
	return neighbors.count
}

func CheckCell(cell bool, neighbors int) bool {
	birth := (neighbors == 3)
	persistAtTwo := cell && (neighbors == 2)

	return birth || persistAtTwo
}

func IterateSystem(system *System) *System {
	var wg sync.WaitGroup
	newSystem := new(System)
	newSystem.Size = system.Size
	newSystem.Alive = make([]Pair, 0)

	wg.Add(system.Size * system.Size)
	for j := 0; j < system.Size; j++ {
		for i := 0; i < system.Size; i++ {
			go func(x int, y int) {
				defer wg.Done()
				localPair := Pair{x: x, y: y}
				count := CountNeighbors(localPair, system)
				isAlive := CheckCell(system.contains(localPair), count)
				if isAlive {
					newSystem.mu.Lock()
					newSystem.append(localPair)
					newSystem.mu.Unlock()
				}
			}(i, j)
		}
	}

	wg.Wait()
	return newSystem
}

/*
func RenderSystem2(system *System) string {
	boxSize := 15
	output := fmt.Sprintf("<div style='position: relative; width: %vpx; height: %vpx'>", boxSize*system.Size, boxSize*system.Size)
	defer func() { output = fmt.Sprintf("%s</div>", output) }()

	for _, pair := range system.Alive {
		cell := fmt.Sprintf("<i id='b' style='top: %vpx; left: %vpx'></i>", pair.x*boxSize, pair.y*boxSize)
		output = fmt.Sprintf("%s%s", output, cell)
	}

	return output
}
*/

func RenderSystem(system *System) string {
	const blackBox string = "<i id='b'></i>"
	const whiteBox string = "<i id='w'></i>"

	var output string
	for y := 0; y < system.Size; y++ {
		for x := 0; x < system.Size; x++ {
			localPair := Pair{x: x, y: y}
			isAlive := system.contains(localPair)
			if isAlive {
				output = fmt.Sprintf("%s%s", output, blackBox)
			} else {
				output = fmt.Sprintf("%s%s", output, whiteBox)
			}
		}
		output = fmt.Sprintf("%s<br>", output)
	}

	return output
}

func RandomSystem(size int) *System {
	newSystem := *new(System)
	newSystem.Size = size

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			isAlive := RandBool()
			if isAlive {
				newSystem.append(Pair{x: x, y: y})
			}
		}
	}

	return &newSystem
}

func RunSystem(s *System) (string, *System) {
	if len(s.Alive) == 0 {
		/*
			s.append(Pair{x: 3, y: 1})
			s.append(Pair{x: 1, y: 2})
			s.append(Pair{x: 3, y: 2})
			s.append(Pair{x: 2, y: 3})
			s.append(Pair{x: 3, y: 3})
			return RenderSystem(s), s
		*/
		newSystem := RandomSystem(s.Size)
		return RenderSystem(newSystem), newSystem
	}

	newSystem := IterateSystem(s)
	return RenderSystem(newSystem), newSystem
}

func RandBool() bool {
	return rand.Intn(2) == 1
}
