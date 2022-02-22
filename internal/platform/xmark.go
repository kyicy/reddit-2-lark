package platform

import "time"

type Mark struct {
	LastPubDate       time.Time `toml:"-" json:"-"`
	LastPubDateString string    `toml:"date" json:"date"`
}

type XMark struct {
	Items map[string]*Mark `toml:"items" json:"items"`
}
