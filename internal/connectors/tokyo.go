package connectors

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/yayuyokitano/livefetcher/internal/core/fetchers"
	"github.com/yayuyokitano/livefetcher/internal/core/htmlquerier"
	"github.com/yayuyokitano/livefetcher/internal/core/util"
	"golang.org/x/net/html"
)

var ShimokitazawaTestFetcher = fetchers.Simple{
	BaseURL:        "http://localhost:9999/",
	InitialURL:     "http://localhost:9999/static/testLive.html",
	LiveSelector:   "//div[@class='live']",
	TitleQuerier:   *htmlquerier.Q("//h2"),
	ArtistsQuerier: *htmlquerier.QAll("//li"),
	PriceQuerier:   *htmlquerier.Q("//p[@class='price']"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//p[@class='date']"),
		MonthQuerier:     *htmlquerier.Q("//p[@class='date']").After("年"),
		DayQuerier:       *htmlquerier.Q("//p[@class='date']").After("月"),
		OpenTimeQuerier:  *htmlquerier.Q("//p[@class='open']"),
		StartTimeQuerier: *htmlquerier.Q("//p[@class='start']"),

		IsYearInLive:  true,
		IsMonthInLive: true,
	},

	PrefectureName: "tokyo",
	AreaName:       "shimokitazawa",
	VenueID:        "shimokitazawa-test",
	Latitude:       35.661563,
	Longitude:      139.666938,

	TestInfo: fetchers.TestInfo{
		IgnoreTest: true,
	},
}

/**************
 *            *
 *  Hachioji  *
 *            *
 **************/

var HachiojiRipsFetcher = fetchers.CreateHachiojiRinkyDinkFetcher("rips", 35.656937, 139.336687, fetchers.TestInfo{
	NumberOfLives:         15,
	FirstLiveTitle:        "“THE MESEEKS JAPAN TOUR 2025”",
	FirstLiveArtists:      []string{"THE MESEEKS(Switzerland)", "LoG", "2STRIKE 3BALL", "clyde", "at field of school", "CAECUM"},
	FirstLivePrice:        "adv:2500yen、door:3000yen",
	FirstLivePriceEnglish: "adv:2500yen、door:3000yen",
	FirstLiveOpenTime:     time.Date(2025, 3, 1, 17, 15, 0, 0, util.JapanTime),
	FirstLiveStartTime:    time.Date(2025, 3, 1, 17, 45, 0, 0, util.JapanTime),
	FirstLiveURL:          util.InsertYearMonth("http://rips.rinkydink.info/rips/schedule?yr=%d&month=%d"),
})

var HachiojiMatchVoxFetcher = fetchers.CreateHachiojiRinkyDinkFetcher("matchvox", 35.656937, 139.336687, fetchers.TestInfo{
	NumberOfLives:         15,
	FirstLiveTitle:        "Match Vox 21st anniversary-初日-",
	FirstLiveArtists:      []string{"SAIHATE", "ザ・ドーベルマンション", "Jacob Jr.", "THE ERIC MARK'S", "サンカクスイ"},
	FirstLivePrice:        "adv.¥2000、/door ¥2500",
	FirstLivePriceEnglish: "adv.¥2000、/door ¥2500",
	FirstLiveOpenTime:     time.Date(2025, 4, 1, 18, 0, 0, 0, util.JapanTime),
	FirstLiveStartTime:    time.Date(2025, 4, 1, 18, 30, 0, 0, util.JapanTime),
	FirstLiveURL:          util.InsertYearMonth("http://matchvox.rinkydink.info/matchvox/schedule?yr=%d&month=%d"),
})

/***************
 *             *
 *  Kichijoji  *
 *             *
 ***************/

type kichijojiBlackAndBlueResponse []struct {
	Title       string `json:"title"`
	Url         string `json:"url"`
	Date        string `json:"start"`
	Description string `json:"description"`
}

var KichijojiBlackAndBlueFetcher = fetchers.Simple{
	BaseURL: "https://blackandblue.tokyo/",
	LiveHTMLFetcher: func(testDocument []byte) (nodes []*html.Node, err error) {
		nodes = make([]*html.Node, 0)
		var list kichijojiBlackAndBlueResponse
		if testDocument == nil {
			if err = util.GetJSON(
				util.InsertYearMonth("https://blackandblue.tokyo/wp-admin/admin-ajax.php?action=eventorganiser-fullcal&start=%d-%02d-01&timeformat=g%3Ai%20A"),
				&list,
			); err != nil {
				return
			}
		} else {
			if err = json.Unmarshal(testDocument, &list); err != nil {
				return
			}
		}
		for _, live := range list {
			if strings.HasSuffix(live.Description, "</br></br>") {
				continue
			}
			n, err := html.Parse(strings.NewReader(fmt.Sprintf(
				"<span id='date'>%s</span><span id='title'>%s</span><span id='body'>%s</span><span id='url'>%s</span>",
				live.Date,
				live.Title,
				live.Description,
				live.Url,
			)))
			if err != nil {
				continue
			}
			nodes = append(nodes, n)
		}
		return
	},
	DetailsLinkSelector: "//span[@id='url']",
	TitleQuerier:        *htmlquerier.Q("//span[@id='title']"),
	ArtistsQuerier:      *htmlquerier.Q("//span[@id='body']").After("【出演】/kk").Before("/kk").Split("/"),
	PriceQuerier:        *htmlquerier.Q("//span[@id='body']").After("charge").Before("/kk").Trim().Prefix("¥"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//span[@id='date']"),
		MonthQuerier:     *htmlquerier.Q("//span[@id='date']").After("-"),
		DayQuerier:       *htmlquerier.Q("//span[@id='date']").After("-").After("-"),
		OpenTimeQuerier:  *htmlquerier.Q("//span[@id='body']").After("open"),
		StartTimeQuerier: *htmlquerier.Q("//span[@id='body']").After("start").After("/"),

		IsYearInLive:  true,
		IsMonthInLive: true,
	},

	PrefectureName: "tokyo",
	AreaName:       "kichijoji",
	VenueID:        "kichijoji-blackandblue",
	Latitude:       35.703563,
	Longitude:      139.581687,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         38,
		FirstLiveTitle:        "“MIX’N ROCK”",
		FirstLiveArtists:      []string{"サンバラッチャ", "真夜華", "Take life easy", "ぶるうまんJAM"},
		FirstLivePrice:        "¥2600＋1d",
		FirstLivePriceEnglish: "¥2600＋1d",
		FirstLiveOpenTime:     time.Date(2025, 3, 1, 18, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 3, 1, 18, 30, 0, 0, util.JapanTime),
	},
}

var KichijojiClubSeataFetcher = fetchers.CreateBassOnTopFetcher(
	"https://seata.jp/",
	"https://seata.jp/schedule/calendar/20%d/%02d/",
	"tokyo", "kichijoji", "kichijoji-clubseata",
	fetchers.TestInfo{
		NumberOfLives:         17,
		FirstLiveTitle:        "Outstanding Vol.2",
		FirstLiveArtists:      []string{"柘榴", "雫", "skip-A", "せな", "月乃", "成宮 亮", "ぽむ", "ゆあ"},
		FirstLivePrice:        "1部：前売券4,200円/当日券4,300円　2部：前売券4,700円/当日券4,800円(1Drink代金￥700別途必要)",
		FirstLivePriceEnglish: "1部：Reservation券4,200円/Door券4,300円　2部：Reservation券4,700円/Door券4,800円(1DrinkPrice￥700Must be purchased separately)",
		FirstLiveOpenTime:     time.Date(2025, 3, 2, 12, 15, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 3, 2, 12, 45, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://seata.jp/schedule/detail/40625",
	},
	35.705813, 139.581813,
)

var KichijojiPlanetKFetcher = fetchers.Simple{
	BaseURL:        "http://inter-planets.net/",
	InitialURL:     "http://inter-planets.net/calendar",
	NextSelector:   "//div[contains(@class, 'ai1ec-pagination')]/a[contains(@class, 'ai1ec-next-page')]",
	LiveSelector:   "//div[@class='ai1ec-agenda-view']/div[contains(@class, 'ai1ec-date')][not(contains(.//span[@class='ai1ec-event-title'], 'ホールレンタル'))]",
	TitleQuerier:   *htmlquerier.Q("//span[@class='ai1ec-event-title']"),
	ArtistsQuerier: *htmlquerier.QAll("//div[@class='ai1ec-event-description']/p[1]//strong"),
	DetailQuerier:  *htmlquerier.Q("//div[@class='ai1ec-event-description']"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:  *htmlquerier.Q("//div[@class='ai1ec-year']"),
		MonthQuerier: *htmlquerier.Q("//div[@class='ai1ec-month']"),
		DayQuerier:   *htmlquerier.Q("//div[@class='ai1ec-day']"),

		IsYearInLive:  true,
		IsMonthInLive: true,
	},

	PrefectureName: "tokyo",
	AreaName:       "kichijoji",
	VenueID:        "kichijoji-planetk",
	Latitude:       35.705312,
	Longitude:      139.578688,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         27,
		FirstLiveTitle:        "PLANET ATTACK!!!",
		FirstLiveArtists:      []string{"ザ・ビンビンズ", "CAROLAN’S", "SHAKTPUNK", "カトウマサタカ", "マグロジェット"},
		FirstLivePrice:        "TICKET ¥2,500、/ ¥3,000 (+1D)",
		FirstLivePriceEnglish: "TICKET ¥2,500、/ ¥3,000 (+1D)",
		FirstLiveOpenTime:     time.Date(2025, 3, 1, 17, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 3, 1, 18, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "http://inter-planets.net/calendar",
	},
}

var KichijojiShuffleFetcher = fetchers.Simple{
	BaseURL:              "http://k-shuffle.com/",
	ShortYearIterableURL: "http://k-shuffle.com/schedule/20%d-%02d",
	LiveSelector:         "//article[@class='sched-box'][not(contains(.//h3, 'HALL RENTAL'))]",
	TitleQuerier:         *htmlquerier.Q("//h3/span[@class='sp-br']"),
	ArtistsQuerier:       *htmlquerier.Q("/p[1]").Split(" / "),
	PriceQuerier:         *htmlquerier.Q("/p[2]/text()[1]").After("adv").Prefix("adv"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//h2"),
		MonthQuerier:     *htmlquerier.Q("//h2").After("."),
		DayQuerier:       *htmlquerier.Q("//h3"),
		OpenTimeQuerier:  *htmlquerier.Q("/p[2]"),
		StartTimeQuerier: *htmlquerier.Q("/p[2]").After("/"),
	},

	PrefectureName: "tokyo",
	AreaName:       "kichijoji",
	VenueID:        "kichijoji-shuffle",
	Latitude:       35.701813,
	Longitude:      139.581187,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         33,
		FirstLiveTitle:        "国生優太 presents live「cobalt」",
		FirstLiveArtists:      []string{"国生優太", "あいあむみー", "マエノミドリ", "一色", "後藤凌", "里星来"},
		FirstLivePrice:        "adv ¥3,500/door ¥-,---",
		FirstLivePriceEnglish: "adv ¥3,500/door ¥-,---",
		FirstLiveOpenTime:     time.Date(2025, 3, 1, 18, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 3, 1, 18, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          util.InsertYearMonth("http://k-shuffle.com/schedule/%d-%02d"),
	},
}

var KichijojiWarpFetcher = fetchers.Simple{
	BaseURL:              "http://warp.rinky.info/",
	ShortYearIterableURL: "http://warp.rinky.info/schedules_cat/20%d-%02d",
	LiveSelector:         "//article[@class='schedules-box'][not(contains(.//h4, 'ホールレンタル'))]",
	TitleQuerier:         *htmlquerier.Q("//h4"),
	ArtistsQuerier:       *htmlquerier.QAll("//div[@class='w-flyer']/text()"),
	PriceQuerier:         *htmlquerier.Q("//section[@class='notes-wrapper']/p[2]"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//section[@class='schedules-navigation']//h3/text()"),
		MonthQuerier:     *htmlquerier.Q("//section[@class='schedules-navigation']//h3/span"),
		DayQuerier:       *htmlquerier.Q("//section[@data-aos='fade-up'][1]"),
		OpenTimeQuerier:  *htmlquerier.Q("//section[@class='notes-wrapper']/p[1]/span"),
		StartTimeQuerier: *htmlquerier.Q("//section[@class='notes-wrapper']/p[1]/span").After(" / "),
	},

	PrefectureName: "tokyo",
	AreaName:       "kichijoji",
	VenueID:        "kichijoji-warp",
	Latitude:       35.704563,
	Longitude:      139.583062,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         20,
		FirstLiveTitle:        `Himeuzu 1st. album "Bouquet" release party`,
		FirstLiveArtists:      []string{"ヒメウズ", "場末", "シノエフヒ", "メレ"},
		FirstLivePrice:        "ADV / DOOR¥2400 / ¥2900 (+1D)",
		FirstLivePriceEnglish: "ADV / DOOR¥2400 / ¥2900 (+1D)",
		FirstLiveOpenTime:     time.Date(2025, 3, 1, 11, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 3, 1, 11, 50, 0, 0, util.JapanTime),
		FirstLiveURL:          util.InsertYearMonth("http://warp.rinky.info/schedules_cat/%d-%02d"),
	},
}

/************
 *          *
 *  Koenji  *
 *          *
 ************/

var KoenjiClubRootsFetcher = fetchers.Simple{
	BaseURL:              "http://www.muribushi.jp/",
	ShortYearIterableURL: "http://www.muribushi.jp/schedule_20%d/%d.html",
	LiveSelector:         "//dl/dd[not(contains(., 'HALL RENTAL'))]",
	TitleQuerier:         *htmlquerier.Q("//font[@color='#ffff00']"),
	DetailQuerier:        *htmlquerier.Q("."),
	ArtistsQuerier:       *htmlquerier.Q("/node()[2][not(text())] | /node()[2][text()]/text()[1]").Split(" / "),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:  *htmlquerier.Q("//h3"),
		MonthQuerier: *htmlquerier.Q("//h3").After(" "),
		DayQuerier:   *htmlquerier.Q("/preceding-sibling::dt").After("/"),
	},

	PrefectureName: "tokyo",
	AreaName:       "koenji",
	VenueID:        "koenji-clubroots",
	Latitude:       35.705562,
	Longitude:      139.648937,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         28,
		FirstLiveTitle:        `Club ROOTS ! 20th Anniv. pre "シミズフェス 2025！"`,
		FirstLiveArtists:      []string{"KG", "ムカイナオト", "丸山尊", "そえじまとしたか", "ショーシャンク南"},
		FirstLivePrice:        "前) ￥2,000",
		FirstLivePriceEnglish: "ADV ) ￥2,000",
		FirstLiveOpenTime:     time.Date(2025, 3, 3, 18, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 3, 3, 19, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          util.InsertYearMonth("http://www.muribushi.jp/schedule_%d/%d.html"),
	},
}

var KoenjiHighFetcher = fetchers.Simple{
	BaseURL:              "http://koenji-high.com/",
	ShortYearIterableURL: "http://koenji-high.com/schedule/?sy=20%d&sm=%d",
	LiveSelector:         "//div[@id='events']/div[contains(@class, 'eventlist')]",
	TitleQuerier:         *htmlquerier.Q("//h3"),
	ArtistsQuerier:       *htmlquerier.QAll("//th[contains(., 'LINE UP')]/following-sibling::td/text()"),
	PriceQuerier:         *htmlquerier.Q("//th[contains(., 'ADV/DOOR')]/following-sibling::td"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//span[@class='yearimage']/@id"),
		MonthQuerier:     *htmlquerier.Q("//span[@class='monthnumimage']/@id"),
		DayQuerier:       *htmlquerier.Q("//span[@class='daynum']"),
		OpenTimeQuerier:  *htmlquerier.Q("//th[contains(., 'OPEN/START')]/following-sibling::td"),
		StartTimeQuerier: *htmlquerier.Q("//th[contains(., 'OPEN/START')]/following-sibling::td").After("/"),
	},

	PrefectureName: "tokyo",
	AreaName:       "koenji",
	VenueID:        "koenji-high",
	Latitude:       35.703563,
	Longitude:      139.651063,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         25,
		FirstLiveTitle:        "SUCKER BAWL~Rock by nature 61~",
		FirstLiveArtists:      []string{"THE POGO", "THE PRISONER", "RAISE A FLAG", "ISHIKAWA"},
		FirstLivePrice:        "￥3,500/￥4,000",
		FirstLivePriceEnglish: "￥3,500/￥4,000",
		FirstLiveOpenTime:     time.Date(2025, 3, 1, 17, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 3, 1, 18, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          util.InsertYearMonth("http://koenji-high.com/schedule/?sy=%d&sm=%d"),
	},
}

