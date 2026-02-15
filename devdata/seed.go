package devdata

import (
	"context"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"log/slog"
	"math"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/adinfinit/jamvote/auth"
	"github.com/adinfinit/jamvote/datastoredb"
	"github.com/adinfinit/jamvote/event"
	"github.com/adinfinit/jamvote/user"
)

var userNames = []string{
	"Admin", "Alice", "Bob", "Charlie", "Diana",
	"Eve", "Frank", "Grace", "Hank", "Ivy",
	"Jack", "Karen", "Leo", "Mona", "Nate",
	"Olivia", "Pete", "Quinn", "Ruby", "Sam",
	"Tina", "Ulysses", "Vera", "Walt", "Xena",
	"Yuri", "Zara", "Amber", "Blake", "Cleo",
	"Derek", "Elsa", "Felix", "Gina", "Hugo",
	"Iris", "Jasper", "Kira", "Liam", "Maya",
	"Noah", "Opal", "Piper", "Reed", "Sage",
	"Troy", "Uma", "Vince", "Wren", "Zelda",
}

var teamNames = []string{
	"Pixel Pirates", "Code Wizards", "Byte Busters", "Dream Weavers",
	"Neon Coders", "Glitch Goblins", "Logic Lords", "Turbo Turtles",
	"Data Dragons", "Sprite Smiths", "Shader Sharks", "Vector Vikings",
	"Kernel Knights", "Raster Rebels", "Mesh Mages", "Voxel Vandals",
	"Loot Llamas", "Bug Bashers", "Frame Fighters", "Stack Samurai",
	"Null Knights", "Poly Phantoms", "Render Rangers", "Bit Bandits",
	"Hex Heroes", "Cache Cats", "Loop Legends", "Spawn Scouts",
	"Flux Foxes", "Warp Wolves", "Ping Pandas", "Debug Demons",
	"Sync Serpents", "Node Nomads", "Frag Falcons", "Query Queens",
	"Drift Droids", "Parse Parrots", "Crypt Crows", "Blaze Bots",
	"Tilt Titans", "Rust Raiders", "Jam Jackals", "Core Cobras",
	"Echo Eagles", "Zoom Zombies", "Grind Griffins", "Mana Monkeys",
	"Saga Spiders", "Loot Lizards", "Nova Newts", "Fury Frogs",
	"Blip Bears", "Warp Wasps", "Dash Doves", "Tide Tigers",
	"Glow Gators", "Rift Ravens", "Bolt Bats", "Hack Hawks",
	"Amp Ants", "Void Vipers", "Snap Snakes", "Clad Clams",
	"Poke Pumas", "Fizz Foxes", "Trek Toads", "Cog Crabs",
	"Orbit Owls", "Plume Pigs", "Dusk Dogs", "Flint Flies",
	"Gust Goats", "Haze Hares", "Isle Ibex", "Jinx Jays",
	"Kite Koi", "Lace Lynx", "Mint Moths", "Nook Narwhals",
	"Oat Otters", "Peak Pines", "Quill Quails", "Reef Rams",
	"Salt Stags", "Turf Terns", "Urn Urchins", "Vale Voles",
	"Weld Wrens", "Yarn Yaks", "Zest Zebus", "Arch Apes",
	"Bark Bees", "Curl Cats", "Dew Deer", "Elm Eels",
	"Fern Fish", "Glen Gulls", "Hill Hens", "Inch Imps",
	"Jade Jags", "Kelp Kites", "Leaf Larks", "Moss Moles",
	"Nest Newts", "Oak Oryx", "Palm Pugs", "Quay Quoks",
	"Root Rooks", "Sand Slugs", "Twig Ticks", "Ursa Umbra",
}

type eventDef struct {
	ID           string
	Name         string
	Theme        string
	Registration bool
	Voting       bool
	Closed       bool
	Revealed     bool
	// EndDaysAgo is the number of days ago the jam ended.
	// Negative values mean the jam ends in the future.
	EndDaysAgo int
}

