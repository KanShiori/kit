package units

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	Bit  Size = 1
	Byte Size = 8 * Bit

	KiB Size = 1024 * Byte
	MiB Size = 1024 * KiB
	GiB Size = 1024 * MiB
	TiB Size = 1024 * GiB

	KB Size = 1000 * Byte
	MB Size = 1000 * KB
	GB Size = 1000 * MB
	TB Size = 1000 * GB

	Kb Size = 1000 * Bit
	Mb Size = 1000 * Kb
	Gb Size = 1000 * Mb
)

var (
	suffixValues = map[string]Size{
		"b": Bit, "B": Byte,
		"k": KiB, "K": KiB, "KB": KB, "KiB": KiB,
		"Kb": Kb, "Kb/s": Kb, "Kbps": Kb,
		"m": MiB, "M": MiB, "MB": MB, "MiB": MiB,
		"Mb": Mb, "Mb/s": Mb, "Mbps": Mb,
		"g": GiB, "G": GiB, "GB": GB, "GiB": GiB,
		"Gb": Gb, "Gb/s": Gb, "Gbps": Gb,
		"t": TiB, "T": TiB, "TB": TB, "TiB": TiB,
	}

	// init 函数中初始化
	pairs []*pair

	re           = regexp.MustCompile(`^[[:space:]]*[.[:digit:]]+[bBkKmMgGtTi/s]+[[:space:]]*$`)
	replaceRegex = regexp.MustCompile("[[:space:][:alpha:]/]+")
)

type Size int64

type pair struct {
	suffix string
	val    Size
}

// ParseSize 将 s 解析为 Size, s 为 float 会有精度损失
func ParseSize(s string) (Size, error) {
	if value, exist := suffixValues[s]; exist {
		return Size(value), nil
	}

	if !re.MatchString(s) {
		return 0, fmt.Errorf("bad format 1: %v", s)
	}

	var size Size
	for _, pair := range pairs {
		if strings.HasSuffix(s, pair.suffix) {
			s = replaceRegex.ReplaceAllString(s, "")

			value, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return 0, fmt.Errorf("bad format 2: %s", s)
			}
			size = Size(value) * pair.val
			break
		}
	}

	return size, nil
}

func AsSize(v float64, unit int64) Size {
	if v < 0 {
		return Size(v)
	}
	return Size(v * float64(unit))
}

func MaxSize(ss ...Size) Size {
	if len(ss) == 0 {
		return 0
	}

	max := ss[0]
	for _, s := range ss {
		if s > max {
			max = s
		}
	}
	return max
}

func MinSize(ss ...Size) Size {
	if len(ss) == 0 {
		return 0
	}

	min := ss[0]
	for _, s := range ss {
		if s < min {
			min = s
		}
	}
	return min
}

func (s Size) TiB() float64 {
	return float64(s) / float64(TiB)
}

func (s Size) TB() float64 {
	return float64(s) / float64(TB)
}

func (s Size) GiB() float64 {
	return float64(s) / float64(GiB)
}

func (s Size) GB() float64 {
	return float64(s) / float64(GB)
}

func (s Size) Gb() float64 {
	return float64(s) / float64(Gb)
}

func (s Size) MiB() float64 {
	return float64(s) / float64(MiB)
}

func (s Size) MB() float64 {
	return float64(s) / float64(MB)
}

func (s Size) Mb() float64 {
	return float64(s) / float64(Mb)
}

func (s Size) KiB() float64 {
	return float64(s) / float64(KiB)
}

func (s Size) KB() float64 {
	return float64(s) / float64(KB)
}

func (s Size) Kb() float64 {
	return float64(s) / float64(Kb)
}

func (s Size) Byte() float64 {
	return float64(s) / float64(Byte)
}

func (s Size) Bit() float64 {
	return float64(s)
}

// BinaryHumanSize 解析为 1024 单位的可读字符串
func (s Size) BinaryHumanSize(precision int) string {
	if s < Byte {
		return fmt.Sprintf("%d%s", s, "b")
	}

	unitStr := []string{"B", "KiB", "MiB", "GiB", "TiB"}
	size, unit := getSizeAndUnit(s.Byte(), 1024.0, unitStr)
	return fmt.Sprintf("%.*g%s", precision, size, unit)
}

// DecimalHumanSize 解析为 1000 单位的可读字符串
func (s Size) DecimalHumanSize(precision int) string {
	if s < Byte {
		return fmt.Sprintf("%d%s", s, "b")
	}

	unitStr := []string{"B", "kB", "MB", "GB", "TB"}

	size, unit := getSizeAndUnit(s.Byte(), 1000.0, unitStr)
	return fmt.Sprintf("%.*g%s", precision, size, unit)
}

func (s Size) NetHumanSize(precision int) string {
	unitStr := []string{"b", "Kb", "Mb", "Gb", "Tb"}
	size, unit := getSizeAndUnit(float64(s), 1000.0, unitStr)
	return fmt.Sprintf("%.*g%s", precision, size, unit)
}

func getSizeAndUnit(size float64, base float64, _map []string) (float64, string) {
	i := 0
	unitsLimit := len(_map) - 1
	for size >= base && i < unitsLimit {
		size = size / base
		i++
	}
	return size, _map[i]
}

// func (s Size) HumanSizeWithPrecision(precision int) string {
// 	return ""
// }

func init() {
	pairs = make([]*pair, 0, len(suffixValues))
	for suffix, val := range suffixValues {
		pairs = append(pairs, &pair{
			suffix: suffix,
			val:    val})
	}
	// 排序, 将后缀大的放在前面, 用于 ParseSize 优先匹配
	sort.Slice(pairs, func(i, j int) bool {
		si := pairs[i]
		sj := pairs[j]
		return len(si.suffix) > len(sj.suffix)
	})

}
