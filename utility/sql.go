package utility

import (
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func RequiredNullable(s interface{}) error {
	err := fmt.Errorf("Invalid value")
	switch v := s.(type) {
	case sql.NullString:
		if v.Valid && len(strings.Trim(v.String, " ")) > 0 {
			return nil
		}
	case sql.NullInt64:
		if v.Valid && v.Int64 != 0 {
			return nil
		}
	case sql.NullFloat64:
		if v.Valid && v.Float64 != 0.00 {
			return nil
		}
	case sql.NullBool:
		if v.Valid {
			return nil
		}
	default:
		return fmt.Errorf("Invalid argument type")
	}

	return err
}

func StringWithDefault(s sql.NullString, d string) string {
	if s.Valid && len(strings.Trim(s.String, " ")) > 0 {
		return s.String
	}

	return d
}

func FromNullString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

func FromNullFloat(s sql.NullFloat64) float64 {
	if s.Valid {
		return s.Float64
	}
	return 0.00
}

func FromNullInt(i sql.NullInt64) int64 {
	if i.Valid {
		return i.Int64
	}

	return 0
}
