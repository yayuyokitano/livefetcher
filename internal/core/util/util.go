package util

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if b > a {
		return a
	}
	return b
}

func GetTimeFromString(s string) (hour string, min string, nextDay bool) {
	colon := strings.Index(s, ":")
	if colon == -1 {
		hour = "03"
		min = "24"
		return
	}
	hour = fmt.Sprintf("%02s", stripNonNumeric(s[Max(colon-2, 0):colon]))
	min = fmt.Sprintf("%02s", stripNonNumeric(s[colon+1:Min(colon+3, len(s))]))

	nhour, err := strconv.Atoi(hour)
	if err != nil {
		hour = "03"
		min = "24"
		return
	}
	if nhour >= 24 {
		hour = fmt.Sprintf("%02s", strconv.Itoa(nhour-24))
		nextDay = true
	}

	return
}

func stripNonNumeric(s string) string {
	var result strings.Builder
	for _, b := range s {
		if '0' <= b && b <= '9' {
			result.WriteByte(byte(b))
		}
	}
	return result.String()
}

func getIndex(s []rune, sep rune) int {
	for i, r := range s {
		if r == sep {
			return i
		}
	}
	return -1
}

func GetDate(md []rune, sep rune) (month string, day string, err error) {
	sepIndex := getIndex(md, sep)
	if sepIndex == -1 {
		err = errors.New("No separator found in date string " + string(md))
		return
	}
	month = fmt.Sprintf("%02s", stripNonNumeric(string(md[Max(sepIndex-2, 0):sepIndex])))
	day = fmt.Sprintf("%02s", stripNonNumeric(string(md[sepIndex+1:Min(sepIndex+3, len(md))])))
	return
}

func GetYearMonth(ym []rune, sep rune) (year string, month string, err error) {
	sepIndex := getIndex(ym, sep)
	if sepIndex == -1 {
		err = errors.New("No separator found in date string " + string(ym))
		return
	}
	year = stripNonNumeric(string(ym[:sepIndex]))
	if len(year) == 2 {
		year = "20" + year
	}
	month = fmt.Sprintf("%02s", stripNonNumeric(string((ym[sepIndex+1:]))))
	return
}

func GetYearMonthDay(ymd []rune, sep1 rune, sep2 rune) (year string, month string, day string, err error) {
	sepIndex := getIndex(ymd, sep1)
	if sepIndex == -1 {
		err = errors.New("No separator found in date string " + string(ymd))
		return
	}
	year = stripNonNumeric(string(ymd[:sepIndex]))
	if len(year) == 2 {
		year = "20" + year
	}
	month, day, err = GetDate(ymd[sepIndex+1:], sep2)
	return
}

var bannedArtists = map[string]bool{
	"":              true,
	"and more...":   true,
	"w/":            true,
	"live:":         true,
	"live":          true,
	"［dj］":          true,
	"[live]":        true,
	"[band]":        true,
	"-band-":        true,
	"[dj]":          true,
	"-live-":        true,
	"◉live◉":        true,
	"live :":        true,
	"【出演】":          true,
	"【ゲスト】":         true,
	"【司会】":          true,
	"転換dj":          true,
	"◉転換dj◉":        true,
	"act":           true,
	"act:":          true,
	"-act-":         true,
	"◉act◉":         true,
	"■live":         true,
	"■dj":           true,
	"■food":         true,
	"■one man":      true,
	"■act":          true,
	"■guest":        true,
	"■guest act":    true,
	"■one man show": true,
	"■bar":          true,
	"■vj":           true,
	"■solo":         true,
	"■host":         true,
	"■shop":         true,
	"dj":            true,
	"dj:":           true,
	"-dj-":          true,
	"◉dj◉":          true,
	"dj :":          true,
	"・料金":           true,
	"料金":            true,
	"ライブ情報":         true,
	"＋1d":           true,
	"+1d":           true,
	"host dj:":      true,
	"host dj":       true,
	"judge":         true,
	"judge:":        true,
	"-judge-":       true,
	"mc":            true,
	"mc:":           true,
	"[selectas]":    true,
	"[on stage]":    true,
	"-mc-":          true,
	"-selector-":    true,
	"-mtr live-":    true,
	"- 再入場可 *再入場毎にドリンク代頂きます / a drink ticket fee charged at every re-entry": true,
	"at spread": true,
	"東京都世田谷区北沢2-12-6 リバーストーンビルb1f": true,
}

var bannedSubstrings = []string{
	"http://",
	"https://",
	"shimokitazawa",
	"livehaus",
	"live haus",
	"下北沢",
	"コメント",
	"リリース",
	"album",
	"アルバム",
	"vol.",
	"food:",
	"【food】",
	"出演者",
}

var bannedRegexes = []string{
	`(?:(?:¥[\d,]+)|(?:[\d,]+円))`, // price
	`\d{2}:\d{2}`,                 // time
	`\d{2}：\d{2}`,                 // time
	`adv.*door`,                   // labels
	`door.*adv`,                   // labels
	`open.*start`,                 // labels
	`start.*open`,                 // labels
	`\d{2}[/. ]\s*\d{2}[/. ]`,     // date
	`\d{2}[/. ]\s*\d{2}\s*\(.*\)`, // date
}

var removable = []string{
	"＜ONE MAN＞",
	"( GUEST ACT )",
	"GUEST DJ : ",
	"DJ：",
}

