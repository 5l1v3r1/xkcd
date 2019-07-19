package xkcdrand_test

import (
	"math/rand"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/inlined/xkcdrand"
)

func TestSequences(t *testing.T) {
	for _, test := range []struct {
		tag    string
		input  []int64
		output []int64
	}{
		{
			tag:   "one number",
			input: []int64{1},
		}, {
			tag:   "multiple numbers",
			input: []int64{1, 2, 3},
		}, {
			tag:    "wrapping sequence",
			input:  []int64{1, 2, 3},
			output: []int64{1, 2, 3, 1, 2, 3, 1},
		},
	} {
		t.Run(test.tag, func(t *testing.T) {
			want := test.output
			if want == nil {
				want = test.input
			}

			stream := xkcdrand.Sequence(test.input...)
			rng := rand.New(stream)
			got := make([]int64, len(want))
			for n := 0; n < len(want); n++ {
				got[n] = rng.Int63()
			}
			diff := cmp.Diff(got, want)
			if diff != "" {
				t.Errorf("xckdrand.Sequence(%v): returned unexpected sequence. got=%v; want=%v; diff=%v", test.input, got, want, diff)
			}
		})
	}
}
func TestUint64(t *testing.T) {
	for _, test := range []struct {
		tag string
		num uint64
	}{
		{
			tag: "big numbers",
			num: 0xFEEDFACEBAADF00D,
		},
		{
			tag: "small numbers",
			num: 0xDEADC0DE,
		},
	} {
		t.Run(test.tag, func(t *testing.T) {
			encoded := xkcdrand.Uint64(test.num)
			stream := xkcdrand.Sequence(encoded)
			rng := rand.New(stream)
			got := rng.Uint64()
			if got != test.num {
				t.Errorf("xckdrand.Uint64(%x): failed to decode; got=%x", test.num, got)
			}
		})
	}
}

func TestFloat64(t *testing.T) {
	for _, test := range []struct {
		tag string
		num float64
	}{
		{
			tag: "basic",
			num: 0.5,
		},
	} {
		t.Run(test.tag, func(t *testing.T) {
			encoded := xkcdrand.Float64(test.num)
			stream := xkcdrand.Sequence(encoded)
			rng := rand.New(stream)
			got64 := rng.Float64()
			if got64 != test.num {
				t.Errorf("xckdrand.Float64(%f): failed to decode; got=%f", test.num, got64)
			}
			got32 := rng.Float32()
			if got32 != float32(test.num) {
				t.Errorf("xckdrand.Float32(%f): failed to decode; got=%f", test.num, got32)
			}
		})
	}

}
