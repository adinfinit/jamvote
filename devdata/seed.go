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
	Info         string
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
	{
		"neon-nights-2024", "Neon Nights 2024", "Glow in the Dark",
		"Create a game where light and glow are central mechanics or aesthetics. Think bioluminescence, neon signs, blacklight effects, or glowing creatures. Make the darkness beautiful!",
		true, false, false, false, -21,
	},
	{
		"pixel-odyssey", "Pixel Odyssey", "Retro Revival",
		"Revisit the golden age of gaming! Build something that captures the spirit of classic 8-bit and 16-bit games. Pixel art encouraged but not required — it's the feel that counts.",
		true, false, false, false, -14,
	},
	{
		"clockwork-dreams", "Clockwork Dreams", "Time Manipulation",
		"Time is your toy. Rewind, fast-forward, slow-mo, parallel timelines — explore what happens when the player can bend time. Bonus points for creative time-loop puzzles.",
		true, false, false, false, -7,
	},

	// Voting open (recently ended jams, voting in progress)
	{
		"cosmic-clash", "Cosmic Clash", "Space Battles",
		"Take the fight to the stars! Whether it's fleet command, dogfights, or orbital bombardment, your game should make players feel the vastness and danger of space combat.",
		false, true, false, false, 3,
	},
	{
		"shadow-realm", "Shadow Realm", "Light and Darkness",
		"Play with the boundary between light and shadow. Shadows can hide, reveal, protect, or threaten. Use contrast as a core design element in gameplay or narrative.",
		false, true, false, false, 10,
	},

	// Voting closed (voting just finished, results pending)
	{
		"wild-cards", "Wild Cards", "Randomness",
		"Embrace chaos! Procedural generation, dice rolls, shuffled decks, random mutations — let unpredictability drive the fun. The best entries make randomness feel fair yet surprising.",
		false, true, true, false, 20,
	},

	// Completed/revealed (past jams with results)
	{
		"ocean-depths", "Ocean Depths", "Underwater Adventure",
		"Dive beneath the waves and explore the mysterious deep. Pressure, oxygen, currents, and strange sea life — the ocean is an alien world right here on Earth.",
		false, true, true, true, 35,
	},
	{
		"robot-uprising", "Robot Uprising", "AI Gone Wrong",
		"The machines have turned. Build a game exploring rogue AI, rebellious robots, or the moment technology slips out of human control. Comedy or horror — your call.",
		false, true, true, true, 60,
	},
	{
		"mystic-forest", "Mystic Forest", "Nature Magic",
		"An ancient forest hums with magic. Craft a game where nature itself is powerful — enchanted groves, talking animals, druidic spells, or ecosystems with a mind of their own.",
		false, true, true, true, 90,
	},
	{
		"fire-and-ice", "Fire and Ice", "Elemental Forces",
		"Harness the raw power of the elements. Fire melts ice, ice freezes water, steam rises — build a game where elemental interactions are at the heart of every challenge.",
		false, true, true, true, 120,
	},
	{
		"tiny-worlds", "Tiny Worlds", "Microscopic",
		"Shrink down and explore worlds invisible to the naked eye. Cells, atoms, insects, dust particles — find the epic in the minuscule.",
		false, true, true, true, 180,
	},
	{
		"last-stand", "Last Stand", "Survival",
		"Everything is against you. Limited resources, relentless enemies, a ticking clock. Make a game about holding on just a little bit longer when all hope seems lost.",
		false, true, true, true, 365,
	},
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
			Info:         def.Info,
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

			grng := teamRNG(tName, gName)
			team := &event.Team{
				Name:    tName,
				Members: members,
				Game: event.Game{
					Name: gName,
					Info: generateGameInfo(grng),
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

var commentsByAspect = map[string][]string{
	"Theme": {
		"Great interpretation of the theme!",
		"The theme connection felt a bit loose.",
		"Really creative take on the theme — wasn't expecting this angle at all.",
		"Theme is present but doesn't really drive the gameplay.",
		"Loved how the theme was woven into every mechanic.",
		"Solid theme usage.",
		"The theme felt like an afterthought unfortunately.",
		"One of the best theme interpretations I've seen this jam. The way you tied it into the narrative was chef's kiss.",
		"Theme is there.",
		"Could have leaned into the theme more, but what's there works.",
	},
	"Enjoyment": {
		"Had a blast playing this!",
		"Fun but a bit short.",
		"I played through it twice — that's rare for a jam game. The core loop is really satisfying.",
		"Controls feel nice and responsive.",
		"Got stuck on level 3 and couldn't figure out what to do.",
		"Enjoyable!",
		"Not my genre but I can see the appeal.",
		"The difficulty curve is perfect. Started easy, ramped up gradually, and the final challenge felt earned.",
		"A bit repetitive after the first few minutes.",
		"Really fun concept, would love to see this expanded post-jam.",
	},
	"Aesthetics": {
		"Love the art style!",
		"Visuals are clean and cohesive.",
		"The art direction is stunning. The color palette alone tells a story. Sound design complements it perfectly.",
		"Could use some polish on the UI.",
		"The music is a banger.",
		"Nice aesthetic.",
		"Placeholder art but the game underneath is solid.",
		"Sound effects are satisfying, especially the jump sound. Small detail but it makes a big difference.",
		"The particle effects are gorgeous.",
		"Visually it works, nothing groundbreaking but competent.",
	},
	"Innovation": {
		"Never seen a mechanic like this before!",
		"Interesting twist on a classic formula.",
		"This is genuinely original. I've played hundreds of jam games and this mechanic is new to me. Patent it.",
		"Pretty standard gameplay loop.",
		"The combination of mechanics is clever.",
		"Fresh idea!",
		"Feels very similar to [other game] but that's not necessarily bad.",
		"The core innovation is the way the two systems interact — separately they're simple, together they're deep.",
		"Not particularly innovative but executed well.",
		"I like the experimental approach even if it doesn't fully land.",
	},
	"Bonus": {
		"Extra polish for a jam game!",
		"The tutorial was really well done.",
		"Incredible amount of content for 72 hours. Four levels, a boss fight, AND a story? How.",
		"Nice attention to detail.",
		"The credits sequence made me smile.",
		"Goes above and beyond.",
	},
}

var gamePremises = []string{
	"You play as a tiny robot trying to escape a crumbling factory before it collapses.",
	"A mysterious signal draws your submarine deeper into an uncharted trench.",
	"Take control of a street cat navigating rooftops to reunite with its owner.",
	"You're a courier in a cyberpunk city — deliver packages while dodging drones and rival gangs.",
	"Guide a spark of light through a dark cavern, illuminating forgotten murals along the way.",
	"A wizard's apprentice accidentally unleashes chaos and must fix each room of the tower.",
	"Play as a chef in a haunted restaurant where the ingredients fight back.",
	"You wake up on a space station with no memory. The AI says everything is fine. It isn't.",
	"Control a paper airplane through a child's imagination — the classroom is your world.",
	"An old lighthouse keeper must keep the flame burning through one final, impossible storm.",
	"Pilot a seed pod through the wind, searching for the perfect place to take root.",
	"A shapeshifting blob must mimic objects in a museum to avoid the security guards.",
	"Defend your anthill against invading beetles by coordinating worker, soldier, and scout ants.",
	"You're a ghost trying to scare enough people out of your house before it gets demolished.",
	"Run a potion shop by day, explore the enchanted forest for ingredients by night.",
}

var gameMadeWith = []string{
	"Made with Godot in 72 hours.",
	"Built using Unity. First jam for two of our members!",
	"Made with Godot 4.3. All assets created during the jam.",
	"Created with PICO-8.",
	"Built in Unreal Engine 5. We bit off more than we could chew but we're happy with the result.",
	"Made from scratch in C++ and SDL2.",
	"Built with Godot. Music composed in LMMS.",
	"Made with Love2D. Art done in Aseprite.",
	"Created in Game Maker. This was a solo project.",
	"Built with Bevy (Rust). Our first game jam using an ECS framework.",
}

var gameControls = []string{
	"WASD to move, Space to jump, Mouse to aim.",
	"Arrow keys to move, Z to interact, X to dash.",
	"Mouse only — click to move, drag to interact.",
	"WASD for movement, E to interact, Q to use ability. Gamepad supported.",
	"Arrow keys or WASD. Space to confirm, Escape to pause.",
	"Point and click. Right-click to examine objects.",
	"WASD + Mouse. Left click to shoot, right click for shield.",
}

var gameExtras = []string{
	"",
	"Art and music made by our team during the jam. SFX from freesound.org.",
	"There's a secret ending if you collect all the hidden stars.",
	"Sound warning: the final level has loud effects.",
	"Tip: you can wall-jump by pressing jump while sliding against a wall.",
	"We ran out of time for a tutorial, so: the glowing things are good, the spiky things are bad.",
	"The game saves automatically between levels.",
	"Best played with headphones for the full audio experience.",
	"Known bug: the pause menu sometimes doesn't unpause. Press Escape twice if that happens.",
	"",
	"",
}

func generateGameInfo(rng *rand.Rand) string {
	premise := gamePremises[rng.IntN(len(gamePremises))]
	madeWith := gameMadeWith[rng.IntN(len(gameMadeWith))]
	controls := gameControls[rng.IntN(len(gameControls))]
	extra := gameExtras[rng.IntN(len(gameExtras))]

	info := premise + "\n\n" + madeWith + "\n\nControls:\n" + controls
	if extra != "" {
		info += "\n\n" + extra
	}
	return info
}

func pickComment(rng *rand.Rand, aspect string) string {
	comments := commentsByAspect[aspect]
	if len(comments) == 0 {
		return ""
	}
	return comments[rng.IntN(len(comments))]
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
		themeMean := 1.5 + trng.Float64()*3.0   // [1.5, 4.5]
		enjoyMean := 1.5 + trng.Float64()*3.0   // [1.5, 4.5]
		aesthetMean := 1.5 + trng.Float64()*3.0 // [1.5, 4.5]
		innovMean := 1.5 + trng.Float64()*3.0   // [1.5, 4.5]
		bonusMean := trng.Float64() * 1.5       // [0, 1.5]

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

			// ~40% chance of commenting on each aspect.
			if vrng.Float64() < 0.4 {
				aspects.Theme.Comment = pickComment(vrng, "Theme")
			}
			if vrng.Float64() < 0.4 {
				aspects.Enjoyment.Comment = pickComment(vrng, "Enjoyment")
			}
			if vrng.Float64() < 0.4 {
				aspects.Aesthetics.Comment = pickComment(vrng, "Aesthetics")
			}
			if vrng.Float64() < 0.4 {
				aspects.Innovation.Comment = pickComment(vrng, "Innovation")
			}
			if vrng.Float64() < 0.2 {
				aspects.Bonus.Comment = pickComment(vrng, "Bonus")
			}

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
