package main

import (
	"fmt";
	"errors";
	"time";
	"math/rand";
	"flag";
	"github.com/inancgumus/screen"
)

// Any live cell with fewer than two live neighbors dies, as if by underpopulation. -> false
// Any live cell with more than three live neighbors dies, as if by overpopulation. -> false

// Any live cell with two or three live neighbors lives on to the next generation. -> true
// Any dead cell with exactly three live neighbors becomes a live cell, as if by reproduction. -> true

func CountNeighborhood(neighborhood []bool) int32 {
	var count int32
	count = 0
	for _,val := range neighborhood {
		if val {
			count++
		}
	}
	return count
}

func clamp(val int32, minimum int32, maximum int32) int32 {
	if val < minimum {
		return minimum
	}

	if val > maximum {
		return maximum
	}

	return val
}

func GetCoord(x int, y int, system [][]bool) (bool, error) {
	x_out_bound := (x < 0 || x >= len(system[0])) 
	y_out_bound := (y < 0 || y >= len(system))
	if x_out_bound || y_out_bound {
		return false, errors.New("Coordinate out of bounds")	
	}

	return system[y][x], nil
}

func ConstructNeighborhood(x int, y int, system [][]bool) []bool {
	coords := [][]int{ 
		{ x - 1, y - 1 }, 
		{ x - 1, y },
		{ x - 1, y + 1 },
		{ x, y - 1 },
		{ x, y + 1 },
		{ x + 1, y - 1 },
		{ x + 1, y },
		{ x + 1, y + 1 },
	}

	neighbors := make([]bool, 0)
	for _,c := range coords {
		val, err := GetCoord(c[0], c[1], system)
		if err != nil {
			continue
		}
		neighbors =	append(neighbors, val) 
	}

	return neighbors
}

func CheckCell(cell bool, neighborhood []bool) bool {
	count := CountNeighborhood(neighborhood)

	deadAndTrio := (!cell && count == 3)
	liveAnd2or3 := (cell && (count == 2 || count == 3))

	return (deadAndTrio || liveAnd2or3)	
}

func UpdateSystem(system [][]bool) [][]bool {
	buffer := make([][]bool, len(system))
	for i := range buffer {
		buffer[i] = make([]bool, len(system[i]))
	}
	copy(buffer, system)

	for y := 0; y < len(system); y++ {
		for x := 0; x < len(system[y]); x++ {
			neighborhood := ConstructNeighborhood(x, y, system)
			cell := system[y][x]
			buffer[y][x] = CheckCell(cell, neighborhood)
		}
	}
	return buffer
}

func PrintConway(system [][]bool) {
	for _, arr := range system {
		for _, val := range arr {
			if val {
				fmt.Print(" # ")
			} else {
				fmt.Print("   ")
			}
		}
		fmt.Println("")
	}
}

func RandBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
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

func main() {
	var sFlag int
	var lFlag int
	flag.IntVar(&sFlag, "speed", 1, "Speed of the game.")
	flag.IntVar(&lFlag, "length", 32, "Length of the system.")
	flag.Parse()

	system := RandomSystem(lFlag)

	for {
		screen.Clear()
		PrintConway(system)
		system = UpdateSystem(system)
		time.Sleep( time.Second / time.Duration(sFlag) )
	}
}
