package models

type Event struct {
	ChangeType string `json:"type"`
	NameSpace  int    `json:"namespace"`
	Title      string `json:"title"`
	TitleUrl   string `json:"title_url"`
	User       string `json:"user"`
	Bot        bool   `json:"bot"`
}

type WikiEdit struct {
	Title       string `json:"TITLE"`
	EditsCount  int    `json:"EDITS_COUNT"`
	WindowStart int    `json:"WINDOW_START"`
	WindowEnd   int    `json:"WINDOW_END"`
}
