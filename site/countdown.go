package site

import (
	"math/rand"
	"time"
)

func IsValidTime(t time.Time) bool {
	zero := time.Unix(0, 0)
	return t.After(zero)
}

type Countdown struct {
	ID   int
	Time time.Time
	Left time.Duration
}

func NewCountdown(target time.Time) Countdown {
	return Countdown{
		ID:   rand.Int(),
		Time: target,
		Left: target.Sub(time.Now()),
	}
}

func (c Countdown) Hours() int {
	return int(c.Left.Hours())
}
func (c Countdown) Minutes() int {
	return int(c.Left.Minutes() - float64(int(c.Left.Hours())*60))
}
func (c Countdown) Seconds() int {
	return int(c.Left.Seconds() - float64(int(c.Left.Minutes())*60))
}
