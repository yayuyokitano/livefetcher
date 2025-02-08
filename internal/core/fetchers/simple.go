// Package fetchers contains site fetchers/scrapers.
//
// simple.go contains the simple regular fetcher.
package fetchers

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/yayuyokitano/livefetcher/internal/core/htmlquerier"
	"github.com/yayuyokitano/livefetcher/internal/core/util"
	"github.com/yayuyokitano/livefetcher/internal/core/util/datastructures"
	"golang.org/x/net/html"
)

// TimeHandler is a struct that specifies some queriers and properties relating to getting the
// opening and start times for lives.
type TimeHandler struct {
	// YearQuerier is a querier that returns the year of the live.
	YearQuerier htmlquerier.Querier

	// MonthQuerier is a querier that returns the month of the live.
	MonthQuerier htmlquerier.Querier

	// DayQuerier is a querier that returns the day of the live.
	DayQuerier htmlquerier.Querier

	// OpenTimeQuerier is a querier that returns the open time of the live in format xx:xx
	//
	// The core handles hours >= 24, incrementing day and subtracting hours appropriately.
	// The core also will automatically remove any extra characters not part of a time.
	OpenTimeQuerier htmlquerier.Querier

	// OpenTimeQuerier is a querier that returns the start time of the live in format xx:xx
	//
	// The core handles hours >= 24, incrementing day and subtracting hours appropriately.
	// The core also will automatically remove any extra characters not part of a time.
	StartTimeQuerier htmlquerier.Querier

	// IsYearInLive specifies whether each live has their own year element,
	// or if there is a single shared element for all lives in a page.
	//
	// If IsYearInLive is true, YearQuerier will execute in the context of LiveSelector.
	// If IsYearInLive is false, YearQuerier will execute in the context of document.
	IsYearInLive bool

	// IsMonthInLive specifies whether each live has their own month element,
	// or if there is a single shared element for all lives in a page.
	//
	// If IsMonthInLive is true, MonthQuerier will execute in the context of LiveSelector.
	// If IsMonthInLive is false, MonthQuerier will execute in the context of document.
	IsMonthInLive bool
}

// TestInfo specifies some information relating to connector tests.
//
// Each connector should have a test document named after its connector ID in the test/ folder.
type TestInfo struct {
	// NumberOfLives specifies the expected number of lives in the test document.
	NumberOfLives int
	// FirstLiveTitle specifies the expected title of the first live in the test document.
	FirstLiveTitle string
	// FirstLiveArtists specifies an array of the expected artists of the first live in the test document.
	FirstLiveArtists []string
	// FirstLivePrice specifies the expected price of the first live in the test document.
	FirstLivePrice string
	// FirstLivePriceEnglish specifies the expected translated price of the first live in the test document.
	FirstLivePriceEnglish string
	// FirstLiveOpenTime specifies the expected opening timestamp of the first live in the test document.
	FirstLiveOpenTime time.Time
	// FirstLiveStartTime specifies the expected starting timestamp of the first live in the test document.
	FirstLiveStartTime time.Time
	// FirstLiveURL specifies the expected URL of the first live in the test document.
	FirstLiveURL string
	// KnownEmpty is a workaround property, specifying that we expect that one of the live entries in the live test will be empty.
	//
	// Set this to true if testIsEmpty fails if you confirm that an empty result for a live is correct.
	//
	// TODO: Find a better way to handle this.
	KnownEmpty bool
}

