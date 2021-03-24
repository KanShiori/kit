package units

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAsSize(t *testing.T) {
	req := require.New(t)

	units := []Size{
		Bit, Byte,
		KiB, MiB, GiB, TiB,
		KB, MB, GB, TB,
		Kb, Mb, Gb,
	}

	for i := 0.0; i < 100; i = i + 0.1 {
		for _, unit := range units {
			size := AsSize(i, unit)
			req.Equal(Size(i*float64(unit)), size, "not same")
		}
	}
}

func TestParseSize(t *testing.T) {
	req := require.New(t)

	for i := 0; i < 100; i = i + 1 {
		for suffix, unit := range suffixValues {
			str := fmt.Sprintf("%d%s", i, suffix)
			size, err := ParseSize(str)
			req.Nil(err, "parse failed")

			req.Equal(Size(i*int(unit)), size, "parse %s failed", str)
		}
	}
}