var eventDefs = []eventDef{
	// Registration stage (upcoming jams, end dates in the future)
	{"neon-nights-2024", "Neon Nights 2024", "Glow in the Dark", true, false, false, false, -21},
	{"pixel-odyssey", "Pixel Odyssey", "Retro Revival", true, false, false, false, -14},
	{"clockwork-dreams", "Clockwork Dreams", "Time Manipulation", true, false, false, false, -7},

	// Voting open (recently ended jams, voting in progress)
	{"cosmic-clash", "Cosmic Clash", "Space Battles", false, true, false, false, 3},
	{"shadow-realm", "Shadow Realm", "Light and Darkness", false, true, false, false, 10},

	// Voting closed (voting just finished, results pending)
	{"wild-cards", "Wild Cards", "Randomness", false, false, true, false, 20},

	// Completed/revealed (past jams with results)
	{"ocean-depths", "Ocean Depths", "Underwater Adventure", false, false, true, true, 35},
	{"robot-uprising", "Robot Uprising", "AI Gone Wrong", false, false, true, true, 60},
	{"mystic-forest", "Mystic Forest", "Nature Magic", false, false, true, true, 90},
	{"fire-and-ice", "Fire and Ice", "Elemental Forces", false, false, true, true, 120},
	{"tiny-worlds", "Tiny Worlds", "Microscopic", false, false, true, true, 180},
	{"last-stand", "Last Stand", "Survival", false, false, true, true, 365},
}

var gameNames = []string{
	"Starbound Escape", "Dungeon Pulse", "Chrono Drift", "Shadow Sprint",
	"Neon Blitz", "Frost Forge", "Pixel Storm", "Void Walker",
	"Flame Dash", "Crystal Caves", "Astro Hop", "Lava Loop",
	"Cyber Slice", "Dream Dash", "Ether Edge", "Fury Flight",
	"Glyph Guard", "Hex Hunt", "Ion Ignite", "Jade Jump",
	"Kinetic Keep", "Luna Lash", "Mist March", "Nova Nudge",
	"Orb Orbit", "Prism Prowl", "Quake Quest", "Rift Run",
	"Spark Surge", "Terra Twist", "Ultra Unity", "Volt Vault",
	"Wave Whirl", "Xenon Xing", "Yonder Yell", "Zephyr Zone",
	"Amber Arc", "Blaze Bolt", "Coral Crash", "Dune Dive",
	"Echo Emit", "Flare Flip", "Gale Grip", "Halo Hurl",
	"Ink Inlet", "Jewel Jolt", "Karma Kick", "Lyric Lift",
	"Magma Maze", "Nebula Nap", "Onyx Oath", "Pyre Plunge",
	"Quill Quake", "Rune Rush", "Shard Shift", "Thorn Trail",
	"Umbra Undo", "Venom Vent", "Wisp Warp", "Xeno Xalt",
	"Yawn Yoke", "Zinc Zap", "Aero Arch", "Brine Bump",
	"Cinder Curl", "Dusk Drop", "Ember Etch", "Fume Furl",
	"Grit Glaze", "Horn Haze", "Isle Itch", "Jest Jab",
	"Knot Knit", "Lime Loom", "Mire Meld", "Nook Nock",
	"Oat Ogle", "Peat Plow", "Quag Quip", "Reed Rile",
	"Silt Sway", "Tuft Turn", "Urge Undo", "Vale Veer",
	"Woad Wilt", "Yarn Yawl", "Zeal Zing", "Acre Axle",
	"Bark Bask", "Clay Clap", "Dell Dint", "Elm Edge",
	"Fen Fray", "Glen Gust", "Husk Hew", "Isle Iris",
	"Jute Jive", "Kelp Keen", "Loam Lurk", "Malt Mend",
	"Nub Nip", "Ore Ooze", "Pulp Prod", "Quay Quash",
	"Rime Roam", "Slag Slew", "Talc Tamp", "Umber Urge",
	"Vat Vow", "Whin Wade", "Xyst Xeno", "Yew Yore",
}

