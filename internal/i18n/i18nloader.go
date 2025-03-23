package i18nloader

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

type SimplifiedLocalizer struct {
	localizer *i18n.Localizer
}

func (l SimplifiedLocalizer) Localize(str string, subs ...string) string {
	subMap := make(map[string]string)

	for i := 0; i < len(subs)-1; i += 2 {
		subMap[subs[i]] = subs[i+1]
	}

	str, err := l.localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    str,
		TemplateData: subMap,
	})
	if err != nil {
		return err.Error()
	}
	return str
}

func Init() (err error) {
	bundle = i18n.NewBundle(language.English)

	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	_, err = bundle.LoadMessageFile("internal/i18n/locales/en_US.toml")
	if err != nil {
		return
	}
	_, err = bundle.LoadMessageFile("internal/i18n/locales/ja_JP.toml")
	if err != nil {
		return
	}
	return
}

func TimeIsIndeterminate(t time.Time) bool {
	return t.Hour() == 3 && t.Minute() == 24
}

func FormatDate(t time.Time, langs []string, includeTime bool) string {
	for _, lang := range langs {
		if strings.HasPrefix(lang, "ja") {
			hour := t.Hour()
			if hour <= 5 && !TimeIsIndeterminate(t) {
				hour += 24
				t = t.AddDate(0, 0, -1)
			}
			weekdayInt := int(t.Weekday())
			var weekday = ""
			switch weekdayInt {
			case 0:
				weekday = "日"
			case 1:
				weekday = "月"
			case 2:
				weekday = "火"
			case 3:
				weekday = "水"
			case 4:
				weekday = "木"
			case 5:
				weekday = "金"
			case 6:
				weekday = "土"
			}

			if !includeTime || TimeIsIndeterminate(t) {
				return fmt.Sprintf("%d年%d月%d日（%s）", t.Year(), int(t.Month()), t.Day(), weekday)
			}
			return fmt.Sprintf("%d年%d月%d日（%s）%02d:%02d", t.Year(), int(t.Month()), t.Day(), weekday, hour, t.Minute())
		}
		if strings.HasPrefix(lang, "en") {
			if !includeTime || TimeIsIndeterminate(t) {
				return t.Format("Mon 2 Jan 2006")
			}
			return t.Format("Mon 2 Jan 2006 03:04 PM")
		}
	}
	if !includeTime || TimeIsIndeterminate(t) {
		return t.Format("Mon 2 Jan 2006")
	}
	return t.Format("Mon 2 Jan 2006 03:04 PM")
}

func FormatOpenStartTime(openTime time.Time, startTime time.Time, langs []string) string {
	dateString := FormatDate(startTime, langs, false)
	openHour := openTime.Hour()
	startHour := startTime.Hour()
	for _, lang := range langs {
		if strings.HasPrefix(lang, "ja") {
			if startHour <= 5 && !TimeIsIndeterminate(startTime) {
				startHour += 24
				if !TimeIsIndeterminate(openTime) {
					openHour += 24
				}
			}
			break
		}
		if strings.HasPrefix(lang, "en") {
			break
		}
	}
	if !TimeIsIndeterminate(openTime) {
		dateString += fmt.Sprintf(" OPEN: %02d:%02d", openHour, openTime.Minute())
	}
	if !TimeIsIndeterminate(startTime) {
		dateString += fmt.Sprintf(" START: %02d:%02d", startHour, startTime.Minute())
	}
	return dateString
}

func GetLanguages(r *http.Request) (langs []string) {
	lang := r.FormValue("lang")
	if lang != "" {
		langs = append(langs, lang)
	}
	accept, _, err := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	if err != nil {
		return
	}
	if len(accept) == 0 {
		return
	}
	for _, lang := range accept {
		langs = append(langs, lang.String())
	}
	return
}

func GetMainLanguage(r *http.Request) string {
	lang := r.FormValue("lang")
	accept := r.Header.Get("Accept-Language")
	_, tag, err := i18n.NewLocalizer(bundle, lang, accept).LocalizeWithTag(&i18n.LocalizeConfig{
		MessageID: "general.artists",
	})
	if err != nil {
		return "en-US"
	}
	return tag.String()
}

func GetMainLanguageFromLangs(langs []string) string {
	_, tag, err := i18n.NewLocalizer(bundle, langs...).LocalizeWithTag(&i18n.LocalizeConfig{
		MessageID: "general.artists",
	})
	if err != nil {
		return "en-US"
	}
	return tag.String()
}

func LocalizerFromLangs(langs []string) SimplifiedLocalizer {
	return SimplifiedLocalizer{i18n.NewLocalizer(bundle, langs...)}
}

func GetLocalizer(r *http.Request) SimplifiedLocalizer {
	lang := r.FormValue("lang")
	accept := r.Header.Get("Accept-Language")
	return SimplifiedLocalizer{i18n.NewLocalizer(bundle, lang, accept)}
}
