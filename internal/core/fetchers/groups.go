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
		ShortYearIterableURL: "https://www-shibuya.jp/schedule/20%d%02d.php",
		LiveSelector:         fmt.Sprintf("//div[@id='eventList']//article[%s]", articleCondition),
		ExpandedLiveSelector: "//a[@class='pageLink']",
		TitleQuerier:         *htmlquerier.QAll(`//*[self::h2[@class='title'] or self::p[@class='exp']]`).FilterTitle(" / ", 1),
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
		Latitude:       35.661537,
		Longitude:      139.698734,

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
	latitude float64,
	longitude float64,
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
		Latitude:       latitude,
		Longitude:      longitude,

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
		Latitude:       35.664363,
		Longitude:      139.699203,

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
		Latitude:       35.656838,
		Longitude:      139.667547,

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
		Latitude:       35.664563,
		Longitude:      139.666938,

		TestInfo: testInfo,
	}
}

func CreateDaisyBarFetcher(
	baseUrl string,
	shortYearIterableURL string,
	prefecture string,
	area string,
	venue string,
	testInfo TestInfo,
) Simple {
	return Simple{
		BaseURL:              baseUrl,
		ShortYearIterableURL: shortYearIterableURL,
		LiveSelector:         "//article[@class='schedule-ticket']",
		TitleQuerier:         *htmlquerier.Q("//h2"),
		ArtistsQuerier:       *htmlquerier.Q("//p[@class='artist']").Split("／").Split(" ゲスト:").Before("【ONE MAN】"),
		PriceQuerier:         *htmlquerier.QAll("//div[@class='schedule-list-content_mid']/div[position()=3 or position()=4]").Join(" / "),

		TimeHandler: TimeHandler{
			YearQuerier:      *htmlquerier.Q("//span[@class='year']"),
			MonthQuerier:     *htmlquerier.Q("//span[@class='month']"),
			DayQuerier:       *htmlquerier.Q("//span[@class='day']"),
			OpenTimeQuerier:  *htmlquerier.Q("//div[@class='schedule-list-content_mid']/div[1]"),
			StartTimeQuerier: *htmlquerier.Q("//div[@class='schedule-list-content_mid']/div[2]"),
		},

		PrefectureName: prefecture,
		AreaName:       area,
		VenueID:        venue,
		Latitude:       35.659562,
		Longitude:      139.668063,

		TestInfo: testInfo,
	}
}

func CreateOldDaisyBarFetcher(
	baseUrl string,
	shortYearIterableURL string,
	prefecture string,
	area string,
	venue string,
	yearColor string,
	testInfo TestInfo,
) Simple {
	return Simple{
		BaseURL:              baseUrl,
		ShortYearIterableURL: shortYearIterableURL,
		LiveSelector:         "//article[@class='single-article']",
		TitleQuerier:         *htmlquerier.Q("//p[contains(@class, 'h4')]"),
		ArtistsQuerier:       *htmlquerier.Q("//div[contains(@class, 'artist')]").Split("／").Before("【ONE MAN】"),
		PriceQuerier:         *htmlquerier.Q("//div[contains(@class, 'liveinfo')]/p/span[2]"),

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
		Latitude:       35.659562,
		Longitude:      139.668063,

		TestInfo: testInfo,
	}

}

