package fetchers

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/yayuyokitano/livefetcher/lib/core/htmlquerier"
	"github.com/yayuyokitano/livefetcher/lib/core/util"
	"golang.org/x/net/html"
)

type TimeHandler struct {
	YearQuerier      htmlquerier.Querier
	MonthQuerier     htmlquerier.Querier
	DayQuerier       htmlquerier.Querier
	OpenTimeQuerier  htmlquerier.Querier
	StartTimeQuerier htmlquerier.Querier

	IsYearInLive  bool
	IsMonthInLive bool
}

type TestInfo struct {
	NumberOfLives         int
	FirstLiveTitle        string
	FirstLiveArtists      []string
	FirstLivePrice        string
	FirstLivePriceEnglish string
	FirstLiveOpenTime     time.Time
	FirstLiveStartTime    time.Time
	FirstLiveURL          string
	KnownEmpty            bool
}

type Simple struct {
	BaseURL              string
	InitialURL           string
	LiveSelector         string
	MultiLiveDaySelector string
	ExpandedLiveSelector string

	ShortYearIterableURL string
	NextSelector         string

	TitleQuerier   htmlquerier.Querier
	ArtistsQuerier htmlquerier.Querier
	DetailQuerier  htmlquerier.Querier

	PriceQuerier htmlquerier.Querier

	DetailsLink         string
	DetailsLinkSelector string

	TimeHandler TimeHandler

	PrefectureName string
	AreaName       string

	VenueID string

	TestInfo TestInfo

	Lives     []util.Live
	isTesting bool
}

func (s *Simple) Fetch() (err error) {
	if s.InitialURL != "" && s.NextSelector != "" {
		err = s.iterateUsingNextLink()
		if err != nil {
			return
		}
	} else if s.ShortYearIterableURL != "" {
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
	l, err := s.fetchLives(n, initialURL)
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
	l, err := s.fetchLives(n, initialURL)
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
		var appL []util.Live
		appL, err = s.fetchLives(n, nextURL)
		if err != nil {
			break
		}
		if len(appL) == 0 {
			break
		}
		l = append(l, appL...)
	}
	s.Lives = l
	return
}

func (s *Simple) iterateUsingShortYear() (err error) {
	t := time.Now()
	year := t.Year() % 100
	month := int(t.Month())
	var l []util.Live
	for err == nil {
		var n *html.Node
		var newURL *url.URL
		newURL, err = url.Parse(fmt.Sprintf(s.ShortYearIterableURL, year, month))
		if err != nil {
			break
		}
		n, err = htmlquery.LoadURL(newURL.String())
		if err != nil {
			break
		}
		var appL []util.Live
		appL, err = s.fetchLives(n, newURL)
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

func (s *Simple) fetchLives(n *html.Node, overviewURL *url.URL) (l []util.Live, err error) {
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
		for i := 0; i < util.Min(10, len(overview)); i++ {
			wg.Add(1)
			go fetchLiveConcurrent(overviewURL, queue, s.ExpandedLiveSelector, &wg)
		}
		wg.Wait()
		for _, liveDetails := range res {
			lives = append(lives, LiveContext{
				n:   liveDetails.res,
				url: liveDetails.url,
			})
		}
	} else {
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

func (s *Simple) fetchDetails(live *html.Node, overviewURL *url.URL, year string, month string, day string) (l util.Live, err error) {
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

	l = util.Live{
		Title:        title,
		Artists:      artists,
		OpenTime:     open,
		StartTime:    start,
		Price:        strings.TrimSpace(price),
		PriceEnglish: strings.TrimSpace(util.EnglishPriceHandler(price)),
		Venue: util.LiveHouse{
			ID: s.VenueID,
			Area: util.Area{
				Prefecture: s.PrefectureName,
				Area:       s.AreaName,
			},
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
		now := time.Now()
		if month < int(now.Month()) {
			res = append(res, strconv.Itoa(time.Now().Year()+1))
		} else {
			res = append(res, strconv.Itoa(time.Now().Year()))
		}

	}

	year = res[0]
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
	month = res[0]
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
	day = res[0]
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
	} else {
		prices, err = s.DetailQuerier.Execute(n)
		if err != nil || len(prices) == 0 {
			return
		}
		price = util.FindPrice(prices[0])
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
	} else {
		arr, err = s.DetailQuerier.Execute(n)
		if err != nil || arr[0] == "" {
			open, err = util.ParseTime(date, "03:24")
			return
		}
		open, err = util.ParseTime(date, util.FindTime(arr[0], "open"))
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
	} else {
		arr, err = s.DetailQuerier.Execute(n)
		if err != nil || arr[0] == "" {
			start, err = util.ParseTime(date, "03:24")
			return
		}
		start, err = util.ParseTime(date, util.FindTime(arr[0], "start"))
	}
	return
}
