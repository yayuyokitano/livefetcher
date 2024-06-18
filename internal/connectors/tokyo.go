package connectors

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/yayuyokitano/livefetcher/internal/core/fetchers"
	"github.com/yayuyokitano/livefetcher/internal/core/htmlquerier"
	"github.com/yayuyokitano/livefetcher/internal/core/util"
	"golang.org/x/net/html"
)

/*************
 *           *
 *  Shibuya  *
 *           *
 *************/

var ShibuyaClubQuattroFetcher = fetchers.Simple{
	BaseURL:              "https://www.club-quattro.com/shibuya/schedule/",
	ShortYearIterableURL: "https://www.club-quattro.com/shibuya/schedule/?ym=20%d%02d",
	LiveSelector:         "//div[@class='schedule-list']/div",
	DetailsLinkSelector:  "//a",
	TitleQuerier:         *htmlquerier.QAll("//p[@class='event-text' or @class='event-ttl']").FilterTitle(`[/\n]`, 1),
	ArtistsQuerier:       *htmlquerier.QAll("//p[@class='event-text' or @class='event-ttl']").FilterArtist(`[/\n]`, 0).SplitRegex(`[/\n]`),
	PriceQuerier:         *htmlquerier.Q("//div[contains(./text(), '料金：')]").After("料金："),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//section[@id='schedule-list-top']/p[@class='ttl']"),
		MonthQuerier:     *htmlquerier.Q("//li[@class='schedule-text slick-current']/p[@class='date']"),
		DayQuerier:       *htmlquerier.Q("//p[@class='date']"),
		OpenTimeQuerier:  *htmlquerier.Q("//div[contains(./text(), 'Open:')]").After("Open:"),
		StartTimeQuerier: *htmlquerier.Q("//div[contains(./text(), 'Start:')]").After("Start:"),
	},

	PrefectureName: "tokyo",
	AreaName:       "shibuya",
	VenueID:        "shibuya-clubquattro",

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         27,
		FirstLiveTitle:        "Rei Release Tour 2024 “VOICE MESSAGE”",
		FirstLiveArtists:      []string{"Rei", "澤村一平(dr)", "真船勝博(ba)", "TAIHEI(kb)", "須原杏(violin)"},
		FirstLivePrice:        "前売 ￥5,500",
		FirstLivePriceEnglish: "Reservation ￥5,500",
		FirstLiveOpenTime:     time.Date(2024, 3, 1, 18, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2024, 3, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://www.club-quattro.com/shibuya/schedule/detail.php?id=15380",
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

	TestInfo: fetchers.TestInfo{
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
	PriceQuerier:              *htmlquerier.Q("//img[contains(@src, 'sche_adv_ttl.gif')]/parent::td/following-sibling::td").ReplaceAllRegex(`\s+`, " "),

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

var ShibuyaLaDonnaFetcher = fetchers.Simple{
	BaseURL:              "https://www.la-donna.jp/",
	ShortYearIterableURL: "https://www.la-donna.jp/schedules/?ym=20%d-%02d",
	LiveSelector:         "//div[@class='sec01']/div[.//dd[@class='bigTxt'][.!='電話受付' and .!='店舗休業日' and .!='企業様イベントご利用']]",
	TitleQuerier:         *htmlquerier.Q("//dd[@class='bigTxt']"),
	ArtistsQuerier:       *htmlquerier.QAll("//dt[.='出演アーティスト']/following-sibling::dd/text()").SplitIgnoreWithin("・", '【', '】'),
	PriceQuerier:         *htmlquerier.Q("//dt[.='前売り / 当日']/parent::*").ReplaceAllRegex(`\s+`, " "),

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
			var n *html.Node
			n, err = html.Parse(strings.NewReader(fmt.Sprintf(
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

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         17,
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

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         23,
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
		FirstLivePriceEnglish: "VIP：¥5,000 / 撮影：¥3,000 / 後方：¥1,500 / Students・Women：¥500 / 新規：¥0 (Incl. Tax / 各Drink代Separately)",
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
		FirstLivePriceEnglish: "Door ¥3,500 (Incl. Tax｜Standing｜Drink代Separately)Reservation ¥3,000 (Incl. Tax｜Standing｜Drink代Separately)",
		FirstLiveOpenTime:     time.Date(2023, 11, 7, 17, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 11, 7, 17, 30, 0, 0, util.JapanTime),
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
		FirstLivePriceEnglish: "¥5,000 / ¥5,500 (Incl. Tax / Drink代Separately)※Door券は19:00~、ReservationのEntryが落ち着き次第 WWW Xにて¥5,500+D代で販売いたします。",
		FirstLiveOpenTime:     time.Date(2023, 11, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 11, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://www-shibuya.jp/schedule/017258.php",
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
		FirstLiveOpenTime:     time.Date(2023, 10, 1, 18, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 10, 1, 18, 30, 0, 0, util.JapanTime),
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
		FirstLivePriceEnglish: "Reservation 2500円(Drinks sold separately)／Door 3000円(Drinks sold separately)",
		FirstLiveOpenTime:     time.Date(2023, 11, 2, 18, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 11, 2, 19, 0, 0, 0, util.JapanTime),
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
		FirstLiveOpenTime:     time.Date(2023, 11, 2, 18, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 11, 2, 18, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          util.InsertYearMonth("http://s-era.jp/schedule_cat/%d-%02d/"),
	},
}

var ShimokitazawaFlowersLoftFetcher = fetchers.CreateLoftFetcher(
	"https://www.loft-prj.co.jp/",
	"https://www.loft-prj.co.jp/schedule/flowersloft/date/20%d/%d",
	"tokyo",
	"shimokitazawa",
	"shimokitazawa-flowersloft",
	fetchers.TestInfo{
		NumberOfLives:         30,
		FirstLiveTitle:        "バニーの日！ハロウィンSP",
		FirstLiveArtists:      []string{"tumiki", "はっちゃん", "テルミー", "ShinK", "クロマティーゆうや(NeoN)", "ディスク百合おん", "まみくん3歳", "マル", "有田清幸", "BIDA"},
		FirstLivePrice:        "男性：2,000 (2drink込み）\n女性：入場無料\n・仮装女子は1Dプレゼント\n・バニー...",
		FirstLivePriceEnglish: "Men：2,000 (2drinkIncluded）\nWomen：EntryFree\n・仮装女子は1Dプレゼント\n・バニー...",
		FirstLiveOpenTime:     time.Date(2023, 11, 1, 18, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 11, 1, 18, 30, 0, 0, util.JapanTime),
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
		FirstLiveOpenTime:     time.Date(2023, 11, 1, 18, 30, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 11, 1, 19, 0, 0, 0, util.JapanTime),
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
	PriceQuerier:   *htmlquerier.Q("//div[@class='post-content-content']/p[3]").After("START ").After("\n").ReplaceAllRegex(`\s+`, " "),

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
	"https://www.loft-prj.co.jp/schedule/shelter/date/20%d/%02d",
	"tokyo",
	"shimokitazawa",
	"shimokitazawa-shelter",
	fetchers.TestInfo{
		NumberOfLives:         35,
		FirstLiveTitle:        "SHELTER & 山﨑 presents 「あなたのドラム、詳しく聞かせて？ Vol.14」",
		FirstLiveArtists:      []string{"GUEST:松本誠治（the telephones）", "司会＆進行：山﨑聖之（CONFVSE / fam / The Firewood Project / LOW IQ 01 & THE RHYTHM MAKERS）"},
		FirstLivePrice:        "ADV¥2000＋1D / DOOR¥2400＋1D\n【発売日】\nSHELTER予約",
		FirstLivePriceEnglish: "ADV¥2000＋1D / DOOR¥2400＋1D\n【発売日】\nSHELTERReservation",
		FirstLiveOpenTime:     time.Date(2023, 6, 1, 19, 0, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 6, 1, 19, 30, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://www.loft-prj.co.jp/schedule/shelter/247893",
	},
)

var ShimokitazawaSpreadFetcher = fetchers.Simple{
	BaseURL:        "https://spread.tokyo/",
	InitialURL:     "https://spread.tokyo/schedule.html",
	LiveSelector:   "//div[@id='c7']/div[@class='box'][position()<200]", // not sure why 200 leads to 99 matches but it does
	TitleQuerier:   *htmlquerier.Q("//u/b").ReplaceAllRegex(`\s+`, " ").CutWrapper(`"`, `"`),
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

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         32,
		FirstLiveTitle:        "『百々和宏presents ～Drunk51～』",
		FirstLiveArtists:      []string{"百々和宏（MO’SOME TONEBENDER）", "ホリエアツシ（STRAIGHTENER）", "佐々木亮介（a flood of circle）", "ヤマジカズヒデ（dip）", "有江嘉典（VOLA&THE ORIENTAL MACHINE）", "ウエノコウジ（the HIATUS, Radio Caroline）", "クハラカズユキ（The Birthday）", "有松益男（BACK DROP BOMB）"},
		FirstLivePrice:        "ADV ￥4800 (+1drink) THANK YOU SOLD OUT!!!\n※1drink ￥600",
		FirstLivePriceEnglish: "ADV ￥4800 (+1drink) THANK YOU SOLD OUT!!!\n※1drink ￥600",
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

var ShinjukuLoftFetcher = fetchers.CreateLoftFetcher(
	"https://www.loft-prj.co.jp/",
	"https://www.loft-prj.co.jp/schedule/loft/date/20%d/%02d",
	"tokyo",
	"shinjuku",
	"shinjuku-loft",
	fetchers.TestInfo{
		NumberOfLives:         38,
		FirstLiveTitle:        "エクストロメ!!",
		FirstLiveArtists:      []string{"tipToe.", "airattic", "Finger Runs", "われらがプワプワプーワプワ"},
		FirstLivePrice:        "ADV¥1500 / DOOR¥未定(DRINK代別¥600)\n[発売]\nLive Pocket 5月26日(金)22:00〜5...",
		FirstLivePriceEnglish: "ADV¥1500 / DOOR¥TBA(DRINK代Separately¥600)\n[発売]\nLive Pocket 5月26日(金)22:00〜5...",
		FirstLiveOpenTime:     time.Date(2023, 6, 1, 18, 40, 0, 0, util.JapanTime),
		FirstLiveStartTime:    time.Date(2023, 6, 1, 19, 10, 0, 0, util.JapanTime),
		FirstLiveURL:          "https://www.loft-prj.co.jp/schedule/loft/252531",
	},
)

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
)
