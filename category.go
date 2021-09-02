package fditch

type Category string

// Enum for all categories that can be found on itch.io
const (
	GameAssets    Category = "game-assets"
	Books                  = "books"
	Comics                 = "comics"
	Tools                  = "tools"
	Games                  = "games"
	PhysicalGames          = "physical-games"
	Soundtracks            = "soundtracks"
	GameMods               = "game-mods"
	Misc                   = "misc"
)

// Array containing all categories.
var Categories = []Category{
	GameAssets,
	Books,
	Comics,
	Tools,
	Games,
	PhysicalGames,
	Soundtracks,
	GameMods,
	Misc,
}
