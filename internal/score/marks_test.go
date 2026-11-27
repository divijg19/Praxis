package score

import "testing"

func TestMarksMonotone(t *testing.T) {
	cases := []struct {
		keys, par, hints, want int
	}{
		{5, 5, 0, 3},
		{4, 5, 1, 3},
		{6, 5, 0, 2},
		{7, 5, 0, 2},
		{8, 5, 0, 1},
		{5, 5, 4, 1},
		{6, 5, 4, 1},
		{0, 0, 0, 1},
	}
	for _, c := range cases {
		if got := Marks(c.keys, c.par, c.hints); got != c.want {
			t.Errorf("Marks(%d,%d,%d) = %d, want %d", c.keys, c.par, c.hints, got, c.want)
		}
	}
}

func TestMarksNeverGatesCompletion(t *testing.T) {
	for keys := 0; keys < 50; keys++ {
		m := Marks(keys, 5, 0)
		if m < 1 || m > 3 {
			t.Fatalf("Marks(%d,5,0) = %d, want 1..3", keys, m)
		}
	}
}

func TestMarksMonotoneProperty(t *testing.T) {
	for par := 1; par <= 20; par++ {
		prev := 4
		for keys := 1; keys <= par+20; keys++ {
			m := Marks(keys, par, 0)
			if m > prev {
				t.Fatalf("marks increased with more keys: par=%d keys=%d %d > %d", par, keys, m, prev)
			}
			prev = m
		}
	}
	for par := 1; par <= 20; par++ {
		prev := 4
		for hints := 0; hints <= 10; hints++ {
			m := Marks(par, par, hints)
			if m > prev {
				t.Fatalf("marks increased with more hints: par=%d hints=%d", par, hints)
			}
			prev = m
		}
	}
}
