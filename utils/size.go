package utils

import (
	"fmt"
	"strconv"
	"strings"
)

type ByteSize uint64

const (
	B  = 1
	KB = B << 10
	MB = KB << 10
	GB = MB << 10
	TB = GB << 10
	PB = TB << 10
	EB = PB << 10
)

var unit bool = false
var format int = 1

func (b ByteSize) String() string {
	switch {
	case b%EB == 0:
		return fmt.Sprintf("%d E%sB", b/EB, b.Unit())
	case b%PB == 0:
		return fmt.Sprintf("%d P%sB", b/PB, b.Unit())
	case b%TB == 0:
		return fmt.Sprintf("%d T%sB", b/TB, b.Unit())
	case b%GB == 0:
		return fmt.Sprintf("%d G%sB", b/GB, b.Unit())
	case b%MB == 0:
		return fmt.Sprintf("%d M%sB", b/MB, b.Unit())
	case b%KB == 0:
		return fmt.Sprintf("%d K%sB", b/KB, b.Unit())
	default:
		return fmt.Sprintf("%d B", b)
	}
}

func (b ByteSize) StringFormat() string {
	floatValue := float64(b)
	_formatArray := []string{
		"%.", b.Format(), "f",
	}
	switch {
	case b < KB:
		return fmt.Sprintf("%d", b)
	case b < MB:
		_formatArray = append(_formatArray, " K", b.Unit(), "B")
		_format := strings.Join(_formatArray, "")
		return fmt.Sprintf(_format, floatValue/KB)
	case b < GB:
		_formatArray = append(_formatArray, " M", b.Unit(), "B")
		_format := strings.Join(_formatArray, "")
		return fmt.Sprintf(_format, floatValue/MB)
	case b < TB:
		_formatArray = append(_formatArray, " G", b.Unit(), "B")
		_format := strings.Join(_formatArray, "")
		return fmt.Sprintf(_format, floatValue/GB)
	case b < PB:
		_formatArray = append(_formatArray, " T", b.Unit(), "B")
		_format := strings.Join(_formatArray, "")
		return fmt.Sprintf(_format, floatValue/TB)
	case b < EB:
		_formatArray = append(_formatArray, " P", b.Unit(), "B")
		_format := strings.Join(_formatArray, "")
		return fmt.Sprintf(_format, floatValue/PB)
	default:
		_formatArray = append(_formatArray, " E", b.Unit(), "B")
		_format := strings.Join(_formatArray, "")
		return fmt.Sprintf(_format, floatValue/EB)
	}
}

func (b ByteSize) Setunit(_unit bool) ByteSize {
	unit = _unit
	return b
}

func (b ByteSize) Unit() string {
	if unit {
		return "i"
	} else {
		return ""
	}
}

func (b ByteSize) SetFormat(i int) ByteSize {
	format = i
	return b
}
func (b ByteSize) Format() string {
	return strconv.Itoa(format)
}
