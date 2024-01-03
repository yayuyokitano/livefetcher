package connectors

import (
	"time"

	"github.com/yayuyokitano/livefetcher/lib/core/fetchers"
	"github.com/yayuyokitano/livefetcher/lib/core/htmlquerier"
	"github.com/yayuyokitano/livefetcher/lib/core/util"
)

/*************
 *           *
 *  Shibuya  *
 *           *
 *************/

var ShibuyaCycloneFetcher = fetchers.Simple{
	BaseURL:              "http://www.cyclone1997.com/",
	ShortYearIterableURL: "http://www.cyclone1997.com/schedule/20%dschedule_%d.html",
	LiveSelector:         "//body/table",
	TitleQuerier:         *htmlquerier.Q("//td/p/span[1]").ReplaceAllRegex(`\s+`, " "),
	ArtistsQuerier:       *htmlquerier.Q("//span/strong").SplitIgnoreWithin("[\n/]", '(', ')'),
	PriceQuerier:         *htmlquerier.Q("//dl[@class='schedule-content__ticket']//p/text()[2]"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//span[@class='date']").Before("/"),
		MonthQuerier:     *htmlquerier.Q("//span[@class='date']").SplitIndex("/", 1),
		DayQuerier:       *htmlquerier.Q("//span[@class='date']").SplitIndex("/", 2).Before("("),
		OpenTimeQuerier:  *htmlquerier.Q("//dl[@class='schedule-content__openstart']//p").Before("／"),
		StartTimeQuerier: *htmlquerier.Q("//dl[@class='schedule-content__openstart']//p").After("／"),
	},

	PrefectureName: "tokyo",
	AreaName:       "shimokitazawa",
	VenueID:        "shimokitazawa-440",

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         29,
		FirstLiveTitle:        "“Stormy Wednesday”",
		FirstLiveArtists:      []string{"王将&The Guv’nor Brothers", "Weiyao Band", "前後のカルマ"},
		FirstLivePrice:        "ADV.￥3,000／DOOR.￥3,400 [1D別]",
		FirstLivePriceEnglish: "ADV.￥3,000／DOOR.￥3,400 [1 Drink purchase required]",
		FirstLiveOpenTime:     time.Unix(1698831000, 0),
		FirstLiveStartTime:    time.Unix(1698832800, 0),
		FirstLiveURL:          "https://clubque.net/schedule/2202/",
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
		FirstLiveOpenTime:     time.Unix(1698829200, 0),
		FirstLiveStartTime:    time.Unix(1698831000, 0),
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
		FirstLiveOpenTime:     time.Unix(1699538400, 0),
		FirstLiveStartTime:    time.Unix(1699543800, 0),
		FirstLiveURL:          "http://eggman.jp/schedule/reraise-house-1on1-battle-season4-vol-10/",
	},
)

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
		FirstLivePriceEnglish: "優先 ¥3,000 Ordinary Ticket ¥0 Door ¥0 （ごReservation時Drink代Separately600円）",
		FirstLiveOpenTime:     time.Unix(1698832800, 0),
		FirstLiveStartTime:    time.Unix(1698834600, 0),
		FirstLiveURL:          "https://shibuya-o.com/crest/schedule/iiiiiiidiom_23-11-1/",
	},
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
		FirstLivePriceEnglish: "ADV 9,500（ごEntry時Drink代Separately600円）",
		FirstLiveOpenTime:     time.Unix(1698829200, 0),
		FirstLiveStartTime:    time.Unix(1698831900, 0),
		FirstLiveURL:          "https://shibuya-o.com/east/schedule/the-dead-daisies/",
	},
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
		FirstLivePriceEnglish: "ADV ¥2,500 DOOR ¥3,000 （ごEntry時Drink代Separately600円）",
		FirstLiveOpenTime:     time.Unix(1698831000, 0),
		FirstLiveStartTime:    time.Unix(1698832800, 0),
		FirstLiveURL:          "https://shibuya-o.com/nest/schedule/may-in-film%e4%b8%bb%e5%82%acmayday-%e6%98%9f%e3%81%ae%e9%99%8d%e3%82%8b%e3%83%8d%e3%82%b9%e3%83%88%e3%81%a7-day1/",
	},
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
		FirstLivePriceEnglish: "ADV ¥4,500 (ごEntry時Drink代Separately600円)",
		FirstLiveOpenTime:     time.Unix(1698831000, 0),
		FirstLiveStartTime:    time.Unix(1698832800, 0),
		FirstLiveURL:          "https://shibuya-o.com/west/schedule/%e6%84%9f%e8%a6%9a%e3%83%94%e3%82%a8%e3%83%ad-10th-anniversary%e3%80%8c%e6%84%9f%e8%a6%9a%e3%83%94%e3%82%a8%e3%83%ad%e3%81%a7%e3%81%99%e3%81%8c%e3%81%aa%e3%81%ab%e3%81%8b%e3%80%8d%e3%83%84%e3%82%a2/",
	},
)

