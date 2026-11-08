package rng

import "testing"

func TestSeedDeterministic(t *testing.T) {
	a := New(Seed([]byte("abc")))
	b := New(Seed([]byte("abc")))
	for i := 0; i < 100; i++ {
		if a.Uint64() != b.Uint64() {
			t.Fatal("sequence is not deterministic")
		}
	}
}

func TestSeedDiffers(t *testing.T) {
	a := New(Seed([]byte("abc")))
	b := New(Seed([]byte("abd")))
	same := true
	for i := 0; i < 16; i++ {
		if a.Uint64() != b.Uint64() {
			same = false
			break
		}
	}
	if same {
		t.Fatal("different seeds produced the same sequence")
	}
}

func TestFloat64Range(t *testing.T) {
	r := New(Seed([]byte("x")))
	for i := 0; i < 1000; i++ {
		f := r.Float64()
		if f < 0 || f >= 1 {
			t.Fatalf("Float64 out of [0,1): %v", f)
		}
	}
}

func TestIntnRange(t *testing.T) {
	r := New(Seed([]byte("x")))
	for i := 0; i < 1000; i++ {
		n := r.Intn(7)
		if n < 0 || n >= 7 {
			t.Fatalf("Intn out of range: %d", n)
		}
	}
}

func TestPerm(t *testing.T) {
	r := New(Seed([]byte("x")))
	p := r.Perm(10)
	seen := map[int]bool{}
	for _, v := range p {
		if v < 0 || v >= 10 || seen[v] {
			t.Fatalf("not a permutation: %v", p)
		}
		seen[v] = true
	}
}
