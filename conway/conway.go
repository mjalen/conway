package conway

// Any live cell with fewer than two live neighbors dies, as if by underpopulation. -> false
// Any live cell with more than three live neighbors dies, as if by overpopulation. -> false

// Any live cell with two or three live neighbors lives on to the next generation. -> true
// Any dead cell with exactly three live neighbors becomes a live cell, as if by reproduction. -> true

import (
	"fmt";
	"math/rand";
	"time";
	"log"
)

// the current state of the system depends ONLY on the previous state

type Pair struct {
	x int
	y int
}

type System2 struct {
	Alive []Pair
	Size int 
}

func (s *System2) contains(pair Pair) bool {
	for _,p := range s.Alive {
		if p.x == pair.x && p.y == pair.y {
			return true
		}
	}
	return false
}

func (s *System2) append(pair Pair) {
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



func CountNeighbors(p Pair, system System2) int {
	xMin :=	Wrap(p.x - 1, 0, system.Size - 1)
	xMax := Wrap(p.x + 1, 0, system.Size - 1) 
	yMin := Wrap(p.y - 1, 0, system.Size - 1)
	yMax := Wrap(p.y + 1, 0, system.Size - 1) 

	xRange := [3]int{xMin, p.x, xMax}
	yRange := [3]int{yMin, p.y, yMax}

	count := 0
	for _,y := range yRange {
		for _,x := range xRange {
			if p.x == x && p.y == y {
				continue
			}
			testPair := Pair{x: x, y: y} 
			isAlive := system.contains(testPair)
			if isAlive {
				count = count + 1
			}
		}
	}

	return count
}

func CheckCell(cell bool, neighbors int) bool {
	birth := (neighbors == 3)
	persistAtTwo := cell && (neighbors == 2)

	return birth || persistAtTwo	
}

func IterateSystem(system System2) System2 {
	newSystem := System2{ Size: system.Size, Alive: make([]Pair, 0) } 

	for y := 0; y < system.Size; y++ {
		for x := 0; x < system.Size; x++ {
			localPair := Pair{x: x, y: y}	
			count := CountNeighbors(localPair, system)
			isAlive := CheckCell(system.contains(localPair), count)
			if isAlive {
				newSystem.append(localPair)
			}
		}
	}

	return newSystem
}

func RenderSystem2(system System2) string {
	const blackBox string = "<span id='b'></span>" 
	const whiteBox string = "<span id='w'></span>"

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

func RandomSystem2(size int) System2 {
	newSystem := *new(System2)
	newSystem.Size = size

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			isAlive := RandBool()
			if isAlive {
				newSystem.append(Pair{ x: x, y: y })
			}
		}
	}

	return newSystem
}

func RunSystem2(s System2) (string, System2) {
	if len(s.Alive) == 0 {
		/*
		s.append(Pair{ x: 3, y: 1 })
		s.append(Pair{ x: 1, y: 2 })
		s.append(Pair{ x: 3, y: 2 })
		s.append(Pair{ x: 2, y: 3 })
		s.append(Pair{ x: 3, y: 3 })
		return RenderSystem2(s), s 
		*/
		newSystem := RandomSystem2(s.Size)
		return RenderSystem2(newSystem), newSystem
	}
	
	newSystem := IterateSystem(s)
	return RenderSystem2(newSystem), newSystem 
}

func GetCoord(x int, y int, system [][]bool) bool {
	y = Wrap(y, 0, len(system) - 1)
	x = Wrap(x, 0, len(system[y]) - 1)

	return system[y][x]
}

func CountNeighborhood(x int, y int, system [][]bool) int {
	coords := [][]int{ 
		{ x - 1, y - 1 }, 
		{ x, y - 1 },
		{ x + 1, y - 1},
		{ x - 1, y },
		{ x + 1, y },
		{ x - 1, y + 1 },
		{ x, y + 1 },
		{ x + 1, y + 1 },
	}

	count := 0
	for _, c := range coords {
		val := GetCoord(c[0], c[1], system)
		if val {
			count = count + 1
		}
	}

	return count 
}

func CopySystem(system [][]bool) [][]bool {
	buffer := make([][]bool, len(system))
	for i, _ := range buffer {
		buffer[i] = make([]bool, len(system[i]))
	}
	copy(buffer, system)
	return buffer
}

func UpdateSystem(system [][]bool) [][]bool {
	buffer := CopySystem(system)
	for y := 0; y < len(system); y++ {
		for x := 0; x < len(system[y]); x++ {
			neighbors := CountNeighborhood(x, y, system)
			buffer[y][x] = CheckCell(system[y][x], neighbors)
		}
	}

	return buffer
}

func RenderConway(system [][]bool) string {
	const blackBox string = "<span id='b'></span>" 
	const whiteBox string = "<span id='w'></span>"

	var output string

	for _, arr := range system {
		for _, val := range arr {
			if val {
				output = fmt.Sprintf("%s%s", output, blackBox)
			} else {
				output = fmt.Sprintf("%s%s", output, whiteBox)
			}
		}
		output = fmt.Sprintf("%s<br>", output) 
	}

	return output
}

func RandBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(10) == 1
}

func EmptySystem(length int) [][]bool {
	system := make([][]bool, length)

	for i := range system {
		system[i] = make([]bool, length)
		for j := range system[i] {
			system[i][j] = false 
		}
	}

	return system
}

type Cell struct {
	x int
	y int
	val bool
}

func InjectPattern(length int, pattern []Cell) [][]bool {
	system := EmptySystem(length)

	for _, cell := range pattern {
		log.Println("Injecting pattern to (%v, %v)", cell.x, cell.y)
		system[cell.y][cell.x] = cell.val	
	}
	return system
}


func RandomSystem(length int) [][]bool {
	system := make([][]bool, length)

	for i := range system {
		system[i] = make([]bool, length)
		for j := range system[i] {
			system[i][j] = RandBool()
		}
	}

	return system
}


type System struct {
	Size int
	Speed int
	System [][]bool
}

func (s *System) Run() string {
	if s.System == nil {
		pattern := make([]Cell, 0)
		pattern = append(pattern, Cell{ x: 3, y: 1, val: true })
		pattern = append(pattern, Cell{ x: 1, y: 2, val: true })
		pattern = append(pattern, Cell{ x: 3, y: 2, val: true })
		pattern = append(pattern, Cell{ x: 2, y: 3, val: true })
		pattern = append(pattern, Cell { x: 3, y: 3, val: true })

		s.System = InjectPattern(s.Size, pattern)
		return RenderConway(s.System)
	}

	copy(s.System, UpdateSystem(s.System))
	return RenderConway(s.System)
}


