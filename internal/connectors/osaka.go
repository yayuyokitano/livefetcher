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

var ShinsaibashiAnimaFetcher = fetchers.Simple{
	BaseURL:             "https://liveanima.jp/",
	InitialURL:          "https://liveanima.jp/?page_id=174",
	LiveSelector:        "//div[contains(@class, 'eo-events')]/div[@class='container']",
	DetailsLinkSelector: "//a",
	TitleQuerier:        *htmlquerier.Q("//a/span"),
	ArtistsQuerier:      *htmlquerier.Q("//span[text()='出演者']/following-sibling::text()").SplitIgnoreWithin("[/\n]", '(', ')'),
	PriceQuerier:        *htmlquerier.Q("//span[text()='PRICE']/following-sibling::text()").ReplaceAllRegex(`(\s| )+`, " "),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//h5/span[@class='animabadge mont']"),
		MonthQuerier:     *htmlquerier.Q("//h5/span[@class='animabadge mont']").After("/"),
		DayQuerier:       *htmlquerier.Q("//h5/span[@class='animabadge mont']").After("/").After("/"),
		OpenTimeQuerier:  *htmlquerier.Q("//span[text()='OPEN']/following-sibling::text()"),
		StartTimeQuerier: *htmlquerier.Q("//span[text()='START']/following-sibling::text()"),

		IsMonthInLive: true,
		IsYearInLive:  true,
	},

	PrefectureName: "osaka",
	AreaName:       "shinsaibashi",
	VenueID:        "shinsaibashi-anima",

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         36,
		FirstLiveTitle:        "EMERGENCY CALL",
		FirstLiveArtists:      []string{"Absopetus-アブソプ-", "CYCLONISTA", "MAGMAZ", "MiNO", "mistress", "ココロシンドローム", "処刑台のシンデレラ"},
		FirstLivePrice:        "ADV ¥1500 DOOR ¥2000 ＋1drink(¥600)",
		FirstLivePriceEnglish: "ADV ¥1500 DOOR ¥2000 ＋1drink(¥600)",
		FirstLiveOpenTime:     time.Date(2024, 3, 6, 18, 45, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 6, 19, 15, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://liveanima.jp/live/gig/emergency-call",
	},
}

