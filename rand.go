package xkcd

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/inlined/rand"
)

// fakeRand implements rand.Rand with a fixed sequence of numbers.
// If called more times than the seeded sequence is long, the source
// sequence will repeat.
type fakeRand struct {
	seq []uint64
	pos int
}

func (r *fakeRand) ExpFloat64() float64 {
	return r.Float64()
}

func (r *fakeRand) Float32() float32 {
	return float32(r.Float64())
}

func (r *fakeRand) Float64() float64 {
	f := float64(r.Int63()) / (1 << 63)
	return f
}

func (r *fakeRand) Int() int {
	return int(r.Uint64())
}

func (r *fakeRand) Int31() int32 {
	return int32(r.Uint64())
}

func (r *fakeRand) Int31n(n int32) int32 {
	return int32(r.Uint64()) % n
}

func (r *fakeRand) Int63() int64 {
	return int64(r.Uint64())
}

func (r *fakeRand) Int63n(n int64) int64 {
	return int64(r.Uint64()) % n
}

func (r *fakeRand) Intn(n int) int {
	return r.Int() % n
}

func (r *fakeRand) NormFloat64() float64 {
	return r.Float64()
}

func (r *fakeRand) Perm(n int) []int {
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = int(r.Int31n(int32(n)))
	}
	return a
}

// Read is not implemented because it is likely to be ambiguous regarding
// users' expectations (e.g. if I push an int64, is that 8 bytes or one?)
func (r *fakeRand) Read(p []byte) (n int, err error) {
	return 0, errors.New("xkcd.Rand.Read not implemented")
}

// Seed is unimplemented because this is a fake that is explicitly seeded.
func (r *fakeRand) Seed(seed int64) { /* ignored */ }

// Shuffle is not implemented because it is very unlikely that users realize
// how many pieces from the "tape" of fake numbers this will consume (the
// built-in implementation is the Fisher-Yates shuffle)
func (r *fakeRand) Shuffle(n int, swap func(i, j int)) {
	/* ignored */
}

func (r *fakeRand) Uint32() uint32 {
	return uint32(r.Uint64())
}

func (r *fakeRand) Uint64() uint64 {
	i := r.seq[r.pos]
	r.pos = (r.pos + 1) % len(r.seq)
	return i
}

// Rand returns an `inlined/rand.Rand` implementation that will return
// the provided pattern cyclically. Sequence assumes that random values
// will be provided in order of copmatible types for which they will be used.
// Rand will automatically detect int, int32, int64, float32, and float64.
// Note: Since the `math/rand` library does not accept negative numbers, neither
// does Rand(). Seeding a negative number will cause a panic.
func Rand(seq ...interface{}) rand.Rand {
	transformed := make([]uint64, len(seq))
	for n, i := range seq {
		transformed[n] = toU64(i)
	}
	return &fakeRand{seq: transformed}
}

func toU64(i interface{}) uint64 {
	switch x := i.(type) {
	// For these types we need to
	case int8:
		return uint64(x)
	case uint8:
		return uint64(x)
	case int16:
		return uint64(x)
	case uint16:
		return uint64(x)
	case int:
		return uint64(x)
	case uint:
		return uint64(x)
	case int32:
		return uint64(x)
	case uint32:
		return uint64(x)
	case int64:
		return uint64(x)
	case uint64:
		return x
	case float32:
		return floatBytes(float64(x))
	case float64:
		return floatBytes(x)
	default:
		panic(fmt.Sprintf("xkcdran.Seed(%v): cannot use seed of type %s", x, reflect.TypeOf(x).Kind()))
	}
}

// rand.Rand.Float64 creates a float between [0,1) based on the proportion
// an int63 with max int63
func floatBytes(f float64) uint64 {
	if f >= 1 {
		panic(fmt.Sprintf("xkcd.Rand(%f): cannot use floats >= 1.0", f))
	}
	f = f * (1 << 63)
	return uint64(f)
}
