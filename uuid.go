package uuid

import (
	"math"
	"math/rand"
	"time"
)

var safeChars string = "23456789ABCDEFGHJKLMNPQRSTWXYZabcdefghijkmnopqrstuvwxyz"
var safeCharsIdx map[byte]float64 = map[byte]float64{}
var timeFrame float64 = 86400000 * 7

func init() {
	for i, c := range safeChars {
		safeCharsIdx[byte(c)] = float64(i)
	}
}

func checkDigit(what string) uint8 {
	var hash int32 = 0
	n := len(what)

	for i := 0; i < n-1; i++ {
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
	if len(what) != 20 && len(what) != 22 {
		return false
	}
	check := checkDigit(what)
	digit := what[len(what)-1]
	return check == digit
}

func Date(what string) *time.Time {
	if !Valid(what) {
		return nil
	}

	a := safeCharsIdx[byte(what[0])]
	b := safeCharsIdx[byte(what[1])]
	t := (a*55 + b + 2000) * timeFrame

	m := float64(55)
	l := 3 + (len(what) - 20)
	for i := 1; i <= l; i++ {
		c := safeCharsIdx[byte(what[1+i])]
		t += c * timeFrame / m
		m *= 55
	}

	t = math.Ceil(t)
	tm := time.Unix(int64(math.Floor(t/1000)), (int64(t)%1000)*int64(time.Millisecond))
	return &tm
}

func Code(what string) string {
	if !Valid(what) {
		return ""
	}

	return string(what[len(what)-15])
}

func create(c string, t time.Time, zeroed bool, milli bool) string {
	if len(c) != 1 {
		c = "1"
	}

	var (
		now    float64
		weeks  float64
		offset float64
		scale  float64
		last   int = 3
		d      []byte
	)

	if milli {
		last += 2
	}

	d = make([]byte, 20+(last-3))

	now = float64(t.UnixNano()) / float64(time.Millisecond)
	weeks = math.Floor(now / timeFrame)
	offset = now - weeks*timeFrame
	scale = 55
	weeks -= 2000

	d[0] = safeChars[int64(math.Floor(weeks/55))%55]
	d[1] = safeChars[int64(weeks)%55]

	for i := 1; i <= last; i++ {
		remain := math.Floor(offset / timeFrame * scale)
		offset -= remain * timeFrame / scale
		scale *= 55
		d[1+i] = safeChars[int64(remain)]
	}

	d[2+last] = c[0]

	if !zeroed {
		rand.Seed(time.Now().UnixNano())
	}
	for i := last + 3; i < 16+last; i++ {
		if zeroed {
			d[i] = c[0]
		} else {
			d[i] = safeChars[rand.Int31n(54)]
		}
	}
	d[16+last] = 0

	d[16+last] = checkDigit(string(d))

	return string(d)
}

func New(code string) string            { return create(code, time.Now(), false, false) }
func Before(date time.Time) string      { return create("0", date, true, false) }
func NewMilli(code string) string       { return create(code, time.Now(), false, true) }
func BeforeMilli(date time.Time) string { return create("0", date, true, true) }
