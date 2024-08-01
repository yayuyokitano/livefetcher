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

type AuthUser struct {
	ID         int64  `json:"id"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	IsVerified bool   `json:"is_verified"`
	Avatar     string `json:"avatar"`
}

type Live struct {
	Id            int64
	Title         string
	Artists       []string
	OpenTime      time.Time
	StartTime     time.Time
	Price         string
	PriceEnglish  string
	Venue         LiveHouse
	URL           string
	IsFavorited   bool
	FavoriteCount int
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

type FavoriteButtonInfo struct {
	Id            int
	IsFavorited   bool
	FavoriteCount int
}