// Seed populates the database with test data if it is empty.
func Seed(log *slog.Logger, db *datastoredb.DB) {
	ctx := context.Background()
	users := db.Users(ctx)

	existing, err := users.List()
	if err != nil {
		log.Error("seed: failed to list users", "error", err)
		return
	}
	if len(existing) > 0 {
		log.Info("seed: database already has users, skipping")
		return
	}

	log.Info("seed: populating database with test data")

	// Create users.
	userIDs := make([]user.UserID, len(userNames))
	for i, name := range userNames {
		cred := &auth.Credentials{
			Provider: "development",
			ID:       auth.DevelopmentUserID(name),
			Email:    strings.ToLower(name) + "@example.com",
			Name:     name,
		}
		u := &user.User{
			Name:  name,
			Email: cred.Email,
			Admin: i == 0,
		}
		id, err := users.Create(cred, u)
		if err != nil {
			log.Error("seed: failed to create user", "name", name, "error", err)
			return
		}
		userIDs[i] = id
	}

	log.Info("seed: created users", "count", len(userIDs))

	adminID := userIDs[0]
	events := db.Events(ctx)
	teamIdx := 0

	for i, def := range eventDefs {
		ev := &event.Event{
			ID:           event.EventID(def.ID),
			Name:         def.Name,
			Theme:        def.Theme,
			Created:      time.Now().AddDate(0, 0, -def.EndDaysAgo-7),
			StartTime:    time.Now().AddDate(0, 0, -def.EndDaysAgo-2),
			EndTime:      time.Now().AddDate(0, 0, -def.EndDaysAgo),
			Registration: def.Registration,
			Voting:       def.Voting,
			Closed:       def.Closed,
			Revealed:     def.Revealed,
			Organizers:   []user.UserID{adminID},
		}

		// Assign jammers: pick 30 users starting at offset based on event index.
		var jammers []user.UserID
		for j := 0; j < 30; j++ {
			uid := userIDs[(i*7+j)%len(userIDs)]
			jammers = append(jammers, uid)
		}
		ev.Jammers = jammers

		// Assign 3 judges per event.
		var judges []user.UserID
		for j := 0; j < 3; j++ {
			uid := userIDs[(i*3+j+40)%len(userIDs)]
			judges = append(judges, uid)
		}
		ev.Judges = judges

		if err := events.Create(ev); err != nil {
			log.Error("seed: failed to create event", "id", def.ID, "error", err)
			return
		}

		// Create 6-12 teams per event.
		erng := eventRNG(def.ID)
		teamCount := 6 + erng.IntN(7)
		var teams []*event.Team
		for t := 0; t < teamCount; t++ {
			memberCount := 1 + erng.IntN(5) // 1-5 members
			var members []event.Member
			for m := 0; m < memberCount; m++ {
				uid := userIDs[(i*10+t*3+m)%len(userIDs)]
				members = append(members, event.Member{
					ID:   uid,
					Name: userNames[(i*10+t*3+m)%len(userNames)],
				})
			}

			tName := teamNames[teamIdx%len(teamNames)]
			gName := gameNames[teamIdx%len(gameNames)]
			teamIdx++

			team := &event.Team{
				Name:    tName,
				Members: members,
				Game: event.Game{
					Name: gName,
					Info: fmt.Sprintf("A game created for %s by team %s.", def.Name, tName),
				},
			}
			team.Game.Link.Download = fmt.Sprintf("https://example.com/games/%s", strings.ReplaceAll(strings.ToLower(gName), " ", "-"))

			teamID, err := events.CreateTeam(ev.ID, team)
			if err != nil {
				log.Error("seed: failed to create team", "event", def.ID, "team", tName, "error", err)
				return
			}
			team.ID = teamID
			teams = append(teams, team)
		}

		// Generate ballots for events past registration.
		if def.Voting || def.Closed || def.Revealed {
			ballotCount := seedBallots(events, ev, teams, jammers)
			log.Info("seed: created event", "id", def.ID, "teams", len(teams), "ballots", ballotCount)
		} else {
			log.Info("seed: created event", "id", def.ID, "teams", len(teams))
		}
	}

	log.Info("seed: done")
}

