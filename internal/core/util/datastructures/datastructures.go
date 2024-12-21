package datastructures

import (
	"slices"
	"time"
)

type User struct {
	ID           int64  `form:"-"`
	Email        string `form:"email"`
	Username     string `form:"username"`
	Nickname     string `form:"nickname"`
	PasswordHash string `form:"-"`
	Bio          string `form:"bio"`
	Location     string `form:"location"`
	IsVerified   bool   `form:"-"`
	Avatar       string `form:"-"`
}

type AuthUser struct {
	ID         int64  `json:"id"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	IsVerified bool   `json:"isVerified"`
	Avatar     string `json:"avatar"`
}

type LiveWithGeoJSON struct {
	Lives   []Live        `json:"lives"`
	GeoJson []LiveGeoJSON `json:"geoJson"`
}

type LiveGeoJSON struct {
	Type       string            `json:"type"`
	Properties GeoJSONProperties `json:"properties"`
	Geometry   GeoJSONGeometry   `json:"geometry"`
}

type GeoJSONProperties struct {
	Name         string `json:"name"`
	ID           int    `json:"id"`
	PopupContent string `json:"popupContent"`
}

type GeoJSONGeometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Live struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Artists       []string  `json:"artists"`
	OpenTime      time.Time `json:"opentime"`
	StartTime     time.Time `json:"starttime"`
	Price         string    `json:"price"`
	PriceEnglish  string    `json:"priceEn"`
	Venue         LiveHouse `json:"venue"`
	URL           string    `json:"url"`
	IsFavorited   bool      `json:"isFavorited"`
	FavoriteCount int       `json:"favoriteCount"`

	// only used for livelists
	LiveListLiveID  int64  `json:"-"`
	LiveListOwnerID int64  `json:"-"`
	Desc            string `json:"-"`
}

type Area struct {
	ID         int    `json:"id"`
	Prefecture string `json:"prefecture"`
	Area       string `json:"area"`
}

type LiveHouse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Url         string  `json:"url"`
	Description string  `json:"description"`
	Area        Area    `json:"area"`
	Longitude   float64 `json:"longitude"`
	Latitude    float64 `json:"latitude"`
}

type FavoriteButtonInfo struct {
	ID            int
	IsFavorited   bool
	FavoriteCount int
}

type LiveListWriteRequest struct {
	ID     int64
	UserID int64
	Title  string
	Desc   string
}

type LiveList struct {
	ID            int64
	Title         string
	Desc          string
	LiveDesc      string
	User          User
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Lives         []Live
	IsFavorited   bool
	FavoriteCount int
}

type AddToLiveListTemplateParams struct {
	LiveID            int64
	LiveLiveLists     []LiveList
	PersonalLiveLists []LiveList
}

type AddToLiveListParameters struct {
	LiveDesc           string
	LiveID             int
	ExistingLiveListID int    // specified if existing live
	NewLiveListTitle   string // specified if new live
	AdditionType       string
}

type NotificationType int16

const (
	NotificationTypeEdited NotificationType = iota + 1
	NotificationTypeDeleted
	NotificationTypeAdded
)

type NotificationFieldType int16

const (
	NotificationFieldTitle NotificationFieldType = iota + 1
	NotificationFieldOpenTime
	NotificationFieldStartTime
	NotificationFieldPrice
	NotificationFieldPriceEnglish
	NotificationFieldURL
	NotificationFieldVenue
	NotificationFieldArtists
)

func (nt NotificationFieldType) String() string {
	return [...]string{
		"",
		"notifications.title",
		"notifications.open",
		"notifications.start",
		"notifications.price",
		"notifications.price",
		"notifications.url",
		"notifications.venue",
		"notifications.artists",
	}[nt]
}

func compareNotificationFields(a, b NotificationField) int {
	if a.Type == b.Type {
		return 0
	}
	order := []NotificationFieldType{
		NotificationFieldTitle,
		NotificationFieldVenue,
		NotificationFieldOpenTime,
		NotificationFieldStartTime,
		NotificationFieldArtists,
		NotificationFieldPrice,
		NotificationFieldPriceEnglish,
		NotificationFieldURL,
	}
	for _, fieldType := range order {
		if a.Type == fieldType {
			return -1
		}
		if b.Type == fieldType {
			return 1
		}
	}
	return 0
}

type NotificationFields []NotificationField

func (nf NotificationFields) Sort() NotificationFields {
	slices.SortFunc(nf, func(a, b NotificationField) int {
		return compareNotificationFields(a, b)
	})
	return nf
}

type NotificationField struct {
	Type     NotificationFieldType
	OldValue string
	NewValue string
}

type NotificationsWrapper struct {
	UnseenCount   int
	Notifications []Notification
}

type Notification struct {
	ID                 int64
	Type               NotificationType
	LiveID             *int64
	LiveTitle          string
	Seen               bool
	CreatedAt          time.Time
	NotificationFields NotificationFields
}

type FieldLineItem struct {
	InnerText     string
	IsHighlighted bool
}

type FieldLine struct {
	Old FieldLineItem
	New FieldLineItem
}

type SavedSearch struct {
	Id         int64
	UserId     int64
	TextSearch string
}