// Simple is the basic fetcher, which currently all fetchers base themselves off of.
type Simple struct {
	// BaseURL is the base URL of the live website.
	// Fetchers often do much href parsing, often requiring us to know the base url in advance.
	//
	// TODO: Remove this, we can infer this from whichever URL is being used for first fetch.
	BaseURL string

	// InitialURL specifies a starting URL.
	// For InitialURL to work, all lives must be in the same page, or NextSelector must be specified.
	//
	// If there are multiple pages, and a usable NextSelector isn't available, IterableURL must be used.
	InitialURL string

	// LiveSelector specifies an xpath selector for one live.
	// Livefetcher will query for all instances of this selector, and treat every match as a separate live.
	LiveSelector string

	// LiveHTMLFetcher specifies a function that returns an array of html nodes corresponding to lives
	// Do not use this unless absolutely necessary
	LiveHTMLFetcher func([]byte) ([]*html.Node, error)

	// MultiLiveDaySelector provides a selector for a more complicated case of multiple lives in same day.
	//
	// Some websites will group multiple lives occurring on the same day under one singular wrapper,
	// and not provide the day of the live inside both elements.
	// In this case, we need to get the date using wrapper element,
	// but get all other info inside each of the live elements.
	//
	// On such a site, Liveselector should be the selector for the entire day wrapper,
	// and MultiLiveDaySelector should be the selector for each of the individual lives.
	//
	// See the Loft fetcher in groups.go for an example of this,
	// with some examples of the relevant page layout on https://www.loft-prj.co.jp/schedule/loft/date/2024/04
	MultiLiveDaySelector string

	// ExpandedLiveSelector is a selector for an anchor element leading to the full live details of a given live.
	//
	// In some cases, all the info needed for a live is not available on the schedule page,
	// and you need to navigate to a separate page for every single live to get the correct details.
	//
	// In this case, use ExpandedLiveSelector.
	//
	// If ExpandedLiveSelector is specified:
	//
	// 1. LiveSelector is used to fetch all lives on schedule page.
	//
	// 2. ExpandedLiveSelector is used within the scope of each live gotten using LiveSelector.
	//
	// 3. href of ExpandedLiveSelector element is navigated to, and all live-context detail queriers executed within that page.
	ExpandedLiveSelector string

	// In some rare cases ExpandedLiveSelector might lead to an article containing multiple lives.
	//
	// ExpandedLiveGroupSelector returns all those individual lives for further use.
	ExpandedLiveGroupSelector string

	// ShortYearIterableURL is a URL with two %d formatters specifying year and month in that order.
	//
	// year and motnh are given without leading zero, if leading zero is needed, provide this yourself using %02d.
	//
	// TODO: either make LongYearIterableURL or expand this to work in both cases. For now just use 20%02d in this case.
	ShortYearIterableURL string

	// ShortYearReverseIterableURL is the same as ShortYearIterableURL, except the order of the format strings is changed, so month is before year.
	ShortYearReverseIterableURL string
	// NextSelector is the selector of a link to the next page of schedule, showing newer lives than the current page.
	//
	// Livefetcher will follow the href of the element specified by NextSelector to get more lives, until no more lives are found.
	//
	// Must be specified along with an InitialURL.
	NextSelector string

	// TitleQuerier is a Querier that returns the title of the live.
	TitleQuerier htmlquerier.Querier
	// ArtistsQuerier is a Querier that returns an array of the artists of the live.
	ArtistsQuerier htmlquerier.Querier
	// DetailQuerier is a Querier that will return an unstructured blob of text,
	// which can be used to replace ArtistsQuerier, PriceQuerier, OpenTimeQuerier, and/or StartTimeQuerier.
	//
	// DetailQuerier is significantly less accurate,
	// and should only be used if the above queriers cannot be reliably created,
	// but can often make decent guesses.
	//
	// DetailQuerier will be overridden by the above queriers,
	// and you can choose to for instance only specify PriceQuerier and DetailQuerier,
	// which will cause PriceQuerier to be used for fetching price,
	// and DetailQuerier to be used for fetching artists, open time, and start time.
	//
	// Avoid using this if possible.
	DetailQuerier htmlquerier.Querier

	// PriceQuerier is a querier that returns the price of the live, including any details about the price.
	PriceQuerier htmlquerier.Querier

	// DetailsLink, if specified, will be the link for all lives returned by connector.
	// This is only useful if lives have no individual links, AND you are fetching from some hidden API.
	DetailsLink string
	// DetailsLinkSelector is the selector within a live for a link to details about the live.
	//
	// This will set the link for each live to the href of the element of the DetailsLinkSelector.
	//
	// Note that this does not need to be used if ExpandedLiveSelector is used.
	DetailsLinkSelector string

	// TimeHandler is a TimeHandler struct used to fetch time details about a live.
	// See TimeHandler documentation for details.
	TimeHandler TimeHandler

	// PrefectureName is the prefecture name for the connector.
	// These are standardized, and you must use the same as all other connectors within same prefecture, CASE SENSITIVE!
	//
	// If a new prefecture is added, locale must also be added to internal/i18n/locales toml files as well.
	PrefectureName string
	// AreaName is the area name for the connector.
	// These are standardized, and you must use the same as all other connectors within same area, CASE SENSITIVE!
	//
	// If a new area is added, locale must also be added to internal/i18n/locales toml files as well.
	//
	// Multiple prefectures may have identically named areas, and they will be treated as entirely separate.
	AreaName string

	// VenueID is the ID of the venue.
	// This must be globally unique.
	//
	// Do not change the ID of a venue unless there is a VERY strong reason to do so.
	//
	// A venue renaming is in itself not reason to change VenueID, only change locales files in this case.
	VenueID string

	// Longitude is the east/west longitude coordinate of livehouse, -180/180
	Longitude float64

	// Latitude is the north/south latitude coordinate of livehouse, -90/90
	Latitude float64

	// TestInfo is a struct specifying expected values for some tests for the connector.
	// See TestInfo documentation for details.
	TestInfo TestInfo

	// Lives is used internally in the core for processing lives.
	// Do not use this in connectors.
	Lives []datastructures.Live
	// isTesting is used internally in the core for processing lives.
	// Do not use this in connectors.
	isTesting bool
}