var ShibuyaWWWFetcher = fetchers.CreateWWWFetcher(
	"@data-place='www' or @data-place='wwwxwww'",
	"shibuya-www",
	fetchers.TestInfo{
		NumberOfLives:         30,
		FirstLiveTitle:        "#ピコリフ 3rdワンマンライブ「Shiny Sparkle」",
		FirstLiveArtists:      []string{"#ピコリフ"},
		FirstLivePrice:        "VIP：¥5,000 / 撮影：¥3,000 / 後方：¥1,500 / 学生・女性：¥500 / 新規：¥0 (税込 / 各ドリンク代別)",
		FirstLivePriceEnglish: "VIP：¥5,000 / 撮影：¥3,000 / 後方：¥1,500 / Students・Women：¥500 / 新規：¥0 (Incl. Tax / 各Drink代Separately)",
		FirstLiveOpenTime:     time.Unix(1698832800, 0),
		FirstLiveStartTime:    time.Unix(1698834600, 0),
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
		FirstLivePriceEnglish: "Door ¥3,500 (Incl. Tax｜Standing｜Drink代Separately)Reservation ¥3,000 (Incl. Tax｜Standing｜Drink代Separately)",
		FirstLiveOpenTime:     time.Unix(1699345800, 0),
		FirstLiveStartTime:    time.Unix(1699345800, 0),
		FirstLiveURL:          "https://www-shibuya.jp/schedule/017137.php",
	},
)

var ShibuyaWWWXFetcher = fetchers.CreateWWWFetcher(
	"@data-place='wwwx'",
	"shibuya-wwwx",
	fetchers.TestInfo{
		NumberOfLives:         27,
		FirstLiveTitle:        "BIG ROMANTIC RECORDS presents Carsick Cars live in Tokyo",
		FirstLiveArtists:      []string{"Carsick Cars", "Hello Shitty (a.k.a Sophia from UptownRecords)"},
		FirstLivePrice:        "¥5,000 / ¥5,500 (税込 / ドリンク代別)※当日券は19:00~、前売りの入場が落ち着き次第 WWW Xにて¥5,500+D代で販売いたします。",
		FirstLivePriceEnglish: "¥5,000 / ¥5,500 (Incl. Tax / Drink代Separately)※Door券は19:00~、ReservationりのEntryが落ち着き次第 WWW Xにて¥5,500+D代で販売いたします。",
		FirstLiveOpenTime:     time.Unix(1698832800, 0),
		FirstLiveStartTime:    time.Unix(1698832800, 0),
		FirstLiveURL:          "https://www-shibuya.jp/schedule/017258.php",
	},
)

/*******************
 * 								 *
 *	Shimokitazawa	 *
 *								 *
 *******************/