var KoenjiShowBoatFetcher = fetchers.Simple{
	BaseURL:              "https://www.showboat1993.com/",
	ShortYearIterableURL: "https://www.showboat1993.com/2024/20%d-%d",
	LiveSelector:         "//div[@data-testid='mesh-container-content']/section[2]/div[@data-testid='inline-content']/div[@data-testid='mesh-container-content']/div/div[@data-testid='inline-content']/div[@data-testid='mesh-container-content']",
	TitleQuerier:         *htmlquerier.Q("/div[2]/p[1]"),
	ArtistsQuerier:       *htmlquerier.QAll("//p[contains(., '出演')]/following-sibling::p[1]").Split(" / "),
	PriceQuerier:         *htmlquerier.Q("//p[contains(., '前売')]").After("前売").Prefix("前売"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.QAll("//span[@style='font-family:brandon-grot-w01-light,sans-serif;']").KeepIndex(2),
		MonthQuerier:     *htmlquerier.Q("//span[@style='font-family:brandon-grot-w01-light,sans-serif;']"),
		DayQuerier:       *htmlquerier.Q("/div[1]"),
		OpenTimeQuerier:  *htmlquerier.Q("//p[contains(., '開場')]"),
		StartTimeQuerier: *htmlquerier.Q("//p[contains(., '開演')]").After(" / "),
	},

	PrefectureName: "tokyo",
	AreaName:       "koenji",
	VenueID:        "koenji-showboat",
	Latitude:       35.703563,
	Longitude:      139.651063,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         22,
		FirstLiveTitle:        "Magic V PRESENTS Rock’n Freedom Bash VOL.5",
		FirstLiveArtists:      []string{"44 REVOLVER", "Lady Ray-X", "NORRA LOOSE"},
		FirstLivePrice:        "前売￥3,000 / 当日￥3,500",
		FirstLivePriceEnglish: "Reservation￥3,000 / Door￥3,500",
		FirstLiveOpenTime:     time.Date(2025, 3, 1, 18, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 3, 1, 18, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          util.InsertYearMonth("https://www.showboat1993.com/2024/%d-%d"),
	},
}

/*************
 *           *
 *  Toshima  *
 *           *
 *************/