func (s *Simple) Fetch() (err error) {
	if s.InitialURL != "" && s.NextSelector != "" {
		err = s.iterateUsingNextLink()
		if err != nil {
			return
		}
	} else if s.ShortYearIterableURL != "" || s.ShortYearReverseIterableURL != "" {
		err = s.iterateUsingShortYear()
		if err != nil {
			return
		}
	} else {
		err = s.fetchSingle()
		if err != nil {
			return
		}
	}
	return
}

func (s *Simple) fetchSingle() (err error) {
	n, err := htmlquery.LoadURL(s.InitialURL)
	if err != nil {
		return
	}
	initialURL, err := url.Parse(s.InitialURL)
	if err != nil {
		return
	}
	l, err := s.fetchLives(n, initialURL, nil)
	if err != nil {
		return
	}
	s.Lives = l
	return
}

func (s *Simple) iterateUsingNextLink() (err error) {
	base, err := url.Parse(s.BaseURL)
	if err != nil {
		return
	}

	n, err := htmlquery.LoadURL(s.InitialURL)
	if err != nil {
		return
	}
	initialURL, err := url.Parse(s.InitialURL)
	if err != nil {
		return
	}
	l, err := s.fetchLives(n, initialURL, nil)
	if err != nil {
		return
	}
	prevURL := initialURL

	for next, err2 := htmlquery.Query(n, s.NextSelector); next != nil && err2 == nil; next, err2 = htmlquery.Query(n, s.NextSelector) {
		var nextURL *url.URL
		nextURL, err = base.Parse(htmlquery.SelectAttr(next, "href"))
		if err != nil {
			break
		}
		if strings.HasPrefix(prevURL.String(), nextURL.String()) || strings.HasPrefix(nextURL.String(), prevURL.String()) {
			break
		}
		prevURL = nextURL
		n, err = htmlquery.LoadURL(nextURL.String())
		if err != nil {
			break
		}
		var appL []datastructures.Live
		appL, err = s.fetchLives(n, nextURL, nil)
		if err != nil {
			break
		}
		if len(appL) == 0 {
			break
		}
		// also leave if all lives are from a previous year
		hasCurrentYearLive := false
		for _, live := range appL {
			if live.OpenTime.Year() >= time.Now().Year() {
				hasCurrentYearLive = true
				break
			}
		}
		if !hasCurrentYearLive {
			break
		}

		l = append(l, appL...)
	}
	s.Lives = l
	return
}

func (s *Simple) getNewIterableURL(year, month int) (*url.URL, error) {
	if s.ShortYearIterableURL != "" {
		return url.Parse(fmt.Sprintf(s.ShortYearIterableURL, year, month))
	} else {
		return url.Parse(fmt.Sprintf(s.ShortYearReverseIterableURL, month, year))
	}
}