var ShimokitazawaArtistFetcher = fetchers.Simple{
	BaseURL:        "http://www.c-artist.com/",
	InitialURL:     util.InsertYearMonth("http://www.c-artist.com/schedule/list/%d%d.txt"),
	LiveSelector:   "//div[@class='sche']",
	TitleQuerier:   *htmlquerier.Q("//p[contains(@class, 'guestname')]/text()[1]"),
	ArtistsQuerier: *htmlquerier.Q("//p[contains(@class, 'guestname')]/text()[last()]").Split("\u00A0").Before("』〜one-man_live〜").Before("』〜two-man_live〜").TrimPrefix("『"),
	PriceQuerier:   *htmlquerier.Q("//p[@class='ex']").After(" / "),
	DetailsLink:    "http://www.c-artist.com/schedule/",

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

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         17,
		FirstLiveTitle:        "『ex_jikkatsu_a:side』",
		FirstLiveArtists:      []string{"パぁ", "かわぐちシンゴ", "モリタケル"},
		FirstLivePrice:        "2,000yen（1ドリンク込み)",
		FirstLivePriceEnglish: "2,000yen（1DrinkIncluded)",
		FirstLiveOpenTime:     time.Unix(1698832800, 0),
		FirstLiveStartTime:    time.Unix(1698832800, 0),
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
		NumberOfLives:         29,
		FirstLiveTitle:        "Get the Crispy! vol.2",
		FirstLiveArtists:      []string{"Crispy Camera Club", "Kamisado", "Salan", "Rumi Nagasawa(LIGHTERS)"},
		FirstLivePrice:        "ADV￥2,400-／DOOR￥2,900-（+1D）",
		FirstLivePriceEnglish: "ADV￥2,400-／DOOR￥2,900-（+1D）",
		FirstLiveOpenTime:     time.Unix(1698831000, 0),
		FirstLiveStartTime:    time.Unix(1698832800, 0),
		FirstLiveURL:          "https://toos.co.jp/basementbar/ev/get-the-crispy-vol-2/",
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

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         26,
		FirstLiveTitle:        "地下教室 vol.13",
		FirstLiveArtists:      []string{"uyuu", "アウフヘーベン", "シンクマクラ", "憂牡丹"},
		FirstLivePrice:        "¥2,000-/¥2,500-(＋1D)",
		FirstLivePriceEnglish: "¥2,000-/¥2,500-(＋1D)",
		FirstLiveOpenTime:     time.Unix(1698829200, 0),
		FirstLiveStartTime:    time.Unix(1698831000, 0),
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
		FirstLiveOpenTime:     time.Unix(1698918300, 0),
		FirstLiveStartTime:    time.Unix(1698920100, 0),
		FirstLiveURL:          "https://chikamichi-otemae.com/chikamichi/852/",
	},
)

