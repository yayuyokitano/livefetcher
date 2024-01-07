package fetchers

import (
	"fmt"

	"github.com/yayuyokitano/livefetcher/internal/core/htmlquerier"
)

func CreateWWWFetcher(
	articleCondition string,
	venue string,
	testInfo TestInfo,
) Simple {
	return Simple{
		BaseURL:              "https://www-shibuya.jp/",
		InitialURL:           "https://www-shibuya.jp/schedule/",
		LiveSelector:         fmt.Sprintf("//div[@id='eventList']//article[%s]", articleCondition),
		NextSelector:         "//ul[@class='navigation']/li[@class='next']/a",
		ExpandedLiveSelector: "//a[@class='pageLink']",
		TitleQuerier:         *htmlquerier.Q("//header//div[@class='event']/p"),
		ArtistsQuerier:       *htmlquerier.QAll("//dt[text()='LINE UP']/following-sibling::dd/a | //dt[text()='LINE UP']/following-sibling::dd/text()").ReplaceAllRegex(`^｜.*｜$`, "").ReplaceAllRegex("^◾️WWW(.*)", ""),
		PriceQuerier:         *htmlquerier.Q("//dt[text()='ADV./DOOR']/following-sibling::dd"),

		TimeHandler: TimeHandler{
			YearQuerier:      *htmlquerier.Q("//header//div[@class='date']/p[@class='year']"),
			MonthQuerier:     *htmlquerier.Q("//header//div[@class='date']/p[@class='month']/text()"),
			DayQuerier:       *htmlquerier.Q("//header//div[@class='date']/p[@class='day']"),
			OpenTimeQuerier:  *htmlquerier.Q("//dd[@class='openstart']").Before("/").Before("｜"),
			StartTimeQuerier: *htmlquerier.Q("//dd[@class='openstart']").After("/").After("｜"),

			IsYearInLive:  true,
			IsMonthInLive: true,
		},

		PrefectureName: "tokyo",
		AreaName:       "shibuya",
		VenueID:        venue,

		TestInfo: testInfo,
	}
}

func CreateOFetcher(
	baseURL string,
	initialURL string,
	prefecture string,
	area string,
	venue string,
	testInfo TestInfo,
) Simple {
	return Simple{
		BaseURL:              baseURL,
		InitialURL:           initialURL,
		LiveSelector:         "//div[@class='p-schedule__list']/div[contains(@class, 'p-scheduled-card')]",
		NextSelector:         "//div[contains(@class, 'p-schedule__nav-item--next')]/a",
		ExpandedLiveSelector: "//a",
		TitleQuerier:         *htmlquerier.Q("//span[@class='p-schedule-detail__title-main']"),
		ArtistsQuerier:       *htmlquerier.Q("//ul[@class='p-schedule-detail__artist']").Split("\n"),
		PriceQuerier:         *htmlquerier.Q("//div[@class='p-schedule-detail__blcok'][2]").After("OPEN").After("START").ReplaceAllRegex(`(\s+)|(\d{2}:\d{2})`, " "),

		TimeHandler: TimeHandler{
			MonthQuerier:     *htmlquerier.Q("//span[@class='p-schedule-detail__date-item']").Before("/"),
			DayQuerier:       *htmlquerier.Q("//span[@class='p-schedule-detail__date-item']").After("/"),
			OpenTimeQuerier:  *htmlquerier.Q("//div[@class='p-schedule-detail__dl'][1]//div[@class='p-schedule-detail__dd']"),
			StartTimeQuerier: *htmlquerier.Q("//div[@class='p-schedule-detail__dl'][2]//div[@class='p-schedule-detail__dd']"),

			IsYearInLive:  true,
			IsMonthInLive: true,
		},

		PrefectureName: prefecture,
		AreaName:       area,
		VenueID:        venue,

		TestInfo: testInfo,
	}
}

