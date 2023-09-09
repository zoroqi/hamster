package data

import "time"

type Manga struct {
	Source     string    `json:"source"`
	Title      string    `json:"title"`
	Link       string    `json:"link"`
	Type       string    `json:"type"`
	Cover      string    `json:"cover"`
	Last       string    `json:"last"`
	LastUpdate time.Time `json:"lastUpdate"`
}