var ShimokitazawaClub251Fetcher = fetchers.Simple{
	BaseURL:              "http://www.club251.com/",
	ShortYearIterableURL: "http://www.club251.com/schedule/schedule-%d%02d.html",
	LiveSelector:         "//div[contains(@class, 'Schedule-1day')]",
	TitleQuerier:         *htmlquerier.Q("//div[contains(@class, 'EVENT-TITLE')]"),
	ArtistsQuerier:       *htmlquerier.Q("//p[contains(@class, 'Performer')]").Split("／"),
	PriceQuerier:         *htmlquerier.Q("//p[contains(@class, 'DETAIL')][1]/text()[2]"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//p[contains(@class, 'month')]").After("/"),
		MonthQuerier:     *htmlquerier.Q("//div[contains(@class, 'DATE')]/p").Before("/"),
		DayQuerier:       *htmlquerier.Q("//div[contains(@class, 'DATE')]/p").After("/").Before("."),
		OpenTimeQuerier:  *htmlquerier.Q("//p[contains(@class, 'DETAIL')][1]/text()[1]").Before("/"),
		StartTimeQuerier: *htmlquerier.Q("//p[contains(@class, 'DETAIL')][1]/text()[1]").After("/"),

		IsMonthInLive: true,
	},

	PrefectureName: "tokyo",
	AreaName:       "shimokitazawa",
	VenueID:        "shimokitazawa-club251",

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         28,
		FirstLiveTitle:        "\"9DayzGlitchClubTokyo pre.“太陽と月光のクリシェ” 下北沢CLUB251 30th ANNIVERSARY！！\"",
		FirstLiveArtists:      []string{"9DayzGlitchClubTokyo", "AZ-ON", "RONLON", "GTRA"},
		FirstLivePrice:        "adv¥2,900/door¥3,400(別途1D600円)",
		FirstLivePriceEnglish: "adv¥2,900/door¥3,400(Separately1D600円)",
		FirstLiveOpenTime:     time.Unix(1696150800, 0),
		FirstLiveStartTime:    time.Unix(1696152600, 0),
		FirstLiveURL:          "http://www.club251.com/schedule/schedule-%d%02d.html",
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

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         29,
		FirstLiveTitle:        "“Stormy Wednesday”",
		FirstLiveArtists:      []string{"王将&The Guv’nor Brothers", "Weiyao Band", "前後のカルマ"},
		FirstLivePrice:        "ADV.￥3,000／DOOR.￥3,400 [1D別]",
		FirstLivePriceEnglish: "ADV.￥3,000／DOOR.￥3,400 [1 Drink purchase required]",
		FirstLiveOpenTime:     time.Unix(1698831000, 0),
		FirstLiveStartTime:    time.Unix(1698832800, 0),
		FirstLiveURL:          "https://clubque.net/schedule/2202/",
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

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         29,
		FirstLiveTitle:        "“Stormy Wednesday”",
		FirstLiveArtists:      []string{"王将&The Guv’nor Brothers", "Weiyao Band", "前後のカルマ"},
		FirstLivePrice:        "ADV.￥3,000／DOOR.￥3,400 [1D別]",
		FirstLivePriceEnglish: "ADV.￥3,000／DOOR.￥3,400 [1 Drink purchase required]",
		FirstLiveOpenTime:     time.Unix(1698831000, 0),
		FirstLiveStartTime:    time.Unix(1698832800, 0),
		FirstLiveURL:          "https://clubque.net/schedule/2202/",
	},
}

var ShimokitazawaDaisyBarFetcher = fetchers.CreateDaisyBarFetcher(
	"https://daisybar.jp/",
	"https://daisybar.jp/events/event/",
	"tokyo",
	"shimokitazawa",
	"shimokitazawa-daisybar",
	"color-blue2",
	fetchers.TestInfo{
		NumberOfLives:         35,
		FirstLiveTitle:        "From The Beginning",
		FirstLiveArtists:      []string{"Heap", "DreadNought", "髙橋 翼(Paper Bag)", "michi(アバウトチルドレン)", "ワタナベイベーマコト(etymon)"},
		FirstLivePrice:        "前売り 2500円(D別)／当日 3000円(D別)",
		FirstLivePriceEnglish: "Reservationり 2500円(Drinks sold separately)／Door 3000円(Drinks sold separately)",
		FirstLiveOpenTime:     time.Unix(1698917400, 0),
		FirstLiveStartTime:    time.Unix(1698919200, 0),
		FirstLiveURL:          "https://daisybar.jp/events/event/",
	},
)

var ShimokitazawaDyCubeFetcher = fetchers.Simple{
	BaseURL:              "https://dycube.tokyo/",
	ShortYearIterableURL: "https://dycube.tokyo/schedule/?ext_num-year=20%d&ext_num-month=%02d",
	LiveSelector:         "//article[@class='schedule-article']",
	TitleQuerier:         *htmlquerier.Q("//h3").ReplaceAllRegex(`\s+`, " "),
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

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         34,
		FirstLiveTitle:        "vessy 初企画 「diamonddust」",
		FirstLiveArtists:      []string{"nene", "栢本ての", "藍谷凪", "vessy"},
		FirstLivePrice:        "TICKET ¥2,500(+1drink¥600)",
		FirstLivePriceEnglish: "TICKET ¥2,500(+1drink¥600)",
		FirstLiveOpenTime:     time.Unix(1698832800, 0),
		FirstLiveStartTime:    time.Unix(1698834600, 0),
		FirstLiveURL:          "https://dycube.tokyo/schedule/?ext_num-year=20%d&ext_num-month=%02d",
	},
}