func CreateEggmanFetcher(
	baseURL string,
	shortYearIterableURL string,
	prefecture string,
	area string,
	venue string,
	testInfo TestInfo,
) Simple {
	return Simple{
		BaseURL:              baseURL,
		ShortYearIterableURL: shortYearIterableURL,
		LiveSelector:         "//article[@class='scheduleList']",
		DetailsLinkSelector:  "//a",
		TitleQuerier:         *htmlquerier.Q("//h1"),
		ArtistsQuerier:       *htmlquerier.Q("//div[@class='act']").SplitIgnoreWithin("[\n/]", '(', ')'),
		PriceQuerier:         *htmlquerier.Q("//div[@class='scheListBody']/ul").After("START ").ReplaceAllRegex(`(\s+)|(\d{2}:\d{2})`, " "),

		TimeHandler: TimeHandler{
			YearQuerier:      *htmlquerier.Q("//div[contains(@class, 'monthHeader')]/h1").Before("."),
			MonthQuerier:     *htmlquerier.Q("//div[contains(@class, 'monthHeader')]/h1").After("."),
			DayQuerier:       *htmlquerier.Q("//div[@class='scheListHeader']/time/strong"),
			OpenTimeQuerier:  *htmlquerier.Q("//div[@class='scheListBody']//li[1]"),
			StartTimeQuerier: *htmlquerier.Q("//div[@class='scheListBody']//li[2]"),
		},

		PrefectureName: prefecture,
		AreaName:       area,
		VenueID:        venue,

		TestInfo: testInfo,
	}
}

func CreateToosFetcher(
	baseURL string,
	shortYearIterableURL string,
	prefecture string,
	area string,
	venue string,
	testInfo TestInfo,
) Simple {
	return Simple{
		BaseURL:              baseURL,
		ShortYearIterableURL: shortYearIterableURL,
		LiveSelector:         "//article[contains(@class, 'type-event')]",
		ExpandedLiveSelector: "//a",
		TitleQuerier:         *htmlquerier.Q("//div[@class='main_title']"),
		ArtistsQuerier:       *htmlquerier.Q("//div[@class='box']/div[contains(@class, 'title')][text()='ACT']/following-sibling::div").Split("\n"),
		PriceQuerier:         *htmlquerier.Q("//div[@class='sub_detail']").SplitIndex("\n", 1),

		TimeHandler: TimeHandler{
			YearQuerier:      *htmlquerier.Q("//div[@class='date']").Before("年").HalfWidth(),
			MonthQuerier:     *htmlquerier.Q("//div[@class='date']").Before("月").After("年").HalfWidth(),
			DayQuerier:       *htmlquerier.Q("//div[@class='date']").Before("日").After("月").HalfWidth(),
			OpenTimeQuerier:  *htmlquerier.Q("//div[@class='sub_detail']").Before("\n").After("OPEN").Before("／"),
			StartTimeQuerier: *htmlquerier.Q("//div[@class='sub_detail']").Before("\n").After("START").Before("／"),
			IsMonthInLive:    true,
			IsYearInLive:     true,
		},

		PrefectureName: prefecture,
		AreaName:       area,
		VenueID:        venue,

		TestInfo: testInfo,
	}
}

func CreateChikamichiFetcher(
	baseURL string,
	initialURL string,
	prefecture string,
	area string,
	venue string,
	testInfo TestInfo,
) Simple {
	return Simple{
		BaseURL:              baseURL,
		InitialURL:           initialURL,
		LiveSelector:         "//article",
		ExpandedLiveSelector: "//a",
		TitleQuerier:         *htmlquerier.Q("//h3"),
		ArtistsQuerier:       *htmlquerier.Q("//dt[text()='LINE UP']/following-sibling::dd").Split(" / "),
		PriceQuerier:         *htmlquerier.Q("//dt[text()='PRICE']/following-sibling::dd"),

		TimeHandler: TimeHandler{
			YearQuerier:      *htmlquerier.Q("//span[@class='lh-1']/span[1]").Before("."),
			MonthQuerier:     *htmlquerier.Q("//span[@class='lh-1']/span[1]").SplitIndex(".", 1),
			DayQuerier:       *htmlquerier.Q("//span[@class='lh-1']/span[1]").SplitIndex(".", 2),
			OpenTimeQuerier:  *htmlquerier.Q("//dt[text()='OPEN / START']/following-sibling::dd").Before("/"),
			StartTimeQuerier: *htmlquerier.Q("//dt[text()='OPEN / START']/following-sibling::dd").After("/"),
			IsMonthInLive:    true,
			IsYearInLive:     true,
		},

		PrefectureName: prefecture,
		AreaName:       area,
		VenueID:        venue,

		TestInfo: testInfo,
	}
}