var ShinsaibashiBeyondFetcher = fetchers.CreateBassOnTopFetcher(
	"https://beyond-osaka.jp/",
	"https://beyond-osaka.jp/schedule/calendar/20%d/%02d/",
	"osaka",
	"shinsaibashi",
	"shinsaibashi-beyond",
	fetchers.TestInfo{
		NumberOfLives:         33,
		FirstLiveTitle:        "おかんに見られたくない いや、見て欲しいVol.4 らくだのレコ発",
		FirstLiveArtists:      []string{"らくだのこぶX", "セックスマシーン!!", "百回中百回", "Blow the instability(O.A)"},
		FirstLivePrice:        "ADV/DOOR ￥3,600/￥4,000（別途1Drink代金¥600-必要）",
		FirstLivePriceEnglish: "ADV/DOOR ￥3,600/￥4,000（Separately1DrinkPrice¥600-Necessary）",
		FirstLiveOpenTime:     time.Date(2024, 3, 1, 18, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 1, 18, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://beyond-osaka.jp/schedule/detail/29388",
	},
)

var ShinsaibashiBigcatFetcher = fetchers.Simple{
	BaseURL:              "https://bigcat-live.com/",
	ShortYearIterableURL: "https://bigcat-live.com/20%d/%d",
	LiveSelector:         "//div[contains(@class, 'archive_block')]",
	TitleQuerier:         *htmlquerier.Q("//h3[@class='ttl']"),
	ArtistsQuerier:       *htmlquerier.Q("//dt[text()='LIVE INFO']/following-sibling::dd/p").Split("/"),
	PriceQuerier:         *htmlquerier.QAll("//dt[text()='ADV' or text()='DOOR']/ancestor::dl").Join(" ").ReplaceAllRegex(`\s+`, " "),

	TimeHandler: fetchers.TimeHandler{
		MonthQuerier:     *htmlquerier.Q("//span[@class='date_txt']"),
		DayQuerier:       *htmlquerier.Q("//span[@class='date_txt']").After("."),
		OpenTimeQuerier:  *htmlquerier.Q("//dt[text()='OPEN']/following-sibling::dd"),
		StartTimeQuerier: *htmlquerier.Q("//dt[text()='START']/following-sibling::dd"),

		IsMonthInLive: true,
	},

	PrefectureName: "osaka",
	AreaName:       "shinsaibashi",
	VenueID:        "shinsaibashi-bigcat",

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         23,
		FirstLiveTitle:        "押忍フェス in BIGCAT",
		FirstLiveArtists:      []string{"KRD8", "WT☆Egret", "すたんぴっ！", "シンセカイヒーロー", "森ふうか", "HIGH SPY DOLL", "ミケネコガールズ", "Mellow giRLs", "Vress", "LOViSH", "caprice", "frecia", "いつでも夢を", "link start", "REBEL REBEL", "EVERYTHING IS WONDER", "Lunouir Tiara", "イロハサクラ"},
		FirstLivePrice:        "ADV 優先：￥2,400一般：￥1,000 DOOR 優先：￥3,400一般：￥2,000",
		FirstLivePriceEnglish: "ADV Priority entry：￥2,400Ordinary Ticket：￥1,000 DOOR Priority entry：￥3,400Ordinary Ticket：￥2,000",
		FirstLiveOpenTime:     time.Date(util.GetRelevantYear(3), 3, 1, 15, 15, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(util.GetRelevantYear(3), 3, 1, 15, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://bigcat-live.com/20%d/%d",
	},
}

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

var ShinsaibashiClapperFetcher = fetchers.Simple{
	BaseURL:              "https://clapper.jp/",
	ShortYearIterableURL: "https://clapper.jp/data/category/20%d-%02d/",
	LiveSelector:         "//ul[@id='scheduleList']/li",
	TitleQuerier:         *htmlquerier.Q("//h4[@class='event_name']").CutWrapper("『", "』"),
	ArtistsQuerier:       *htmlquerier.QAll("//h5[text()='出演']/following-sibling::p[1]/text()"),
	PriceQuerier:         *htmlquerier.QAll("//h5[text()='料金']/following-sibling::text()").Join(" "),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//span[@class='ev_date']"),
		MonthQuerier:     *htmlquerier.Q("//span[@class='ev_date']").After("."),
		DayQuerier:       *htmlquerier.Q("//span[@class='ev_date']").After(".").After("."),
		OpenTimeQuerier:  *htmlquerier.QAll("//h5[text()='OPEN／START']/following-sibling::text()").Join(""),
		StartTimeQuerier: *htmlquerier.QAll("//h5[text()='OPEN／START']/following-sibling::text()").Join("").After(":"),

		IsYearInLive:  true,
		IsMonthInLive: true,
	},

	PrefectureName: "osaka",
	AreaName:       "shinsaibashi",
	VenueID:        "shinsaibashi-clapper",

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         23,
		FirstLiveTitle:        "大阪最終単独公演『集諦』",
		FirstLiveArtists:      []string{"NIGAI"},
		FirstLivePrice:        "前売¥5,000-(1D別)　当日¥0-(1D別)",
		FirstLivePriceEnglish: "Reservation¥5,000-(1 Drink purchase required)　Door¥0-(1 Drink purchase required)",
		FirstLiveOpenTime:     time.Date(2024, 3, 7, 17, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 7, 18, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://clapper.jp/data/category/2024-03/",
	},
}

var ShinsaibashiClubVijonFetcher = fetchers.CreateBassOnTopFetcher(
	"https://vijon.jp/",
	"https://vijon.jp/schedule/calendar/20%d/%02d/",
	"osaka",
	"shinsaibashi",
	"shinsaibashi-clubvijon",
	fetchers.TestInfo{
		NumberOfLives:         38,
		FirstLiveTitle:        "FREE! FREE! FREE!",
		FirstLiveArtists:      []string{"Louve noir", "ZERO", "sp1at", "CARKLAND", "LOCO"},
		FirstLivePrice:        "ADV/DOOR ￥0 別途2Drink代￥1,200",
		FirstLivePriceEnglish: "ADV/DOOR ￥0 Separately2Drink代￥1,200",
		FirstLiveOpenTime:     time.Date(2024, 3, 1, 17, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 1, 18, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://vijon.jp/schedule/detail/33505",
	},
)

var ShinsaibashiConpassFetcher = fetchers.Simple{
	BaseURL:              "https://www.conpass.jp/",
	ShortYearIterableURL: "https://www.conpass.jp/?cat=4&m=20%d%02d",
	LiveSelector:         "//div[@id='main']//ul/div[not(@class) and not(@id)]",
	ExpandedLiveSelector: "//a",
	TitleQuerier:         *htmlquerier.Q("//p[@class='event_tittle']"),
	ArtistsQuerier:       *htmlquerier.QAll("//p[text()='LINEUP:']/following-sibling::p[1]/text()"),
	PriceQuerier:         *htmlquerier.QAll("//p[text()='CHARGE:']/following-sibling::p[1]/text()").Join(" "),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//span[@class='event_day_in']"),
		MonthQuerier:     *htmlquerier.Q("//span[@class='event_day_in']").After("."),
		DayQuerier:       *htmlquerier.Q("//span[@class='event_day_in']").After(".").After("."),
		OpenTimeQuerier:  *htmlquerier.Q("//span[text()='INFORMATION:']/following-sibling::text()[1]"),
		StartTimeQuerier: *htmlquerier.Q("//span[text()='INFORMATION:']/following-sibling::text()[1]").After(":"),

		IsYearInLive:  true,
		IsMonthInLive: true,
	},

	PrefectureName: "osaka",
	AreaName:       "shinsaibashi",
	VenueID:        "shinsaibashi-conpass",

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         16,
		FirstLiveTitle:        "West Side Unity presents. 『LEAVE YOUTH HERE -EXTRA PARTY-』",
		FirstLiveArtists:      []string{"Demonstration Of Power (UK)", "Despize (UK)", "SAND", "Decasion", "UNMASK aLIVE", "RESENTMENT", "waterweed", "ReVERSE BOYZ", "UNHOLY11", "Fallen Grace", "CE$", "MOON SHOW Fr. JAH WORKS", "DJ ACE Fr. JAH WORKS"},
		FirstLivePrice:        "前売 ¥3,500(D別) 当日 ¥4,000(D別)",
		FirstLivePriceEnglish: "Reservation ¥3,500(Drinks sold separately) Door ¥4,000(Drinks sold separately)",
		FirstLiveOpenTime:     time.Date(2024, 3, 3, 14, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 3, 14, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://www.conpass.jp/7168.html",
	},
}

var ShinsaibashiDropFetcher = fetchers.CreateBassOnTopFetcher(
	"https://vijon.jp/",
	"https://vijon.jp/schedule/calendar/20%d/%02d/",
	"osaka",
	"shinsaibashi",
	"shinsaibashi-drop",
	fetchers.TestInfo{
		NumberOfLives:         39,
		FirstLiveTitle:        "世史久祭り大阪編vol.23 春のドドスコベイベーナイト",
		FirstLiveArtists:      []string{"世史久", "MEGAHORN", "ELBRUNCH", "イチゼンバッカー", "飛太", "ほーDK", "T-face", "浦田哲也", "ウルトラソウル", "10ripeee", "田中佑生大", "竹歳みずほ", "Mifuyu", "林奈恵"},
		FirstLivePrice:        "♢(来場)￥3.800 別途1D代要 ♢(配信)3000円",
		FirstLivePriceEnglish: "♢(In Person)￥3.800 1 Drink must be purchased separately ♢(Livestream)3000円",
		FirstLiveOpenTime:     time.Date(2024, 3, 1, 16, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 1, 17, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://clubdrop.jp/schedule/detail/32308",
	},
)

var ShinsaibashiFanjtwiceFetcher = fetchers.Simple{
	BaseURL:        "http://www.fanj-twice.com/",
	InitialURL:     "http://www.fanj-twice.com/sch_twice/sch000.html",
	LiveSelector:   "//div[contains(@class, 'cssskin-_block_main_news_m')]",
	TitleQuerier:   *htmlquerier.Q("//h4"),
	ArtistsQuerier: *htmlquerier.Q("//p[@class='c-body']/span[2]//span[not(*)]").Split(" / "),
	PriceQuerier:   *htmlquerier.Q("//p[@class='c-body']/span[3]//span[@class='d-bold']/following-sibling::text()"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//p[@class='c-lead']"),
		MonthQuerier:     *htmlquerier.Q("//p[@class='c-lead']").After("."),
		DayQuerier:       *htmlquerier.Q("//p[@class='c-lead']").After(".").After("."),
		OpenTimeQuerier:  *htmlquerier.Q("//p[@class='c-body']/span[3]//span[@class='d-bold']/preceding-sibling::text()"),
		StartTimeQuerier: *htmlquerier.Q("//p[@class='c-body']/span[3]//span[@class='d-bold']/preceding-sibling::text()").After("/"),

		IsYearInLive:  true,
		IsMonthInLive: true,
	},

	PrefectureName: "osaka",
	AreaName:       "shinsaibashi",
	VenueID:        "shinsaibashi-fanjtwice",

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         12,
		FirstLiveTitle:        "GEM JAM FES！",
		FirstLiveArtists:      []string{"meluQ", "SITRA.", "青のメロディー", "caprice", "AsteRythm", "Stella!", "ネテルダイヤ", "ゆめポケ", "めいてん"},
		FirstLivePrice:        "前売り 1,900円 / 前方エリア 3,000円（別途1DRINK）",
		FirstLivePriceEnglish: "Reservation 1,900円 / Front area 3,000円（Separately1DRINK）",
		FirstLiveOpenTime:     time.Date(2024, 3, 2, 10, 15, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 2, 10, 35, 0, 0, util.JapanTime),
		FirstLiveURL:          "http://www.fanj-twice.com/sch_twice/sch000.html",
	},
}

var ShinsaibashiHokageFetcher = fetchers.Simple{
	BaseURL:                     "http://musicbarhokage.net/",
	ShortYearReverseIterableURL: "http://musicbarhokage.net/schedule%d_20%d.htm",
	LiveSelector:                "//table[@bordercolor='#FF0000']/tbody/tr/td/div/table/tbody",
	TitleQuerier:                *htmlquerier.Q("/tr[3]//strong"),
	ArtistsQuerier:              *htmlquerier.Q("/tr[4]//strong").Split("\n"),
	PriceQuerier:                *htmlquerier.Q("/tr[6]//strong"),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//span[@class='style15']"),
		MonthQuerier:     *htmlquerier.Q("//span[@class='style15']").After("."),
		DayQuerier:       *htmlquerier.Q("//span[@class='style15']").After(".").After("."),
		OpenTimeQuerier:  *htmlquerier.Q("/tr[5]//strong").After("OPEN:"),
		StartTimeQuerier: *htmlquerier.Q("/tr[5]//strong").After("START:"),

		IsYearInLive:  true,
		IsMonthInLive: true,
	},

	PrefectureName: "osaka",
	AreaName:       "shinsaibashi",
	VenueID:        "shinsaibashi-hokage",

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         24,
		FirstLiveTitle:        "[DUB UNDERGROUND vol.4]",
		FirstLiveArtists:      []string{"DUBBY BON", "Manabu Dub", "Black Warriyah", "BIG \"DUB\" HEAD (fr.Medical Tempo)"},
		FirstLivePrice:        "Adv.1000yen Door.1000yen (+Drink fee)",
		FirstLivePriceEnglish: "Adv.1000yen Door.1000yen (+Drink fee)",
		FirstLiveOpenTime:     time.Date(2024, 3, 1, 21, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 1, 21, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "http://musicbarhokage.net/schedule%d_20%d.htm",
	},
}