func CreateBassOnTopFetcher(
	baseURL string,
	shortYearIterableURL string,
	prefecture string,
	area string,
	venueID string,
	testInfo TestInfo,
	latitude float64,
	longitude float64,
) Simple {
	return Simple{
		BaseURL:              baseURL,
		ShortYearIterableURL: shortYearIterableURL,
		LiveSelector:         "//div[@class='container scheduleList']/ul/li[.//h1/text()!='HALL RENTAL']",
		ExpandedLiveSelector: "//a[@class='btnStyle01']",
		TitleQuerier:         *htmlquerier.Q("//div[@class='scheduleCnt']/h1").CutWrapper("【", "】"),
		ArtistsQuerier:       *htmlquerier.Q("//dl[@class='act']//span").SplitIgnoreWithin("(/|(【MC】))", '(', ')'),
		PriceQuerier:         *htmlquerier.Q("//dl[@class='price']/dd"),

		TimeHandler: TimeHandler{
			YearQuerier:      *htmlquerier.Q("//p[@class='day']"),
			MonthQuerier:     *htmlquerier.Q("//p[@class='day']").After("."),
			DayQuerier:       *htmlquerier.Q("//p[@class='day']").After(".").After("."),
			OpenTimeQuerier:  *htmlquerier.Q("//dl[@class='openTime']/dd"),
			StartTimeQuerier: *htmlquerier.Q("//dl[@class='openTime']/dd").After("/"),

			IsYearInLive:  true,
			IsMonthInLive: true,
		},

		PrefectureName: prefecture,
		AreaName:       area,
		VenueID:        venueID,
		Latitude:       latitude,
		Longitude:      longitude,

		TestInfo: TestInfo{
			NumberOfLives:         testInfo.NumberOfLives,
			FirstLiveTitle:        testInfo.FirstLiveTitle,
			FirstLiveArtists:      testInfo.FirstLiveArtists,
			FirstLivePrice:        testInfo.FirstLivePrice,
			FirstLivePriceEnglish: testInfo.FirstLivePriceEnglish,
			FirstLiveOpenTime:     testInfo.FirstLiveOpenTime,
			FirstLiveStartTime:    testInfo.FirstLiveStartTime,
			SkipOfflineTest:       true,
		},
	}
}

func CreateCycloneFetcher(
	baseURL string,
	shortYearIterableURL string,
	prefecture string,
	area string,
	venueID string,
	dayImageSubstring string,
	testInfo TestInfo,
) Simple {
	return Simple{
		BaseURL:              baseURL,
		ShortYearIterableURL: shortYearIterableURL,
		LiveSelector:         "//body/table",
		// trust
		TitleQuerier:   *htmlquerier.QAll("//td/p/span[1]/descendant::text()[not(preceding-sibling::*[self::span or self::strong]) and normalize-space(.)!='' and (not(ancestor::strong or ancestor::a) or ancestor::span[last()]/text()[1]/preceding-sibling::*[self::span or self::strong])]").Join(" "),
		ArtistsQuerier: *htmlquerier.QAll("//span/strong").Trim().TrimSuffix("PRESENTS").SplitIgnoreWithin("[\n/]", '(', ')'),
		PriceQuerier:   *htmlquerier.Q("//td/p/span[1]/text()[last()-1]"),

		TimeHandler: TimeHandler{
			YearQuerier:      *htmlquerier.Q("//b"),
			MonthQuerier:     *htmlquerier.Q("//b/strong"),
			DayQuerier:       *htmlquerier.Q(fmt.Sprintf("//img[contains(@src, '%s')]/@src", dayImageSubstring)),
			OpenTimeQuerier:  *htmlquerier.Q("//td/p/span[1]/text()[last()-2]"),
			StartTimeQuerier: *htmlquerier.Q("//td/p/span[1]/text()[last()-2]").After("|"),
		},

		PrefectureName: prefecture,
		AreaName:       area,
		VenueID:        venueID,
		Longitude:      139.698562,
		Latitude:       35.661563,

		TestInfo: testInfo,
	}
}

func CreateLoftFetcher(
	baseUrl string,
	shortYearIterableURL string,
	prefecture string,
	area string,
	venue string,
	testInfo TestInfo,
	latitude float64,
	longitude float64,
) Simple {
	return Simple{
		BaseURL:              baseUrl,
		ShortYearIterableURL: shortYearIterableURL,
		LiveSelector:         "//section[contains(@class, 'block_schedule_list')]/div/div[contains(@class, 'list')]/div",
		TitleQuerier:         *htmlquerier.Q("//h1").CutWrapper("『", "』"),
		ArtistsQuerier:       *htmlquerier.QAll("//div[contains(@class, 'entry')]/p").SplitIgnoreWithin("(\n)|( / )", '（', '）'),
		PriceQuerier:         *htmlquerier.Q("//dd[contains(@class, 'ticket_detail_box')]/text()[1]"),
		ExpandedLiveSelector: "//a[contains(@class, 'js-cursor-elm')]",

		TimeHandler: TimeHandler{
			YearQuerier:      *htmlquerier.Q("//nav[contains(@class, 'navi_month')]//div[contains(@class, 'year')]"),
			MonthQuerier:     *htmlquerier.Q("//nav[contains(@class, 'navi_month')]//div[contains(@class, 'month')]/text()[1]"),
			DayQuerier:       *htmlquerier.Q("//time/div[contains(@class, 'day')]"),
			OpenTimeQuerier:  *htmlquerier.Q("//div[contains(@class, 'open')]"),
			StartTimeQuerier: *htmlquerier.Q("//div[contains(@class, 'open')]").After("START"),
		},

		PrefectureName: prefecture,
		AreaName:       area,
		VenueID:        venue,
		Latitude:       latitude,
		Longitude:      longitude,

		TestInfo: testInfo,
	}
}

