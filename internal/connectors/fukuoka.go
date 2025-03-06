package connectors

import (
	"time"

	"github.com/yayuyokitano/livefetcher/internal/core/fetchers"
	"github.com/yayuyokitano/livefetcher/internal/core/htmlquerier"
	"github.com/yayuyokitano/livefetcher/internal/core/util"
)

var FukuokaGrafFetcher = fetchers.CreateGrafFetcher("https://fukuoka-graf.com/", "fukuoka-graf", 33.593063, 130.394813, fetchers.TestInfo{
	NumberOfLives:         12,
	FirstLiveTitle:        `sea's line × FLAGS`,
	FirstLiveArtists:      []string{"sea's line", "Amsterdamned", "the seadays", "elephant", "futurina", "犬のやすらぎ", "奏人心", "竹崎彰悟", "gn8mykitten", "藤山拓", "Etranger"},
	FirstLivePrice:        "￥3000 / ￥3500 / 1DRINK ORDER",
	FirstLivePriceEnglish: "￥3000 / ￥3500 / 1DRINK ORDER",
	FirstLiveOpenTime:     time.Date(2025, 3, 1, 15, 0, 0, 0, util.JapanTime),
	FirstLiveStartTime:    time.Date(2025, 3, 1, 15, 30, 0, 0, util.JapanTime),
	FirstLiveURL:          util.InsertShortYearMonth("https://fukuoka-graf.com/20%d%02d.html"),
}, *htmlquerier.Q("//div[@class='title']/text()[contains(., 'Schedule')]"))

var FukuokaOpsFetcher = fetchers.CreateProjectFamiryFetcher("https://op-s.info/", "fukuoka-ops", 33.592062, 130.395562, fetchers.TestInfo{
	NumberOfLives:         18,
	FirstLiveTitle:        "hyakki pre. 百鬼夜行vol.3",
	FirstLiveArtists:      []string{"hyakki", "IrisaVior", "文明廻花", "hyper luck 2", "Penny Lane"},
	FirstLivePrice:        "ADV: ¥2500 / DOOR: ¥3000",
	FirstLivePriceEnglish: "ADV: ¥2500 / DOOR: ¥3000",
	FirstLiveOpenTime:     time.Date(2025, 3, 1, 18, 0, 0, 0, util.JapanTime),
	FirstLiveStartTime:    time.Date(2025, 3, 1, 18, 30, 0, 0, util.JapanTime),
	FirstLiveURL:          "https://op-s.info/schedule/hyakki-pre-%%e7%%99%%be%%e9%%ac%%bc%%e5%%a4%%9c%%e8%%a1%%8cvol-3/",
})

var FukuokaPublicBarBassicFetcher = fetchers.Simple{
	BaseURL:              "https://freecalend.com/",
	ShortYearIterableURL: "https://freecalend.com/open/mem82502",
	LiveSelector:         "//div[@id='schedule']/div[contains(@class, 'days')]",
	TitleQuerier:         *htmlquerier.Q("//div[@class='cat' and ./text()='TITLE']/following-sibling::div[@class='desc']"),
	ArtistsQuerier:       *htmlquerier.Q("//div[@class='cat' and ./text()='CAST']/following-sibling::div[@class='desc']").Split(" / "),
	PriceQuerier:         *htmlquerier.Q("//div[@class='cat' and ./text()='ADV / DOOR']/following-sibling::div[@class='desc']"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//div[@class='title']/text()[contains(., 'Schedule')]"),
		MonthQuerier:     *htmlquerier.Q("//div[@class='title']/text()[contains(., 'Schedule')]").After("/"),
		DayQuerier:       *htmlquerier.Q("//div[@class='date2']").After("/"),
		OpenTimeQuerier:  *htmlquerier.Q("//div[@class='cat' and ./text()='OPEN / START']/following-sibling::div[@class='desc']"),
		StartTimeQuerier: *htmlquerier.Q("//div[@class='cat' and ./text()='OPEN / START']/following-sibling::div[@class='desc']").After("/"),
	},

	PrefectureName: "fukuoka",
	AreaName:       "tenjin",
	VenueID:        "fukuoka-graf",
	Latitude:       33.593063,
	Longitude:      130.394813,
	RequireArtists: true,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         12,
		FirstLiveTitle:        `sea's line × FLAGS`,
		FirstLiveArtists:      []string{"sea's line", "Amsterdamned", "the seadays", "elephant", "futurina", "犬のやすらぎ", "奏人心", "竹崎彰悟", "gn8mykitten", "藤山拓", "Etranger"},
		FirstLivePrice:        "￥3000 / ￥3500 / 1DRINK ORDER",
		FirstLivePriceEnglish: "￥3000 / ￥3500 / 1DRINK ORDER",
		FirstLiveOpenTime:     time.Date(2025, 3, 1, 15, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 3, 1, 15, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          util.InsertShortYearMonth("https://fukuoka-graf.com/20%d%02d.html"),
	},
}

var FukuokaVoodooLoungeFetcher = fetchers.CreateGrafFetcher("https://voodoolounge.jp/", "fukuoka-voodoolounge", 33.593063, 130.394687, fetchers.TestInfo{
	NumberOfLives:         12,
	FirstLiveTitle:        `香野子東名阪福ツアー`,
	FirstLiveArtists:      []string{"香野子"},
	FirstLivePrice:        "【優先入場料金】 ￥3000 / 1DRINK ORDER 【一般前売料金】 ￥1500 / 1DRINK ORDER",
	FirstLivePriceEnglish: "【Priority entryEntry料金】 ￥3000 / 1DRINK ORDER 【Ordinary TicketReservation料金】 ￥1500 / 1DRINK ORDER",
	FirstLiveOpenTime:     time.Date(2025, 4, 4, 18, 50, 0, 0, util.JapanTime),
	FirstLiveStartTime:    time.Date(2025, 4, 4, 19, 20, 0, 0, util.JapanTime),
	FirstLiveURL:          util.InsertShortYearMonth("https://voodoolounge.jp/20%d%02d.html"),
}, htmlquerier.Querier{})

var FukuokaQueblickFetcher = fetchers.CreateProjectFamiryFetcher("https://queblick.com/", "fukuoka-queblick", 33.589562, 130.393562, fetchers.TestInfo{
	NumberOfLives:         22,
	FirstLiveTitle:        "Mix Box",
	FirstLiveArtists:      []string{"ちょこはち", "ABYSS", "mm.", "Writers Sky"},
	FirstLivePrice:        "ADV: ¥2500 / DOOR: ¥3000",
	FirstLivePriceEnglish: "ADV: ¥2500 / DOOR: ¥3000",
	FirstLiveOpenTime:     time.Date(2025, 3, 1, 18, 0, 0, 0, util.JapanTime),
	FirstLiveStartTime:    time.Date(2025, 3, 1, 18, 30, 0, 0, util.JapanTime),
	FirstLiveURL:          "https://queblick.com/schedule/mix-box-60/",
})
