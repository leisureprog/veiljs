package rng

type RNG struct {
	state uint64
}

func New(seed uint64) *RNG {
	r := &RNG{state: seed}
	r.Uint64()
	return r
}

func (r *RNG) Uint64() uint64 {
	r.state += 0x9E3779B97F4A7C15
	z := r.state
	z = (z ^ (z >> 30)) * 0xBF58476D1CE4E5B9
	z = (z ^ (z >> 27)) * 0x94D049BB133111EB
	return z ^ (z >> 31)
}

func (r *RNG) Intn(n int) int {
	if n <= 1 {
		return 0
	}
	return int(r.Uint64() % uint64(n))
}

func (r *RNG) Float64() float64 {
	return float64(r.Uint64()>>11) / float64(uint64(1)<<53)
}

func (r *RNG) Perm(n int) []int {
	p := make([]int, n)
	for i := range p {
		p[i] = i
	}
	for i := n - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		p[i], p[j] = p[j], p[i]
	}
	return p
}

func Seed(parts ...[]byte) uint64 {
	h := uint64(14695981039346656037)
	for _, p := range parts {
		for _, b := range p {
			h ^= uint64(b)
			h *= 1099511628211
		}
	}
	return h
}
