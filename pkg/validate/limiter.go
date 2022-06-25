package validate

import (
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

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