// eventRNG creates a deterministic RNG seeded from the event ID.
func eventRNG(eventID string) *rand.Rand {
	h := fnv.New64a()
	h.Write([]byte(eventID))
	return rand.New(rand.NewPCG(h.Sum64(), 0))
}

// teamRNG creates a deterministic RNG seeded from the team and game name.
func teamRNG(teamName, gameName string) *rand.Rand {
	h := fnv.New64a()
	h.Write([]byte(teamName))
	h.Write([]byte{0})
	h.Write([]byte(gameName))
	return rand.New(rand.NewPCG(h.Sum64(), 0))
}

// voterRNG creates a deterministic RNG seeded from voter, team, and event.
func voterRNG(eventID string, voterID user.UserID, teamName string) *rand.Rand {
	h := fnv.New64a()
	h.Write([]byte(eventID))
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(voterID))
	h.Write(buf[:])
	h.Write([]byte(teamName))
	return rand.New(rand.NewPCG(h.Sum64(), 0))
}

// normalScore generates a normally distributed score clamped to [min, max].
func normalScore(rng *rand.Rand, mean, stddev, min, max float64) float64 {
	v := rng.NormFloat64()*stddev + mean
	// Round to nearest 0.1.
	v = math.Round(v*10) / 10
	return math.Max(min, math.Min(max, v))
}

// seedBallots generates ballots for an event. Each team gets a hash-derived
// mean quality, and each voter's scores are normally distributed around it.
func seedBallots(events event.Repo, ev *event.Event, teams []*event.Team, voters []user.UserID) int {
	count := 0
	for _, team := range teams {
		// Derive per-team mean scores from team+game name hash.
		trng := teamRNG(team.Name, team.Game.Name)
		themeMean := 1.5 + trng.Float64()*3.0     // [1.5, 4.5]
		enjoyMean := 1.5 + trng.Float64()*3.0     // [1.5, 4.5]
		aesthetMean := 1.5 + trng.Float64()*3.0   // [1.5, 4.5]
		innovMean := 1.5 + trng.Float64()*3.0     // [1.5, 4.5]
		bonusMean := trng.Float64() * 1.5          // [0, 1.5]

		for _, voterID := range voters {
			if team.HasMemberID(voterID) {
				continue
			}

			vrng := voterRNG(string(ev.ID), voterID, team.Name)

			// ~70% of voters submit a ballot.
			if vrng.Float64() > 0.7 {
				continue
			}

			aspects := event.Aspects{
				Theme:      event.Aspect{Score: normalScore(vrng, themeMean, 0.7, 1, 5)},
				Enjoyment:  event.Aspect{Score: normalScore(vrng, enjoyMean, 0.7, 1, 5)},
				Aesthetics: event.Aspect{Score: normalScore(vrng, aesthetMean, 0.7, 1, 5)},
				Innovation: event.Aspect{Score: normalScore(vrng, innovMean, 0.7, 1, 5)},
				Bonus:      event.Aspect{Score: normalScore(vrng, bonusMean, 0.5, 0, 2.5)},
			}
			aspects.Overall = event.Aspect{Score: aspects.Total()}

			ballot := &event.Ballot{
				Voter:     voterID,
				Team:      team.ID,
				Index:     int64(count),
				Completed: true,
				Aspects:   aspects,
			}

			if err := events.SubmitBallot(ev.ID, ballot); err != nil {
				continue
			}
			count++
		}
	}
	return count
}
