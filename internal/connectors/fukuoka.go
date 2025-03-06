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

var ZeppFukuokaFetcher = fetchers.CreateZeppFetcher("fukuoka", "fukuoka", "fukuoka", "fukuoka-zepp", 33.593563, 130.363187, fetchers.TestInfo{
	NumberOfLives:         9,
	FirstLiveTitle:        "04 Limited Sazabys「MOON tour 2025」",
	FirstLiveArtists:      []string{"04Limited Sazabys"},
	FirstLivePrice:        "1Fスタンディング（整理番号付）/ SOLD OUT 2F指定席/ SOLD OUT",
	FirstLivePriceEnglish: "1FStanding（Numbered tickets (may affect entry order)）/ SOLD OUT 2FReserved Seating/ SOLD OUT",
	FirstLiveOpenTime:     time.Date(2025, 4, 4, 18, 0, 0, 0, util.JapanTime),
	FirstLiveStartTime:    time.Date(2025, 4, 4, 19, 0, 0, 0, util.JapanTime),
	FirstLiveURL:          "https://www.zepp.co.jp/hall/fukuoka/schedule/single/?rid=146915",
})

var KokuraFuseFetcher = fetchers.Simple{
	BaseURL:              "https://kokurafuse.com/",
	ShortYearIterableURL: "https://kokurafuse.com/monthly/?d=20%d-%02d-01",
	LiveSelector:         "//article[@class='schedule-item']",
	DetailsLinkSelector:  "//a",
	TitleQuerier:         *htmlquerier.Q("//h2"),
	ArtistsQuerier:       *htmlquerier.Q("//dl[@class='event__cast']/dd").Split(" / "),
	PriceQuerier:         *htmlquerier.Q("//dl[@class='event__price']/dd").NormalizeWhitespace(),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//h3[@class='content-title']/small"),
		MonthQuerier:     *htmlquerier.Q("//h3[@class='content-title']/small").After("年"),
		DayQuerier:       *htmlquerier.Q("//span[@class='event__date-day']"),
		OpenTimeQuerier:  *htmlquerier.Q("//dl[@class='event__time']//dd"),
		StartTimeQuerier: *htmlquerier.Q("//dl[@class='event__time']//dd").After("開演"),
	},

	PrefectureName: "fukuoka",
	AreaName:       "kokura",
	VenueID:        "kokura-fuse",
	Latitude:       33.886437,
	Longitude:      130.879937,
	RequireArtists: true,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         24,
		FirstLiveTitle:        `MUZIC PUMP`,
		FirstLiveArtists:      []string{"PSYCO LOGIC BOX", "BADWHY’s", "reo goble and his band", "珊々瑚々", "dop", "THEBIGDIPPER"},
		FirstLivePrice:        "前売 2,500円 / 当日 3,000円",
		FirstLivePriceEnglish: "Reservation 2,500円 / Door 3,000円",
		FirstLiveOpenTime:     time.Date(2025, 3, 1, 17, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 3, 1, 18, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://kokurafuse.com/schedule/schedule2697/",
	},
}