var ShinsaibashiJanusFetcher = fetchers.Simple{
	BaseURL:              "https://janusosaka.com/",
	ShortYearIterableURL: "https://janusosaka.com/schedule/20%d-%02d/",
	LiveSelector:         "//article[@class='c-scheduleList']",
	TitleQuerier:         *htmlquerier.Q("//div[@class='c-scheduleList__head--title']"),
	ArtistsQuerier:       *htmlquerier.Q("//div[@class='c-scheduleList__head--act']").SplitIgnoreWithin(`( / )|(【?((opening)|(Opening)|(OPENING))?\s*((Guest)|(guest)|(GUEST)|(ゲスト)|(artist)|(Artist)|(ARTIST)|(act)|(Act)|(ACT))\s*((artist)|(Artist)|(ARTIST)|(act)|(Act)|(ACT))?】?((\s*):)?)|(O.A.(\s*):)|(【DJ/MC】)|(【LIVE】)|(\(O.A.\))`, '(', ')'), // dont worry about it
	PriceQuerier:         *htmlquerier.Q("//dt[text()='ADV/DOOR']/following-sibling::dd").Prefix("ADV/DOOR: ").ReplaceAllRegex(`\s+`, " "),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//div[@id='c-breadcrumb']//li[last()]"),
		MonthQuerier:     *htmlquerier.Q("//div[@class='c-scheduleList__date--month']"),
		DayQuerier:       *htmlquerier.Q("//div[@class='c-scheduleList__date--date']"),
		OpenTimeQuerier:  *htmlquerier.Q("//dt[text()='OPEN/START']/following-sibling::dd"),
		StartTimeQuerier: *htmlquerier.Q("//dt[text()='OPEN/START']/following-sibling::dd").After("/"),

		IsMonthInLive: true,
	},

	PrefectureName: "osaka",
	AreaName:       "shinsaibashi",
	VenueID:        "shinsaibashi-janus",

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         28,
		FirstLiveTitle:        "帝国喫茶ワンマンツアー2024 「きみの待つ場所へ春のメロディーを」",
		FirstLiveArtists:      []string{"帝国喫茶"},
		FirstLivePrice:        "ADV/DOOR: ￥3,800 / 未定",
		FirstLivePriceEnglish: "ADV/DOOR: ￥3,800 / TBA",
		FirstLiveOpenTime:     time.Date(2024, 3, 1, 18, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://janusosaka.com/schedule/20%d-%02d/",
	},
}