func CreateOldLoftFetcher(
	baseUrl string,
	shortYearIterableURL string,
	prefecture string,
	area string,
	venue string,
	testInfo TestInfo,
	latitude float64,
	longitude float64,
) Simple {
	return Simple{
		BaseURL:              baseUrl,
		ShortYearIterableURL: shortYearIterableURL,
		LiveSelector:         "//table[contains(@class, 'timetable')]/tbody/tr",
		TitleQuerier:         *htmlquerier.Q("//h3").CutWrapper("『", "』"),
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
		Latitude:       latitude,
		Longitude:      longitude,

		TestInfo: testInfo,
	}
}

func CreateOmatsuriFetcher(baseUrl string, prefectureName, areaName, venueID string, latitude, longitude float64, testInfo TestInfo) Simple {
	return Simple{
		BaseURL:              baseUrl,
		ShortYearIterableURL: baseUrl + "events?date=20%d%%2F%02d",
		LiveSelector:         "//div[@class='blog-entry event-wrap row'][not(contains(.//h2, 'HALL RENTAL') or contains(.//h2, 'Hall Rental') or contains(.//h2, 'business holiday') or contains(.//h2, 'Charter Recording'))]",
		DetailsLinkSelector:  "//h2/a",
		TitleQuerier:         *htmlquerier.Q("//h2"),
		ArtistsQuerier:       *htmlquerier.QAll("//p[@class='act']/span"),
		PriceQuerier:         *htmlquerier.QAll("//p[@class='charge']"),

		TimeHandler: TimeHandler{
			YearQuerier:      *htmlquerier.Q("//ul[@class='pagination']/li[@class='active']"),
			MonthQuerier:     *htmlquerier.Q("//ul[@class='pagination']/li[@class='active']").After("."),
			DayQuerier:       *htmlquerier.Q("//span[@class='day']"),
			OpenTimeQuerier:  *htmlquerier.Q("//p[@class='time']"),
			StartTimeQuerier: *htmlquerier.Q("//p[@class='time']").After("START"),
		},

		PrefectureName: prefectureName,
		AreaName:       areaName,
		VenueID:        venueID,
		Latitude:       latitude,
		Longitude:      longitude,

		TestInfo: testInfo,
	}
}

func CreateProjectFamiryFetcher(baseUrl string, venueId string, latitude, longitude float64, testInfo TestInfo) Simple {
	return Simple{
		BaseURL:              baseUrl,
		ShortYearIterableURL: baseUrl + "schedule/20%d/%02d/",
		LiveSelector:         "//div[@class='sche-archives-flex']/a",
		ExpandedLiveSelector: ".",
		TitleQuerier:         *htmlquerier.Q("//h2"),
		ArtistsQuerier:       *htmlquerier.Q("//dt[./text() = 'ACT:']/following-sibling::dd").Split(" / "),
		PriceQuerier:         *htmlquerier.QAll("//dl[@class='Bbox'][contains(./dt/text(), 'ADV:') or contains(./dt/text(), 'DOOR:')]").Join(" / "),

		TimeHandler: TimeHandler{
			YearQuerier:      *htmlquerier.Q("//li[@class='achives-year']"),
			MonthQuerier:     *htmlquerier.Q("//li[@class='achives-year']").After("/"),
			DayQuerier:       *htmlquerier.Q("//span[@class='date-text']").After("/").After("/"),
			OpenTimeQuerier:  *htmlquerier.Q("//dt[contains(./text(), 'OPEN:')]/following-sibling::dd"),
			StartTimeQuerier: *htmlquerier.Q("//dt[contains(./text(), 'START:')]/following-sibling::dd"),
		},

		PrefectureName: "fukuoka",
		AreaName:       "tenjin",
		VenueID:        venueId,
		Latitude:       latitude,
		Longitude:      longitude,
		RequireArtists: true,

		TestInfo: testInfo,
	}
}
