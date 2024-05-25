package vector

import (
	"testing"
	"slices"
)

func TestWrap(t *testing.T) {
	cases := []struct {
		pair Pair
		min Pair
		max Pair
		want Pair 
	}{
		{ pair: Pair{0, 0}, min: Pair{0, 0}, max: Pair{2,2}, want: Pair{0,0} },
		{ pair: Pair{2, 0}, min: Pair{0, 0}, max: Pair{2,2}, want: Pair{0,0} },
		{ pair: Pair{0, 2}, min: Pair{0, 0}, max: Pair{2,2}, want: Pair{0,0} },
		{ pair: Pair{-1, 0}, min: Pair{0, 0}, max: Pair{2,2}, want: Pair{1,0} },
		{ pair: Pair{0, -1}, min: Pair{0, 0}, max: Pair{2,2}, want: Pair{0,1} },
		{ pair: Pair{-1, -1}, min: Pair{0, 0}, max: Pair{2,3}, want: Pair{1, 2} },
		{ pair: Pair{-1, -1}, min: Pair{-1, 0}, max: Pair{2,3}, want: Pair{-1,2} },
		{ pair: Pair{5, 0}, min: Pair{0, 1}, max: Pair{2,3}, want: Pair{0,2} },
	}

	for _, c := range cases {
		got := c.pair.Wrap(c.min, c.max)

		if got != c.want {
			t.Errorf("%+v.Wrap(%+v, %+v) == %+v, want %+v.", c.pair, c.min, c.max, got, c.want)
		}
	}
}

func TestNeighbors(t *testing.T) {
	cmpF := func(p, r Pair) int {
		if p.X == r.X {
			return (p.Y - r.Y) 
		}

		return (p.X - r.X) 
	}


	cases := []struct {
		pair Pair
		size int
		want []Pair
	}{
		{
			pair: Pair{1,1}, 
			size: 5,	
			want: []Pair{
				{0,0},
				{0,1},
				{0,2},
				{1,0},
				{1,2},
				{2,0},
				{2,1},
				{2,2},
			},
		},
		{
			pair: Pair{0,0},
			size: 5,
			want: []Pair{
				{4,4},
				{4,0},
				{4,1},
				{0,4},
				{1,4},
				{0,1},
				{1,1},
				{1,0},
			},
		},
	}

	for _, c := range cases {
		got := c.pair.Neighbors(c.size)
		slices.SortFunc(got, cmpF)
		slices.SortFunc(c.want, cmpF)

		if !slices.Equal(got, c.want) {
			t.Errorf("%+v.Neighbors(%v) == %+v, want %+v.", c.pair, c.size, got, c.want)
		}
	}
}