var ShinsaibashiKingCobraFetcher = fetchers.Simple{
	BaseURL:              "http://king-cobra.net/",
	ShortYearIterableURL: "http://king-cobra.net/schedule/20%d_%d.html",
	LiveSelector:         "//font[@color='#00CCFF' and string-length(normalize-space(text())) > 10]",
	TitleQuerier:         *htmlquerier.Q("/.").Trim().CutWrapper("『", "』").ReplaceAllRegex(`\s+`, " "),
	ArtistsQuerier:       *htmlquerier.QAll("/ancestor::tr[1]/following-sibling::tr[1]/td[1]//text()").DeleteFrom("[FOOD]"),
	PriceQuerier:         *htmlquerier.QAll("/ancestor::tr[1]/following-sibling::tr[1]/td[3]//text()").Join("").ReplaceAllRegex(`\s+`, " "),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//font[@color='#FF33CC']"),
		MonthQuerier:     *htmlquerier.Q("/ancestor::td[1]/preceding-sibling::td[1]//text()[contains(., '月')]"),
		DayQuerier:       *htmlquerier.Q("/ancestor::td[1]/preceding-sibling::td[1]//text()[contains(., '月')]").After("月"),
		OpenTimeQuerier:  *htmlquerier.Q("/ancestor::tr[1]/following-sibling::tr[1]/td[2]//text()[contains(., '開場')]"),
		StartTimeQuerier: *htmlquerier.Q("/ancestor::tr[1]/following-sibling::tr[1]/td[2]//text()[contains(., '開演')]"),

		IsMonthInLive: true,
	},

	PrefectureName: "osaka",
	AreaName:       "shinsaibashi",
	VenueID:        "shinsaibashi-kingcobra",

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         10,
		FirstLiveTitle:        "ヘルスマパンク 春の電波ジャック!!",
		FirstLiveArtists:      []string{"ギターパンダ", "THE FLYING PANTS", "アンモニアンズ", "JOKE?!", "THE MAYUCHIX", "ラティーノ山口", "大義"},
		FirstLivePrice:        "・ADV.3,500 ・DOOR.4,000",
		FirstLivePriceEnglish: "・ADV.3,500 ・DOOR.4,000",
		FirstLiveOpenTime:     time.Date(2024, 3, 2, 17, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 2, 17, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "http://king-cobra.net/schedule/20%d_%d.html",
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
		FirstLivePriceEnglish: "Reservation3,500円/Door4,000円(DrinkNot included in ticket600円)/Camera fee＋1,000円",
		FirstLiveOpenTime:     time.Date(2024, 3, 10, 18, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 10, 18, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://livehouse-kurage.com/schedule/%e5%a4%a9%e5%a5%b3%e7%a5%9e%e6%a8%82-%e7%a5%9e%e6%a5%bd%e7%a5%ad%e3%80%8e%e3%80%9c%e8%8a%b1%e3%80%9c%e3%80%8f/",
	},
}

