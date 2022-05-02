package uid_utils

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	isLowerCase = true
)

func getChar(value int) string {
	if value < 10 {
		return fmt.Sprintf("%1d", value)
	}
	asciiA := 65 //based on aschii table
	if isLowerCase {
		asciiA = 97
	}
	if value < 36 {
		return string(rune(asciiA + value - 10))
	} else {
		return fmt.Sprintf("%2d", value)
	}
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}
func GetUid(prefix string, isLower bool) string {
	isLowerCase = isLower
	dt := time.Now().UTC()
	yr := dt.Year() % 2000
	mm := int(dt.Month())
	h := dt.Hour()
	m := dt.Minute()
	s := dt.Second()
	rand.Seed(time.Now().UnixNano())
	r := RandomInt(1, 999)
	return fmt.Sprintf("%s%s%s%s%s%s%d", prefix, getChar(yr), getChar(mm), getChar(h), getChar(m), getChar(s), r)
}
