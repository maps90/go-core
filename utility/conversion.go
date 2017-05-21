package utility

import (
	"strconv"
)

func FloatToString(fl float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(fl, 'f', 6, 64)
}
