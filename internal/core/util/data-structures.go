package util

import (
	"time"
)

type User struct {
	ID           int64
	Email        string
	Username     string
	Nickname     string
	PasswordHash string
	Bio          string
	Location     string
	IsVerified   bool
	Avatar       string
}

type Live struct {
	Title        string
	Artists      []string
	OpenTime     time.Time
	StartTime    time.Time
	Price        string
	PriceEnglish string
	Venue        LiveHouse
	URL          string
	IsBookmarked bool
	Bookmarks    int
}

type LiveWithID struct {
	ID int64
	Live
}

type Area struct {
	ID         int    `db:"id"`
	Prefecture string `db:"prefecture"`
	Area       string `db:"area"`
}

type LiveHouse struct {
	ID          string  `db:"id"`
	Url         string  `db:"url"`
	Description string  `db:"description"`
	Area        Area    `db:"areas_id"`
	Longitude   float64 `db:"longitude"`
	Latitude    float64 `db:"latitude"`
}