func (s *Simple) iterateUsingShortYear() (err error) {
	t := time.Now()
	year := t.Year() % 100
	month := int(t.Month())
	var l []datastructures.Live
	for err == nil {
		var n *html.Node
		var newURL *url.URL
		newURL, err = s.getNewIterableURL(year, month)
		if err != nil {
			break
		}
		n, err = htmlquery.LoadURL(newURL.String())
		if err != nil {
			break
		}
		var appL []datastructures.Live
		appL, err = s.fetchLives(n, newURL, nil)
		if err != nil {
			break
		}
		if len(appL) == 0 {
			break
		}
		l = append(l, appL...)
		month++
		if month > 12 {
			month = 1
			year++
		}
	}
	s.Lives = l
	return
}

type LiveQueueElement struct {
	live *html.Node
	res  *html.Node
	url  *url.URL
}

func fetchLiveConcurrent(baseURL *url.URL, queue chan *LiveQueueElement, expandedLiveSelector string, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range queue {
		liveAnchor, err := htmlquery.Query(job.live, expandedLiveSelector)
		if err != nil || liveAnchor == nil {
			return
		}

		var liveDetails *html.Node
		var url *url.URL
		url, err = baseURL.Parse(htmlquery.SelectAttr(liveAnchor, "href"))
		if err != nil {
			return
		}
		job.url = url
		liveDetails, err = htmlquery.LoadURL(url.String())
		if err != nil || liveDetails == nil {
			return
		}
		job.res = liveDetails
	}
}

type LiveContext struct {
	n   *html.Node
	url *url.URL
}

func (s *Simple) fetchLives(n *html.Node, overviewURL *url.URL, testDocument []byte) (l []datastructures.Live, err error) {
	var lives []LiveContext
	if s.ExpandedLiveSelector != "" {
		var overview []*html.Node
		overview, err = htmlquery.QueryAll(n, s.LiveSelector)
		if err != nil {
			return
		}
		if overview == nil {
			err = fmt.Errorf("fetching overview returned nil from %s", overviewURL)
			return
		}

		// fetch lives concurrently, limit it to 10 at a time, make sure responses are in the correct order
		var wg sync.WaitGroup
		queue := make(chan *LiveQueueElement, len(overview))
		var res []*LiveQueueElement
		for _, live := range overview {
			job := &LiveQueueElement{live: live}
			res = append(res, job)
			queue <- job
		}
		close(queue)
		for i := 0; i < min(10, len(overview)); i++ {
			wg.Add(1)
			go fetchLiveConcurrent(overviewURL, queue, s.ExpandedLiveSelector, &wg)
		}
		wg.Wait()
		for _, liveDetails := range res {
			if s.ExpandedLiveGroupSelector == "" {
				lives = append(lives, LiveContext{
					n:   liveDetails.res,
					url: liveDetails.url,
				})
			} else {
				var liveNodes []*html.Node
				liveNodes, err = htmlquery.QueryAll(liveDetails.res, s.ExpandedLiveGroupSelector)
				if err != nil || len(liveNodes) == 0 {
					continue
				}
				for _, liveNode := range liveNodes {
					lives = append(lives, LiveContext{
						n:   liveNode,
						url: liveDetails.url,
					})
				}
			}
		}
	} else if s.LiveSelector != "" {
		var rawLives []*html.Node
		rawLives, err = htmlquery.QueryAll(n, s.LiveSelector)
		if err != nil {
			return
		}
		if rawLives == nil {
			err = errors.New("raw live query returned nil")
			return
		}
		for _, live := range rawLives {
			lives = append(lives, LiveContext{
				n:   live,
				url: overviewURL,
			})
		}
	} else if s.LiveHTMLFetcher != nil {
		var rawLives []*html.Node
		rawLives, err = s.LiveHTMLFetcher(testDocument)
		if err != nil {
			return
		}
		if rawLives == nil {
			err = errors.New("raw live query returned nil")
			return
		}
		for _, live := range rawLives {
			lives = append(lives, LiveContext{
				n:   live,
				url: overviewURL,
			})
		}
	}

	if len(lives) == 0 {
		return
	}

	var year string
	if !s.TimeHandler.IsYearInLive {
		year, err = s.getYear(n)
		if err != nil {
			return
		}
	}

	var month string
	if !s.TimeHandler.IsMonthInLive {
		month, err = s.getMonth(n)
		if err != nil {
			return
		}
	}

	var day string

	for _, live := range lives {

		if s.TimeHandler.IsYearInLive {
			year, err = s.getYear(live.n)
			if err != nil {
				fmt.Println(err)
				err = nil
				continue
			}
		}

		prevMonth := month
		if s.TimeHandler.IsMonthInLive {
			month, err = s.getMonth(live.n)
			if err != nil || month == "" {
				err = nil
				month = prevMonth
			}
		}

		prevDay := day
		day, err = s.getDay(live.n)
		if err != nil || day == "" {
			err = nil
			day = prevDay
		}

		timeCutoff := time.Now().AddDate(0, -1, 0)

		if s.MultiLiveDaySelector == "" {
			appL, err := s.fetchDetails(live.n, live.url, year, month, day)
			if err != nil {
				fmt.Println(err)
				err = nil
				continue
			}
			if !s.isTesting && appL.StartTime.Before(timeCutoff) {
				continue
			}
			l = append(l, appL)
		} else {
			dailyLives, err := htmlquery.QueryAll(live.n, s.MultiLiveDaySelector)
			if err != nil || dailyLives == nil {
				continue
			}
			for _, dailyLive := range dailyLives {
				appL, err := s.fetchDetails(dailyLive, live.url, year, month, day)
				if err != nil {
					fmt.Println(err)
					err = nil
					continue
				}
				if !s.isTesting && appL.StartTime.Before(timeCutoff) {
					continue
				}
				l = append(l, appL)
			}
		}
	}
	return
}

