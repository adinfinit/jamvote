# Scoring

Games are scored on the following aspects:

<dl>
<dt>Theme (1 &ndash; 5)</dt>
<dd>How well does the game interpret the jam theme?</dd>

<dt>Enjoyment (1 &ndash; 5)</dt>
<dd>How does the game generally feel to play?</dd>

<dt>Aesthetics (1 &ndash; 5)</dt>
<dd>How well is the story, art and audio executed?</dd>

<dt>Innovation (1 &ndash; 5)</dt>
<dd>Is there something novel in the game?</dd>

<dt>Bonus (0 &ndash; 2.5)</dt>
<dd>Is there anything exceptionally special about it?</dd>
</dl>

The **Overall** score is the weighted average of all
aspects: 5 \* (Theme + Enjoyment + Aesthetics + Innovation + Bonus) /
4.5, clamped to the range 1 – 5.

## Why have aspects?

The voting could be done with a single slider from 0 to 5, but having
multiple aspects allows for a more nuanced and comprehensive evaluation
of games. Each aspect provides a different perspective on the game,
allowing voters to consider various facets of its quality. This
approach also helps to ensure that the voting process is fair and
balanced, as it takes into account a range of factors that contribute
to a game's overall quality. Having people write specific feedback
about different aspects of the game allows jammers to improve their
games and learn from each other.

## Why these aspects?

Many aspects were considered, but these were the most important ones
and seemed to cover the key qualities of a game. These scoring
categories act as goals for people to think about. They also try to
balance abilities between different groups, so that the "programmers
team" won't necessarily win over the "artists team".

**"Theme"** is a crucial aspect for game jams. It discourages people
from coming with premeditated games and encourages creativity.
Themes in general make jams more interesting, because you can see
how other people interpret the theme and create unique games.

**"Enjoyment"** is another crucial aspect. It could also have been called
"Fun", "Engagement", "Creativity", or "Playability". However, these
don't fully capture the breadth of what games can be. Games can be
stressful or sad and still be enjoyable. Overall, this tries to
capture two ideas:

- Do you want to play this game again?
- How much did the game make you feel?

**"Aesthetics"** rewards games that are visually cohesive, have good writing,
have suitable sound design, and feel like they belong together.
There are many ways that games can excel besides gameplay, such as visuals,
audio, narrative, character design, and world-building. Aesthetics tries to
capture all of these. It could have been called "Graphics" or
"Sound", but those labels can make people focus on looking "pretty".
Games can be simple or ugly and still suit the experience. For
example, a grungy and sketchy aesthetic can create a gritty and
oppressive atmosphere.

**"Innovation"** rewards people for trying something new. Coming up with
new things in a short time often means that the gameplay mechanics
cannot be balanced, hurting the enjoyment of the game. However, novel
games are more memorable than a remake of a classic. Innovation can
happen in many ways, such as new gameplay mechanics, new art styles,
or new sound design. It's important to remember that innovation
doesn't always mean creating something completely new, but also
improving upon existing ideas.

The **"Bonus"** category captures everything else, including things that
are not covered by the other categories. Sometimes there's a
[quality without a name](https://en.wikipedia.org/wiki/The_Timeless_Way_of_Building#Summary),
something that is difficult to describe but is important to the game.
It can be how all the pieces fit together, some category that's missing,
or some excellence in a specific area where "5" just isn't
enough. However, this is a very subjective category, so it is
weighted at 0.5 compared to the other categories.

## Scoring statistics

The scoring tries to be as fair as possible with the number of votes
available. We use a simple mean to calculate the score, because it is
easy to understand and fair enough. A median could be a better
indicator of overall quality, as it is less influenced by extreme
values; similarly, excluding outliers could improve accuracy.
However, either approach would make the scoring more complicated to understand,
or make it more likely to have ties.

We use 10 as the minimum number of votes to consider a score
reliable. At that point a single vote can change the score by at most
~5%, which can still affect the ranking, but there is no easy way to
get significantly better accuracy with that few votes.

Then the next stage is 30 votes, which is considered to be the
"minimum observations for a large sample". At that point we can be fairly
confident that the final result represents the opinion of everyone involved.

In practice, most games receive significantly more votes than that.