var IkebukuroAdmFetcher = fetchers.Simple{
	BaseURL:              "https://adm-rock.com/",
	InitialURL:           "https://adm-rock.com/schedule/%e3%83%aa%e3%82%b9%e3%83%88/",
	LiveSelector:         "//div[@class='tribe-events-calendar-list__event-wrapper tribe-common-g-col'][not(contains(.//h3, 'HALL RENTAL'))]",
	ExpandedLiveSelector: "//a",
	TitleQuerier:         *htmlquerier.Q("//h1"),
	ArtistsQuerier:       *htmlquerier.QAll("//div[@class='tribe-events-single-event-description tribe-events-content']/p[1]/text()"),
	PriceQuerier:         *htmlquerier.Q("//p[contains(., 'ADV/DOOR')]"),

	TimeHandler: fetchers.TimeHandler{
		MonthQuerier:     *htmlquerier.Q("//span[@class='tribe-event-date-start']"),
		DayQuerier:       *htmlquerier.Q("//span[@class='tribe-event-date-start']").After("月"),
		OpenTimeQuerier:  *htmlquerier.Q("//p[contains(., 'OP/ST')]"),
		StartTimeQuerier: *htmlquerier.Q("//p[contains(., 'OP/ST')]").After("/").After("/"),
		IsMonthInLive:    true,
	},

	PrefectureName: "tokyo",
	AreaName:       "toshima",
	VenueID:        "ikebukuro-adm",
	Latitude:       35.729563,
	Longitude:      139.716063,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         53,
		FirstLiveTitle:        "南無阿部陀仏×Adm presents「千里の道は池袋から！！」",
		FirstLiveArtists:      []string{"南無阿部陀仏", "輪廻", "wise man", "友達博物舘"},
		FirstLivePrice:        "■ADV/DOOR ¥2500+1D/¥3000+1D",
		FirstLivePriceEnglish: "■ADV/DOOR ¥2500+1D/¥3000+1D",
		FirstLiveOpenTime:     time.Date(util.GetRelevantYear(3), 3, 12, 18, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(util.GetRelevantYear(3), 3, 12, 18, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://adm-rock.com/schedule/%e5%8d%97%e7%84%a1%e9%98%bf%e9%83%a8%e9%99%80%e4%bb%8fxadm-presents%e3%80%8c%e5%8d%83%e9%87%8c%e3%81%ae%e9%81%93%e3%81%af%e6%b1%a0%e8%a2%8b%e3%81%8b%e3%82%89%ef%bc%81%ef%bc%81%e3%80%8d/",
	},
}

var OtsukaDeepaFetcher = fetchers.CreateBassOnTopFetcher(
	"https://otsukadeepa.jp/",
	"https://otsukadeepa.jp/schedule/calendar/20%d/%02d/",
	"tokyo", "toshima", "otsuka-deepa",
	fetchers.TestInfo{
		NumberOfLives:         15,
		FirstLiveTitle:        "23歳を正しく始める方法",
		FirstLiveArtists:      []string{"kilaku｡", "白昼夢", "THIRSTY RAGE", "ネコノヒタイ", "atLestyRe", "CITY OVER", "Nothing Neverminds"},
		FirstLivePrice:        "ADV/DOOR ￥2,000-/￥2,500- (別途1d￥700-)",
		FirstLivePriceEnglish: "ADV/DOOR ￥2,000-/￥2,500- (Separately1d￥700-)",
		FirstLiveOpenTime:     time.Date(2025, 3, 4, 16, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 3, 4, 16, 50, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://otsukadeepa.jp/schedule/detail/37282",
	},
	35.730438, 139.728063,
)

var OtsukaMeetsFetcher = fetchers.CreateOmatsuriFetcher(
	"https://meets.rinky.info/",
	"tokyo",
	"toshima",
	"otsuka-meets",
	35.730187,
	139.728187,
	fetchers.TestInfo{
		NumberOfLives:         29,
		FirstLiveTitle:        "Part 1] WOYATSU BIRTHDAY LIVE~2025",
		FirstLiveArtists:      []string{"をやつ", "ちワ", "ニケ", "あつや"},
		FirstLivePrice:        "￥5000(+Drink)",
		FirstLivePriceEnglish: "￥5000(+Drink)",
		FirstLiveOpenTime:     time.Date(2025, 3, 1, 12, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 3, 1, 12, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://meets.rinky.info/events/25956",
	},
)

/*************
 *           *
 *  Shibuya  *
 *           *
 *************/

var ShibuyaChelseaHotelFetcher = fetchers.Simple{
	BaseURL:              "http://www.chelseahotel.jp/",
	ShortYearIterableURL: "http://www.chelseahotel.jp/%d%02d.html",
	LiveSelector:         "//body/p[@class='bold' and .//a]",
	TitleQuerier:         *htmlquerier.QAll("./text()[not(./preceding-sibling::a) and position() > 1]").Join(" "),
	ArtistsQuerier:       *htmlquerier.QAll("./a"),
	PriceQuerier:         *htmlquerier.Q("./font/text()[not(./preceding-sibling::a) and position() > 1]"),

	TimeHandler: fetchers.TimeHandler{
		MonthQuerier:     *htmlquerier.Q("./text()[1]"),
		DayQuerier:       *htmlquerier.Q("./text()[1]").After("/"),
		OpenTimeQuerier:  *htmlquerier.Q("./font/text()[1]"),
		StartTimeQuerier: *htmlquerier.Q("./font/text()[1]").After("START"),
		IsMonthInLive:    true,
	},

	PrefectureName: "tokyo",
	AreaName:       "shibuya",
	VenueID:        "shibuya-chelseahotel",
	Latitude:       35.662437,
	Longitude:      139.697938,
	RequireArtists: true,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         22,
		FirstLiveTitle:        `Cha'R ONE-MAN Live in TOKYO 「New Glow」`,
		FirstLiveArtists:      []string{"Cha'R"},
		FirstLivePrice:        "前売り￥3,500/当日￥4,000(1ドリンク別)",
		FirstLivePriceEnglish: "前売り￥3,500/当日￥4,000(1ドリンク別)",
		FirstLiveOpenTime:     time.Date(2025, 3, 1, 16, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 3, 1, 17, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          util.InsertShortYearMonth("http://www.chelseahotel.jp/%d%02d.html"),
	},
}

var ShibuyaClubQuattroFetcher = fetchers.Simple{
	BaseURL:              "https://www.club-quattro.com/shibuya/schedule/",
	ShortYearIterableURL: "https://www.club-quattro.com/shibuya/schedule/?ym=20%d%02d",
	LiveSelector:         "//div[@class='event-box']",
	DetailsLinkSelector:  "//a",
	TitleQuerier:         *htmlquerier.QAll("//p[@class='txt-01' or @class='txt-02']").FilterTitle(`[/\n]`, 1),
	ArtistsQuerier:       *htmlquerier.QAll("//p[@class='txt-01' or @class='txt-02']").FilterArtist(`[/\n]`, 0).SplitRegex(`[/\n]`),
	PriceQuerier:         *htmlquerier.Q("//dt[contains(./text(), '料金')]/following-sibling::dd"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//div[@class='head']/p[contains(@class, 'year')]"),
		MonthQuerier:     *htmlquerier.Q("//div[@class='head']//div[contains(@class, 'month-list')]/a[contains(@class, 'current')]"),
		DayQuerier:       *htmlquerier.Q("//p[@class='day']"),
		OpenTimeQuerier:  *htmlquerier.Q("//dt[contains(./text(), '開場/開演')]/following-sibling::dd"),
		StartTimeQuerier: *htmlquerier.Q("//dt[contains(./text(), '開場/開演')]/following-sibling::dd").After("/"),
	},

	PrefectureName: "tokyo",
	AreaName:       "shibuya",
	VenueID:        "shibuya-clubquattro",
	Longitude:      139.697563,
	Latitude:       35.661062,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         30,
		FirstLiveTitle:        `The Biscats TOUR 2024~2025「ロカビリーナイト」`,
		FirstLiveArtists:      []string{"The Biscats"},
		FirstLivePrice:        "前売 ￥5,000",
		FirstLivePriceEnglish: "Reservation ￥5,000",
		FirstLiveOpenTime:     time.Date(2025, 2, 1, 17, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 2, 1, 18, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://www.club-quattro.com/shibuya/schedule/detail/?cd=016227",
	},
}

var ShibuyaCycloneFetcher = fetchers.CreateCycloneFetcher(
	"http://www.cyclone1997.com/",
	"http://www.cyclone1997.com/schedule/20%dschedule_%d.html",
	"tokyo",
	"shibuya",
	"shibuya-cyclone",
	"cyclone_day",
	fetchers.TestInfo{
		NumberOfLives:         34,
		FirstLiveTitle:        "SWD Japan Proudly pre. EMMURE JAPAN TOUR 2024",
		FirstLiveArtists:      []string{"EMMURE", "DVRK", "Sailing Before The Wind", "VICTIMOFDECEPTION", "HAILROSE"},
		FirstLivePrice:        "ADV ¥7000 | DOOR TBA (+1D)",
		FirstLivePriceEnglish: "ADV ¥7000 | DOOR TBA (+1D)",
		FirstLiveOpenTime:     time.Date(2024, 3, 1, 18, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 1, 18, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "http://www.cyclone1997.com/schedule/20%dschedule_%d.html",
	},
)

var ShibuyaDiveFetcher = fetchers.Simple{
	BaseURL:              "https://shibuya-dive.com/",
	ShortYearIterableURL: "https://shibuya-dive.com/schedule/?date=20%d-%02d",
	LiveSelector:         "//article",
	DetailsLinkSelector:  "//a",
	TitleQuerier:         *htmlquerier.Q("//h3"),
	ArtistsQuerier:       *htmlquerier.Q("//th[.='ACT']/following-sibling::td").Split("/"),
	PriceQuerier:         *htmlquerier.QAll("//th[.='ADV' or .='DOOR']/following-sibling::td").Join("、DOOR: ").Prefix("ADV: "),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//p[@class='schedule-date']"),
		MonthQuerier:     *htmlquerier.Q("//p[@class='schedule-date']").After("."),
		DayQuerier:       *htmlquerier.Q("//p[@class='schedule-date']").After(".").After("."),
		OpenTimeQuerier:  *htmlquerier.Q("//th[.='OPEN']/following-sibling::td"),
		StartTimeQuerier: *htmlquerier.Q("//th[.='START']/following-sibling::td"),

		IsYearInLive:  true,
		IsMonthInLive: true,
	},

	PrefectureName: "tokyo",
	AreaName:       "shibuya",
	VenueID:        "shibuya-dive",
	Latitude:       35.653937,
	Longitude:      139.708562,

	TestInfo: fetchers.TestInfo{
		// Currently, there are no lives available
		IgnoreTest:            true,
		NumberOfLives:         3,
		FirstLiveTitle:        "HEY! BIRTHDAY LIVE & PARTY! ～ゆーりーとあーいーの秘密の花園～ 秋元悠里　作島藍　生誕イベント開催！！！",
		FirstLiveArtists:      []string{"Hey!Mommy!"},
		FirstLivePrice:        "ADV: VIPチケット：¥6,000 / 一般チケット：¥2,000 各+D代別、DOOR: 一般チケット：¥3,000 +D代別",
		FirstLivePriceEnglish: "ADV: VIP Ticket：¥6,000 / Ordinary Ticket Ticket：¥2,000 各+Drink must be purchased separately、DOOR: Ordinary Ticket Ticket：¥3,000 +Drink must be purchased separately",
		FirstLiveOpenTime:     time.Date(2024, 3, 20, 11, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 20, 12, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://shibuya-dive.com/schedule/hey-birthday-live-party-%%ef%%bd%%9e%%e3%%82%%86%%e3%%83%%bc%%e3%%82%%8a%%e3%%83%%bc%%e3%%81%%a8%%e3%%81%%82%%e3%%83%%bc%%e3%%81%%84%%e3%%83%%bc%%e3%%81%%ae%%e7%%a7%%98%%e5%%af%%86%%e3%%81%%ae%%e8%%8a%%b1%%e5%%9c%%92%%ef%%bd%%9e-%%e7%%a7%%8b%%e5%%85%%83/",
	},
}

var ShibuyaEggmanDayFetcher = fetchers.CreateEggmanFetcher(
	"http://eggman.jp/",
	"http://eggman.jp/schedule-cat/daytime/?syear=20%d&smonth=%02d",
	"tokyo",
	"shibuya",
	"shibuya-eggmanday",
	fetchers.TestInfo{
		NumberOfLives:         21,
		FirstLiveTitle:        "[ 2 0 0 A ] – INT. PENTUMNUS OMNIS –",
		FirstLiveArtists:      []string{"BLVELY", "DIVE TO THE 2ND", "MELODRiVE", "SpendyMily", "été"},
		FirstLivePrice:        "学生 : ¥2400 / 一般 : ¥3400",
		FirstLivePriceEnglish: "Students : ¥2400 / Ordinary Ticket : ¥3400",
		FirstLiveOpenTime:     time.Date(2023, 11, 1, 18, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 11, 1, 18, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "http://eggman.jp/schedule/2-0-0-a-int-pentumnus-omnis/",
	},
)

var ShibuyaEggmanNightFetcher = fetchers.CreateEggmanFetcher(
	"http://eggman.jp/",
	"http://eggman.jp/schedule-cat/nighttime/?syear=20%d&smonth=%02d",
	"tokyo",
	"shibuya",
	"shibuya-eggmannight",
	fetchers.TestInfo{
		NumberOfLives:         7,
		FirstLiveTitle:        "RE:raise house 1on1 battle season4 vol.10",
		FirstLiveArtists:      []string{"Akari", "Suthoom(SYMBOL-ISM)", "HIRO(ALMA/DANCE FUSION/Novel Nextus)", "SACHIO", "Takky", "TOSHIYA SGSD"},
		FirstLivePrice:        "ENTRY　事前： 2000円/1D別　当日：3000円/1D別 観戦：1500円/1D別",
		FirstLivePriceEnglish: "ENTRY　Reservation： 2000円/1 Drink purchase required　Door：3000円/1 Drink purchase required 観戦：1500円/1 Drink purchase required",
		FirstLiveOpenTime:     time.Date(2023, 11, 9, 23, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 11, 10, 0, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "http://eggman.jp/schedule/reraise-house-1on1-battle-season4-vol-10/",
	},
)

var ShibuyaFowsFetcher = fetchers.CreateOmatsuriFetcher(
	"https://shibuya-fows.jp/",
	"tokyo",
	"shibuya",
	"shibuya-fows",
	35.661563,
	139.697938,
	fetchers.TestInfo{
		NumberOfLives:         3,
		FirstLiveTitle:        "Hype The Rock Vol.3",
		FirstLiveArtists:      []string{"Dannie May", "NOMELON NOLEMON", "紫 今"},
		FirstLivePrice:        "",
		FirstLivePriceEnglish: "",
		FirstLiveOpenTime:     time.Date(2025, 3, 2, 18, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 3, 2, 18, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://shibuya-fows.jp/events/25457",
	},
)

var ShibuyaGarretFetcher = fetchers.CreateCycloneFetcher(
	"http://www.cyclone1997.com/",
	"http://www.cyclone1997.com/garret/g_schedule/garret_20%dschedule_%d.html",
	"tokyo",
	"shibuya",
	"shibuya-garret",
	"garret_day",
	fetchers.TestInfo{
		NumberOfLives:         30,
		FirstLiveTitle:        "T.M.Music pre. Ten56. Japan Tour Tokyo Show",
		FirstLiveArtists:      []string{"Ten56.(France)", "Suggestions", "Evilgloom", "Nimbus", "DIVINITIST"},
		FirstLivePrice:        "ADV ¥6900 | DOOR ¥7900 (+1D)",
		FirstLivePriceEnglish: "ADV ¥6900 | DOOR ¥7900 (+1D)",
		FirstLiveOpenTime:     time.Date(2024, 3, 3, 17, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 3, 17, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "http://www.cyclone1997.com/garret/g_schedule/garret_20%dschedule_%d.html",
	},
)

var ShibuyaGeeGeFetcher = fetchers.Simple{
	BaseURL:                   "http://www.gee-ge.net/",
	ShortYearIterableURL:      "http://www.gee-ge.net/schedule/20%d/%02d/",
	LiveSelector:              "//div[contains(@class, 'sche_box')]//img/parent::a",
	ExpandedLiveSelector:      ".",
	ExpandedLiveGroupSelector: "//div[@class='box']/div[contains(@class, 'sche_box') and child::table]",
	TitleQuerier:              *htmlquerier.Q("//strong").ReplaceAll("\n", "").Trim().CutWrapper("『", "』").CutWrapper("【", "】").CutWrapper("「", "」"), // lol
	ArtistsQuerier:            *htmlquerier.Q("//img[contains(@src, 'artist_sche_ttl.gif')]/parent::td/following-sibling::td").SplitIgnoreWithin("[\n、/]", '(', ')').After("◎").After("・"),
	PriceQuerier:              *htmlquerier.Q("//img[contains(@src, 'sche_adv_ttl.gif')]/parent::td/following-sibling::td"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//img[contains(@src, 'sche_date_ttl.gif')]/parent::td/following-sibling::td"),
		MonthQuerier:     *htmlquerier.Q("//img[contains(@src, 'sche_date_ttl.gif')]/parent::td/following-sibling::td").After("/"),
		DayQuerier:       *htmlquerier.Q("//img[contains(@src, 'sche_date_ttl.gif')]/parent::td/following-sibling::td").After("/").After("/"),
		OpenTimeQuerier:  *htmlquerier.Q("//img[contains(@src, 'sche_detail_ttl.gif')]/parent::td/following-sibling::td"),
		StartTimeQuerier: *htmlquerier.Q("//img[contains(@src, 'sche_detail_ttl.gif')]/parent::td/following-sibling::td").After("/"),

		IsYearInLive:  true,
		IsMonthInLive: true,
	},

	PrefectureName: "tokyo",
	AreaName:       "shibuya",
	VenueID:        "shibuya-geege",
	Latitude:       35.662463,
	Longitude:      139.698734,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         34,
		FirstLiveTitle:        "ウダガワガールズコレクション vol.770",
		FirstLiveArtists:      []string{"衿衣", "ナカノユウキ", "sanoha", "杏珠", "菜々姫", "鈴木里咲"},
		FirstLivePrice:        "ADV ¥2800 / DOOR ¥3300 + 1DRINK ¥700",
		FirstLivePriceEnglish: "ADV ¥2800 / DOOR ¥3300 + 1DRINK ¥700",
		FirstLiveOpenTime:     time.Date(2024, 3, 2, 11, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 2, 12, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "http://www.gee-ge.net/detail/2024/03/02/#6131",
	},
}

// TODO: this is in harajuku
var ShibuyaLaDonnaFetcher = fetchers.Simple{
	BaseURL:              "https://www.la-donna.jp/",
	ShortYearIterableURL: "https://www.la-donna.jp/schedules/?ym=20%d-%02d",
	LiveSelector:         "//div[@class='sec01']/div[.//dd[@class='bigTxt'][.!='電話受付' and .!='店舗休業日' and .!='企業様イベントご利用']]",
	TitleQuerier:         *htmlquerier.Q("//dd[@class='bigTxt']"),
	ArtistsQuerier:       *htmlquerier.QAll("//dt[.='出演アーティスト']/following-sibling::dd/text()").SplitIgnoreWithin("・", '【', '】'),
	PriceQuerier:         *htmlquerier.Q("//dt[.='前売り / 当日']/parent::*"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//div[@class='monthly']/div[@class='date']"),
		MonthQuerier:     *htmlquerier.Q("//div[@class='monthly']/div[@class='date']/span"),
		DayQuerier:       *htmlquerier.Q("//div[@class='date']").After("/"),
		OpenTimeQuerier:  *htmlquerier.Q("//dt[.='OPEN / START']/following-sibling::dd"),
		StartTimeQuerier: *htmlquerier.Q("//dt[.='OPEN / START']/following-sibling::dd").After("/"),
	},

	PrefectureName: "tokyo",
	AreaName:       "shibuya",
	VenueID:        "shibuya-ladonna",
	Latitude:       35.669163,
	Longitude:      139.706891,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         26,
		FirstLiveTitle:        "【5/9へ延期】Thank You for All of You !!　〜村瀬由衣2nd Album「幸せのスケッチ」Complete Live〜",
		FirstLiveArtists:      []string{"村瀬由衣(Vo)", "鈴木雄大(Gt&Cho)", "浜田美樹(Cho)", "安部潤(Key)", "鎌田清(Dr)", "山口克彦(Gt)", "遠山陽介(Ba)", "杉真理"},
		FirstLivePrice:        "前売り / 当日 6,000円 / 6,500円",
		FirstLivePriceEnglish: "Reservation / Door 6,000円 / 6,500円",
		FirstLiveOpenTime:     time.Date(2024, 3, 1, 18, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 1, 19, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://www.la-donna.jp/schedules/?ym=20%d-%02d",
	},
}

var ShibuyaOCrestFetcher = fetchers.CreateOFetcher(
	"https://shibuya-o.com/",
	"https://shibuya-o.com/crest/schedule/",
	"tokyo",
	"shibuya",
	"shibuya-ocrest",
	fetchers.TestInfo{
		NumberOfLives:         41,
		FirstLiveTitle:        "IIIIIIIDIOM FREE LIVE Road to -勇往邁進-",
		FirstLiveArtists:      []string{"IIIIIIIDIOM"},
		FirstLivePrice:        "優先 ¥3,000 一般 ¥0 当日 ¥0 （ご予約時ドリンク代別途600円）",
		FirstLivePriceEnglish: "Priority entry ¥3,000 Ordinary Ticket ¥0 Door ¥0 （ごReservation時DrinkNot included in ticket600円）",
		FirstLiveOpenTime:     time.Date(util.GetRelevantYear(11), 11, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(util.GetRelevantYear(11), 11, 1, 19, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://shibuya-o.com/crest/schedule/iiiiiiidiom_23-11-1/",
	},
	35.658687,
	139.695562,
)

var ShibuyaOEastFetcher = fetchers.CreateOFetcher(
	"https://shibuya-o.com/",
	"https://shibuya-o.com/east/schedule/",
	"tokyo",
	"shibuya",
	"shibuya-oeast",
	fetchers.TestInfo{
		NumberOfLives:         32,
		FirstLiveTitle:        "THE DEAD DAISIES",
		FirstLiveArtists:      []string{"THE DEAD DAISIES", "THE MIDNIGHT ROSES（Opening Act）"},
		FirstLivePrice:        "ADV 9,500（ご入場時ドリンク代別途600円）",
		FirstLivePriceEnglish: "ADV 9,500（ごWhen enteringDrinkNot included in ticket600円）",
		FirstLiveOpenTime:     time.Date(util.GetRelevantYear(11), 11, 1, 18, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(util.GetRelevantYear(11), 11, 1, 18, 45, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://shibuya-o.com/east/schedule/the-dead-daisies/",
	},
	35.658713,
	139.695609,
)

var ShibuyaONestFetcher = fetchers.CreateOFetcher(
	"https://shibuya-o.com/",
	"https://shibuya-o.com/nest/schedule/",
	"tokyo",
	"shibuya",
	"shibuya-onest",
	fetchers.TestInfo{
		NumberOfLives:         35,
		FirstLiveTitle:        "mayday-星が降るネストで-",
		FirstLiveArtists:      []string{"may in film", "PRSMIN", "zanka"},
		FirstLivePrice:        "ADV ¥2,500 DOOR ¥3,000 （ご入場時ドリンク代別途600円）",
		FirstLivePriceEnglish: "ADV ¥2,500 DOOR ¥3,000 （ごWhen enteringDrinkNot included in ticket600円）",
		FirstLiveOpenTime:     time.Date(util.GetRelevantYear(11), 11, 1, 18, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(util.GetRelevantYear(11), 11, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://shibuya-o.com/nest/schedule/may-in-film%e4%b8%bb%e5%82%acmayday-%e6%98%9f%e3%81%ae%e9%99%8d%e3%82%8b%e3%83%8d%e3%82%b9%e3%83%88%e3%81%a7-day1/",
	},
	35.658563,
	139.695313,
)

var ShibuyaOWestFetcher = fetchers.CreateOFetcher(
	"https://shibuya-o.com/",
	"https://shibuya-o.com/west/schedule/",
	"tokyo",
	"shibuya",
	"shibuya-owest",
	fetchers.TestInfo{
		NumberOfLives:         29,
		FirstLiveTitle:        "感覚ピエロ 10th ANNIVERSARY「感覚ピエロですがなにか」ツアー",
		FirstLiveArtists:      []string{"感覚ピエロ"},
		FirstLivePrice:        "ADV ¥4,500 (ご入場時ドリンク代別途600円)",
		FirstLivePriceEnglish: "ADV ¥4,500 (ごWhen enteringDrinkNot included in ticket600円)",
		FirstLiveOpenTime:     time.Date(util.GetRelevantYear(11), 11, 1, 18, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(util.GetRelevantYear(11), 11, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://shibuya-o.com/west/schedule/%e6%84%9f%e8%a6%9a%e3%83%94%e3%82%a8%e3%83%ad-10th-anniversary%e3%80%8c%e6%84%9f%e8%a6%9a%e3%83%94%e3%82%a8%e3%83%ad%e3%81%a7%e3%81%99%e3%81%8c%e3%81%aa%e3%81%ab%e3%81%8b%e3%80%8d%e3%83%84%e3%82%a2/",
	},
	35.658488,
	139.695297,
)

type shibuyaStrobeListResponse struct {
	Schedule struct {
		Year  string `json:"year"`
		Month string `json:"month"`
		List  map[string]struct {
			Event []struct {
				Artist   string `json:"artist"`
				Charge   string `json:"charge"`
				Id       string `json:"id"`
				OpenTime string `json:"opentime"`
				Title    string `json:"title"`
			} `json:"event"`
		} `json:"list"`
	} `json:"schedule"`
}

type shibuyaStrobeFlatElement struct {
	Year   string
	Month  string
	Day    string
	Time   string
	Artist string
	Price  string
	Id     string
	Title  string
}

func appendFlattenedStrobeList(a *[]shibuyaStrobeFlatElement, listResponse shibuyaStrobeListResponse) {
	keyInts := make([]int, 0)
	keys := make([]string, 0)
	for i := range listResponse.Schedule.List {
		n, err := strconv.Atoi(i)
		if err != nil {
			continue
		}
		keyInts = append(keyInts, n)
	}
	slices.Sort(keyInts)
	for i := range keyInts {
		keys = append(keys, strconv.Itoa(i))
	}

	for _, day := range keys {
		content := listResponse.Schedule.List[day]
		for _, event := range content.Event {
			if event.Id == "" {
				continue
			}
			*a = append(*a, shibuyaStrobeFlatElement{
				Year:   listResponse.Schedule.Year,
				Month:  listResponse.Schedule.Month,
				Day:    day,
				Time:   event.OpenTime,
				Artist: event.Artist,
				Price:  event.Charge,
				Id:     event.Id,
				Title:  event.Title,
			})
		}
	}
}

func getShibuyaStrobeLiveList(testDocument []byte) (list []shibuyaStrobeFlatElement) {
	list = make([]shibuyaStrobeFlatElement, 0)
	if testDocument != nil {
		var res shibuyaStrobeListResponse
		if err := json.Unmarshal(testDocument, &res); err != nil {
			return
		}
		appendFlattenedStrobeList(&list, res)
	} else {
		t := time.Now()
		year := t.Year() % 100
		month := int(t.Month())
		prevLength := -1

		for len(list) != prevLength {
			prevLength = len(list)
			var res shibuyaStrobeListResponse
			if err := util.GetJSON(
				fmt.Sprintf("https://www.strobe-cafe.com/schedule/get-data-schedule.php?y=20%d&m=%02d", year, month),
				&res,
			); err != nil {
				break
			}
			appendFlattenedStrobeList(&list, res)

			month++
			if month > 12 {
				month = 1
				year++
			}
		}
	}
	return
}

func createShibuyaStrobeHtml(live shibuyaStrobeFlatElement) string {
	content := "<body>"
	content += fmt.Sprintf(`<a href="https://www.strobe-cafe.com/schedule/%s/%s/%s.html"></a>`, live.Year, live.Month, live.Id)
	content += fmt.Sprintf(`<p id="date">%s/%s/%s</p>`, live.Year, live.Month, live.Day)
	content += fmt.Sprintf(`<p id="artist">%s</p>`, live.Artist)
	content += fmt.Sprintf(`<p id="price">%s</p>`, live.Price)
	content += fmt.Sprintf(`<p id="title">%s</p>`, live.Title)
	content += fmt.Sprintf(`<p id="time">%s</p>`, live.Time)
	content += "</body>"
	return content
}

// TODO: this is in harajuku
var ShibuyaStrobeFetcher = fetchers.Simple{
	BaseURL: "https://www.strobe-cafe.com/",
	LiveHTMLFetcher: func(testDocument []byte) (nodes []*html.Node, err error) {
		nodes = make([]*html.Node, 0)
		list := getShibuyaStrobeLiveList(testDocument)
		for _, live := range list {
			node, err := html.Parse(strings.NewReader(
				createShibuyaStrobeHtml(live),
			))
			if err != nil || node == nil {
				continue
			}
			nodes = append(nodes, node)
		}
		return
	},
	DetailsLinkSelector: "//a",
	TitleQuerier:        *htmlquerier.Q("//p[@id='title']"),
	ArtistsQuerier:      *htmlquerier.Q("//p[@id='artist']").SplitIgnoreWithin("/|(【.*?】)", '(', ')'),
	PriceQuerier:        *htmlquerier.Q("//p[@id='price']"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//p[@id='date']"),
		MonthQuerier:     *htmlquerier.Q("//p[@id='date']").After("/"),
		DayQuerier:       *htmlquerier.Q("//p[@id='date']").After("/").After("/"),
		OpenTimeQuerier:  *htmlquerier.Q("//p[@id='time']"),
		StartTimeQuerier: *htmlquerier.Q("//p[@id='time']").After("/"),

		IsYearInLive:  true,
		IsMonthInLive: true,
	},

	PrefectureName: "tokyo",
	AreaName:       "shibuya",
	VenueID:        "shibuya-strobe",
	Latitude:       35.671563,
	Longitude:      139.704437,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         16,
		FirstLiveTitle:        "古川本舗FC限定ライブ　昼公演【re:ROLE】",
		FirstLiveArtists:      []string{"古川本舗"},
		FirstLivePrice:        "前売 4500円 / 当日 5500円 (各+1drink order)",
		FirstLivePriceEnglish: "Reservation 4500円 / Door 5500円 (各+1drink order)",
		FirstLiveOpenTime:     time.Date(2024, 3, 2, 13, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 2, 13, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://www.strobe-cafe.com/schedule/2024/03/3128.html",
	},
}

type shibuyaTokioTokyoListResponse []struct {
	Document struct {
		Fields struct {
			Default struct {
				MapValue struct {
					Fields struct {
						List struct {
							ArrayValue struct {
								Values []struct {
									ReferenceValue string `json:"referenceValue"`
								} `json:"values"`
							} `json:"arrayValue"`
						} `json:"list"`
					} `json:"fields"`
				} `json:"mapValue"`
			} `json:"default"`
		} `json:"fields"`
	} `json:"document"`
}

type shibuyaTokioTokyoLiveJSON struct {
	Fields struct {
		Default struct {
			MapValue struct {
				Fields struct {
					Title struct {
						StringValue string `json:"stringValue"`
					} `json:"title"`
					Body struct {
						StringValue string `json:"stringValue"`
					} `json:"body"`
					Slug struct {
						StringValue string `json:"stringValue"`
					} `json:"slug"`
					Date struct {
						StringValue string `json:"stringValue"`
					} `json:"hMPvAUXl"` //wtf?
				} `json:"fields"`
			} `json:"mapValue"`
		} `json:"default"`
	} `json:"fields"`
}

func shibuyaTokioTokyoGetBody(body, title string) string {
	if body != "" {
		return body
	}
	return title
}

// please never break i do not want to ever deal with this again what the fuck is this abomination
var ShibuyaTokioTokyoFetcher = fetchers.Simple{
	BaseURL: "https://tokio.world/",
	LiveHTMLFetcher: func(testDocument []byte) (nodes []*html.Node, err error) {
		nodes = make([]*html.Node, 0)
		var list shibuyaTokioTokyoListResponse
		if testDocument == nil {
			if err = util.GetJSON(
				"https://api.cms.studiodesignapp.com/documents:runQuery?q=eyJzdHJ1Y3R1cmVkUXVlcnkiOnsiZnJvbSI6W3siY29sbGVjdGlvbklkIjoicHVibGlzaGVkIiwiYWxsRGVzY2VuZGFudHMiOnRydWV9XSwid2hlcmUiOnsiY29tcG9zaXRlRmlsdGVyIjp7Im9wIjoiQU5EIiwiZmlsdGVycyI6W3siZmllbGRGaWx0ZXIiOnsiZmllbGQiOnsiZmllbGRQYXRoIjoiX21ldGEucHJvamVjdC5pZCJ9LCJvcCI6IkVRVUFMIiwidmFsdWUiOnsic3RyaW5nVmFsdWUiOiIyNGMyMTZkOTUwY2U0OTY5YWU2ZiJ9fX0seyJmaWVsZEZpbHRlciI6eyJmaWVsZCI6eyJmaWVsZFBhdGgiOiJfbWV0YS5zY2hlbWEua2V5In0sIm9wIjoiRVFVQUwiLCJ2YWx1ZSI6eyJzdHJpbmdWYWx1ZSI6InplQ2FyaG5yIn19fV19fSwib3JkZXJCeSI6W3siZmllbGQiOnsiZmllbGRQYXRoIjoiX21ldGEucHVibGlzaGVkQXQifSwiZGlyZWN0aW9uIjoiREVTQ0VORElORyJ9XSwibGltaXQiOjF9fQ%3D%3D",
				&list,
			); err != nil {
				return
			}
		} else {
			if err = json.Unmarshal(testDocument, &list); err != nil {
				return
			}
		}
		for _, liveReference := range list[0].Document.Fields.Default.MapValue.Fields.List.ArrayValue.Values {
			var live shibuyaTokioTokyoLiveJSON
			if err = util.GetJSON(
				strings.Replace(liveReference.ReferenceValue, "projects/studio-7e371/databases/(default)/", "https://api.cms.studiodesignapp.com/", 1),
				&live,
			); err != nil {
				continue
			}
			if strings.Contains(live.Fields.Default.MapValue.Fields.Title.StringValue, "ホールレンタル") {
				continue
			}
			n, err := html.Parse(strings.NewReader(fmt.Sprintf(
				"<span id='date'>%s</span><span id='title'>%s</span><span id='body'>%s</span><span id='url'>%s</span>",
				live.Fields.Default.MapValue.Fields.Date.StringValue,
				live.Fields.Default.MapValue.Fields.Title.StringValue,
				shibuyaTokioTokyoGetBody(live.Fields.Default.MapValue.Fields.Body.StringValue, live.Fields.Default.MapValue.Fields.Title.StringValue),
				fmt.Sprintf("https://tokio.world/posts/%s", live.Fields.Default.MapValue.Fields.Slug.StringValue),
			)))
			if err != nil {
				continue
			}
			nodes = append(nodes, n)
		}
		return
	},
	DetailsLinkSelector: "//span[@id='url']",
	TitleQuerier:        *htmlquerier.Q("//span[@id='body']"),
	ArtistsQuerier:      *htmlquerier.Q("//span[@id='title']").After("】").Split("/"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:  *htmlquerier.Q("//span[@id='date']"),
		MonthQuerier: *htmlquerier.Q("//span[@id='date']").After("/"),
		DayQuerier:   *htmlquerier.Q("//span[@id='date']").After("/").After("/"),

		IsYearInLive:  true,
		IsMonthInLive: true,
	},

	PrefectureName: "tokyo",
	AreaName:       "shibuya",
	VenueID:        "shibuya-tokiotokyo",
	Latitude:       35.662562,
	Longitude:      139.698937,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         18,
		FirstLiveTitle:        "MUSIC NEXUS Presents Live 導 -SHIRUBE- Vol.0",
		FirstLiveArtists:      []string{"BESPER", "LUKA", "佐野諒太", "浜野はるき", "灯橙あか"},
		FirstLivePrice:        "このライブハウスのイベントの値段にアクセスできません。ライブのリンクをチェックしてください。",
		FirstLivePriceEnglish: "Cannot access prices for lives at this venue. Please check live link.",
		FirstLiveOpenTime:     time.Date(2024, 3, 12, 03, 24, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 12, 03, 24, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://tokio.world/posts/OMwpVbLU",
	},
}

var ShibuyaVeatsFetcher = fetchers.Simple{
	BaseURL:              "https://veats.jp/",
	ShortYearIterableURL: "https://veats.jp/schedule/?param=20%d%02d",
	LiveSelector:         "//div[@class='today-contents']/a",
	ExpandedLiveSelector: ".",
	TitleQuerier:         *htmlquerier.QAll("//p[@class='ttl']//text()[1]").KeepIndex(-1),
	ArtistsQuerier:       *htmlquerier.Q("//dt[.='LINE UP']/following-sibling::dd").SplitIgnoreWithin("/", '(', ')'),
	PriceQuerier:         *htmlquerier.Q("//dt[.='ADV /  DOOR']/following-sibling::dd"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//p[@class='year']"),
		MonthQuerier:     *htmlquerier.Q("//div[@class='head']/p"),
		DayQuerier:       *htmlquerier.Q("//div[@class='head']/p").After("/"),
		OpenTimeQuerier:  *htmlquerier.Q("//dt[.='OPEN / START']/following-sibling::dd"),
		StartTimeQuerier: *htmlquerier.Q("//dt[.='OPEN / START']/following-sibling::dd").After("/"),
		IsMonthInLive:    true,
	},

	PrefectureName: "tokyo",
	AreaName:       "shibuya",
	VenueID:        "shibuya-veats",
	Latitude:       35.660613,
	Longitude:      139.697641,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         32,
		FirstLiveTitle:        "BLACK IRIS WEEKLY LIVE",
		FirstLiveArtists:      []string{"BLACK IRIS"},
		FirstLivePrice:        "前方エリア¥4,000・後方エリア¥1,000 / 前方エリア¥4,500・後方エリア¥1,000 (D代別)",
		FirstLivePriceEnglish: "Front area¥4,000・Rear area¥1,000 / Front area¥4,500・Rear area¥1,000 (Drink must be purchased separately)",
		FirstLiveOpenTime:     time.Date(2024, 3, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 1, 19, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://veats.jp/schedule/5667/",
	},
}

var ShibuyaWWWFetcher = fetchers.CreateWWWFetcher(
	"@data-place='www' or @data-place='wwwxwww'",
	"shibuya-www",
	fetchers.TestInfo{
		NumberOfLives:         30,
		FirstLiveTitle:        "#ピコリフ 3rdワンマンライブ「Shiny Sparkle」",
		FirstLiveArtists:      []string{"#ピコリフ"},
		FirstLivePrice:        "VIP：¥5,000 / 撮影：¥3,000 / 後方：¥1,500 / 学生・女性：¥500 / 新規：¥0 (税込 / 各ドリンク代別)",
		FirstLivePriceEnglish: "VIP：¥5,000 / 撮影：¥3,000 / 後方：¥1,500 / Students・Women：¥500 / 新規：¥0 (Incl. Tax / 各DrinkSeparately)",
		FirstLiveOpenTime:     time.Date(2023, 11, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 11, 1, 19, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://www-shibuya.jp/schedule/017256.php",
	},
)

var ShibuyaWWWBetaFetcher = fetchers.CreateWWWFetcher(
	"@data-place='www_beta' or @data-place='wwwxwww'",
	"shibuya-wwwbeta",
	fetchers.TestInfo{
		NumberOfLives:         4,
		FirstLiveTitle:        "MONK WORK BASE",
		FirstLiveArtists:      []string{"仙人掌 & S-kaine", "LIl' Leise But Gold", "NEI", "anddytoystore", "AI Jacky", "KEYTOTHECITY", "Yusef Imamura & Sano", "Palmpark", "YU", "cut skateboards", "kizukush", "JON (UGLY WEAPON)", "Joe cupertino", "CHAPAH", "VOLOJZA", "interplay", "凸凹。", "AI.U", "Kazuhiko Fujita", "SEX 山口", "Wide Escapes", "MONK", "Hanaboo", "stedman", "YAMAPIZZA", "stadman"},
		FirstLivePrice:        "当日 ¥3,500 (税込｜スタンディング｜ドリンク代別)前売 ¥3,000 (税込｜スタンディング｜ドリンク代別)",
		FirstLivePriceEnglish: "Door ¥3,500 (Incl. Tax｜Standing｜DrinkSeparately)Reservation ¥3,000 (Incl. Tax｜Standing｜DrinkSeparately)",
		FirstLiveOpenTime:     time.Date(2023, 11, 7, 17, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 11, 7, 17, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://www-shibuya.jp/schedule/017137.php",
		IgnoreTest:            true,
	},
)

var ShibuyaWWWXFetcher = fetchers.CreateWWWFetcher(
	"@data-place='wwwx' or @data-place='www_x'",
	"shibuya-wwwx",
	fetchers.TestInfo{
		NumberOfLives:         27,
		FirstLiveTitle:        "BIG ROMANTIC RECORDS presents Carsick Cars live in Tokyo",
		FirstLiveArtists:      []string{"Carsick Cars", "Hello Shitty (a.k.a Sophia from UptownRecords)"},
		FirstLivePrice:        "¥5,000 / ¥5,500 (税込 / ドリンク代別)※当日券は19:00~、前売りの入場が落ち着き次第 WWW Xにて¥5,500+D代で販売いたします。",
		FirstLivePriceEnglish: "¥5,000 / ¥5,500 (Incl. Tax / DrinkSeparately)※Door券は19:00~、ReservationのEntryが落ち着き次第 WWW Xにて¥5,500+D代で販売いたします。",
		FirstLiveOpenTime:     time.Date(2023, 11, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 11, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://www-shibuya.jp/schedule/017258.php",
	},
)

var YoyogiBarbaraFetcher = fetchers.CreateOmatsuriFetcher(
	"https://barbara.omatsuri.tech/",
	"tokyo",
	"shibuya",
	"yoyogi-barbara",
	35.682688,
	139.699563,
	fetchers.TestInfo{
		NumberOfLives:         20,
		FirstLiveTitle:        "GOTCHA MIX!!!!",
		FirstLiveArtists:      []string{"森口しゅな", "ゆみちぃ", "みほりょうすけ", "スギムラリョウイチ", "金子TKO"},
		FirstLivePrice:        "ADV/DOOR 3,000/3,500 配信2,000",
		FirstLivePriceEnglish: "ADV/DOOR 3,000/3,500 Livestream2,000",
		FirstLiveOpenTime:     time.Date(2025, 3, 1, 18, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 3, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://barbara.omatsuri.tech/events/26505",
	},
)

/*******************
 * 								 *
 *	Shimokitazawa	 *
 *								 *
 *******************/

var ShimokitazawaArtistFetcher = fetchers.Simple{
	BaseURL:              "http://www.c-artist.com/",
	ShortYearIterableURL: "http://www.c-artist.com/schedule/list/20%d%d.txt",
	LiveSelector:         "//div[@class='sche']",
	TitleQuerier:         *htmlquerier.Q("//p[contains(@class, 'guestname')]/text()[1]"),
	ArtistsQuerier:       *htmlquerier.Q("//p[contains(@class, 'guestname')]/text()[last()]").Split("\u00A0").Before("』〜").TrimPrefix("『"),
	PriceQuerier:         *htmlquerier.Q("//p[@class='ex']").After(" / "),
	DetailsLink:          "http://www.c-artist.com/schedule/",

	TimeHandler: fetchers.TimeHandler{
		MonthQuerier:    *htmlquerier.Q("//p[@class='day']").Before("/"),
		DayQuerier:      *htmlquerier.Q("//p[@class='day']").After("/").Before("("),
		OpenTimeQuerier: *htmlquerier.Q("//p[@class='ex']").Before(" / ").After("open"),
		IsMonthInLive:   true,
		IsYearInLive:    true,
	},

	PrefectureName: "tokyo",
	AreaName:       "shimokitazawa",
	VenueID:        "shimokitazawa-artist",
	Latitude:       35.663562,
	Longitude:      139.668609,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         17,
		FirstLiveTitle:        "『jikkatsu_c:side』〜 two-man_live 〜",
		FirstLiveArtists:      []string{"いちろう", "酒井勇也"},
		FirstLivePrice:        "2,000yen（1ドリンク込み）",
		FirstLivePriceEnglish: "2,000yen（1DrinkIncluded）",
		FirstLiveOpenTime:     time.Date(util.GetRelevantYear(3), 3, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(util.GetRelevantYear(3), 3, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "http://www.c-artist.com/schedule/",
	},
}

var ShimokitazawaBasementBarFetcher = fetchers.CreateToosFetcher(
	"https://www.toos.co.jp/",
	"https://toos.co.jp/basementbar/event/on/20%d/%02d/",
	"tokyo",
	"shimokitazawa",
	"shimokitazawa-basementbar",
	fetchers.TestInfo{
		NumberOfLives:         31,
		FirstLiveTitle:        "sickufo fanclub",
		FirstLiveArtists:      []string{"あずき過保護セット", "小林壮侍", "年齢バンド", "まほろば", "Bambi club", "GRASAM ANIMAL", "Laget’s Jam Stack", "sokkuriufo (from Khaki)", "sickufo(60min)"},
		FirstLivePrice:        "ADV￥1,800-（+2D）",
		FirstLivePriceEnglish: "ADV￥1,800-（+2D）",
		FirstLiveOpenTime:     time.Date(2024, 3, 1, 17, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 1, 17, 50, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://toos.co.jp/basementbar/ev/sickufo-fanclub/",
	},
)

var ShimokitazawaChikamatsuFetcher = fetchers.Simple{
	BaseURL:              "https://chikamatsu-nite.com/",
	InitialURL:           "https://chikamatsu-nite.com/schedule/",
	LiveSelector:         "//ul[@id='event-list']/li",
	ExpandedLiveSelector: "//a",
	TitleQuerier:         *htmlquerier.Q("//div[@id='event-title']//h2"),
	ArtistsQuerier:       *htmlquerier.Q("//dl[@id='event-view']/dd[3]").Split("/"),
	PriceQuerier:         *htmlquerier.Q("//dl[@id='event-view']/dd[2]"),

	TimeHandler: fetchers.TimeHandler{
		MonthQuerier:     *htmlquerier.Q("//div[@id='event-title']//time").Before("/"),
		DayQuerier:       *htmlquerier.Q("//div[@id='event-title']//time").After("/"),
		OpenTimeQuerier:  *htmlquerier.Q("//dl[@id='event-view']/dd[1]").Before("/"),
		StartTimeQuerier: *htmlquerier.Q("//dl[@id='event-view']/dd[1]").After("/"),
		IsMonthInLive:    true,
		IsYearInLive:     true,
	},

	PrefectureName: "tokyo",
	AreaName:       "shimokitazawa",
	VenueID:        "shimokitazawa-chikamatsu",
	Latitude:       35.656813,
	Longitude:      139.667562,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         26,
		FirstLiveTitle:        "地下教室 vol.13",
		FirstLiveArtists:      []string{"uyuu", "アウフヘーベン", "シンクマクラ", "憂牡丹"},
		FirstLivePrice:        "¥2,000-/¥2,500-(＋1D)",
		FirstLivePriceEnglish: "¥2,000-/¥2,500-(＋1D)",
		FirstLiveOpenTime:     time.Date(util.GetRelevantYear(11), 11, 1, 18, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(util.GetRelevantYear(11), 11, 1, 18, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://chikamatsu-nite.com/schedule/2023/11/01-id2504.php",
	},
}

var ShimokitazawaChikamichiFetcher = fetchers.CreateChikamichiFetcher(
	"https://chikamichi-otemae.com/",
	"https://chikamichi-otemae.com/chikamichi/",
	"tokyo",
	"shimokitazawa",
	"shimokitazawa-chikamichi",
	fetchers.TestInfo{
		NumberOfLives:         24,
		FirstLiveTitle:        "VINTAGE ROCK × チケットぴあ presents FUTURE RESERVE vol.5",
		FirstLiveArtists:      []string{"omeme tenten", "すなお", "daisansei", "カラコルムの山々"},
		FirstLivePrice:        "¥2,900(+1D¥600)",
		FirstLivePriceEnglish: "¥2,900(+1D¥600)",
		FirstLiveOpenTime:     time.Date(2023, 11, 2, 18, 45, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 11, 2, 19, 15, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://chikamichi-otemae.com/chikamichi/852/",
	},
)

var ShimokitazawaClub251Fetcher = fetchers.Simple{
	BaseURL:              "https://www.club251.com/",
	ShortYearIterableURL: "https://club251.com/schedule/?yr=20%d&mo=%d",
	LiveSelector:         "//div[@class='schedule-in']/div[@class=' eventful' or @class=' eventful-today']/table",
	MultiLiveDaySelector: "//table[@class='about-in' and .//h2/text()!='PRIVATE']",
	TitleQuerier:         *htmlquerier.Q("//h2"),
	ArtistsQuerier:       *htmlquerier.Q("//h2/following-sibling::p[@class='fw-bold']").Split("／"),
	PriceQuerier:         *htmlquerier.Q("//text()[contains(., 'CHARGE : ')]").After("CHARGE : "),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//div[@class='schedulebox-solo']/div[@class='clearfix']//h3"),
		MonthQuerier:     *htmlquerier.Q("//div[@class='schedulebox-solo']/div[@class='clearfix']//h3").After("年"),
		DayQuerier:       *htmlquerier.Q("//th"),
		OpenTimeQuerier:  *htmlquerier.Q("//text()[contains(., 'OPEN ') and contains(., 'START ')]"),
		StartTimeQuerier: *htmlquerier.Q("//text()[contains(., 'OPEN ') and contains(., 'START ')]").After("START"),
	},

	PrefectureName: "tokyo",
	AreaName:       "shimokitazawa",
	VenueID:        "shimokitazawa-club251",
	Latitude:       35.658313,
	Longitude:      139.667312,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         32,
		FirstLiveTitle:        "INFINITY LIVE presents 『ULTLA POP BUBBLE』",
		FirstLiveArtists:      []string{"遥か、彼方。", "Goodbye for First kiss", "開歌-かいか-", "月刊PAM", "美味しい水玉", "かわいいからって甘くみないで", "Noreco", "ponderosa may bloom"},
		FirstLivePrice:        "前売 ¥2,300- / 当日 ¥2,800- (+1D)",
		FirstLivePriceEnglish: "Reservation ¥2,300- / Door ¥2,800- (+1D)",
		FirstLiveOpenTime:     time.Date(2025, 2, 1, 10, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 2, 1, 10, 15, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://club251.com/schedule/?yr=20%d&mo=%d",
	},
}

var Shimokitazawa440Fetcher = fetchers.Simple{
	BaseURL:              "http://440.tokyo/",
	InitialURL:           "http://440.tokyo/schedule/",
	NextSelector:         "//div[@class='month']/div[@class='next']/a",
	ExpandedLiveSelector: "//a",
	LiveSelector:         "//div[@class='schedule-list__list']/article",
	TitleQuerier:         *htmlquerier.Q("//h2"),
	ArtistsQuerier:       *htmlquerier.Q("//h1").Split("｜"),
	PriceQuerier:         *htmlquerier.Q("//dl[@class='schedule-content__ticket']//p/text()[2]"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//span[@class='date']").Before("/"),
		MonthQuerier:     *htmlquerier.Q("//span[@class='date']").SplitIndex("/", 1),
		DayQuerier:       *htmlquerier.Q("//span[@class='date']").SplitIndex("/", 2).Before("("),
		OpenTimeQuerier:  *htmlquerier.Q("//dl[@class='schedule-content__openstart']//p").Before("／"),
		StartTimeQuerier: *htmlquerier.Q("//dl[@class='schedule-content__openstart']//p").After("／"),
		IsMonthInLive:    true,
		IsYearInLive:     true,
	},

	PrefectureName: "tokyo",
	AreaName:       "shimokitazawa",
	VenueID:        "shimokitazawa-440",
	Latitude:       35.658313,
	Longitude:      139.667391,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         29,
		FirstLiveTitle:        "“Stormy Wednesday”",
		FirstLiveArtists:      []string{"王将&The Guv’nor Brothers", "Weiyao Band", "前後のカルマ"},
		FirstLivePrice:        "ADV.￥3,000／DOOR.￥3,400 [1D別]",
		FirstLivePriceEnglish: "ADV.￥3,000／DOOR.￥3,400 [1 Drink purchase required]",
		FirstLiveOpenTime:     time.Date(2023, 11, 1, 18, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 11, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "http://440.tokyo/events/231101/",
	},
}

var ShimokitazawaClubQueFetcher = fetchers.Simple{
	BaseURL:              "https://clubque.net/",
	InitialURL:           "https://clubque.net/schedule/",
	NextSelector:         "//div[@class='month']/div[@class='next']/a",
	ExpandedLiveSelector: "//a",
	LiveSelector:         "//div[@class='schedule-list__list']/article",
	TitleQuerier:         *htmlquerier.Q("//h2"),
	ArtistsQuerier:       *htmlquerier.Q("//h1").Split("｜"),
	PriceQuerier:         *htmlquerier.Q("//dl[@class='schedule-content__ticket']//p/text()[2]"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//span[@class='date']").Before("/"),
		MonthQuerier:     *htmlquerier.Q("//span[@class='date']").SplitIndex("/", 1),
		DayQuerier:       *htmlquerier.Q("//span[@class='date']").SplitIndex("/", 2).Before("("),
		OpenTimeQuerier:  *htmlquerier.Q("//dl[@class='schedule-content__openstart']//p").Before("／"),
		StartTimeQuerier: *htmlquerier.Q("//dl[@class='schedule-content__openstart']//p").After("／"),
		IsMonthInLive:    true,
		IsYearInLive:     true,
	},

	PrefectureName: "tokyo",
	AreaName:       "shimokitazawa",
	VenueID:        "shimokitazawa-clubque",
	Latitude:       35.660938,
	Longitude:      139.668813,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         29,
		FirstLiveTitle:        "“Stormy Wednesday”",
		FirstLiveArtists:      []string{"王将&The Guv’nor Brothers", "Weiyao Band", "前後のカルマ"},
		FirstLivePrice:        "ADV.￥3,000／DOOR.￥3,400 [1D別]",
		FirstLivePriceEnglish: "ADV.￥3,000／DOOR.￥3,400 [1 Drink purchase required]",
		FirstLiveOpenTime:     time.Date(2023, 11, 1, 18, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 11, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://clubque.net/schedule/2202/",
	},
}

var ShimokitazawaDaisyBarFetcher = fetchers.CreateDaisyBarFetcher(
	"https://daisybar.jp/",
	"https://daisybar.jp/schedule/20%d/%d/",
	"tokyo",
	"shimokitazawa",
	"shimokitazawa-daisybar",
	fetchers.TestInfo{
		NumberOfLives:         33,
		FirstLiveTitle:        "Paper Bag presents 「自分ってなんなんだっけ」今日くらいはいいでしょ爆酒Night",
		FirstLiveArtists:      []string{"Paper Bag", "fountin", "Maju2", "詩野"},
		FirstLivePrice:        "前売 ¥2500 / 当日 ¥3000",
		FirstLivePriceEnglish: "Reservation ¥2500 / Door ¥3000",
		FirstLiveOpenTime:     time.Date(2025, 2, 1, 18, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 2, 1, 18, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          util.InsertYearMonth("https://daisybar.jp/schedule/%d/%d/"),
	},
)

var ShimokitazawaDyCubeFetcher = fetchers.Simple{
	BaseURL:              "https://dycube.tokyo/",
	ShortYearIterableURL: "https://dycube.tokyo/schedule/?ext_num-year=20%d&ext_num-month=%02d",
	LiveSelector:         "//article[@class='schedule-article']",
	TitleQuerier:         *htmlquerier.Q("//h3"),
	ArtistsQuerier:       *htmlquerier.Q("//div[@class='schedule-article-body__txt']/p[1]").Split("\n"),
	PriceQuerier:         *htmlquerier.Q("//div[@class='schedule-article-body__txt']/p[2]/text()[2]"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//a[contains(@class, 'schedule-head-year__btn')][contains(@class, '-active')]"),
		MonthQuerier:     *htmlquerier.Q("//a[contains(@class, 'schedule-head-month__btn')][contains(@class, '-active')]"),
		DayQuerier:       *htmlquerier.Q("//p[@class='schedule-article-head__day']"),
		OpenTimeQuerier:  *htmlquerier.Q("//div[@class='schedule-article-body__txt']/p[2]/text()[1]").Before("/"),
		StartTimeQuerier: *htmlquerier.Q("//div[@class='schedule-article-body__txt']/p[2]/text()[1]").After("/"),
	},

	PrefectureName: "tokyo",
	AreaName:       "shimokitazawa",
	VenueID:        "shimokitazawa-dycube",
	Latitude:       35.662037,
	Longitude:      139.666609,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         34,
		FirstLiveTitle:        "vessy 初企画 「diamonddust」",
		FirstLiveArtists:      []string{"nene", "栢本ての", "藍谷凪", "vessy"},
		FirstLivePrice:        "TICKET ¥2,500(+1drink¥600)",
		FirstLivePriceEnglish: "TICKET ¥2,500(+1drink¥600)",
		FirstLiveOpenTime:     time.Date(2023, 11, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 11, 1, 19, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://dycube.tokyo/schedule/?ext_num-year=20%d&ext_num-month=%02d",
	},
}

var ShimokitazawaEraFetcher = fetchers.Simple{
	BaseURL:        "http://s-era.jp/",
	InitialURL:     util.InsertYearMonth("http://s-era.jp/schedule_cat/%d-%02d/"),
	NextSelector:   "//section[contains(@class, 'schedule-navigation')]/div[2]/p[2]/a",
	LiveSelector:   "//article[contains(@class, 'schedule-box')]",
	TitleQuerier:   *htmlquerier.Q("//h4"),
	ArtistsQuerier: *htmlquerier.Q("//div[contains(@class, 'w-flyer')]").BeforeSelector("//div[contains(@class, 'detail-texts')]").Before("\n\n").Before("/OPEN ").SplitIgnoreWithin("/", '（', '）').Trim().TrimSuffix("/"),
	PriceQuerier:   *htmlquerier.Q("//div[contains(@class, 'detail-grid')]/p[2]"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//section[contains(@class, 'schedule-navigation')]//h3").Before("."),
		MonthQuerier:     *htmlquerier.Q("//time/text()").Before("."),
		DayQuerier:       *htmlquerier.Q("//time/text()").After("."),
		OpenTimeQuerier:  *htmlquerier.Q("//div[contains(@class, 'detail-grid')]/p[1]/span[2]"),
		StartTimeQuerier: *htmlquerier.Q("//div[contains(@class, 'detail-grid')]/p[1]/span[4]"),
		IsMonthInLive:    true,
	},

	PrefectureName: "tokyo",
	AreaName:       "shimokitazawa",
	VenueID:        "shimokitazawa-era",
	Latitude:       35.663313,
	Longitude:      139.668313,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         38,
		FirstLiveTitle:        "ERA presents. A Day in The Life Vol.159",
		FirstLiveArtists:      []string{"paddy isle", "Akarumeno brown", "Blueberry Mondays", "Qurukuma", "Spinning Plums"},
		FirstLivePrice:        "ADV ¥2000DOOR ¥2500 (+1D¥600)",
		FirstLivePriceEnglish: "ADV ¥2000DOOR ¥2500 (+1D¥600)",
		FirstLiveOpenTime:     time.Date(2023, 11, 2, 18, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 11, 2, 18, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          util.InsertYearMonth("http://s-era.jp/schedule_cat/%d-%02d/"),
	},
}

var ShimokitazawaFlowersLoftFetcher = fetchers.CreateLoftFetcher(
	"https://www.loft-prj.co.jp/",
	"https://www.loft-prj.co.jp/schedule/flowersloft/schedule?scheduleyear=20%d&schedulemonth=%d",
	"tokyo",
	"shimokitazawa",
	"shimokitazawa-flowersloft",
	fetchers.TestInfo{
		NumberOfLives:         34,
		FirstLiveTitle:        "Terrafirma",
		FirstLiveArtists:      []string{"aruga", "くぐり", "Yellow mo", "乙女絵画", "Sorry No Camisole", "[DJ] myein", "KARIN"},
		FirstLivePrice:        "ADV.DOOR ¥2,400(+1D)",
		FirstLivePriceEnglish: "ADV.DOOR ¥2,400(+1D)",
		FirstLiveOpenTime:     time.Date(2025, 2, 1, 23, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 2, 1, 23, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://www.loft-prj.co.jp/schedule/flowersloft/300436",
	},
	35.662188,
	139.667937,
)

var ShimokitazawaLagunaFetcher = fetchers.CreateOldDaisyBarFetcher(
	"https://s-laguna.jp/",
	"https://s-laguna.jp/events/event/on/20%d/%02d/",
	"tokyo",
	"shimokitazawa",
	"shimokitazawa-laguna",
	"color-white",
	fetchers.TestInfo{
		NumberOfLives:         38,
		FirstLiveTitle:        "BRIDGE",
		FirstLiveArtists:      []string{"ミアナシメント", "金田すなほ", "リョマル", "ANAIS"},
		FirstLivePrice:        "前売 3000円(D別) 当日 3500円(D別)",
		FirstLivePriceEnglish: "Reservation 3000円(Drinks sold separately) Door 3500円(Drinks sold separately)",
		FirstLiveOpenTime:     time.Date(2023, 11, 1, 18, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 11, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://s-laguna.jp/events/event/on/20%d/%02d/",
	},
)

var ShimokitazawaLiveHausFetcher = fetchers.Simple{
	BaseURL:             "https://livehaus.jp/",
	InitialURL:          "https://livehaus.jp/schedule/",
	NextSelector:        "//a[contains(@class, 'tribe-events-c-nav__next')]",
	LiveSelector:        "//article[contains(@class, 'tribe-events-calendar-list__event')]",
	TitleQuerier:        *htmlquerier.Q("//h3"),
	DetailQuerier:       *htmlquerier.Q("//div[contains(@class, 'tribe-events-calendar-list__event-description')]").PreserveWhitespace(),
	DetailsLinkSelector: "//a[contains(@class, 'tribe-events-calendar-list__event-featured-image-link')]",

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:   *htmlquerier.Q("//time[contains(@class, 'tribe-events-calendar-list__month-separator-text')]").After("月"),
		MonthQuerier:  *htmlquerier.Q("//span[contains(@class, 'tribe-event-date-start')]").Before("/"),
		DayQuerier:    *htmlquerier.Q("//span[contains(@class, 'tribe-event-date-start')]").After("/"),
		IsMonthInLive: true,
	},

	PrefectureName: "tokyo",
	AreaName:       "shimokitazawa",
	VenueID:        "shimokitazawa-livehaus",
	Latitude:       35.659562,
	Longitude:      139.667562,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         19,
		FirstLiveTitle:        "突撃! 隣のトシクラブ vol.4 in 下北沢",
		FirstLiveArtists:      []string{"ドトキンズ", "Los Rancheros", "The Silver Sonics", "SHOWA BOYZ feat.DJ NIGARA", "SHJ (LONDON NITE)", "TAGO!", "KAZU SUDO (Caribbean Dandy)", "内藤啓介 (Chingcame)", "NIGARA (GARA)", "ITA (Nat Records)", "TSUNE (Lewis Leathers Japan)", "Mr.X(A&Y)", "ciibo", "KOSEI", "Yuuki(Tip Clothing&co.)", "RYO", "Shunsuke", "Shima Volume(104club)", "TOSHI & VOW"},
		FirstLivePrice:        "ADM¥2000",
		FirstLivePriceEnglish: "ADM¥2000",
		FirstLiveOpenTime:     time.Date(2023, 11, 3, 15, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 11, 3, 15, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://livehaus.jp/event/%e7%aa%81%e6%92%83-%e9%9a%a3%e3%81%ae%e3%83%88%e3%82%b7%e3%82%af%e3%83%a9%e3%83%96-vol-4-in-%e4%b8%8b%e5%8c%97%e6%b2%a2/",
	},
}

var ShimokitazawaLiveHolicFetcher = fetchers.Simple{
	BaseURL:             "https://liveholic.jp/",
	InitialURL:          "https://liveholic.jp/schedule/",
	LiveSelector:        "//div[@class='schedulegroup']/dl",
	DetailsLinkSelector: "//h2/a",
	TitleQuerier:        *htmlquerier.Q("//h2"),
	ArtistsQuerier:      *htmlquerier.Q("//p[@class='detail'][1]").Split(" / "),
	PriceQuerier:        *htmlquerier.Q("//p[@class='detail'][3]"),

	TimeHandler: fetchers.TimeHandler{
		MonthQuerier:     *htmlquerier.Q("//p[@class='date']").Before("."),
		DayQuerier:       *htmlquerier.Q("//p[@class='date']/text()").After("."),
		OpenTimeQuerier:  *htmlquerier.Q("//p[@class='detail'][2]/text()[1]").Before("/"),
		StartTimeQuerier: *htmlquerier.Q("//p[@class='detail'][2]/text()[2]").Before("/"),

		IsMonthInLive: true,
		IsYearInLive:  true,
	},

	PrefectureName: "tokyo",
	AreaName:       "shimokitazawa",
	VenueID:        "shimokitazawa-liveholic",
	Latitude:       35.661412,
	Longitude:      139.669078,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         52,
		FirstLiveTitle:        "オレンジバンドシャトル",
		FirstLiveArtists:      []string{"ザストーンマスター", "おじさん的思考ディスク", "BamblueSiA", "Sonne", "Route4th", "GO SEE REGRET", "※EiNyは出演キャンセル"},
		FirstLivePrice:        "一般¥2,000 学生¥1,400 (D別)",
		FirstLivePriceEnglish: "Ordinary Ticket¥2,000 Students¥1,400 (Drinks sold separately)",
		FirstLiveOpenTime:     time.Date(util.GetRelevantYear(11), 11, 6, 17, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(util.GetRelevantYear(11), 11, 6, 18, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://liveholic.jp/schedule/2023/11/post-455.php",
	},
}

var ShimokitazawaMonaRecordsFetcher = fetchers.Simple{
	BaseURL:              "https://www.mona-records.com/",
	ShortYearIterableURL: "https://www.mona-records.com/date/20%d/%02d/?category_name=livespace",
	LiveSelector:         "//div[contains(@class, 'sidepostlist-item')]",
	ExpandedLiveSelector: "//a",
	TitleQuerier:         *htmlquerier.Q("//h1"),
	ArtistsQuerier:       *htmlquerier.Q("//td[text()='出演']/following-sibling::td").Split(" / "),
	PriceQuerier:         *htmlquerier.Q("//td[text()='料金']/following-sibling::td"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//time").Before("/"),
		MonthQuerier:     *htmlquerier.Q("//time").SplitIndex("/", 1),
		DayQuerier:       *htmlquerier.Q("//time").SplitIndex("/", 2).Before("("),
		OpenTimeQuerier:  *htmlquerier.Q("//td[text()='時間']/following-sibling::td").Before("/"),
		StartTimeQuerier: *htmlquerier.Q("//td[text()='時間']/following-sibling::td").After("/"),

		IsYearInLive:  true,
		IsMonthInLive: true,
	},

	PrefectureName: "tokyo",
	AreaName:       "shimokitazawa",
	VenueID:        "shimokitazawa-monarecords",
	Latitude:       35.660463,
	Longitude:      139.667516,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         30,
		FirstLiveTitle:        "LOLIPOP",
		FirstLiveArtists:      []string{"じゆりぴ", "若杉果歩", "イロハマイ"},
		FirstLivePrice:        "前売：¥2,900+1Drink",
		FirstLivePriceEnglish: "Reservation：¥2,900+1Drink",
		FirstLiveOpenTime:     time.Date(2023, 11, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 11, 1, 19, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://www.mona-records.com/livespace/17026/",
	},
}

var ShimokitazawaMosaicFetcher = fetchers.Simple{
	BaseURL:        "http://mu-seum.co.jp",
	InitialURL:     "http://mu-seum.co.jp/schedule.html",
	NextSelector:   "//span[contains(@class, 'calendar_next')]/a",
	LiveSelector:   "//table[contains(@class, 'listCal')]",
	TitleQuerier:   *htmlquerier.Q("//tr[1]/td/p[1]"),
	ArtistsQuerier: *htmlquerier.Q("//td[contains(@class, 'live_menu')]/strong").SplitIgnoreWithin(" / ", '（', '）'),
	PriceQuerier:   *htmlquerier.Q("//td[contains(@class, 'live_menu')]/text()[5]"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//td[contains(@class, 'month_eng')]/text()[1]"),
		MonthQuerier:     *htmlquerier.Q("//tr[1]/th").Before("/"),
		DayQuerier:       *htmlquerier.Q("//tr[1]/th").After("/"),
		OpenTimeQuerier:  *htmlquerier.Q("//td[contains(@class, 'live_menu')]/text()[4]").After("OPEN "),
		StartTimeQuerier: *htmlquerier.Q("//td[contains(@class, 'live_menu')]/text()[4]").After("START "),

		IsMonthInLive: true,
	},

	PrefectureName: "tokyo",
	AreaName:       "shimokitazawa",
	VenueID:        "shimokitazawa-mosaic",
	Latitude:       35.659688,
	Longitude:      139.668563,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         31,
		FirstLiveTitle:        "『MOSAiC iDOL INFINITY - One hundred -』",
		FirstLiveArtists:      []string{"エレクトリックリボン", "月刊PAM", "セカイシティ", "LUNCH KIDS", "AKUMATICA"},
		FirstLivePrice:        "前売 ¥100 / 当日 ¥1,100(+1Drink ¥600)",
		FirstLivePriceEnglish: "Reservation ¥100 / Door ¥1,100(+1Drink ¥600)",
		FirstLiveOpenTime:     time.Date(2023, 6, 1, 18, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 6, 1, 18, 15, 0, 0, util.JapanTime),
		FirstLiveURL:          "http://mu-seum.co.jp/schedule.html",
	},
}

var ShimokitazawaOtemaeFetcher = fetchers.CreateChikamichiFetcher(
	"https://chikamichi-otemae.com/",
	"https://chikamichi-otemae.com/otemae/",
	"tokyo",
	"shimokitazawa",
	"shimokitazawa-otemae",
	fetchers.TestInfo{
		NumberOfLives:         6,
		FirstLiveTitle:        "Otemae Vol.1",
		FirstLiveArtists:      []string{"ELLY", "実香", "弓詩"},
		FirstLivePrice:        "ADV : ¥2,000(+1D¥600) / DAY : ¥2,500(+1D¥600)",
		FirstLivePriceEnglish: "ADV : ¥2,000(+1D¥600) / DAY : ¥2,500(+1D¥600)",
		FirstLiveOpenTime:     time.Date(2023, 10, 1, 18, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 10, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://chikamichi-otemae.com/otemae/963/",
		KnownEmpty:            true,
	},
)

var ShimokitazawaRegFetcher = fetchers.Simple{
	BaseURL:        "https://www.reg-r2.com/",
	InitialURL:     "https://www.reg-r2.com/?page_id=7250",
	NextSelector:   "//a[contains(text(), '次の月')]",
	LiveSelector:   "//table[@id='live_date']//tr",
	TitleQuerier:   *htmlquerier.Q("//div[@class='live_title']"),
	ArtistsQuerier: *htmlquerier.Q("//div[@class='performer_name']").Split("  "), // no idea why \n doesnt work here
	PriceQuerier:   *htmlquerier.Q("//div[@class='time_price']/text()[2]"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//h2").Before("年"),
		MonthQuerier:     *htmlquerier.Q("//h2").After("年").Before("月"),
		DayQuerier:       *htmlquerier.Q("//td[@class='date_weekday']/text()[1]").After("/"),
		OpenTimeQuerier:  *htmlquerier.Q("//div[@class='time_price']/text()[1]").SplitIndex("/", 1),
		StartTimeQuerier: *htmlquerier.Q("//div[@class='time_price']/text()[1]").SplitIndex("/", 2),
	},

	PrefectureName: "tokyo",
	AreaName:       "shimokitazawa",
	VenueID:        "shimokitazawa-reg",
	Latitude:       35.658187,
	Longitude:      139.667937,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         27,
		FirstLiveTitle:        "HAWK ANARCHY 2nd Album 宇宙侵略紀行release tour「日本侵略紀行」",
		FirstLiveArtists:      []string{"HAWK ANARCHY", "New Age Core", "OFELIA", "ALL I WANT"},
		FirstLivePrice:        "前売り￥2,500-（ドリンク別）/当日￥3,000-（ドリンク別）",
		FirstLivePriceEnglish: "Reservation￥2,500-（Drinks sold separately）/Door￥3,000-（Drinks sold separately）",
		FirstLiveOpenTime:     time.Date(2023, 11, 2, 18, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 11, 2, 18, 45, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://www.reg-r2.com/?page_id=7250",
	},
}

var ShimokitazawaShangrilaFetcher = fetchers.Simple{
	BaseURL:        "https://www.shan-gri-la.jp/",
	InitialURL:     "https://www.shan-gri-la.jp/tokyo/category/schedule/",
	LiveSelector:   "//div[@id='content']/div[contains(@class, 'hentry')]",
	TitleQuerier:   *htmlquerier.Q("//strong"),
	ArtistsQuerier: *htmlquerier.Q("//div[@class='post-content-content']/p[2]").Split("\n"),
	PriceQuerier:   *htmlquerier.Q("//div[@class='post-content-content']/p[3]").After("START ").After("\n"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//table[@id='wp-calendar']/caption").Before("年"),
		MonthQuerier:     *htmlquerier.Q("//h2").Before("/"),
		DayQuerier:       *htmlquerier.Q("//h2").After("/").Before("("),
		OpenTimeQuerier:  *htmlquerier.Q("//div[@class='post-content-content']/p[3]/text()[1]"),
		StartTimeQuerier: *htmlquerier.Q("//div[@class='post-content-content']/p[3]/text()[1]").After("/"),
		IsMonthInLive:    true,
	},

	PrefectureName: "tokyo",
	AreaName:       "shimokitazawa",
	VenueID:        "shimokitazawa-shangrila",
	Latitude:       35.660488,
	Longitude:      139.668516,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         5,
		FirstLiveTitle:        "The Cheserasera 2024 春の邂逅ツーマンライブ 其の壱",
		FirstLiveArtists:      []string{"The Cheserasera", "ABSTRACT MASH"},
		FirstLivePrice:        "前売￥4,400 / 当日￥4,900 （1ドリンク￥600別）",
		FirstLivePriceEnglish: "Reservation￥4,400 / Door￥4,900 （1Drink￥600Separately）",
		FirstLiveOpenTime:     time.Date(2024, 4, 4, 19, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 4, 4, 19, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://www.shan-gri-la.jp/tokyo/category/schedule/",
	},
}

var ShimokitazawaShelterFetcher = fetchers.CreateLoftFetcher(
	"https://www.loft-prj.co.jp",
	"https://www.loft-prj.co.jp/schedule/shelter/schedule?scheduleyear=20%d&schedulemonth=%d",
	"tokyo",
	"shimokitazawa",
	"shimokitazawa-shelter",
	fetchers.TestInfo{
		NumberOfLives:         37,
		FirstLiveTitle:        "Get into gear",
		FirstLiveArtists:      []string{"LIQUID SCREEN", "TREEBERRYS", "KOJI OZAKI メンバー:", "尾﨑 浩治(ex.Samantha’s Favourite, ex.BOYCE)", "嶋村 輝之(赤い夕陽, ex.Samantha’s Favourite )", "岩渕 尚史(Sloppy Joe, ex.BOYCE)", "加嶋 幸平(the SUN, ex.ROCKBOTTOM)", "DJ: 9232atfr"},
		FirstLivePrice:        "ADV¥3000 / DOOR¥3500",
		FirstLivePriceEnglish: "ADV¥3000 / DOOR¥3500",
		FirstLiveOpenTime:     time.Date(2025, 2, 1, 18, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 2, 1, 18, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://www.loft-prj.co.jp/schedule/shelter/301422",
	},
	35.661488,
	139.669453,
)

var ShimokitazawaSpreadFetcher = fetchers.Simple{
	BaseURL:        "https://spread.tokyo/",
	InitialURL:     "https://spread.tokyo/schedule.html",
	LiveSelector:   "//div[@id='c7']/div[@class='box'][position()<200]", // not sure why 200 leads to 99 matches but it does
	TitleQuerier:   *htmlquerier.Q("//u/b").CutWrapper(`"`, `"`),
	ArtistsQuerier: *htmlquerier.QAll("//div/span[last()]/text()"),
	PriceQuerier:   *htmlquerier.Q("//div/span[4]"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//div/span[1]"),
		MonthQuerier:     *htmlquerier.Q("//div/span[1]").After("."),
		DayQuerier:       *htmlquerier.Q("//div/span[1]").After(".").After("."),
		OpenTimeQuerier:  *htmlquerier.Q("//div/span[3]"),
		StartTimeQuerier: *htmlquerier.Q("//div/span[3]").After("START"),

		IsYearInLive:  true,
		IsMonthInLive: true,
	},

	PrefectureName: "tokyo",
	AreaName:       "shimokitazawa",
	VenueID:        "shimokitazawa-spread",
	Latitude:       35.660838,
	Longitude:      139.667734,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         99,
		FirstLiveTitle:        `techmoris presents "feel free...!I feel free...! Release Party"`,
		FirstLiveArtists:      []string{"ラジカセ狂気", "Fanny Hill", "techmoris", "139????"},
		FirstLivePrice:        "ADV. ¥2,000 | DOOR. ¥2,500 | U23. ¥1,500 (+1D)",
		FirstLivePriceEnglish: "ADV. ¥2,000 | DOOR. ¥2,500 | U23. ¥1,500 (+1D)",
		FirstLiveOpenTime:     time.Date(2024, 4, 5, 19, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 4, 5, 20, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://spread.tokyo/schedule.html",
	},
}

var ShimokitazawaThreeFetcher = fetchers.CreateToosFetcher(
	"https://www.toos.co.jp/",
	"https://www.toos.co.jp/3/events/event/on/20%d/%02d/",
	"tokyo",
	"shimokitazawa",
	"shimokitazawa-three",
	fetchers.TestInfo{
		NumberOfLives:         31,
		FirstLiveTitle:        "D.B.Inches 1st EP ”Instinct, Filter Bubble” Release Party",
		FirstLiveArtists:      []string{"aruga", "HOT DOG LOVE", "Nenne", "Suzuki Ryuto", "D.B.Inches"},
		FirstLivePrice:        "ADV￥2,500-／DOOR￥3,000-（+1D）",
		FirstLivePriceEnglish: "ADV￥2,500-／DOOR￥3,000-（+1D）",
		FirstLiveOpenTime:     time.Date(2023, 11, 2, 18, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 11, 2, 19, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://www.toos.co.jp/3/events/d-b-inches-1st-ep-instinct-filter-bubble-release-party/",
	},
)

var ShimokitazawaWaverFetcher = fetchers.Simple{
	BaseURL:              "https://waverwaver.net/",
	ShortYearIterableURL: "https://waverwaver.net/category/schedule/20%d年%02d月/",
	ExpandedLiveSelector: "//h1/a",
	LiveSelector:         "//div[contains(@class, 'postList')]/article[contains(@class, 'media')]",
	TitleQuerier:         *htmlquerier.Q("//ul[contains(@class, 'schedule_list')]//i[contains(@class, 'icon_title')]/following-sibling::text()"),
	ArtistsQuerier:       *htmlquerier.Q("//ul[contains(@class, 'schedule_list')]//i[contains(@class, 'icon_lineup')]/following-sibling::text()").Split("/"),
	PriceQuerier:         *htmlquerier.Q("//ul[contains(@class, 'schedule_list')]//i[contains(@class, 'icon_adv')]/following-sibling::text()"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//h1").Before("年"),
		MonthQuerier:     *htmlquerier.Q("//h1").Before("月").After("年"),
		DayQuerier:       *htmlquerier.Q("//h1").Before("日").After("月"),
		OpenTimeQuerier:  *htmlquerier.Q("//ul[contains(@class, 'schedule_list')]//i[contains(@class, 'icon_open')]/following-sibling::text()").Before(" / "),
		StartTimeQuerier: *htmlquerier.Q("//ul[contains(@class, 'schedule_list')]//i[contains(@class, 'icon_open')]/following-sibling::text()").After(" / ").Before(" "),
		IsMonthInLive:    true,
		IsYearInLive:     true,
	},

	PrefectureName: "tokyo",
	AreaName:       "shimokitazawa",
	VenueID:        "shimokitazawa-waver",
	Latitude:       35.660187,
	Longitude:      139.667937,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         49,
		FirstLiveTitle:        "『笑うはこには福来る』",
		FirstLiveArtists:      []string{"小山恭代", "高田えぬひろ", "矢部りんご", "ナカザワ"},
		FirstLivePrice:        "adv¥2,200 / door¥2,500 各＋1drink代別途(¥600)",
		FirstLivePriceEnglish: "adv¥2,200 / door¥2,500 各＋1drinkNot included in ticket(¥600)",
		FirstLiveOpenTime:     time.Date(2023, 11, 2, 18, 45, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 11, 2, 19, 15, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://waverwaver.net/2023/11/02/2023%%e5%%b9%%b411%%e6%%9c%%8802%%e6%%97%%a5%%e6%%9c%%a8/",
	},
}

/***************
 *						 *
 *	Setagaya  *
 *						 *
 ***************/

var NishieifukuJamFetcher = fetchers.CreateOmatsuriFetcher(
	"https://jam.rinky.info/",
	"tokyo",
	"setagaya",
	"nishieifuku-jam",
	35.678313,
	139.635063,
	fetchers.TestInfo{
		NumberOfLives:         34,
		FirstLiveTitle:        "DAY EVENT] Momo♡Ai Solo Performance -Thank you for meeting me",
		FirstLiveArtists:      []string{"もも♡あい"},
		FirstLivePrice:        "Sチケット¥8,000 / Aチケット¥5,000 / 一般¥2,000 / 学生¥1,500 ＋ DRINK¥600",
		FirstLivePriceEnglish: "S-Ticket¥8,000 / A-Ticket¥5,000 / Ordinary Ticket¥2,000 / Students¥1,500 ＋ DRINK¥600",
		FirstLiveOpenTime:     time.Date(2025, 3, 1, 9, 45, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 3, 1, 10, 45, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://jam.rinky.info/events/26435",
	},
)

var SangenjayaHeavensDoorFetcher = fetchers.Simple{
	BaseURL:        "https://heavens-door-music.com/",
	InitialURL:     "https://heavens-door-music.com/schedule/",
	LiveSelector:   "//div[@class='clearfix sche_box']",
	TitleQuerier:   *htmlquerier.Q("//h2"),
	ArtistsQuerier: *htmlquerier.QAll("//div[contains(@class, 'tribe-events-list-event-description')]/p[1] | //div[contains(@class, 'tribe-events-list-event-description')]/h1").Split(" ／ ").Split("\n").After("O.A.："),
	DetailQuerier:  *htmlquerier.Q("//div[contains(@class, 'tribe-events-list-event-description')]"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:  *htmlquerier.Q("//h2[@class='tribe-events-list-separator-month']"),
		MonthQuerier: *htmlquerier.Q("//h2[@class='tribe-events-list-separator-month']").After("/"),
		DayQuerier:   *htmlquerier.Q("//span[@class='tribe-event-date-start']").After("/"),
	},

	PrefectureName: "tokyo",
	AreaName:       "setagaya",
	VenueID:        "sangenjaya-heavensdoor",
	Latitude:       35.642187,
	Longitude:      139.671688,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         18,
		FirstLiveTitle:        "fuck’in shit! White day",
		FirstLiveArtists:      []string{"ヒゲと味噌汁", "教祖仮面", "堕天使Project", "平成墓嵐", "Brain Stupid"},
		FirstLivePrice:        "前売/当日　￥2,500/￥2,800",
		FirstLivePriceEnglish: "Reservation/Door　￥2,500/￥2,800",
		FirstLiveOpenTime:     time.Date(2025, 3, 14, 18, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 3, 14, 19, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://heavens-door-music.com/schedule/",
	},
}

var ShindaitaFeverFetcher = fetchers.Simple{
	BaseURL:        "https://www.fever-popo.com/",
	InitialURL:     util.InsertYearMonth("https://www.fever-popo.com/schedule/%d/%02d/"),
	NextSelector:   "//div[@id='mekuri']/a[2]",
	LiveSelector:   "//div[contains(@class, 'hentry')]",
	TitleQuerier:   *htmlquerier.Q("//h2[contains(@class, 'eventtitle')]").After(")\u00A0"),
	ArtistsQuerier: *htmlquerier.Q("//h3/p").SplitIgnoreWithin("(\n)|( / )", '(', ')'),
	PriceQuerier:   *htmlquerier.Q("//div[2]/div[1]/div[3]"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//h2[contains(@class, 'eventtitle')]").Before("."),
		MonthQuerier:     *htmlquerier.Q("//h2[contains(@class, 'eventtitle')]").SplitIndex(".", 1),
		DayQuerier:       *htmlquerier.Q("//h2[contains(@class, 'eventtitle')]").SplitIndex(".", 2).Before(" "),
		OpenTimeQuerier:  *htmlquerier.Q("//div[2]/div[1]/div[2]").Before(" / "),
		StartTimeQuerier: *htmlquerier.Q("//div[2]/div[1]/div[2]").After(" / "),
		IsMonthInLive:    true,
		IsYearInLive:     true,
	},

	PrefectureName: "tokyo",
	AreaName:       "setagaya",
	VenueID:        "shindaita-fever",
	Latitude:       35.662813,
	Longitude:      139.660062,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         32,
		FirstLiveTitle:        "『百々和宏presents ～Drunk51～』",
		FirstLiveArtists:      []string{"百々和宏（MO’SOME TONEBENDER）", "ホリエアツシ（STRAIGHTENER）", "佐々木亮介（a flood of circle）", "ヤマジカズヒデ（dip）", "有江嘉典（VOLA&THE ORIENTAL MACHINE）", "ウエノコウジ（the HIATUS, Radio Caroline）", "クハラカズユキ（The Birthday）", "有松益男（BACK DROP BOMB）"},
		FirstLivePrice:        "ADV ￥4800 (+1drink) THANK YOU SOLD OUT!!! ※1drink ￥600",
		FirstLivePriceEnglish: "ADV ￥4800 (+1drink) THANK YOU SOLD OUT!!! ※1drink ￥600",
		FirstLiveOpenTime:     time.Date(2023, 11, 1, 18, 45, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 11, 1, 19, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          util.InsertYearMonth("https://www.fever-popo.com/schedule/%d/%02d/"),
	},
}

/**************
 * 					  *
 *	Shinjuku  *
 *						*
 **************/

var ShinjukuAcbHallFetcher = fetchers.Simple{
	BaseURL:              "https://acb-hall.jp/",
	ShortYearIterableURL: "https://acb-hall.jp/schedule.php?year=20%d&month=%d",
	LiveSelector:         "//article[@class='cal']/table/tbody/tr[.//li[@class='band']/text()!='TBA']",
	TitleQuerier:         *htmlquerier.Q("//li/b"),
	ArtistsQuerier:       *htmlquerier.Q("//li[@class='band']").Split(" / "),
	PriceQuerier:         *htmlquerier.QAll("//li[@class='adv-door']/span").Join(" / "),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//ul[@class='pager']/li[2]"),
		MonthQuerier:     *htmlquerier.Q("//ul[@class='pager']/li[2]").After("年"),
		DayQuerier:       *htmlquerier.Q("//th/text()[1]"),
		OpenTimeQuerier:  *htmlquerier.Q("//li[@class='open-start']/span[1]").After("OPEN:"),
		StartTimeQuerier: *htmlquerier.Q("//li[@class='open-start']/span[2]").After("START:"),
	},

	PrefectureName: "tokyo",
	AreaName:       "shinjuku",
	VenueID:        "shinjuku-acbhall",
	Latitude:       35.695937,
	Longitude:      139.702437,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         10,
		FirstLiveTitle:        "Overwhelming LIBERTY TOUR 2024 Final series",
		FirstLiveArtists:      []string{"OwL", "GOOD4NOTHING"},
		FirstLivePrice:        "前売: 3,000円 / 当日: 3,500円",
		FirstLivePriceEnglish: "Reservation: 3,000円 / Door: 3,500円",
		FirstLiveOpenTime:     time.Date(2025, 2, 1, 18, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 2, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          util.InsertShortYearMonth("https://acb-hall.jp/schedule.php?year=20%d&month=%d"),
	},
}

var ShinjukuHeistFetcher = fetchers.Simple{
	BaseURL:              "https://heist.tokyo/",
	ShortYearIterableURL: "https://heist.tokyo/schedule/20%d/%02d/",
	LiveSelector:         "//div[@class='sche-archives-flex']/a",
	ExpandedLiveSelector: ".",
	TitleQuerier:         *htmlquerier.Q("//h2[@class='detail-single']"),
	ArtistsQuerier:       *htmlquerier.Q("//h4[@class='act-name']").Split("\n").Split(" / "),
	PriceQuerier:         *htmlquerier.QAll("//td[@class='table-title'][./text()='ADV:' or ./text()='DOOR:']/following-sibling::td[@class='table-detail']").Join(" / DOOR: ").Prefix("ADV: "),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//li[@class='achives-year']"),
		MonthQuerier:     *htmlquerier.Q("//li[@class='achives-year']").After("/"),
		DayQuerier:       *htmlquerier.Q("//span[@class='date-text']").After("/").After("/"),
		OpenTimeQuerier:  *htmlquerier.Q("//td[@class='table-title'][./text()='OPEN:']/following-sibling::td[@class='table-detail']"),
		StartTimeQuerier: *htmlquerier.Q("//td[@class='table-title'][./text()='START:']/following-sibling::td[@class='table-detail']"),
	},

	PrefectureName: "tokyo",
	AreaName:       "shinjuku",
	VenueID:        "shinjuku-heist",
	Latitude:       35.695437,
	Longitude:      139.703562,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         24,
		FirstLiveTitle:        "『アイプラアイドルパーティー #17』",
		FirstLiveArtists:      []string{"ギミラブ！", "シェリー"},
		FirstLivePrice:        "ADV: 2400 / DOOR: 2900",
		FirstLivePriceEnglish: "ADV: 2400 / DOOR: 2900",
		FirstLiveOpenTime:     time.Date(2025, 2, 1, 10, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 2, 1, 10, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://heist.tokyo/schedule/schedule-3042/",
	},
}

var ShinjukuLoftFetcher = fetchers.CreateLoftFetcher(
	"https://www.loft-prj.co.jp/",
	"https://www.loft-prj.co.jp/schedule/loft/schedule?scheduleyear=20%d&schedulemonth=%d",
	"tokyo",
	"shinjuku",
	"shinjuku-loft",
	fetchers.TestInfo{
		NumberOfLives:         35,
		FirstLiveTitle:        "邂逅 2025",
		FirstLiveArtists:      []string{"ジュウ", "Haze", "らそんぶる", "夕方と猫", "毎晩揺れてスカート", "ウマシカて", "パキルカ", "pinfu", "天", "アンと私", "JIGDRESS", "Cody・Lee(李)", "チョーキューメイ", "超☆社会的サンダル", "おとなりにぎんが計画"},
		FirstLivePrice:        "ADV:通常¥4400・学生¥3400 / DOOR:通常¥4900(DRINK代別¥600)",
		FirstLivePriceEnglish: "ADV:Regular¥4400・Students¥3400 / DOOR:Regular¥4900(DRINKSeparately¥600)",
		FirstLiveOpenTime:     time.Date(2025, 2, 1, 12, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 2, 1, 12, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://www.loft-prj.co.jp/schedule/loft/schedule/301392",
	},
	35.695538,
	139.702578,
)

var ShinjukuMarbleFetcher = fetchers.Simple{
	BaseURL: "https://shinjuku-marble.com/",
	LiveHTMLFetcher: func(testDocument []byte) (nodes []*html.Node, err error) {
		nodes = make([]*html.Node, 0)
		var newNodes []*html.Node
		offset := 0
		endDate := time.Now().Format("2006-01-02")
		hasMoreEvent := true
		if testDocument == nil {
			for i := 0; i < 20 && hasMoreEvent; i++ {
				client := &http.Client{Timeout: 10 * time.Second}
				reqBody := fmt.Sprintf("action=mec_grid_load_more&mec_start_date=%s&mec_offset=%d&atts%%5Bsk-options%%5D%%5Bgrid%%5D%%5Bstyle%%5D=classic&atts%%5Bid%%5D=314&apply_sf_date=0", endDate, offset)
				var req *http.Request
				req, err = http.NewRequest("POST", "https://shinjuku-marble.com/wp-admin/admin-ajax.php", strings.NewReader(reqBody))
				if err != nil {
					return
				}
				req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
				var res *http.Response
				res, err = client.Do(req)
				if err != nil {
					return
				}
				defer res.Body.Close()

				var b []byte
				b, err = io.ReadAll(res.Body)
				if err != nil {
					return
				}
				newNodes, offset, endDate, hasMoreEvent = parseShinjukuMarbleResponse(b)
				if newNodes != nil {
					nodes = append(nodes, newNodes...)
				}
			}
		} else {
			newNodes, _, _, _ = parseShinjukuMarbleResponse(testDocument)
			nodes = newNodes
		}
		return
	},
	TitleQuerier:        *htmlquerier.Q("//h1"),
	DetailQuerier:       *htmlquerier.QAll("//div[contains(@class, 'mec-single-event-description')]/p/text()").Join("\n"),
	ArtistsQuerier:      *htmlquerier.QAll("//div[contains(@class, 'mec-single-event-description')]/p/text()").DeleteUntil("●出演"),
	DetailsLinkSelector: "//link[@rel='canonical']",

	PrefectureName: "tokyo",
	AreaName:       "shinjuku",
	VenueID:        "shinjuku-marble",
	Latitude:       35.696313,
	Longitude:      139.700562,

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:   *htmlquerier.Q("//span[@class='mec-start-date-label']").Split(" ").KeepIndex(-1),
		MonthQuerier:  *htmlquerier.Q("//span[@class='mec-start-date-label']").Before("月"),
		DayQuerier:    *htmlquerier.Q("//span[@class='mec-start-date-label']").After("月").Trim().Before(" "),
		IsYearInLive:  true,
		IsMonthInLive: true,
	},

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         12,
		FirstLiveTitle:        "THURSDAY’S YOUTH 8th Anniversary Live”Just after dark”",
		FirstLiveArtists:      []string{"THURSDAY’S YOUTH"},
		FirstLivePrice:        "前売り¥3500、当日¥3900",
		FirstLivePriceEnglish: "Reservation¥3500、Door¥3900",
		FirstLiveOpenTime:     time.Date(2025, 3, 9, 17, 00, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 3, 9, 17, 43, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://shinjuku-marble.com/calender/thursdays-youth-8th-anniversary-livejust-after-dark/",
	},
}

type ShinjukuMarbleResponse struct {
	EndDate      string `json:"end_date"`
	HasMoreEvent int    `json:"has_more_event"`
	Html         string `json:"html"`
	Offset       int    `json:"offset"`
}

func parseShinjukuMarbleResponse(s []byte) (nodes []*html.Node, offset int, endDate string, hasMoreEvent bool) {
	nodes = make([]*html.Node, 0)
	var res ShinjukuMarbleResponse
	if err := json.Unmarshal(s, &res); err != nil {
		return
	}

	offset = res.Offset
	hasMoreEvent = res.HasMoreEvent == 1
	endDate = res.EndDate

	n, err := htmlquery.Parse(strings.NewReader(res.Html))
	if err != nil {
		return
	}
	links, err := htmlquery.QueryAll(n, "//article")
	if err != nil {
		return
	}
	baseUrl, err := url.Parse("https://shinjuku-marble.com/")
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	queue := make(chan *fetchers.LiveQueueElement, len(links))
	var liveSlice []*fetchers.LiveQueueElement
	for _, live := range links {
		job := &fetchers.LiveQueueElement{Live: live}
		liveSlice = append(liveSlice, job)
		queue <- job
	}
	close(queue)
	for i := 0; i < min(10, len(links)); i++ {
		wg.Add(1)
		go fetchers.FetchLiveConcurrent(baseUrl, queue, "//a", &wg)
	}
	wg.Wait()
	for _, liveDetails := range liveSlice {
		nodes = append(nodes, liveDetails.Res)
	}
	return
}

var ShinjukuMarzFetcher = fetchers.Simple{
	BaseURL:              "http://www.marz.jp/",
	ShortYearIterableURL: "http://www.marz.jp/schedule/20%d/%02d/",
	LiveSelector:         "//article",
	TitleQuerier:         *htmlquerier.Q("//h1"),
	ArtistsQuerier:       *htmlquerier.QAll("//div[@class='entrybody']/a/p/text()"),
	PriceQuerier:         *htmlquerier.Q("//div[@class='entryex']/p/text()[2]"),
	DetailsLinkSelector:  "//h1/a",

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//div[@id='h']/h1").Before("."),
		MonthQuerier:     *htmlquerier.Q("//div[@id='h']/h1").After("."),
		DayQuerier:       *htmlquerier.Q("//time").After("."),
		OpenTimeQuerier:  *htmlquerier.Q("//div[@class='entryex']/p/text()[1]").Before(" / "),
		StartTimeQuerier: *htmlquerier.Q("//div[@class='entryex']/p/text()[1]").After(" / "),
	},

	PrefectureName: "tokyo",
	AreaName:       "shinjuku",
	VenueID:        "shinjuku-marz",
	Latitude:       35.696313,
	Longitude:      139.700688,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         4,
		FirstLiveTitle:        "少女脱兎１周年ライブ 「惑星探査機、乗車」",
		FirstLiveArtists:      []string{"少女脱兎", "MAPA", "赤いくらげ"},
		FirstLivePrice:        "adv ¥3,000 / door ¥3,500 (+1drink¥600)",
		FirstLivePriceEnglish: "adv ¥3,000 / door ¥3,500 (+1drink¥600)",
		FirstLiveOpenTime:     time.Date(2025, 4, 21, 18, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 4, 21, 18, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "http://www.marz.jp/schedule/2025/04/post_613.html",
	},
}

var ShinjukuSamuraiFetcher = fetchers.Simple{
	BaseURL:              "https://live-samurai.jp/",
	ShortYearIterableURL: "https://live-samurai.jp/20%d/%02d/",
	LiveSelector:         "//table[@class='post']",
	TitleQuerier:         *htmlquerier.Q("//h4[@class='entryTitle']"),
	ArtistsQuerier:       *htmlquerier.Q("//table[@class='post2']/tbody/tr[3]/td[2]").SplitRegex(`/`).Trim(),
	PriceQuerier:         *htmlquerier.QAll("//table[@class='post2']/tbody/tr[2]/td").Join(" "),
	DetailsLinkSelector:  "//h4/a",

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//h1"),
		MonthQuerier:     *htmlquerier.Q("//h1").After("年"),
		DayQuerier:       *htmlquerier.Q("//font[@color='#ffffff']/b/text()[1]").After("/"),
		OpenTimeQuerier:  *htmlquerier.Q("//table[@class='post2']/tbody/tr[1]/td[2]"),
		StartTimeQuerier: *htmlquerier.Q("//table[@class='post2']/tbody/tr[1]/td[2]").After("/"),
	},

	PrefectureName: "tokyo",
	AreaName:       "shinjuku",
	VenueID:        "shinjuku-samurai",
	Latitude:       35.697938,
	Longitude:      139.701813,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         24,
		FirstLiveTitle:        "Shinjuku SAMURAI pre. うてばひびく vol.118",
		FirstLiveArtists:      []string{"歌田真紀", "Spit Lulu's", "花酔い"},
		FirstLivePrice:        "ADV./DOOR ￥2,500/￥3,000　※入場時1ドリンク￥600別",
		FirstLivePriceEnglish: "ADV./DOOR ￥2,500/￥3,000　※When entering1Drink￥600Separately",
		FirstLiveOpenTime:     time.Date(2025, 2, 3, 18, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 2, 3, 19, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://live-samurai.jp/2025/02/03/__trashed-2-34/",
	},
}

var ShinjukuScienceFetcher = fetchers.Simple{
	BaseURL:              "https://club-science.com/",
	ShortYearIterableURL: "https://club-science.com/schedule/20%d/%02d/",
	LiveSelector:         "//div[contains(@class, 'schedule_box')][not(contains(.//div[@class='sche_ttl'], '株式会社'))]",
	MultiLiveDaySelector: "//div[@class='schedule_center_list']",
	TitleQuerier:         *htmlquerier.Q("//div[@class='sche_ttl']"),
	ArtistsQuerier:       *htmlquerier.QAll("//div[@class='sche_detail'][1]//text()").Split("\n").Split(" / ").DeleteFrom("＜FOOD＞").DeleteFrom("【FOOD】").DeleteFrom("【FOOD BOOTH】"),
	PriceQuerier:         *htmlquerier.Q("//div[@class='sche_detail'][2]/text()[2]"),
	DetailsLinkSelector:  "//a",

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//div[@class='next_preview']/span[@class='sche-date']"),
		MonthQuerier:     *htmlquerier.Q("//div[@class='next_preview']/span[@class='sche-date']").After("年"),
		DayQuerier:       *htmlquerier.Q("//div[@class='sche-date']"),
		OpenTimeQuerier:  *htmlquerier.Q("//div[@class='sche_detail'][2]/text()[1]"),
		StartTimeQuerier: *htmlquerier.Q("//div[@class='sche_detail'][2]/text()[1]").After("START"),
	},

	PrefectureName: "tokyo",
	AreaName:       "shinjuku",
	VenueID:        "shinjuku-science",
	Latitude:       35.695437,
	Longitude:      139.703562,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         24,
		FirstLiveTitle:        `OXYMORPHONN × バミューダ★バガボンド × Shinjuku club SCIENCE presents "SCIENTIFIC NIGHTMARE”vol.8`,
		FirstLiveArtists:      []string{"OXYMORPHONN", "バミューダ★バガボンド", "THE DISASTER POINTS", "SOBUT", "STRIKE AGAIN", "kitsunevi", "HOTVOX", "You-suke Hirata(Kikoku/CLOCK CHANNEL/Anti Class/R.I.P CLEAR)", "$HUN (HOMEWARD TATTOO PARLOR)", "AREA"},
		FirstLivePrice:        "ADV 2500 DOOR 3000",
		FirstLivePriceEnglish: "ADV 2500 DOOR 3000",
		FirstLiveOpenTime:     time.Date(2025, 2, 1, 16, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 2, 1, 17, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://club-science.com/detail/2025/02/01",
	},
}

var ShinjukuZircoTokyoFetcher = fetchers.CreateBassOnTopFetcher(
	"https://zirco-tokyo.jp/",
	"https://zirco-tokyo.jp/schedule/calendar/20%d/%02d/",
	"tokyo",
	"shinjuku",
	"shinjuku-zircotokyo",
	fetchers.TestInfo{
		NumberOfLives:         32,
		FirstLiveTitle:        "専門学校東京ビジュアルアーツ　音楽総合学科ミュージシャン専攻卒業公演",
		FirstLiveArtists:      []string{},
		FirstLivePrice:        "ADV/DOOR ￥0- (1Drink代金￥600別途必要)",
		FirstLivePriceEnglish: "ADV/DOOR ￥0- (1DrinkPrice￥600Must be purchased separately)",
		FirstLiveOpenTime:     time.Date(2024, 3, 1, 17, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 1, 18, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://zirco-tokyo.jp/schedule/detail/31777",
	},
	35.693763,
	139.703859,
)

var ZeppShinjukuFetcher = fetchers.CreateZeppFetcher("shinjuku", "tokyo", "shinjuku", "shinjuku-zepp", 35.695937, 139.700688, fetchers.TestInfo{
	NumberOfLives:         17,
	FirstLiveTitle:        "⼩久保柚乃⽣誕ソロライブ「炭素。」",
	FirstLiveArtists:      []string{"小久保柚乃"},
	FirstLivePrice:        "スタンディング/ ¥6,900 カメコ席/ ¥12,000",
	FirstLivePriceEnglish: "Standing/ ¥6,900 カメコ席/ ¥12,000",
	FirstLiveOpenTime:     time.Date(2025, 4, 3, 18, 0, 0, 0, util.JapanTime),
	FirstLiveStartTime:    time.Date(2025, 4, 3, 19, 0, 0, 0, util.JapanTime),
	FirstLiveURL:          "https://www.zepp.co.jp/hall/shinjuku/schedule/single/?rid=147339",
})

/***********
 *         *
 *  Other  *
 *         *
 ***********/

var TakadanobabaClubPhaseFetcher = fetchers.Simple{
	BaseURL:              "https://www.club-phase.com/",
	ShortYearIterableURL: "https://www.club-phase.com/schedule/20%d-%02d",
	LiveSelector:         "//table[@id='sched']/tbody/tr",
	TitleQuerier:         *htmlquerier.Q("/td/p[1]"),
	ArtistsQuerier:       *htmlquerier.QAll("//p[@class='sc_artist']/text()").Split(" / "),
	PriceQuerier:         *htmlquerier.Q("//p[@class='sc_price']/text()[2]"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//section[@id='contents-right']//h2"),
		MonthQuerier:     *htmlquerier.Q("//section[@id='contents-right']//h2").After(" "),
		DayQuerier:       *htmlquerier.Q("//span[@class='sc_date']"),
		OpenTimeQuerier:  *htmlquerier.Q("//p[@class='sc_price']/text()[1]"),
		StartTimeQuerier: *htmlquerier.Q("//p[@class='sc_price']/text()[1]").After("START"),
	},

	PrefectureName: "tokyo",
	AreaName:       "tokyo",
	VenueID:        "takadanobaba-clubphase",
	Latitude:       35.714687,
	Longitude:      139.706313,

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         30,
		FirstLiveTitle:        `高田馬場CLUB PHASE THE NINTH APOLLO pre PHASEの24周年を祝う "TETORAと炙りなタウンの2マン"`,
		FirstLiveArtists:      []string{"TETORA", "炙りなタウン"},
		FirstLivePrice:        "ADV ￥3500/DOOR ￥4500/D別■",
		FirstLivePriceEnglish: "ADV ￥3500/DOOR ￥4500/Drinks sold separately■",
		FirstLiveOpenTime:     time.Date(2025, 3, 1, 17, 15, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2025, 3, 1, 18, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          util.InsertYearMonth("https://www.club-phase.com/schedule/%d-%02d"),
	},
}

var ZeppDiverCityFetcher = fetchers.CreateZeppFetcher("divercity", "tokyo", "tokyo", "tokyo-zeppdivercity", 35.624688, 139.774563, fetchers.TestInfo{
	NumberOfLives:         24,
	FirstLiveTitle:        "DEEN LIVE JOY-Break26 〜ROCK ON!〜 追加公演",
	FirstLiveArtists:      []string{"DEEN"},
	FirstLivePrice:        "全席指定/ ¥8,500",
	FirstLivePriceEnglish: "全席指定/ ¥8,500",
	FirstLiveOpenTime:     time.Date(2025, 3, 1, 16, 30, 0, 0, util.JapanTime),
	FirstLiveStartTime:    time.Date(2025, 3, 1, 17, 30, 0, 0, util.JapanTime),
	FirstLiveURL:          "https://www.zepp.co.jp/hall/divercity/schedule/single/?rid=146730",
})

var ZeppHanedaFetcher = fetchers.CreateZeppFetcher("haneda", "tokyo", "tokyo", "tokyo-zepphaneda", 35.547438, 139.756562, fetchers.TestInfo{
	NumberOfLives:         29,
	FirstLiveTitle:        "aiko Live Tour「Love Like Rock vol.10」",
	FirstLiveArtists:      []string{"aiko"},
	FirstLivePrice:        "1Fスタンディング/ ¥7,500 2F指定席/ ¥7,500 2Fスタンディング/ ¥7,500",
	FirstLivePriceEnglish: "1Fスタンディング/ ¥7,500 2F指定席/ ¥7,500 2Fスタンディング/ ¥7,500",
	FirstLiveOpenTime:     time.Date(2025, 3, 1, 16, 30, 0, 0, util.JapanTime),
	FirstLiveStartTime:    time.Date(2025, 3, 1, 17, 30, 0, 0, util.JapanTime),
	FirstLiveURL:          "https://www.zepp.co.jp/hall/haneda/schedule/single/?rid=140708",
})
