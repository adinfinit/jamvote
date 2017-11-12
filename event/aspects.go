package event

type Aspects struct {
	Theme      float64
	Enjoyment  float64
	Aesthetics float64
	Innovation float64
	Bonus      float64
}

func (aspects *Aspects) EnsureRange() {
	clamp(&aspects.Theme, 1, 5)
	clamp(&aspects.Enjoyment, 1, 5)
	clamp(&aspects.Aesthetics, 1, 5)
	clamp(&aspects.Innovation, 1, 5)
	clamp(&aspects.Bonus, 1, 5)
}

func (aspects *Aspects) Total() float64 {
	return (aspects.Theme +
		aspects.Enjoyment +
		aspects.Aesthetics +
		aspects.Innovation +
		aspects.Bonus*0.5) / (5*4 + 5*0.5)
}

func clamp(v *float64, min, max float64) {
	if *v < min {
		*v = min
	}
	if *v > max {
		*v = max
	}
}

/*

Theme
How well does it interpret the theme
1 Not even close
2 Resembling
3 Related
4 Spot on
5 Novel Interpretation

Enjoyment
How does the game generally feel
1 Boring
2 Not playing again
3 Nice
4 Didn't want to stop
5 Will play later.

Aesthetics
How well is the story, art and audio executed
1 None
2 Needs tweaks
3 Nice
4 Really good
5 Exceptional

Innovation
Something novel in the game
1 Seen it a lot
2 Interesting variation
3 Interesting approach
4 Never seen before
5 Exceptional

Bonus
Anything exceptionally special about * 0,5
1 Nothing special
2 Really liked *
3 Really loved *
4 Loved everything
5 <3

*/