func (s *Simple) fetchDetails(live *html.Node, overviewURL *url.URL, year string, month string, day string) (l datastructures.Live, err error) {
	date := fmt.Sprintf("%s-%s-%s", year, month, day)

	var open time.Time
	var start time.Time
	open, err = s.getOpenTime(live, date)
	if err != nil {
		return
	}

	start, err = s.getStartTime(live, date)
	if err != nil {
		return
	}
	if open.Hour() == 3 && open.Minute() == 24 && !(start.Hour() == 3 && start.Minute() == 24) {
		open = start
	}
	if start.Hour() == 3 && start.Minute() == 24 && !(open.Hour() == 3 && open.Minute() == 24) {
		start = open
	}

	var price string
	price, err = s.getPrice(live)
	if err != nil {
		return
	}

	var title string
	title, err = s.getTitle(live)
	if err != nil {
		return
	}

	var artists []string
	artists, err = s.FetchArtists(live)
	if err != nil {
		return
	}

	detailsURL := overviewURL.String()
	if s.DetailsLinkSelector != "" {
		var detailsLink *html.Node
		detailsLink, err = htmlquery.Query(live, s.DetailsLinkSelector)
		if err == nil && detailsLink != nil {
			var newURL *url.URL
			newURL, err = overviewURL.Parse(htmlquery.SelectAttr(detailsLink, "href"))
			if err == nil && newURL != nil {
				detailsURL = newURL.String()
			}
		}
	}
	if s.DetailsLink != "" {
		var newURL *url.URL
		newURL, err = overviewURL.Parse(s.DetailsLink)
		if err == nil && newURL != nil {
			detailsURL = newURL.String()
		}
	}

	l = datastructures.Live{
		Title:        title,
		Artists:      artists,
		OpenTime:     open,
		StartTime:    start,
		Price:        strings.TrimSpace(price),
		PriceEnglish: strings.TrimSpace(util.EnglishPriceHandler(price)),
		Venue: datastructures.LiveHouse{
			ID: s.VenueID,
			Area: datastructures.Area{
				Prefecture: s.PrefectureName,
				Area:       s.AreaName,
			},
			Latitude:  s.Latitude,
			Longitude: s.Longitude,
		},
		URL: detailsURL,
	}
	return
}

func (s *Simple) FetchArtists(n *html.Node) (a []string, err error) {
	if s.ArtistsQuerier.Initialized {
		a, err = s.ArtistsQuerier.Execute(n)
		a = util.ProcessArtists(a)
	} else {
		a, err = s.DetailQuerier.Execute(n)
		if err != nil || len(a) == 0 {
			return
		}
		a = util.ProcessArtists(strings.Split(a[0], "\n"))
	}
	return
}

