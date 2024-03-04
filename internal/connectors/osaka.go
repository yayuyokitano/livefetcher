package connectors

import (
	"time"

	"github.com/yayuyokitano/livefetcher/internal/core/fetchers"
	"github.com/yayuyokitano/livefetcher/internal/core/htmlquerier"
	"github.com/yayuyokitano/livefetcher/internal/core/util"
)

/******************
 * 								*
 *	Shinsaibashi	*
 *								*
 ******************/

var ShinsaibashiBronzeFetcher = fetchers.Simple{
	BaseURL:              "http://osakabronze.com",
	ShortYearIterableURL: "http://osakabronze.com/schedule20%d%02d.php",
	LiveSelector:         "//div[@class='eventbox']",
	TitleQuerier:         *htmlquerier.Q("//p[@class='midashi']").ReplaceAllRegex(`\s+`, " "), // 10/10
	ArtistsQuerier:       *htmlquerier.Q("//p[@class='bandlist']").SplitIgnoreWithin(`\n|( \/ )`, '(', ')'),
	PriceQuerier:         *htmlquerier.Q("//p[@class='openstart']/text()[2]").After("TICKET "),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//h4").Before("年"),
		MonthQuerier:     *htmlquerier.Q("//h4").After("年").Before("月"),
		DayQuerier:       *htmlquerier.Q("//h4").After("月").Before("日"),
		OpenTimeQuerier:  *htmlquerier.Q("//p[@class='openstart']"),
		StartTimeQuerier: *htmlquerier.Q("//p[@class='openstart']").After("START "),

		IsYearInLive:  true,
		IsMonthInLive: true,
	},

	PrefectureName: "osaka",
	AreaName:       "shinsaibashi",
	VenueID:        "shinsaibashi-bronze",

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         27,
		FirstLiveTitle:        "1st full album Desire of life Game Change Tour final seriesLove the past play the future",
		FirstLiveArtists:      []string{"Hyuga", "DETOX"},
		FirstLivePrice:        "adv ¥2500 door ¥3000(別途1D ¥600)",
		FirstLivePriceEnglish: "adv ¥2500 door ¥3000(Separately1D ¥600)",
		FirstLiveOpenTime:     time.Date(2024, 3, 2, 18, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 2, 19, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "http://osakabronze.com/schedule20%d%02d.php",
	},
}

var ShinsaibashiKurageFetcher = fetchers.Simple{
	BaseURL:              "https://livehouse-kurage.com",
	InitialURL:           "https://livehouse-kurage.com/schedule/",
	LiveSelector:         "//li[@class='archive_li']",
	ExpandedLiveSelector: "//a",
	TitleQuerier:         *htmlquerier.Q("//h4"),
	ArtistsQuerier:       *htmlquerier.Q("//p[@class='schedule_act']").SplitIgnoreWithin(`[\n/、]`, '(', ')'),
	PriceQuerier:         *htmlquerier.Q("//p[@class='schedule_price']"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//span[@class='schedule_year']"),
		MonthQuerier:     *htmlquerier.Q("//p[@class='schedule_day']"),
		DayQuerier:       *htmlquerier.Q("//p[@class='schedule_day']").After("."),
		OpenTimeQuerier:  *htmlquerier.Q("//p[@class='schedule_time']"),
		StartTimeQuerier: *htmlquerier.Q("//p[@class='schedule_time']").After("/"),

		IsYearInLive:  true,
		IsMonthInLive: true,
	},

	PrefectureName: "osaka",
	AreaName:       "shinsaibashi",
	VenueID:        "shinsaibashi-kurage",

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         10,
		FirstLiveTitle:        "天女神樂 神楽祭『〜花〜』",
		FirstLiveArtists:      []string{"天女神樂", "Panic×Panic", "メイビスレーヌ"},
		FirstLivePrice:        "前売3,500円/当日4,000円(ドリンク代別途600円)/カメラ登録料＋1,000円",
		FirstLivePriceEnglish: "Reservation3,500円/Door4,000円(Drink代Separately600円)/Camera fee＋1,000円",
		FirstLiveOpenTime:     time.Date(2024, 3, 10, 18, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 10, 18, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://livehouse-kurage.com/schedule/%e5%a4%a9%e5%a5%b3%e7%a5%9e%e6%a8%82-%e7%a5%9e%e6%a5%bd%e7%a5%ad%e3%80%8e%e3%80%9c%e8%8a%b1%e3%80%9c%e3%80%8f/",
	},
}

// this is insanity
var ShinsaibashiUtausakanaFetcher = fetchers.Simple{
	BaseURL:             "http://utausakana.com/",
	InitialURL:          "http://utausakana.com/menu/",
	NextSelector:        "//div[@class='pager']/a[@class='next']",
	LiveSelector:        "//div[@class='menu_body']/p[text()='act)']",
	TitleQuerier:        *htmlquerier.Q("/preceding::p[1]").TrimPrefix("『").TrimSuffix("』"),
	ArtistsQuerier:      *htmlquerier.QAll("/following-sibling::p").DeleteFrom("\u00A0"),
	PriceQuerier:        *htmlquerier.QAll("/following-sibling::p").DeleteUntil("\u00A0").KeepIndex(1).ReplaceAll("\u00A0", ""),
	DetailsLinkSelector: "/ancestor::div[@class='menu_body']/div[@class='menu_title']/a",

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("/ancestor::div[@class='menu_list']/div[@class='menu_category']"),
		MonthQuerier:     *htmlquerier.Q("/ancestor::div[@class='menu_list']/div[@class='menu_category']").After("年"),
		DayQuerier:       *htmlquerier.Q("/ancestor::div[@class='menu_body']/div[@class='menu_title']").After("/"),
		OpenTimeQuerier:  *htmlquerier.QAll("/following-sibling::p").DeleteUntil("\u00A0").KeepIndex(0),
		StartTimeQuerier: *htmlquerier.QAll("/following-sibling::p").DeleteUntil("\u00A0").KeepIndex(0).After("/").After("/"),

		IsYearInLive:  true,
		IsMonthInLive: true,
	},

	PrefectureName: "osaka",
	AreaName:       "shinsaibashi",
	VenueID:        "shinsaibashi-utausakana",

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         37,
		FirstLiveTitle:        "song日和",
		FirstLiveArtists:      []string{"杉本ラララ", "清原ありさ", "osakana", "雨蘭"},
		FirstLivePrice:        "adv/day ¥2500(1d別)",
		FirstLivePriceEnglish: "adv/day ¥2500(1dSeparately)",
		FirstLiveOpenTime:     time.Date(2024, 3, 1, 18, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "http://utausakana.com/menu/1023584",
	},
}
