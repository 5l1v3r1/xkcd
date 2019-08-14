package xkcd_test

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/inlined/xkcd"
)

func TestSequences(t *testing.T) {
	for _, test := range []struct {
		tag    string
		input  []interface{}
		output []interface{}
	}{
		{
			tag:   "one number",
			input: []interface{}{1},
		}, {
			tag:   "multiple numbers",
			input: []interface{}{1, 2, 3},
		}, {
			tag:    "wrapping sequence",
			input:  []interface{}{1, 2, 3},
			output: []interface{}{1, 2, 3, 1, 2, 3, 1},
		},
	} {
		t.Run(test.tag, func(t *testing.T) {
			want := test.output
			if want == nil {
				want = test.input
			}

			rng := xkcd.Rand(test.input...)
			got := make([]interface{}, len(want))
			for n := 0; n < len(want); n++ {
				got[n] = rng.Int()
			}
			diff := cmp.Diff(got, want)
			if diff != "" {
				t.Errorf("xckdrand.Sequence(%v): returned unexpected sequence. got=%v; want=%v; diff=%v", test.input, got, want, diff)
			}
		})
	}
}

func TestTypes(t *testing.T) {
	vals := []interface{}{
		float64(0.1),
		float32(0.42),
		float64(math.Pi / 10),
		int(42),
		int32(1337),
		int32(1234),
		int64(0xCAFEDEADBEEF),
		int64(0xBAADF00D),
		int(2),
		float64(0.1027),
		uint32(789),
		uint64(6051985),
	}
	rng := xkcd.Rand(vals...)
	got := []interface{}{
		rng.ExpFloat64(),
		rng.Float32(),
		rng.Float64(),
		rng.Int(),
		rng.Int31(),
		rng.Int31n(2096),
		rng.Int63(),
		rng.Int63n(math.MaxUint32),
		rng.Intn(42),
		rng.NormFloat64(),
		rng.Uint32(),
		rng.Uint64(),
	}
	diff := cmp.Diff(got, vals)
	if diff != "" {
		t.Errorf("xckd.Rand(%v): returned unexpected sequence. got=%v; diff=%v", vals, got, diff)
	}
}

func TestNFuncsWraparound(t *testing.T) {
	input := []interface{}{
		52,
		42,
	}
	want := []interface{}{
		int32(10),
		int64(0),
	}

	rng := xkcd.Rand(input...)
	got := []interface{}{
		rng.Int31n(42),
		rng.Int63n(42),
	}
	diff := cmp.Diff(got, want)
	if diff != "" {
		t.Errorf("xckd.Rand(%v).Int(31|63)n(42): returned unexpected sequence. got=%v; want=%v; diff=%v", input, got, want, diff)
	}
}