func (s *Simple) getTitle(n *html.Node) (title string, err error) {
	var a []string
	a, err = s.TitleQuerier.Execute(n)
	if err != nil {
		return
	}
	title = a[0]
	return
}

func isolateFirstNumber(old string) (isolated string, err error) {
	re, err := regexp.Compile(`\d+`)
	if err != nil {
		return
	}
	isolated = re.FindString(old)
	return
}

func (s *Simple) getYear(n *html.Node) (year string, err error) {
	var res []string
	if s.TimeHandler.YearQuerier.Initialized {
		res, err = s.TimeHandler.YearQuerier.Execute(n)
		if err != nil {
			return
		}
	} else {
		var month int
		var monthstr string
		monthstr, err = s.getMonth(n)
		if err != nil {
			return
		}
		month, err = strconv.Atoi(monthstr)
		if err != nil {
			return
		}
		res = append(res, strconv.Itoa(util.GetRelevantYear(month)))
	}

	year, err = isolateFirstNumber(res[0])
	if err != nil {
		return
	}
	if len(year) == 2 {
		year = "20" + year
	}
	return
}

func (s *Simple) getMonth(n *html.Node) (month string, err error) {
	res, err := s.TimeHandler.MonthQuerier.Execute(n)
	if err != nil {
		return
	}
	month, err = isolateFirstNumber(res[0])
	if err != nil {
		return
	}

	if len(month) == 1 {
		month = "0" + month
	}
	return
}

func (s *Simple) getDay(n *html.Node) (day string, err error) {
	res, err := s.TimeHandler.DayQuerier.Execute(n)
	if err != nil {
		return
	}
	day, err = isolateFirstNumber(res[0])
	if err != nil {
		return
	}
	if len(day) == 1 {
		day = "0" + day
	}
	return
}

func (s *Simple) getPrice(n *html.Node) (price string, err error) {
	var prices []string
	if s.PriceQuerier.Initialized {
		prices, err = s.PriceQuerier.Execute(n)
		if err != nil || len(prices) == 0 {
			return
		}
		price = prices[0]
	} else if s.DetailQuerier.Initialized {
		prices, err = s.DetailQuerier.Execute(n)
		if err != nil || len(prices) == 0 {
			return
		}
		price = util.FindPrice(prices)
	} else {
		price = "このライブハウスのイベントの値段にアクセスできません。ライブのリンクをチェックしてください。"
	}
	return
}

func (s *Simple) getOpenTime(n *html.Node, date string) (open time.Time, err error) {
	var arr []string
	if s.TimeHandler.OpenTimeQuerier.Initialized {
		arr, err = s.TimeHandler.OpenTimeQuerier.Execute(n)
		if err != nil || arr[0] == "" {
			open, err = util.ParseTime(date, "03:24")
			return
		}
		open, err = util.ParseTime(date, arr[0])
	} else if s.DetailQuerier.Initialized {
		arr, err = s.DetailQuerier.Execute(n)
		if err != nil || arr[0] == "" {
			open, err = util.ParseTime(date, "03:24")
			return
		}
		open, err = util.ParseTime(date, util.FindTime(strings.Join(arr, ""), "open"))
	} else {
		open, err = util.ParseTime(date, "03:24")
	}
	return
}

func (s *Simple) getStartTime(n *html.Node, date string) (start time.Time, err error) {
	var arr []string
	if s.TimeHandler.StartTimeQuerier.Initialized {
		arr, err = s.TimeHandler.StartTimeQuerier.Execute(n)
		if err != nil || arr[0] == "" {
			start, err = util.ParseTime(date, "03:24")
			return
		}
		start, err = util.ParseTime(date, arr[0])
	} else if s.DetailQuerier.Initialized {
		arr, err = s.DetailQuerier.Execute(n)
		if err != nil || arr[0] == "" {
			start, err = util.ParseTime(date, "03:24")
			return
		}
		start, err = util.ParseTime(date, util.FindTime(strings.Join(arr, ""), "start"))
	} else {
		start, err = util.ParseTime(date, "03:24")
	}
	return
}
