package validate

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

const (
	NoBreakSp       = '\u00A0' // no-break space
	NarrowNoBreakSp = '\u202F' // narrow no-break space
	WideNoBreakSp   = '\uFEFF' // wide no-break space
	EnNoBreakSp     = '\u0020' // no-break space
	ZhWideNoBreakSp = '\u3000' // wide no-break space

	LeftToRightEmbedding = '\u202a' // left-to-right embedding
)

// DeSpace delete space in the val
func DeSpace(val string) string {
	val = strings.ReplaceAll(val, string(NoBreakSp), "")
	val = strings.ReplaceAll(val, string(NarrowNoBreakSp), "")
	val = strings.ReplaceAll(val, string(WideNoBreakSp), "")
	val = strings.ReplaceAll(val, string(EnNoBreakSp), "")
	val = strings.ReplaceAll(val, string(ZhWideNoBreakSp), "")
	return val
}

// EnglishLimiter english limiter
// this limiter is used to check english
// but all special symbols will also be prohibited in the parameter
// Please use with caution
func EnglishLimiter(fl validator.FieldLevel) bool {
	val, ok := fl.Field().Interface().(string)
	if ok {
		for _, r := range val {
			if r < 'A' || r > 'z' {
				return false
			}
		}
	}
	return true
}

// IntegerLimiter positive integer limiter
// string to int if error then return false else return true
func IntegerLimiter(fl validator.FieldLevel) bool {
	val, ok := fl.Field().Interface().(string)
	if ok {
		if _, err := strconv.Atoi(val); err != nil {
			return false
		}
	}
	return true
}

// NumberLimiter positive number limiter
// string to float if error then return false else return true
func NumberLimiter(fl validator.FieldLevel) bool {
	val, ok := fl.Field().Interface().(string)
	if ok {
		if _, err := strconv.ParseFloat(val, 64); err != nil {
			return false
		}
	}
	return true
}

// EqNowDayLimiter equal now day limiter
func EqNowDayLimiter(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if ok {
		now := time.Now()
		if date.Year() == now.Year() && date.Month() == now.Month() && date.Day() == now.Day() {
			return true
		}
	}
	return false
}

// GtNowDayLimiter greater than now day limiter
func GtNowDayLimiter(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if ok {
		now := time.Now()
		if date.Year() >= now.Year() && date.Month() >= now.Month() && date.Day() > now.Day() {
			return true
		}
	}
	return false
}

// LtNowDayLimiter less than now day limiter
func LtNowDayLimiter(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if ok {
		now := time.Now()
		if date.Year() <= now.Year() && date.Month() <= now.Month() && date.Day() < now.Day() {
			return true
		}
	}
	return false
}