var ShinsaibashiLoftPlusOneWestFetcher = fetchers.CreateLoftFetcher(
	"https://www.loft-prj.co.jp/",
	"https://www.loft-prj.co.jp/schedule/west/date/20%d/%02d",
	"osaka",
	"shinsaibashi",
	"shinsaibashi-loftplusonewest",
	fetchers.TestInfo{
		NumberOfLives:         25,
		FirstLiveTitle:        "はっぴー空間",
		FirstLiveArtists:      []string{"ChanceMovement"},
		FirstLivePrice:        "◎観覧について\n前売,当日共に￥1,500(共に1オーダー必須（￥500以上）)\n■観覧チ...",
		FirstLivePriceEnglish: "◎観覧について\nReservation,Door共に￥1,500(共に1オーダー必須（￥500以上）)\n■観覧チ...",
		FirstLiveOpenTime:     time.Date(2024, 3, 2, 12, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 2, 12, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://www.loft-prj.co.jp/schedule/west/277016",
	},
)

var ShinsaibashiMuseFetcher = fetchers.Simple{
	BaseURL:              "http://osaka.muse-live.com/",
	ShortYearIterableURL: "http://osaka.muse-live.com/schedule/?y=20%d&m=%d",
	LiveSelector:         "//article[@class='media schedule']",
	TitleQuerier:         *htmlquerier.Q("//h3"),
	ArtistsQuerier:       *htmlquerier.QAll("//div[@class='schedule_content']/p[1]/a"),
	PriceQuerier:         *htmlquerier.Q("//ul[@class='schedule_info_list']/li[2]/span[2]").ReplaceAllRegex(`\s+`, " "),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//div[@class='schedule_date']/span"),
		MonthQuerier:     *htmlquerier.Q("//span[@class='month']"),
		DayQuerier:       *htmlquerier.Q("//span[@class='month']/following-sibling::text()"),
		OpenTimeQuerier:  *htmlquerier.Q("//ul[@class='schedule_info_list']/li[1]/span[2]"),
		StartTimeQuerier: *htmlquerier.Q("//ul[@class='schedule_info_list']/li[1]/span[2]").After("/"),

		IsMonthInLive: true,
	},

	PrefectureName: "osaka",
	AreaName:       "shinsaibashi",
	VenueID:        "shinsaibashi-muse",

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         28,
		FirstLiveTitle:        "FASTMUSIC CARNIVAL TOUR2024",
		FirstLiveArtists:      []string{"Bentham", "SAKANAMON", "板歯目"},
		FirstLivePrice:        "ADV.¥4,000 入場時DRINK代別途600円必要",
		FirstLivePriceEnglish: "ADV.¥4,000 When enteringDRINKNot included in ticket600円Necessary",
		FirstLiveOpenTime:     time.Date(2024, 3, 1, 18, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "http://osaka.muse-live.com/schedule/?y=20%d&m=%d",
	},
}

