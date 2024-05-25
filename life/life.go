package life

import (
	"fmt"
	"math/rand"
	"slices"
	"sync"

	"conway-http/life/vector"
)

type Rules struct {
	Birth   []int
	Survive []int
}

type System struct {
	Alive []vector.Pair
	Rules Rules
	Size  int
}

func (s *System) contains(v vector.Pair) bool {
	for _, p := range s.Alive {
		if p == v {
			return true
		}
	}
	return false
}

func (s *System) append(ps ...vector.Pair) {
	s.Alive = append(s.Alive, ps...)
}

func (s *System) GetSize() int {
	return s.Size
}

func (s *System) GetAlive() []vector.Pair {
	return s.Alive
}

func (s *System) CountAlive(ps []vector.Pair) int {
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

func (s *System) Next() System {
	next := System{
		Rules: s.Rules,
		Size:  s.Size,
	}

	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, a := range s.Alive {
		ns := a.Neighbors(s.Size)
		for _, p := range append(ns, a) {
			mu.Lock()
			checked := next.contains(p)
			mu.Unlock()
			if checked {
				continue
			}

			wg.Add(1)
			go func(p vector.Pair) {
				c := s.CountAlive(p.Neighbors(s.Size))
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
			isAlive := s.contains(vector.Pair{X: x, Y: y})
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
				s.append(vector.Pair{X: x, Y: y})
			}
		}
	}

	return s
}

func (s System) Inject(ps ...vector.Pair) System {
	s.Alive = ps
	return s
}
