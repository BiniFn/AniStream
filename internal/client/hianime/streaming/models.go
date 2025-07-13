package streaming

type ScrapedEncryptedSources struct {
	Sources   string         `json:"sources"`
	Server    int            `json:"server"`
	Intro     ScrapedSegment `json:"intro"`
	Outro     ScrapedSegment `json:"outro"`
	Tracks    []ScrapedTrack `json:"tracks"`
	Encrypted bool           `json:"encrypted"`
}

type ScrapedUnencryptedSources struct {
	Source     string         `json:"source"`
	ServerName string         `json:"serverName"`
	Type       string         `json:"type"`
	Intro      ScrapedSegment `json:"intro"`
	Outro      ScrapedSegment `json:"outro"`
	Tracks     []ScrapedTrack `json:"tracks"`
}

type ScrapedSegment struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type ScrapedTrack struct {
	File    string `json:"file"`
	Kind    string `json:"kind"`
	Label   string `json:"label,omitempty"`
	Default bool   `json:"default,omitempty"`
}

type ScrapedMegaplaySources struct {
	Sources struct {
		File string `json:"file"`
	} `json:"sources"`
	Server int            `json:"server"`
	Intro  ScrapedSegment `json:"intro"`
	Outro  ScrapedSegment `json:"outro"`
	Tracks []ScrapedTrack `json:"tracks"`
}
