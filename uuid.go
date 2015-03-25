package uuid

import (
	"fmt"
	"time"
	"math/rand"
)

var safeChars string = "23456789ABCDEFGHJKLMNPQRSTWXYZabcdefghijkmnopqrstuvwxyz"
var safeCharsIdx map[rune]int64 = map[rune]int64{}
var timeFrame int64 = 86400000 * 7

func init() {
	for i, c := range safeChars {
		safeCharsIdx[c] = int64(i)
	}
}

func checkDigit(what string) uint8 {
	var hash int32 = 0
	n := len(what)

	for i := 0; i < n - 1; i++ {
		hash = ((hash << 5) - hash) + int32(what[i])
		hash |= 0
	}
	hash = int32(uint32(hash) % 55)

	if hash >= 55 {
		hash = 54
	}

	return safeChars[hash]
}


func Valid(what string) bool {
	check := checkDigit(what)
	digit := what[19]
	return check == digit
}

func Date(what string) *time.Time {
	if len(what) != 20 {
		return nil
	}

	fmt.Println(rune(what[1]), what[2], what[3])

	a := safeCharsIdx[rune(what[1])]
	b := safeCharsIdx[rune(what[2])]
	c := safeCharsIdx[rune(what[3])]

	var t int64
	t = (a * 55 + b + 2000) * timeFrame + ^^(c * timeFrame / 55)
	tm := time.Unix(^^(t/1000), (t % 1000) * 1000)
	return &tm
}

func New(c string) string {
	if len(c) != 1 {
		c = "1"
	}

 	d := make([]byte, 20)
	d[0] = c[0]
	now := ^^(time.Now().UnixNano() / 1000000)
	weeks := ^^(now / timeFrame)
	ff := timeFrame * 55
	over := now - (weeks * timeFrame)
	remain := int64((float64(over) / float64(timeFrame)) * 55)

	weeks -= 2000;

	d[1] = safeChars[^^(weeks / 55) % 55]
	d[2] = safeChars[weeks % 55]
	d[3] = safeChars[remain]

	for i := 4; i < 19; i++ {
		d[i] = safeChars[rand.Int31n(54)]
	}
	d[19] = 0;

	d[19] = checkDigit(string(d))

	return string(d)
}