var ShimokitazawaEraFetcher = fetchers.Simple{
	BaseURL:        "http://s-era.jp/",
	InitialURL:     util.InsertYearMonth("http://s-era.jp/schedule_cat/%d-%02d/"),
	NextSelector:   "//section[contains(@class, 'schedule-navigation')]/div[2]/p[2]/a",
	LiveSelector:   "//article[contains(@class, 'schedule-box')]",
	TitleQuerier:   *htmlquerier.Q("//h4").ReplaceAllRegex(`\s+`, " "),
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

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         38,
		FirstLiveTitle:        "ERA presents. A Day in The Life Vol.159",
		FirstLiveArtists:      []string{"paddy isle", "Akarumeno brown", "Blueberry Mondays", "Qurukuma", "Spinning Plums"},
		FirstLivePrice:        "ADV ¥2000DOOR ¥2500 (+1D¥600)",
		FirstLivePriceEnglish: "ADV ¥2000DOOR ¥2500 (+1D¥600)",
		FirstLiveOpenTime:     time.Unix(1698915600, 0),
		FirstLiveStartTime:    time.Unix(1698917400, 0),
		FirstLiveURL:          util.InsertYearMonth("http://s-era.jp/schedule_cat/%d-%02d/"),
	},
}

var ShimokitazawaFlowersLoftFetcher = fetchers.CreateLoftFetcher(
	"https://www.loft-prj.co.jp/",
	util.InsertYearMonth("https://www.loft-prj.co.jp/schedule/flowersloft/date/%d/%d"),
	"tokyo",
	"shimokitazawa",
	"shimokitazawa-flowersloft",
	fetchers.TestInfo{
		NumberOfLives:         30,
		FirstLiveTitle:        "バニーの日！ハロウィンSP",
		FirstLiveArtists:      []string{"tumiki", "はっちゃん", "テルミー", "ShinK", "クロマティーゆうや(NeoN)", "ディスク百合おん", "まみくん3歳", "マル", "有田清幸", "BIDA"},
		FirstLivePrice:        "男性：2,000 (2drink込み）\n女性：入場無料\n・仮装女子は1Dプレゼント\n・バニー...",
		FirstLivePriceEnglish: "Men：2,000 (2drinkIncluded）\nWomen：EntryFree\n・仮装女子は1Dプレゼント\n・バニー...",
		FirstLiveOpenTime:     time.Unix(1698831000, 0),
		FirstLiveStartTime:    time.Unix(1698831000, 0),
		FirstLiveURL:          "https://www.loft-prj.co.jp/schedule/flowersloft/267929",
	},
)

var ShimokitazawaLagunaFetcher = fetchers.CreateDaisyBarFetcher(
	"https://s-laguna.jp/",
	"https://s-laguna.jp/events/event/",
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
		FirstLiveOpenTime:     time.Unix(1698831000, 0),
		FirstLiveStartTime:    time.Unix(1698832800, 0),
		FirstLiveURL:          "https://s-laguna.jp/events/event/",
	},
)

