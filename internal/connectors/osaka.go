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