var ShinsaibashiPangeaFetcher = fetchers.Simple{
	BaseURL:             "https://liveanima.jp/",
	InitialURL:          "https://livepangea.com/schedule/",
	LiveSelector:        "//div[contains(@class, 'eo-events')]/div",
	DetailsLinkSelector: "//a",
	TitleQuerier:        *htmlquerier.Q("//a").CutWrapper(`"`, `"`),
	ArtistsQuerier:      *htmlquerier.Q("//span[text()='出演者']/following-sibling::p").SplitIgnoreWithin("[/\n]", '(', ')'),
	PriceQuerier:        *htmlquerier.Q("//span[text()='PRICE']/following-sibling::text()").ReplaceAllRegex(`(\s| )+`, " "),

	TimeHandler: fetchers.TimeHandler{
		MonthQuerier:     *htmlquerier.Q("//p[contains(@class,'live_mom')]"),
		DayQuerier:       *htmlquerier.Q("//p[contains(@class,'live_day')]"),
		OpenTimeQuerier:  *htmlquerier.Q("//span[text()='OPEN']/following-sibling::text()"),
		StartTimeQuerier: *htmlquerier.Q("//span[text()='START']/following-sibling::text()"),

		IsMonthInLive: true,
	},

	PrefectureName: "osaka",
	AreaName:       "shinsaibashi",
	VenueID:        "shinsaibashi-pangea",

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         69,
		FirstLiveTitle:        "Lyanas 1st mini Album「Telescope of You」Release Tour “共鳴シンフォニア”",
		FirstLiveArtists:      []string{"Lyanas", "Cleo", "しゃららんベイビーズ", "Serpent Stellar"},
		FirstLivePrice:        "ADV ¥2500 DOOR ¥3000 【＋1drink(¥600)】",
		FirstLivePriceEnglish: "ADV ¥2500 DOOR ¥3000 【＋1drink(¥600)】",
		FirstLiveOpenTime:     time.Date(util.GetRelevantYear(3), 3, 7, 18, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(util.GetRelevantYear(3), 3, 7, 19, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://livepangea.com/live/event-17706",
	},
}

// this is insanity
var ShinsaibashiUtausakanaFetcher = fetchers.Simple{
	BaseURL:             "http://utausakana.com/",
	InitialURL:          "http://utausakana.com/menu/",
	NextSelector:        "//div[@class='pager']/a[@class='next']",
	LiveSelector:        "//div[@class='menu_body']/p[text()='act)']",
	TitleQuerier:        *htmlquerier.Q("/preceding::p[1]").CutWrapper("『", "』"),
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

var ShinsaibashiVaronFetcher = fetchers.CreateBassOnTopFetcher(
	"https://osaka-varon.jp/",
	"https://osaka-varon.jp/schedule/calendar/20%d/%02d/",
	"osaka",
	"shinsaibashi",
	"shinsaibashi-varon",
	fetchers.TestInfo{
		NumberOfLives:         29,
		FirstLiveTitle:        "The Dust’n’Bonez 「20th anniversary」",
		FirstLiveArtists:      []string{"The Dust’n’Bonez"},
		FirstLivePrice:        "ADV/DOOR ￥5000/￥5500(1D別・整理番号付・税込)",
		FirstLivePriceEnglish: "ADV/DOOR ￥5000/￥5500(1 Drink purchase required・Numbered tickets (may affect entry order)・Incl. Tax)",
		FirstLiveOpenTime:     time.Date(2024, 3, 1, 18, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://osaka-varon.jp/schedule/detail/31089",
	},
)