var ShimokitazawaLiveHausFetcher = fetchers.Simple{
	BaseURL:             "https://livehaus.jp/",
	InitialURL:          "https://livehaus.jp/schedule/",
	NextSelector:        "//a[contains(@class, 'tribe-events-c-nav__next')]",
	LiveSelector:        "//article[contains(@class, 'tribe-events-calendar-list__event')]",
	TitleQuerier:        *htmlquerier.Q("//h3"),
	DetailQuerier:       *htmlquerier.Q("//div[contains(@class, 'tribe-events-calendar-list__event-description')]"),
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

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         19,
		FirstLiveTitle:        "突撃! 隣のトシクラブ vol.4 in 下北沢",
		FirstLiveArtists:      []string{"ドトキンズ", "Los Rancheros", "The Silver Sonics", "SHOWA BOYZ feat.DJ NIGARA", "SHJ (LONDON NITE)", "TAGO!", "KAZU SUDO (Caribbean Dandy)", "内藤啓介 (Chingcame)", "NIGARA (GARA)", "ITA (Nat Records)", "TSUNE (Lewis Leathers Japan)", "Mr.X(A&Y)", "ciibo", "KOSEI", "Yuuki(Tip Clothing&co.)", "RYO", "Shunsuke", "Shima Volume(104club)", "TOSHI & VOW"},
		FirstLivePrice:        "ADM¥2000",
		FirstLivePriceEnglish: "ADM¥2000",
		FirstLiveOpenTime:     time.Unix(1698991200, 0),
		FirstLiveStartTime:    time.Unix(1698991200, 0),
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

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         52,
		FirstLiveTitle:        "オレンジバンドシャトル",
		FirstLiveArtists:      []string{"ザストーンマスター", "おじさん的思考ディスク", "BamblueSiA", "Sonne", "Route4th", "GO SEE REGRET", "※EiNyは出演キャンセル"},
		FirstLivePrice:        "一般¥2,000 学生¥1,400 (D別)",
		FirstLivePriceEnglish: "Ordinary Ticket¥2,000 Students¥1,400 (Drinks sold separately)",
		FirstLiveOpenTime:     time.Unix(1699259400, 0),
		FirstLiveStartTime:    time.Unix(1699261200, 0),
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

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         30,
		FirstLiveTitle:        "LOLIPOP",
		FirstLiveArtists:      []string{"じゆりぴ", "若杉果歩", "イロハマイ"},
		FirstLivePrice:        "前売：¥2,900+1Drink",
		FirstLivePriceEnglish: "Reservation：¥2,900+1Drink",
		FirstLiveOpenTime:     time.Unix(1698832800, 0),
		FirstLiveStartTime:    time.Unix(1698834600, 0),
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
		OpenTimeQuerier:  *htmlquerier.Q("//td[contains(@class, 'live_menu')]/text()[4]").After("OPEN ").Before(" "),
		StartTimeQuerier: *htmlquerier.Q("//td[contains(@class, 'live_menu')]/text()[4]").After("START ").Before(" "),

		IsMonthInLive: true,
	},

	PrefectureName: "tokyo",
	AreaName:       "shimokitazawa",
	VenueID:        "shimokitazawa-mosaic",

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         31,
		FirstLiveTitle:        "『MOSAiC iDOL INFINITY - One hundred -』",
		FirstLiveArtists:      []string{"エレクトリックリボン", "月刊PAM", "セカイシティ", "LUNCH KIDS", "AKUMATICA"},
		FirstLivePrice:        "前売 ¥100 / 当日 ¥1,100(+1Drink ¥600)",
		FirstLivePriceEnglish: "Reservation ¥100 / Door ¥1,100(+1Drink ¥600)",
		FirstLiveOpenTime:     time.Unix(1685610000, 0),
		FirstLiveStartTime:    time.Unix(1685610900, 0),
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
		FirstLiveOpenTime:     time.Unix(1696152600, 0),
		FirstLiveStartTime:    time.Unix(1696154400, 0),
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

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         27,
		FirstLiveTitle:        "HAWK ANARCHY 2nd Album 宇宙侵略紀行release tour「日本侵略紀行」",
		FirstLiveArtists:      []string{"HAWK ANARCHY", "New Age Core", "OFELIA", "ALL I WANT"},
		FirstLivePrice:        "前売り￥2,500-（ドリンク別）/当日￥3,000-（ドリンク別）",
		FirstLivePriceEnglish: "Reservationり￥2,500-（Drinks sold separately）/Door￥3,000-（Drinks sold separately）",
		FirstLiveOpenTime:     time.Unix(1698915600, 0),
		FirstLiveStartTime:    time.Unix(1698918300, 0),
		FirstLiveURL:          "https://www.reg-r2.com/?page_id=7250",
	},
}

