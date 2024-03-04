package fetchers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func getCurrentShortURL(shortYearIterableURL string) string {
	t := time.Now()
	year := t.Year() % 100
	month := int(t.Month())
	return strings.Split(fmt.Sprintf(shortYearIterableURL, year, month), "%!(EXTRA")[0]
}

func (s *Simple) Test(t *testing.T) {
	s.isTesting = true
	path := fmt.Sprintf("../../../test/%s/%s/%s.html", s.PrefectureName, s.AreaName, s.VenueID)
	n, err := htmlquery.LoadDoc(path)
	if err != nil {
		t.Error(err)
		return
	}

	if s.InitialURL != "" && s.NextSelector != "" {
		err = s.testHasNextURL(n)
		if err != nil {
			t.Error(err)
			return
		}
		err = s.testRemoteInitialNext()
		if err != nil {
			t.Error(err)
			return
		}
		err = s.testStaticLive(n, s.InitialURL)
		if err != nil {
			t.Error(err)
		}
		return
	}

	if s.ShortYearIterableURL != "" {
		url := getCurrentShortURL(s.ShortYearIterableURL)

		err = s.testRemoteShortYearIterable(url)
		if err != nil {
			t.Error(err)
			return
		}
		err = s.testStaticLive(n, url)
		if err != nil {
			t.Error(err)
		}
		return
	}

	err = s.testRemoteShortYearIterable(s.InitialURL)
	if err != nil {
		t.Error(err)
		return
	}
	err = s.testStaticLive(n, s.InitialURL)
	if err != nil {
		t.Error(err)
	}
}

func (s *Simple) testRemoteShortYearIterable(url string) (err error) {
	n, err := htmlquery.LoadURL(url)
	if err != nil {
		return
	}

	err = s.testNotEmpty(n, url)
	return
}

func (s *Simple) testRemoteInitialNext() (err error) {
	n, err := htmlquery.LoadURL(s.InitialURL)
	if err != nil {
		return
	}

	err = s.testHasNextURL(n)
	if err != nil {
		return
	}

	err = s.testNotEmpty(n, s.InitialURL)
	return
}

func (s *Simple) testHasNextURL(n *html.Node) (err error) {
	next, err := htmlquery.Query(n, s.NextSelector)
	if err != nil {
		return
	}
	if next == nil {
		err = errors.New("next link is nil")
		return
	}
	nextURL := htmlquery.SelectAttr(next, "href")
	if nextURL == "" {
		err = fmt.Errorf("next link is empty")
		return
	}
	return
}

func (s *Simple) testStaticLive(n *html.Node, path string) (err error) {
	var pathURL *url.URL
	pathURL, err = url.Parse(path)
	if err != nil {
		return
	}
	l, err := s.fetchLives(n, pathURL)
	if err != nil {
		return
	}
	if len(l) != s.TestInfo.NumberOfLives {
		err = fmt.Errorf("expected %d lives, got %d", s.TestInfo.NumberOfLives, len(l))
		return
	}
	firstLive := l[0]
	if !reflect.DeepEqual(firstLive.Artists, s.TestInfo.FirstLiveArtists) {
		var expected []byte
		expected, err = json.Marshal(s.TestInfo.FirstLiveArtists)
		if err != nil {
			expected = []byte(fmt.Sprintf("%v", s.TestInfo.FirstLiveArtists))
			err = nil
		}
		var actual []byte
		actual, err = json.Marshal(firstLive.Artists)
		if err != nil {
			actual = []byte(fmt.Sprintf("%v", firstLive.Artists))
			err = nil
		}
		err = fmt.Errorf("expected artists %v, got %v", string(expected), string(actual))
		return
	}
	if firstLive.Title != s.TestInfo.FirstLiveTitle {
		err = fmt.Errorf("expected title %s, got %s", s.TestInfo.FirstLiveTitle, firstLive.Title)
		return
	}
	if firstLive.Price != s.TestInfo.FirstLivePrice {
		err = fmt.Errorf("expected price %s, got %s", s.TestInfo.FirstLivePrice, firstLive.Price)
		return
	}
	if firstLive.PriceEnglish != s.TestInfo.FirstLivePriceEnglish {
		err = fmt.Errorf("expected english price %s, got %s", s.TestInfo.FirstLivePriceEnglish, firstLive.PriceEnglish)
		return
	}
	if firstLive.OpenTime.Unix() != s.TestInfo.FirstLiveOpenTime.Unix() {
		err = fmt.Errorf("expected opentime %s, got %s", s.TestInfo.FirstLiveOpenTime, firstLive.OpenTime)
		return
	}
	if firstLive.StartTime.Unix() != s.TestInfo.FirstLiveStartTime.Unix() {
		err = fmt.Errorf("expected starttime %s, got %s", s.TestInfo.FirstLiveStartTime, firstLive.StartTime)
		return
	}
	if s.InitialURL != "" && firstLive.URL != s.TestInfo.FirstLiveURL {
		err = fmt.Errorf("expected url %s, got %s", s.TestInfo.FirstLiveURL, firstLive.URL)
		return
	}
	if s.ShortYearIterableURL != "" && firstLive.URL != getCurrentShortURL(s.TestInfo.FirstLiveURL) {
		err = fmt.Errorf("expected url %s, got %s", getCurrentShortURL(s.TestInfo.FirstLiveURL), firstLive.URL)
		return
	}
	return
}

func (s *Simple) testNotEmpty(n *html.Node, path string) (err error) {
	if s.TestInfo.KnownEmpty {
		return
	}
	var pathURL *url.URL
	pathURL, err = url.Parse(path)
	if err != nil {
		return
	}
	l, err := s.fetchLives(n, pathURL)
	if err != nil {
		return
	}
	if len(l) == 0 {
		err = errors.New("no lives fetched")
		return
	}

	hasTitle := false
	for _, live := range l {
		if live.Title != "" {
			hasTitle = true
			break
		}
	}
	if !hasTitle {
		err = errors.New("no title fetched from any live")
		return
	}

	hasArtist := false
	for _, live := range l {
		if len(live.Artists) != 0 {
			hasArtist = true
			break
		}
	}
	if !hasArtist {
		err = errors.New("no artists fetched from any live")
		return
	}

	hasPrice := false
	for _, live := range l {
		if live.Price != "" {
			hasPrice = true
			break
		}
	}
	if !hasPrice {
		err = errors.New("no price fetched from any live")
		return
	}

	hasStartTime := false
	for _, live := range l {
		if !live.StartTime.IsZero() {
			hasStartTime = true
			break
		}
	}
	if !hasStartTime {
		err = errors.New("no start time fetched from any live")
		return
	}

	hasOpenTime := false
	for _, live := range l {
		if !live.OpenTime.IsZero() {
			hasOpenTime = true
			break
		}
	}
	if !hasOpenTime {
		err = errors.New("no open time fetched from any live")
		return
	}

	if l[0].Venue.ID != s.VenueID {
		err = fmt.Errorf("benueID is wrong, expected %s, got %s", s.VenueID, l[0].Venue.ID)
		return
	}

	allHaveURL := true
	for _, live := range l {
		if live.URL == "" {
			allHaveURL = false
			break
		}
	}
	if !allHaveURL {
		err = errors.New("url not fetched from all lives")
		return
	}

	return
}
