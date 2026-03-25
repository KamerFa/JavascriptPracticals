package prefs

type Langs struct {
	Arabic          bool `json:"arabic"`
	Transliteration bool `json:"transliteration"`
	Bosnian         bool `json:"bosnian"`
	English         bool `json:"english"`
	Turkish         bool `json:"turkish"`
}

type Prefs struct {
	RevIndex  int    `json:"revIndex"`
	Filter    string `json:"filter"`
	Langs     Langs  `json:"langs"`
	SeerahOpen bool  `json:"seerahOpen"`
	SeerahLang string `json:"seerahLang"`
}

type Bookmark struct {
	ID        int    `json:"id,omitempty"`
	SurahNum  int    `json:"surahNum"`
	AyahNum   int    `json:"ayahNum"`
	RevIndex  int    `json:"revIndex"`
	SurahName string `json:"surahName"`
}