func CreateDaisyBarFetcher(
	baseUrl string,
	initialURL string,
	prefecture string,
	area string,
	venue string,
	yearColor string,
	testInfo TestInfo,
) Simple {
	return Simple{
		BaseURL:        baseUrl,
		InitialURL:     initialURL,
		LiveSelector:   "//article[@class='single-article']",
		TitleQuerier:   *htmlquerier.Q("//p[contains(@class, 'h4')]"),
		ArtistsQuerier: *htmlquerier.Q("//div[contains(@class, 'artist')]").Split("／").Before("【ONE MAN】"),
		PriceQuerier:   *htmlquerier.Q("//div[contains(@class, 'liveinfo')]/p/span[2]"),

		TimeHandler: TimeHandler{
			YearQuerier:      *htmlquerier.Q(fmt.Sprintf("//p[contains(@class, 'h4 %s')]", yearColor)),
			MonthQuerier:     *htmlquerier.Q("//div[contains(@class, 'h3')]").Before("/"),
			DayQuerier:       *htmlquerier.Q("//div[contains(@class, 'h3')]").After("/").Before("("),
			OpenTimeQuerier:  *htmlquerier.Q("//div[contains(@class, 'liveinfo')]/p/span[1]").After("OPEN").Before("START"),
			StartTimeQuerier: *htmlquerier.Q("//div[contains(@class, 'liveinfo')]/p/span[1]").After("START"),
			IsMonthInLive:    true,
		},

		PrefectureName: prefecture,
		AreaName:       area,
		VenueID:        venue,

		TestInfo: testInfo,
	}

}

func CreateLoftFetcher(
	baseUrl string,
	initialURL string,
	prefecture string,
	area string,
	venue string,
	testInfo TestInfo,
) Simple {
	return Simple{
		BaseURL:              baseUrl,
		InitialURL:           initialURL,
		NextSelector:         "//span[@id='next_btn']/a",
		LiveSelector:         "//table[contains(@class, 'timetable')]/tbody/tr",
		TitleQuerier:         *htmlquerier.Q("//h3"),
		MultiLiveDaySelector: "//div[contains(@class, 'event clearfix')]",
		ArtistsQuerier:       *htmlquerier.Q("//p[contains(@class, 'month_content')]").SplitIgnoreWithin("(\n)|( / )", '（', '）'),
		PriceQuerier:         *htmlquerier.Q("//p[contains(@class, 'ticket')]"),
		DetailsLinkSelector:  "//p[contains(@class, 'detail_mono')]/a",

		TimeHandler: TimeHandler{
			YearQuerier:      *htmlquerier.Q("//div[@id='month_top']/h2").Before("年"),
			MonthQuerier:     *htmlquerier.Q("//div[@id='month_top']/h2").After("年").Before("月"),
			DayQuerier:       *htmlquerier.Q("//th[contains(@class, 'day')]/text()[1]"),
			OpenTimeQuerier:  *htmlquerier.Q("//p[contains(@class, 'time_text')]").Before(" / "),
			StartTimeQuerier: *htmlquerier.Q("//p[contains(@class, 'time_text')]").After(" / "),
		},

		PrefectureName: prefecture,
		AreaName:       area,
		VenueID:        venue,
		TestInfo:       testInfo,
	}
}
