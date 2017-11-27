(function(global) {
	"use strict";

	var adjective = ["adamant", "adroit", "amatory", "animistic", "antic", "arcadian", "baleful", "bellicose", "bilious", "boorish", "calamitous", "caustic", "cerulean", "comely", "concomitant", "contumacious", "corpulent", "crapulous", "defamatory", "didactic", "dilatory", "dowdy", "efficacious", "effulgent", "egregious", "endemic", "equanimous", "execrable", "fastidious", "feckless", "fecund", "friable", "fulsome", "garrulous", "guileless", "gustatory", "heuristic", "histrionic", "hubristic", "incendiary", "insidious", "insolent", "intransigent", "inveterate", "invidious", "irksome", "jejune", "jocular", "judicious", "lachrymose", "limpid", "loquacious", "luminous", "mannered", "mendacious", "meretricious", "minatory", "mordant", "munificent", "nefarious", "noxious", "obtuse", "parsimonious", "pendulous", "pernicious", "pervasive", "petulant", "platitudinous", "precipitate", "propitious", "puckish", "querulous", "quiescent", "rebarbative", "recalcitant", "redolent", "rhadamanthine", "risible", "ruminative", "sagacious", "salubrious", "sartorial", "sclerotic", "serpentine", "spasmodic", "strident", "taciturn", "tenacious", "tremulous", "trenchant", "turbulent", "turgid", "ubiquitous", "uxorious", "verdant", "voluble", "voracious", "wheedling", "withering", "zealous", "floofy"];
	var noun = ["ninja", "chair", "pancake", "statue", "unicorn", "rainbows", "laser", "senor", "bunny", "captain", "nibblets", "cupcake", "carrot", "gnomes", "glitter", "potato", "salad", "toejam", "curtains", "beets", "toilet", "exorcism", "stick figures", "mermaid eggs", "sea barnacles", "dragons", "jellybeans", "snakes", "dolls", "bushes", "cookies", "apples", "ice cream", "ukulele", "kazoo", "banjo", "opera singer", "circus", "trampoline", "carousel", "carnival", "locomotive", "hot air balloon", "praying mantis", "animator", "artisan", "artist", "colorist", "inker", "coppersmith", "director", "designer", "flatter", "stylist", "leadman", "limner", "make-up artist", "model", "musician", "penciller", "producer", "scenographer", "set decorator", "silversmith", "teacher", "auto mechanic", "beader", "bobbin boy", "clerk of the chapel", "filling station attendant", "foreman", "maintenance engineering", "mechanic", "miller", "moldmaker", "panel beater", "patternmaker", "plant operator", "plumber", "sawfiler", "shop foreman", "soaper", "stationary engineer", "wheelwright", "woodworkers", "beavers", "cats", "floofs"];
	var verb = ['destroy', 'can’t even', 'win', 'have guns', 'get jokes', 'know things', 'trivialize', 'lift', 'guess well', 'lawyer up', 'care'];

	global.GenerateTeamName = function() {
		return TitleCase(pick(teamFormats)());
	};

	var teamFormats = [
		function() {
			return pick(adjective) + " " + pick(noun);
		},
		function() {
			return pick(adjective) + " " + pick(adjective);
		},
		function() {
			return "The " + pick(adjective);
		},
		function() {
			return "Team " + pick(adjective);
		},
		function() {
			return "Team " + pick(noun);
		},
		function() {
			return pick(noun) + " " + pick(verb);
		},
		function() {
			return pick(noun) + " don't " + pick(verb);
		},
		function() {
			return pick(noun) + " with " + pick(noun);
		},
		function() {
			return pick(verb) + " & " + pick(verb);
		},
		function() {
			return pick(verb) + " " + pick(noun);
		}
	];

	global.GenerateGameName = function() {
		return TitleCase(pick(gameFormats)());
	};

	var gameFormats = [
		function() {
			return pick(adjective) + " " + pick(noun);
		},
		function() {
			return pick(adjective) + " " + pick(adjective);
		},
		function() {
			return "The " + pick(adjective);
		},
		function() {
			return "Fall of " + pick(noun);
		},
		function() {
			return "Tragedy of " + pick(noun);
		},
		function() {
			return pick(noun) + " " + pick(verb);
		},
		function() {
			return pick(noun) + " with " + pick(noun);
		},
		function() {
			return pick(verb) + " & " + pick(verb);
		},
		function() {
			return pick(verb) + " " + pick(noun) + " with " + pick(noun);
		}
	];

	function pick(list) {
		return list[Math.floor(Math.random() * list.length)];
	}

	function TitleCase(str) {
		return str.replace(/\w\S*/g, function(txt) {
			return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();
		}).replace(/\bOr\b/g, "or");
	}
})(window);