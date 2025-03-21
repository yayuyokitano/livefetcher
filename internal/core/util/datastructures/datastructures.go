package datastructures

import (
	"encoding/json"
	"slices"
	"time"
)

type User struct {
	ID                 int                `form:"-"`
	Email              string             `form:"email"`
	Username           string             `form:"username"`
	Nickname           string             `form:"nickname"`
	PasswordHash       string             `form:"-"`
	Bio                string             `form:"bio"`
	Location           string             `form:"location"`
	IsVerified         bool               `form:"-"`
	Avatar             string             `form:"-"`
	CalendarProperties CalendarProperties `form:"-"`
}

type AuthUser struct {
	ID         int    `json:"id"`
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

type GetLivesJsonResponse struct {
	Lives            Lives                     `json:"lives"`
	CalendarEventMap map[string]CalendarEvents `json:"calendarEventMap"`
}

type Paginator struct {
	Total      int `json:"total"`
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
	Page       int `json:"page"`
	TotalPages int `json:"totalPages"`
}

type Lives struct {
	Paginator Paginator `json:"paginator"`
	Lives     []Live    `json:"lives"`
}

type Live struct {
	ID                   int            `json:"id"`
	Title                string         `json:"title"`
	Artists              []string       `json:"artists"`
	OpenTime             time.Time      `json:"opentime"`
	StartTime            time.Time      `json:"starttime"`
	Price                string         `json:"price"`
	PriceEnglish         string         `json:"priceEn"`
	Venue                LiveHouse      `json:"venue"`
	URL                  string         `json:"url"`
	IsFavorited          bool           `json:"isFavorited"`
	FavoriteCount        int            `json:"favoriteCount"`
	ConflictingEvents    CalendarEvents `json:"conflictingEvents"`
	CalendarOpenEventId  string         `json:"calendarOpenEventId"`
	CalendarStartEventId string         `json:"calendarStartEventId"`

	// only used for livelists
	LiveListLiveID  int    `json:"liveListLiveId"`
	LiveListOwnerID int    `json:"liveListOwnerId"`
	Desc            string `json:"desc"`
}

func GetEventEndTime(live Live) time.Time {
	switch len(live.Artists) {
	case 1:
		return live.StartTime.Add(2 * time.Hour)
	case 2:
		return live.StartTime.Add(3 * time.Hour)
	default:
		return live.StartTime.Add(time.Duration(min(len(live.Artists), 10)) * time.Hour)
	}
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
	ID            int  `json:"id"`
	IsFavorited   bool `json:"isFavorited"`
	FavoriteCount int  `json:"favoriteCount"`
}

type LiveListWriteRequest struct {
	ID     int
	UserID int
	Title  string
	Desc   string
}

type LiveList struct {
	ID            int
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
	LiveID            int
	LiveLiveLists     []LiveList
	PersonalLiveLists []LiveList
}

type AddToLiveListParameters struct {
	LiveDesc           string `form:"liveDesc"`
	LiveID             int    `form:"liveId"`
	ExistingLiveListID int    `form:"existingLiveListId"` // specified if existing live
	NewLiveListTitle   string `form:"newLiveListTitle"`   // specified if new live
	AdditionType       string `form:"additionType"`
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
	ID                 int
	Type               NotificationType
	LiveID             *int
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
	Id         int
	UserId     int
	TextSearch string
}

type CalendarType int16

const (
	CalendarTypeGoogle CalendarType = iota + 1
)

type CalendarProperties struct {
	Id    *string `form:"id"`
	Type  *int16  `form:"type"`
	Token *string `form:"token"`
}

type CalendarEvent struct {
	Id    string    `json:"id"`
	Name  string    `json:"name"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type CalendarEvents []CalendarEvent

/*const conflictMargin = 1 * time.Hour

func getHasPrefixIndex(ce CalendarEvents, prefix string) int {
	for i, event := range ce {
		if strings.HasPrefix(event.Name, prefix) {
			return i
		}
	}
	return -1
}*/

func (ce CalendarEvents) ToDataMapString() string {
	b, err := json.Marshal(ce.ToDataMap())
	if err != nil {
		return "{}"
	}
	return string(b)
}

func (ce CalendarEvents) ToDataMap() map[string]CalendarEvents {
	eventMap := make(map[string]CalendarEvents)
	for _, e := range ce {
		if e.Start.Before(time.Now()) && e.End.Before(time.Now()) {
			continue
		}
		neededStrings := make([]string, 0)
		startString := e.Start.Format("2006-01-02")
		endString := e.End.Format("2006-01-02")
		neededStrings = append(neededStrings, startString)
		if e.Start.Hour() < 1 {
			neededStrings = append(neededStrings, e.Start.Add(-2*time.Hour).Format("2006-01-02"))
		}
		if e.End.Hour() > 22 {
			neededStrings = append(neededStrings, e.End.Add(2*time.Hour).Format("2006-01-02"))
		}
		if startString < endString {
			for t := e.End; t.Format("2006-01-02") > startString; t = t.Add(-24 * time.Hour) {
				neededStrings = append(neededStrings, t.Format("2006-01-02"))
			}
		}

		for _, s := range neededStrings {
			if eventMap[s] == nil {
				eventMap[s] = make(CalendarEvents, 0)
			}
			eventMap[s] = append(eventMap[s], e)
		}
	}
	return eventMap
}

/*
func (ce *CalendarEvents) ApplyConflictingEvents(rawLives []Live) []Live {
	slices.SortFunc(*ce, func(a, b CalendarEvent) int {
		return int(a.Start.Unix() - b.Start.Unix())
	})

	lives := make([]Live, len(rawLives))
	copy(lives, rawLives)
	slices.SortFunc(lives, func(a, b Live) int {
		return int(a.StartTime.Unix() - b.StartTime.Unix())
	})

	for i, j := 0, 0; i < len(*ce) && j < len(lives); {
		if !(*ce)[i].End.After(lives[j].StartTime.Add(-conflictMargin)) {
			i++
			continue
		}

		if (*ce)[i].Start.After(GetEventEndTime(lives[j]).Add(conflictMargin)) {
			j++
			continue
		}

		for k := j; (*ce)[i].Start.After(GetEventEndTime(lives[k]).Add(conflictMargin)) && !(*ce)[i].End.After(lives[k].StartTime.Add(-conflictMargin)); k++ {
			if strings.HasPrefix((*ce)[i].Name, "OPEN ") {
				startEventIndex := getHasPrefixIndex(lives[k].ConflictingEvents, "START ")
				if startEventIndex != -1 {
					lives[k].ConflictingEvents[startEventIndex].Name = strings.TrimPrefix(lives[k].ConflictingEvents[startEventIndex].Name, "START ")
					lives[k].ConflictingEvents[startEventIndex].Start = (*ce)[i].Start
					continue
				}
			} else if strings.HasPrefix((*ce)[i].Name, "START ") {
				startEventIndex := getHasPrefixIndex(lives[k].ConflictingEvents, "OPEN ")
				if startEventIndex != -1 {
					lives[k].ConflictingEvents[startEventIndex].Name = strings.TrimPrefix(lives[k].ConflictingEvents[startEventIndex].Name, "OPEN ")
					lives[k].ConflictingEvents[startEventIndex].End = (*ce)[i].End
					continue
				}
			}
			lives[k].ConflictingEvents = append(lives[k].ConflictingEvents, (*ce)[i])
		}
		i++
	}
	return lives
}
*/
