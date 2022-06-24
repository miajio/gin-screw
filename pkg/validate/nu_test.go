package validate_test

import (
	"fmt"
	"testing"
	"time"
)

func TestNu(t *testing.T) {
	date, _ := time.Parse("2006-01-02", "2022-06-21")
	now := time.Now()
	if date.Year() <= now.Year() && date.Month() <= now.Month() && date.Day() < now.Day() {
		fmt.Println("true")
		return
	}
	fmt.Println("false")
}
