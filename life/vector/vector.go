package vector

type Ecosystem interface {
	GetSize() int
	GetAlive() []Pair
}

type Pair struct {
	X int
	Y int
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

func (p *Pair) Neighbors(size int) []Pair { // for now, this INCLUDES p in the neighborhood.
	var n []Pair

	for y := (p.Y - 1); y <= (p.Y + 1); y++ {
		for x := (p.X - 1); x <= (p.X + 1); x++ {
			if p.X == x && p.Y == y {
				continue
			}
			wrap := (&Pair{x, y}).Wrap(Pair{0, 0}, Pair{size, size})
			n = append(n, wrap)
		}
	}

	return n
}

