package connectors

import (
	"time"

	"github.com/yayuyokitano/livefetcher/lib/core/fetchers"
	"github.com/yayuyokitano/livefetcher/lib/core/htmlquerier"
)

/******************
 * 								*
 *	Shinsaibashi	*
 *								*
 ******************/

var ShinsaibashiBronzeFetcher = fetchers.Simple{
	BaseURL:              "http://osakabronze.com",
	ShortYearIterableURL: "http://osakabronze.com/live/?date=%d.%d",
	LiveSelector:         "//div[@class='day']",
	TitleQuerier:         *htmlquerier.Q("//h3"),
	ArtistsQuerier:       *htmlquerier.Q("//h4").SplitIgnoreWithin(`\n|( \/ )`, '（', '）'),
	PriceQuerier:         *htmlquerier.Q("//h5").ReplaceAllRegex(`\s+`, " ").SplitRegexIndex(`START \d+:\d+ `, 1),

	TimeHandler: fetchers.TimeHandler{
		YearQuerier:      *htmlquerier.Q("//h2").Before("."),
		MonthQuerier:     *htmlquerier.Q("//h2").SplitIndex(".", 1),
		DayQuerier:       *htmlquerier.Q("//h2").SplitRegexIndex("[ .]", 2),
		OpenTimeQuerier:  *htmlquerier.Q("//h5").After("OPEN ").Before(" "),
		StartTimeQuerier: *htmlquerier.Q("//h5").After("START ").Before(" "),

		IsYearInLive:  true,
		IsMonthInLive: true,
	},

	PrefectureName: "osaka",
	AreaName:       "shinsaibashi",
	VenueID:        "shinsaibashi-bronze",

	TestInfo: fetchers.TestInfo{
		NumberOfLives:         24,
		FirstLiveTitle:        "「初志貫徹」\n1st E.P \"アイ E.P\" Release Tour \"愛されたいツアー\"",
		FirstLiveArtists:      []string{"21世記少年", "LAURUS NOBILIS", "endroar", "竜也(Genbu)", "青とフーカ"},
		FirstLivePrice:        "ADV ¥2000 DOOR ¥2500 1D別",
		FirstLivePriceEnglish: "ADV ¥2000 DOOR ¥2500 1 Drink purchase required",
		FirstLiveOpenTime:     time.Unix(1685608200, 0),
		FirstLiveStartTime:    time.Unix(1685610000, 0),
		FirstLiveURL:          "http://osakabronze.com/live/?date=%d.%d",
	},
}