var prefixes = []string{
	"•",
	"✰",
}

func ProcessArtists(a []string) (artists []string) {
	artists = []string{}
	for _, artist := range a {
		for _, r := range removable {
			artist = strings.Replace(artist, r, "", -1)
		}
		artist = strings.TrimSpace(artist)
		lower := strings.ToLower(artist)
		if bannedArtists[lower] {
			continue
		}
		var isBanned bool
		for _, substr := range bannedSubstrings {
			if strings.Contains(lower, substr) {
				isBanned = true
				continue
			}
		}
		if isBanned {
			continue
		}

		for _, regex := range bannedRegexes {
			re, err := regexp.Compile(regex)
			if err != nil {
				continue
			}
			if re.MatchString(lower) {
				isBanned = true
				continue
			}
		}
		if isBanned {
			continue
		}

		for _, prefix := range prefixes {
			artist = strings.TrimPrefix(artist, prefix)
		}
		artists = append(artists, artist)
	}
	return
}

func InsertYearMonth(s string) string {
	t := time.Now()
	return fmt.Sprintf(s, t.Year(), int(t.Month()))
}

func InsertShortYearMonth(s string) string {
	t := time.Now()
	return fmt.Sprintf(s, t.Year()%100, int(t.Month()))
}

func SpacedPriceTimeFetcher(d string, s string) (price string, open time.Time, start time.Time, err error) {
	r, err := regexp.Compile(`\s+`)
	if err != nil {
		return
	}
	processed := r.ReplaceAll([]byte(s), []byte(" "))
	split := strings.Split(string(processed), " ")
	for i, v := range split {
		hour, min, nextDay := GetTimeFromString(v)
		if hour == "03" && min == "24" {
			continue
		}
		if open.IsZero() {
			open, err = time.Parse(timeLayout, fmt.Sprintf("%s %s:%s:00 +0900", d, hour, min))
			if err != nil {
				return
			}
			if nextDay {
				open = open.AddDate(0, 0, 1)
			}
		} else {
			start, err = time.Parse(timeLayout, fmt.Sprintf("%s %s:%s:00 +0900", d, hour, min))
			if err != nil {
				return
			}
			if nextDay {
				start = start.AddDate(0, 0, 1)
			}
			price = strings.Join(split[i+1:], " ")
			return
		}
	}
	defaultTime, err := time.Parse(timeLayout, fmt.Sprintf("%s %s:%s:00 +0900", d, "03", "24"))
	if err != nil {
		return
	}
	price, open, start = "", defaultTime, defaultTime
	return
}

func GetUniqueVenues(a []LiveHouse) (b []LiveHouse) {
	m := make(map[string]bool)
	for _, v := range a {
		if !m[v.ID] {
			m[v.ID] = true
			b = append(b, v)
		}
	}
	return
}

func GetUniqueVenueIDs(a []LiveHouse) (b []string) {
	m := make(map[string]bool)
	for _, v := range a {
		if !m[v.ID] {
			m[v.ID] = true
			b = append(b, v.ID)
		}
	}
	return
}

func findNthTime(s string, n int) string {
	re, err := regexp.Compile(`\d{2}:\d{2}`)
	if err != nil {
		return "03:24"
	}
	matches := re.FindAllString(s, n)
	if matches == nil {
		return "03:24"
	}
	return matches[len(matches)-1]
}

func prefixToN(prefix string) int {
	switch prefix {
	case "open":
		return 1
	case "start":
		return 2
	default:
		return -1
	}
}

func FindTime(s string, prefix string) string {
	arr := strings.Split(strings.ToLower(s), prefix)
	if len(arr) < 2 {
		return findNthTime(s, prefixToN(prefix))
	}
	var str string
	if strings.HasSuffix(strings.TrimSpace(arr[0]), "/") {
		tmp := strings.Split(arr[1], "/")
		if len(tmp) < 2 {
			return findNthTime(s, prefixToN(prefix))
		}
		str = strings.TrimSpace(tmp[1])[0:5]
	} else {
		str = strings.TrimSpace(arr[1])[0:5]
	}

	re, err := regexp.Compile(`\d{2}:\d{2}`)
	if err != nil {
		return findNthTime(s, prefixToN(prefix))
	}
	if re.MatchString(str) {
		return str
	}
	return findNthTime(s, prefixToN(prefix))
}

func FindPrice(arr []string) string {
	re, err := regexp.Compile(`[^\s]*\s?(?:(?:[¥￥][\d,]+)|(?:[\d,]+円))(?:\s*\+\d?(?:(?:D)|(?:ドリンク)))?`)
	if err != nil {
		return ""
	}
	for _, s := range arr {
		arr := re.FindAllString(s, 2)
		if arr != nil {
			return strings.Join(arr, "、")
		}
	}
	return ""
}

// GetRelevantYear gets the year for a given month.
// Some connectors have no way to get the year from DOM, so we make a basic set of assumptions:
//
// 1. If the month of the live is equal to or greater than the current month, assume the live is in the current year.
//
// 2. If the month of the live is less than the current month, assume the live is next year.
func GetRelevantYear(month int) int {
	now := time.Now()
	if month < int(now.Month()) {
		return time.Now().Year() + 1
	}
	return time.Now().Year()
}

var JapanTime = time.FixedZone("UTC+9", +9*60*60)
