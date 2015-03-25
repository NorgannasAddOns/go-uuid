package uuid

import (
	"time"
	"math"
	"math/rand"
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
	if len(what) != 20 || !Valid(what) {
		return nil
	}

	a := safeCharsIdx[byte(what[0])]
	b := safeCharsIdx[byte(what[1])]
	c := safeCharsIdx[byte(what[2])]
	d := safeCharsIdx[byte(what[3])]
	e := safeCharsIdx[byte(what[4])]

	t := (a * 55 + b + 2000) * timeFrame
	t += math.Floor(c * timeFrame / 55)
	t += math.Floor(d * timeFrame / (55 * 55))
	t += math.Floor(e * timeFrame / (55 * 55 * 55))
	tm := time.Unix(int64(math.Floor(t/1000)), (int64(t) % 1000) * 1000000)
	return &tm
}

func Code(what string) string {
	if len(what) != 20 || !Valid(what) {
		return ""
	}

	return string(what[5])
}

func create(c string, t time.Time, zeroed bool) string {
	if len(c) != 1 {
		c = "1"
	}

 	d := make([]byte, 20)
	var now, weeks, remain, remain2, remain3, offset, scale float64

	now  = float64(t.UnixNano())/1000000
	weeks = math.Floor(now / timeFrame)
	offset = now - weeks * timeFrame
	scale = 55
	remain = math.Floor(offset / timeFrame * scale)
	offset -= remain * timeFrame / scale
	scale *= 55
	remain2 = math.Floor(offset / timeFrame * scale)
	offset -= remain2 * timeFrame / scale
	scale *= 55
	remain3 = math.Floor(offset / timeFrame * scale)

	weeks -= 2000;

	d[0] = safeChars[int64(math.Floor(weeks / 55)) % 55]
	d[1] = safeChars[int64(weeks) % 55]
	d[2] = safeChars[int64(remain)]
	d[3] = safeChars[int64(remain2)]
	d[4] = safeChars[int64(remain3)]
	d[5] = c[0]

	for i := 6; i < 19; i++ {
		if zeroed {
			d[i] = c[0]
		} else {
			d[i] = safeChars[rand.Int31n(54)]
		}
	}
	d[19] = 0;

	d[19] = checkDigit(string(d))

	return string(d)
}

func New(code string) string {
	return create(code, time.Now(), false)
}

func Before(date time.Time) string {
	return create("0", date, true)
}

