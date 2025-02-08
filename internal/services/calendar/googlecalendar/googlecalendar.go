package googlecalendar

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yayuyokitano/livefetcher/internal/core/util/datastructures"
	"github.com/yayuyokitano/livefetcher/internal/services"
	calendarqueries "github.com/yayuyokitano/livefetcher/internal/services/calendar/queries"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type GoogleCalendar struct{}

func getConfig() (config *oauth2.Config, err error) {
	b, err := os.ReadFile("internal/services/calendar/googlecalendar/credentials.json")
	if err != nil {
		return
	}

	config, err = google.ConfigFromJSON(b, "https://www.googleapis.com/auth/calendar.events.owned", "https://www.googleapis.com/auth/calendar.calendars.readonly")
	return
}

func initService(ctx context.Context, props datastructures.CalendarProperties) (service *calendar.Service, err error) {
	if props.Id == nil || props.Token == nil {
		err = errors.New("no calendar information found")
		return
	}

	config, err := getConfig()
	if err != nil {
		return
	}

	token := &oauth2.Token{}
	err = json.NewDecoder(strings.NewReader(*props.Token)).Decode(token)
	if err != nil {
		return
	}

	return calendar.NewService(ctx, option.WithHTTPClient(config.Client(ctx, token)))
}

func buildEvents(live datastructures.Live) (*calendar.Event, *calendar.Event) {
	fmt.Printf("%+v\n", live)
	openEvent := &calendar.Event{
		Summary:     "OPEN " + live.Venue.Name,
		Description: "<h2>ARTIST</h2><ul>",
		Location:    live.Venue.Name,
		Start: &calendar.EventDateTime{
			DateTime: live.OpenTime.Format(time.RFC3339),
		},
		End: &calendar.EventDateTime{
			DateTime: live.StartTime.Format(time.RFC3339),
		},
	}

	for _, a := range live.Artists {
		openEvent.Description += "<li>" + a + "</li>"
	}
	openEvent.Description += `</ul><a href="https://example.com" target="_blank" rel="noopener noreferrer">LiveRadar</a>`

	startEvent := &calendar.Event{
		Summary:     "START " + live.Venue.Name,
		Description: openEvent.Description,
		Location:    live.Venue.Name,
		Start: &calendar.EventDateTime{
			DateTime: live.StartTime.Format(time.RFC3339),
		},
		End: &calendar.EventDateTime{
			DateTime: datastructures.GetEventEndTime(live).Format(time.RFC3339),
		},
	}

	return openEvent, startEvent
}

func (g *GoogleCalendar) PostEvent(ctx context.Context, props datastructures.CalendarProperties, userId int64, live datastructures.Live) (newLive datastructures.Live, err error) {
	service, err := initService(ctx, props)
	if err != nil {
		return
	}

	openEvent, startEvent := buildEvents(live)

	openCall := service.Events.Insert(*props.Id, openEvent)
	openCalendarEvent, err := openCall.Do()
	if err != nil {
		return
	}

	startCall := service.Events.Insert(*props.Id, startEvent)
	startCalendarEvent, err := startCall.Do()
	if err != nil {
		return
	}

	err = calendarqueries.PostCalendarId(ctx, live.ID, userId, openCalendarEvent.Id, startCalendarEvent.Id)
	if err != nil {
		return
	}
	newLive = live
	newLive.CalendarOpenEventId = openCalendarEvent.Id
	newLive.CalendarStartEventId = startCalendarEvent.Id
	return
}

func (g *GoogleCalendar) PutEvent(ctx context.Context, props datastructures.CalendarProperties, userId int64, live datastructures.Live) (err error) {
	service, err := initService(ctx, props)
	if err != nil {
		return
	}

	openEvent, startEvent := buildEvents(live)

	openEventId, startEventId, err := calendarqueries.GetCalendarId(ctx, live.ID, userId)
	if err != nil {
		return
	}

	openCall := service.Events.Update(*props.Id, openEventId, openEvent)
	_, err = openCall.Do()
	if err != nil {
		return
	}

	startCall := service.Events.Update(*props.Id, startEventId, startEvent)
	_, err = startCall.Do()
	return
}

func (g *GoogleCalendar) DeleteEvent(ctx context.Context, props datastructures.CalendarProperties, userId, liveId int64) (err error) {
	service, err := initService(ctx, props)
	if err != nil {
		return
	}

	openEventId, startEventId, err := calendarqueries.GetCalendarId(ctx, liveId, userId)
	if err != nil {
		return
	}

	openCall := service.Events.Delete(*props.Id, openEventId)
	err = openCall.Do()
	if err != nil {
		return
	}

	startCall := service.Events.Delete(*props.Id, startEventId)
	err = startCall.Do()
	if err != nil {
		return
	}

	err = calendarqueries.DeleteCalendarId(ctx, userId, openEventId, startEventId)
	return
}

func (g *GoogleCalendar) GetAllEvents(ctx context.Context, props datastructures.CalendarProperties, userId int64) (events []datastructures.CalendarEvent, err error) {
	service, err := initService(ctx, props)
	if err != nil {
		return
	}

	serviceEvents, err := service.Events.List(*props.Id).Do()
	if err != nil {
		return
	}

	for _, e := range serviceEvents.Items {
		startTime, err := time.Parse(time.RFC3339, e.Start.DateTime)
		if err != nil {
			// TODO: log
			continue
		}

		endTime, err := time.Parse(time.RFC3339, e.End.DateTime)
		if err != nil {
			// TODO: log
			continue
		}

		events = append(events, datastructures.CalendarEvent{
			Id:    e.Id,
			Name:  e.Summary,
			Start: startTime,
			End:   endTime,
		})
	}
	return
}

func createStateToken(userId int64) (token string, err error) {
	t := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.RegisteredClaims{
		Issuer:    "livefetcher-auth",
		Subject:   fmt.Sprintf("%d", userId),
		Audience:  jwt.ClaimStrings{"https://example.com"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	token, err = t.SignedString(services.PrivateKey)
	return
}

func verifyStateToken(token string) (userId int64, err error) {
	parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return services.PublicKey, nil
	})
	if err != nil {
		return
	}

	if rc, ok := parsedToken.Claims.(*jwt.RegisteredClaims); ok {
		var iuser int
		iuser, err = strconv.Atoi(rc.Subject)
		if err != nil {
			return
		}
		userId = int64(iuser)
	} else {
		err = errors.New("invalid auth token")
	}
	return
}

func GetGoogleAuthCodeUrl(userId int64) (url string, err error) {
	config, err := getConfig()
	if err != nil {
		return
	}
	token, err := createStateToken(userId)
	if err != nil {
		return
	}
	url = config.AuthCodeURL(token, oauth2.AccessTypeOffline)
	fmt.Println(url)
	return
}

type OauthForm struct {
	State string `form:"state"`
	Code  string `form:"code"`
	Scope string `form:"scope"`
}

func ExchangeCode(ctx context.Context, authProps OauthForm) (tok *oauth2.Token, err error) {
	config, err := getConfig()
	if err != nil {
		return
	}
	tok, err = config.Exchange(ctx, authProps.Code)
	return
}
