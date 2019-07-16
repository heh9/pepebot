package api

import (
	"strconv"
)

var Rankings = map[int] map[int] string {
	1: {
		1: "Herald 1",
		2: "Herald 2",
		3: "Herald 3",
		4: "Herald 4",
		5: "Herald 5",
		6: "Herald 6",
		7: "Herald 7",
	},
	2: {
		1: "Guardian 1",
		2: "Guardian 2",
		3: "Guardian 3",
		4: "Guardian 4",
		5: "Guardian 5",
		6: "Guardian 6",
		7: "Guardian 7",
	},
	3: {
		1: "Crusader 1",
		2: "Crusader 2",
		3: "Crusader 3",
		4: "Crusader 4",
		5: "Crusader 5",
		6: "Crusader 6",
		7: "Crusader 7",
	},
	4: {
		1: "Archon 1",
		2: "Archon 2",
		3: "Archon 3",
		4: "Archon 4",
		5: "Archon 5",
		6: "Archon 6",
		7: "Archon 7",
	},
	5: {
		1: "Legend 1",
		2: "Legend 2",
		3: "Legend 3",
		4: "Legend 4",
		5: "Legend 5",
		6: "Legend 6",
		7: "Legend 7",
	},
	6: {
		1: "Ancient 1",
		2: "Ancient 2",
		3: "Ancient 3",
		4: "Ancient 4",
		5: "Ancient 5",
		6: "Ancient 6",
		7: "Ancient 7",
	},
	7: {
		1: "Divine 1",
		2: "Divine 2",
		3: "Divine 3",
		4: "Divine 4",
		5: "Divine 5",
		6: "Divine 6",
		7: "Divine 7",
	},
	8: {
		1: "Immortal 1",
		2: "Immortal 2",
		3: "Immortal 3",
		4: "Immortal 4",
		5: "Immortal 5",
		6: "Immortal 6",
		7: "Immortal 7",
	},
}

func GetPlayerMedalString(rankTier int) string {

	rankTierString := strconv.Itoa(rankTier)

	if rankTier == 0 {
		return "Unknown"
	}

	tier , _ := strconv.Atoi(rankTierString[0:1])
	rank , _ := strconv.Atoi(rankTierString[1:2])

	return Rankings[tier][rank]
}