var ShimokitazawaShangrilaFetcher = fetchers.Simple{
	BaseURL:        "https://www.shan-gri-la.jp/",
	InitialURL:     "https://www.shan-gri-la.jp/tokyo/category/schedule/",
	LiveSelector:   "//div[@id='content']/div[contains(@class, 'hentry')]",
	TitleQuerier:   *htmlquerier.Q("//strong"),
	ArtistsQuerier: *htmlquerier.Q("//div[@class='post-content-content']/p[2]").Split("\n"),
	PriceQuerier:   *htmlquerier.Q("//div[@class='post-content-content']/p[3]/span"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//table[@id='wp-calendar']/caption").Before("年"),
		MonthQuerier:     *htmlquerier.Q("//h2").Before("/"),
		DayQuerier:       *htmlquerier.Q("//h2").After("/").Before("("),
		OpenTimeQuerier:  *htmlquerier.Q("//div[@class='post-content-content']/p[3]/text()[1]").Before("/"),
		StartTimeQuerier: *htmlquerier.Q("//div[@class='post-content-content']/p[3]/text()[1]").After("/"),
		IsMonthInLive:    true,
	},

	PrefectureName: "tokyo",
	AreaName:       "shimokitazawa",
	VenueID:        "shimokitazawa-shangrila",

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         10,
		FirstLiveTitle:        "Girls Face in Shangri-La mini!!",
		FirstLiveArtists:      []string{"Li-V-RAVE", "奏音コレクト", "Fancy Film", "ヒトノユメ", "HALO PALLETE", "SUGAR☆VEGA.\u00adcom"},
		FirstLivePrice:        "予約￥0（1ドリンク￥600別）\n当日￥0（2ドリンク代￥1200別）",
		FirstLivePriceEnglish: "Reservation￥0（1Drink￥600Separately）\nDoor￥0（2Drink代￥1200Separately）",
		FirstLiveOpenTime:     time.Unix(1699153800, 0),
		FirstLiveStartTime:    time.Unix(1699155000, 0),
		FirstLiveURL:          "https://www.shan-gri-la.jp/tokyo/category/schedule/",
	},
}

var ShimokitazawaShelterFetcher = fetchers.CreateLoftFetcher(
	"https://www.loft-prj.co.jp",
	util.InsertYearMonth("https://www.loft-prj.co.jp/schedule/shelter/date/%d/%02d"),
	"tokyo",
	"shimokitazawa",
	"shimokitazawa-shelter",
	fetchers.TestInfo{
		NumberOfLives:         35,
		FirstLiveTitle:        "SHELTER & 山﨑 presents 「あなたのドラム、詳しく聞かせて？ Vol.14」",
		FirstLiveArtists:      []string{"GUEST:松本誠治（the telephones）", "司会＆進行：山﨑聖之（CONFVSE / fam / The Firewood Project / LOW IQ 01 & THE RHYTHM MAKERS）"},
		FirstLivePrice:        "ADV¥2000＋1D / DOOR¥2400＋1D\n【発売日】\nSHELTER予約",
		FirstLivePriceEnglish: "ADV¥2000＋1D / DOOR¥2400＋1D\n【発売日】\nSHELTERReservation",
		FirstLiveOpenTime:     time.Unix(1685613600, 0),
		FirstLiveStartTime:    time.Unix(1685615400, 0),
		FirstLiveURL:          "https://www.loft-prj.co.jp/schedule/shelter/247893",
	},
)

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
		FirstLiveOpenTime:     time.Unix(1698917400, 0),
		FirstLiveStartTime:    time.Unix(1698919200, 0),
		FirstLiveURL:          "https://www.toos.co.jp/3/events/d-b-inches-1st-ep-instinct-filter-bubble-release-party/",
	},
)

