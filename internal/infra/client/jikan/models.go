package jikan

type Image struct {
	ImageURL      string  `json:"image_url"`
	SmallURL      *string `json:"small_image_url,omitempty"`
	LargeImageURL *string `json:"large_image_url,omitempty"`
}

type CharacterImages struct {
	JPG  Image `json:"jpg"`
	Webp Image `json:"webp"`
}

type CharacterData struct {
	MalID  int             `json:"mal_id"`
	Name   string          `json:"name"`
	URL    string          `json:"url"`
	Images CharacterImages `json:"images"`
}

type Character struct {
	Character CharacterData `json:"character"`
	Role      string        `json:"role"`
	Favorites int           `json:"favorites"`
}

type JikanCharactersResponse struct {
	Data []Character `json:"data"`
}

type CharactersNode struct {
	MalID     int    `json:"mal_id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Favorites int    `json:"favorites"`
	Image     Image  `json:"image_url"`
}

type CharactersResponse []CharactersNode

type CharacterFullImages struct {
	JPG  Image `json:"jpg"`
	Webp Image `json:"webp"`
}

type CharacterAnimeImages struct {
	JPG  Image `json:"jpg"`
	Webp Image `json:"webp"`
}

type CharacterAnimeData struct {
	MalID  int                  `json:"mal_id"`
	URL    string               `json:"url"`
	Images CharacterAnimeImages `json:"images"`
	Title  string               `json:"title"`
}

type CharacterAnime struct {
	Role  string             `json:"role"`
	Anime CharacterAnimeData `json:"anime"`
}

type VoiceActorImage struct {
	JPG Image `json:"jpg"`
}

type VoiceActorPerson struct {
	MalID  int             `json:"mal_id"`
	URL    string          `json:"url"`
	Images VoiceActorImage `json:"images"`
	Name   string          `json:"name"`
}

type VoiceActor struct {
	Person   VoiceActorPerson `json:"person"`
	Language string           `json:"language"`
}

type CharacterFullData struct {
	MalID     int                 `json:"mal_id"`
	URL       string              `json:"url"`
	Images    CharacterFullImages `json:"images"`
	Name      string              `json:"name"`
	NameKanji string              `json:"name_kanji"`
	Nicknames []string            `json:"nicknames"`
	Favorites int                 `json:"favorites"`
	About     string              `json:"about"`
	Anime     []CharacterAnime    `json:"anime"`
	Voices    []VoiceActor        `json:"voices"`
}

type JikanCharacterFullDataResponse struct {
	Data CharacterFullData `json:"data"`
}

type PersonImages struct {
	JPG Image `json:"jpg"`
}

type PersonData struct {
	MalID          int          `json:"mal_id"`
	URL            string       `json:"url"`
	WebsiteURL     *string      `json:"website_url"`
	Images         PersonImages `json:"images"`
	Name           string       `json:"name"`
	GivenName      *string      `json:"given_name"`
	FamilyName     *string      `json:"family_name"`
	AlternateNames []string     `json:"alternate_names"`
	Birthday       *string      `json:"birthday"`
	Favorites      int          `json:"favorites"`
	About          *string      `json:"about"`
}

type PersonAnimeData struct {
	MalID  int    `json:"mal_id"`
	URL    string `json:"url"`
	Images struct {
		JPG  Image `json:"jpg"`
		Webp Image `json:"webp"`
	} `json:"images"`
	Title string `json:"title"`
}

type PersonAnime struct {
	Position string          `json:"position"`
	Anime    PersonAnimeData `json:"anime"`
}

type PersonVoiceData struct {
	Role      string          `json:"role"`
	Anime     PersonAnimeData `json:"anime"`
	Character struct {
		MalID  int    `json:"mal_id"`
		URL    string `json:"url"`
		Images struct {
			JPG  Image `json:"jpg"`
			Webp Image `json:"webp"`
		} `json:"images"`
		Name string `json:"name"`
	} `json:"character"`
}

type PersonFullData struct {
	MalID          int               `json:"mal_id"`
	URL            string            `json:"url"`
	WebsiteURL     *string           `json:"website_url"`
	Images         PersonImages      `json:"images"`
	Name           string            `json:"name"`
	GivenName      *string           `json:"given_name"`
	FamilyName     *string           `json:"family_name"`
	AlternateNames []string          `json:"alternate_names"`
	Birthday       *string           `json:"birthday"`
	Favorites      int               `json:"favorites"`
	About          *string           `json:"about"`
	Anime          []PersonAnime     `json:"anime"`
	Voices         []PersonVoiceData `json:"voices"`
}

type JikanPersonFullDataResponse struct {
	Data PersonFullData `json:"data"`
}
