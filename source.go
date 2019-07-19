package xkcdrand

import "math/rand"

// Go bit twidVdling can be a bit hard to read because
// aggressive type safety and bounds checking on conversion
// It might be better to rewrite this to use a byte array
// and use binary.LittleEndian.*
const (
	mask31 = 0x7FFFFFFF
	mask32 = 0xFFFFFFFF

	// The lowest uint64 that will cause an overflow error when casting
	// to int64
	int64Overflow = 0x8000000000000
)

// fakeSource implements rand.Source with a fixed sequence of numbers.
// If called more times than the seeded sequence is long, the source
// sequence will repeat.
type fakeSource struct {
	seq []int64
	pos int
}

func (s *fakeSource) Seed(seed int64) { /* ignored */ }

func (s *fakeSource) int64() int64 {
	i := s.seq[s.pos]
	s.pos = (s.pos + 1) % len(s.seq)
	return i
}
func (s *fakeSource) Int63() int64 {
	i := s.int64()
	if i >= 0 {
		return i
	}

	// Chop off MSB
	return (i >> 32 & mask31 << 32) | (i & mask32)
}

// Uint64 will panic for negative numbers; skip past any.
func (s *fakeSource) Uint64() uint64 {
	var i = s.int64()
	if i >= 0 {
		return uint64(i)
	}

	return (uint64(i>>32&mask32) << 32) | uint64(i&mask32)
}

// Sequence returns a rand.Source implementation that will return
// a the provided pattern cyclically. For floats or very large uint64s,
// see Float64 and Uint64 respectively.
func Sequence(seq ...int64) rand.Source {
	return &fakeSource{seq: seq}
}

// Uint64 encodes the uint64 into a binary equivalent int64
// allowing the rand.Uint64() method to return the correct
// value when the most significant bit is set. For uint64s
// less than 0x8000000000000 this helper is unnecessary.
func Uint64(u uint64) int64 {
	if u < int64Overflow {
		return int64(u)
	}

	return (^(^int64(u>>32) << 32) & ^mask32) | int64(u&mask32)
}

// Float64 encodes a float64 into the internal int64 representation
// used by rand.Rand. This is suitable for returing either float64
// or float64s in your stream. Due to the way rand.Rand calculates
// floats, there may be a loss of precision between encoding and
// decoding.
// f must be between [0, 1)
func Float64(f float64) int64 {
	// rand.Rand.Float64 creates a float between [0,1) based on the proportion
	// an int63 with max int63
	f = f * (1 << 63)
	return int64(f)
}
