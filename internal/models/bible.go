package models

type BibleAPIResponse struct {
	RandomVerse struct {
		Book    string `json:"book"`
		BookID  string `json:"book_id"`
		Chapter int    `json:"chapter"`
		Text    string `json:"text"`
		Verse   int    `json:"verse"`
	} `json:"random_verse"`
	Translation struct {
		Identifier   string `json:"identifier"`
		Language     string `json:"language"`
		LanguageCode string `json:"language_code"`
		License      string `json:"license"`
		Name         string `json:"name"`
	} `json:"translation"`
}

type RandomVerseResponseDTO struct {
	Translation string `json:"translation"`
	Book        string `json:"book"`
	Chapter     int    `json:"chapter"`
	Verse       int    `json:"verse"`
	Text        string `json:"text"`
}
