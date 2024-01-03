package util

import "time"

type Live struct {
	Title        string
	Artists      []string
	OpenTime     time.Time
	StartTime    time.Time
	Price        string
	PriceEnglish string
	Venue        LiveHouse
	URL          string
}

type Area struct {
	ID         int    `db:"id"`
	Prefecture string `db:"prefecture"`
	Area       string `db:"area"`
}

type LiveHouse struct {
	ID          string `db:"id"`
	Url         string `db:"url"`
	Description string `db:"description"`
	Area        Area   `db:"areas_id"`
}
