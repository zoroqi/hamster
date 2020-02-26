package manga

type Manga struct {
	Id     string  `json:"link"`
	Title  string  `json:"title"`
	Titles []Title `json:"titles"`
	//Link           string    `json:"link"`
	WebSites       []WebSite `json:"webSites"`
	Type           string    `json:"type"`
	LastUpdateTime string    `json:"lastUpdateTime"`
}

type Title struct {
	Name     string `json:"name"`
	Language string `json:"language"`
}

type WebSite struct {
	Url   string `json:"url"`
	Brief string `json:"brief"`
}