var ShimokitazawaWaverFetcher = fetchers.Simple{
	BaseURL:              "https://waverwaver.net/",
	InitialURL:           "https://waverwaver.net/category/schedule/",
	NextSelector:         "//ul[contains(@class, 'page-numbers')]//a[contains(@class, 'next')]",
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

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         49,
		FirstLiveTitle:        "『笑うはこには福来る』",
		FirstLiveArtists:      []string{"小山恭代", "高田えぬひろ", "矢部りんご", "ナカザワ"},
		FirstLivePrice:        "adv¥2,200 / door¥2,500 各＋1drink代別途(¥600)",
		FirstLivePriceEnglish: "adv¥2,200 / door¥2,500 各＋1drink代Separately(¥600)",
		FirstLiveOpenTime:     time.Unix(1698918300, 0),
		FirstLiveStartTime:    time.Unix(1698920100, 0),
		FirstLiveURL:          "https://waverwaver.net/2023/11/02/2023%e5%b9%b411%e6%9c%8802%e6%97%a5%e6%9c%a8/",
	},
}

/***************
 *						 *
 *	Shindaita  *
 *						 *
 ***************/

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
	AreaName:       "shindaita",
	VenueID:        "shindaita-fever",

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         32,
		FirstLiveTitle:        "『百々和宏presents ～Drunk51～』",
		FirstLiveArtists:      []string{"百々和宏（MO’SOME TONEBENDER）", "ホリエアツシ（STRAIGHTENER）", "佐々木亮介（a flood of circle）", "ヤマジカズヒデ（dip）", "有江嘉典（VOLA&THE ORIENTAL MACHINE）", "ウエノコウジ（the HIATUS, Radio Caroline）", "クハラカズユキ（The Birthday）", "有松益男（BACK DROP BOMB）"},
		FirstLivePrice:        "ADV ￥4800 (+1drink) THANK YOU SOLD OUT!!!\n※1drink ￥600",
		FirstLivePriceEnglish: "ADV ￥4800 (+1drink) THANK YOU SOLD OUT!!!\n※1drink ￥600",
		FirstLiveOpenTime:     time.Unix(1698831900, 0),
		FirstLiveStartTime:    time.Unix(1698834600, 0),
		FirstLiveURL:          util.InsertYearMonth("https://www.fever-popo.com/schedule/%d/%02d/"),
	},
}

/**************
 * 					  *
 *	Shinjuku  *
 *						*
 **************/

var ShinjukuLoftFetcher = fetchers.CreateLoftFetcher(
	"http://mu-seum.co.jp",
	util.InsertYearMonth("https://www.loft-prj.co.jp/schedule/loft/date/%d/%02d"),
	"tokyo",
	"shinjuku",
	"shinjuku-loft",
	fetchers.TestInfo{
		NumberOfLives:         38,
		FirstLiveTitle:        "エクストロメ!!",
		FirstLiveArtists:      []string{"tipToe.", "airattic", "Finger Runs", "われらがプワプワプーワプワ"},
		FirstLivePrice:        "ADV¥1500 / DOOR¥未定(DRINK代別¥600)\n[発売]\nLive Pocket 5月26日(金)22:00〜5...",
		FirstLivePriceEnglish: "ADV¥1500 / DOOR¥TBA(DRINK代Separately¥600)\n[発売]\nLive Pocket 5月26日(金)22:00〜5...",
		FirstLiveOpenTime:     time.Unix(1685612400, 0),
		FirstLiveStartTime:    time.Unix(1685614200, 0),
		FirstLiveURL:          "https://www.loft-prj.co.jp/schedule/loft/252531",
	},
)