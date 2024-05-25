package life

// Any live cell with fewer than two live neighbors dies, as if by underpopulation. -> false
// Any live cell with more than three live neighbors dies, as if by overpopulation. -> false

// Any live cell with two or three live neighbors lives on to the next generation. -> true
// Any dead cell with exactly three live neighbors becomes a live cell, as if by reproduction. -> true

import (
	"fmt"
	"math/rand"
	"slices"
	"sync"
)

type Pair struct {
	X int
	Y int
}

type Rules struct {
	Birth   []int
	Survive []int
}

type System struct {
	Alive []Pair
	Rules Rules
	Size  int
}

func (s *System) contains(pair Pair) bool {
	for _, p := range s.Alive {
		if p == pair {
			return true
		}
	}

	return false
}

func (s *System) append(ps ...Pair) {
	s.Alive = append(s.Alive, ps...)
}

func Wrap(val int, min int, max int) int {
	if val < min {
		return max - 1
	}

	if val >= max {
		return min
	}

	return val
}

func (p *Pair) Wrap(min Pair, max Pair) Pair {
	return Pair{
		X: Wrap(p.X, min.X, max.X),
		Y: Wrap(p.Y, min.Y, max.Y),
	}
}

func (p *Pair) Neighbors(s *System) []Pair { // for now, this INCLUDES p in the neighborhood.
	var n []Pair

	for y := (p.Y - 1); y <= (p.Y + 1); y++ {
		for x := (p.X - 1); x <= (p.X + 1); x++ {
			if p.X == x && p.Y == y {
				continue
			}
			wrap := (&Pair{x, y}).Wrap(Pair{0, 0}, Pair{s.Size, s.Size})
			n = append(n, wrap)
		}
	}

	return n
}

func (s *System) CountAlive(ps []Pair) int {
	c := 0
	for _, p := range ps {
		if s.contains(p) {
			c++
		}
	}
	return c
}

func (r *Rules) Check(alive bool, neighbors int) bool {
	if alive {
		return slices.Contains(r.Survive, neighbors)
	}

	return slices.Contains(r.Birth, neighbors)
}

func (s *System) Next() System { // perhaps add error handling for empty systems
	next := System{
		Rules: s.Rules,
		Size:  s.Size,
	}

	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, a := range s.Alive {
		ns := a.Neighbors(s)
		for _, p := range append(ns, a) {
			mu.Lock()
			checked := next.contains(p)
			mu.Unlock()
			if checked {
				continue
			}

			wg.Add(1)
			go func(p Pair) {
				c := s.CountAlive(p.Neighbors(s))
				isAlive := next.Rules.Check(s.contains(p), c)
				if isAlive {
					mu.Lock()
					next.append(p)
					mu.Unlock()
				}
				wg.Done()
			}(p)
		}
	}

	/*
		for y := 0; y < s.Size; y++ {
			for x := 0; x < s.Size; x++ {
				wg.Add(1)
				go func(p Pair) {
					c := s.CountAlive(p.Neighbors(s))
					isAlive := next.Rules.Check(s.contains(p), c)
					if isAlive {
						mu.Lock()
						next.append(p)
						mu.Unlock()
					}
					wg.Done()
				}(Pair{x,y})
			}
		}
	*/

	wg.Wait()
	return next
}

func (s *System) ToHTML() string {
	const blackBox string = "<div id='b'></div>"
	const whiteBox string = "<div id='w'></div>"

	var output string
	for y := 0; y < s.Size; y++ {
		output = fmt.Sprintf("%s<div id='r'>", output)
		for x := 0; x < s.Size; x++ {
			p := Pair{x, y}
			isAlive := s.contains(p)
			if isAlive {
				output = fmt.Sprintf("%s%s", output, blackBox)
			} else {
				output = fmt.Sprintf("%s%s", output, whiteBox)
			}
		}
		output = fmt.Sprintf("%s</div>", output)
	}

	return output
}

func (s System) Random(seed int64) System {
	src := rand.New(rand.NewSource(seed))
	for y := 0; y < s.Size; y++ {
		for x := 0; x < s.Size; x++ {
			if src.Intn(2) == 1 {
				s.append(Pair{x, y})
			}
		}
	}

	return s
}

func (s System) Inject(ps ...Pair) System {
	s.Alive = ps
	return s
}
