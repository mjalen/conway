package life

import (
	"math/rand"
	"slices"
	"sync"

	"github.com/mjalen/conway/life/vector"
)

type Rules struct {
	Birth   []int
	Survive []int
}

type System struct {
	Alive []vector.Pair
	Rules Rules
	Size  int
	Seed int64
}

func (s *System) Contains(v vector.Pair) bool {
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
		if s.Contains(p) {
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
			checked := next.Contains(p)
			mu.Unlock()
			if checked {
				continue
			}

			wg.Add(1)
			go func(p vector.Pair) {
				c := s.CountAlive(p.Neighbors(s.Size))
				isAlive := next.Rules.Check(s.Contains(p), c)
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

func (s System) Random() System {
	if s.Seed == 0 {
		s.Seed = rand.Int63n(999999999)
	}
	src := rand.New(rand.NewSource(s.Seed))
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
