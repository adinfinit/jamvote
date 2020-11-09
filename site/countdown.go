package site

import (
	"math/rand"
	"time"
)

// IsValidTime verifies whether t is non-zero.
func IsValidTime(t time.Time) bool {
	zero := time.Unix(0, 0)
	return t.After(zero)
}

// Countdown is a countdown time.
type Countdown struct {
	ID   int
	Time time.Time
	Left time.Duration
}

// NewCountdown creates a countdown with a random identifier.
func NewCountdown(target time.Time) Countdown {
	return Countdown{
		ID:   rand.Int(),
		Time: target,
		Left: target.Sub(time.Now()),
	}
}

// Hours returns how many hours there are left.
func (c Countdown) Hours() int {
	return int(c.Left.Hours())
}

// Minutes returns how many minutes there are left.
func (c Countdown) Minutes() int {
	return int(c.Left.Minutes() - float64(int(c.Left.Hours())*60))
}

// Seconds returns how many seconds there are left.
func (c Countdown) Seconds() int {
	return int(c.Left.Seconds() - float64(int(c.Left.Minutes())*60))
}
