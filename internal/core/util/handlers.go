package util

import (
	"fmt"
	"strings"
	"time"
)

type TimeHandler func(d string, t string) (open int64, start int64, err error)

type PriceHandler func(p string) (price string, err error)

var timeLayout = "2006-01-02 15:04:05 -0700"

func ParseTime(d string, t string) (res time.Time, err error) {
	th, tm, nextDay := GetTimeFromString(t)
	res, err = time.Parse(timeLayout, fmt.Sprintf("%s %s:%s:00 +0900", d, th, tm))
	if err != nil {
		return
	}
	if nextDay {
		res = res.AddDate(0, 0, 1)
	}
	return
}

func EnglishPriceHandler(p string) (price string) {
	price = strings.ReplaceAll(p, "前売り", "Reservation")
	price = strings.ReplaceAll(price, "前売", "Reservation")
	price = strings.ReplaceAll(price, "スタンディング", "Standing")
	price = strings.ReplaceAll(price, "税込", "Incl. Tax")
	price = strings.ReplaceAll(price, "税抜", "Excl. Tax")
	price = strings.ReplaceAll(price, "当日", "Door")
	price = strings.ReplaceAll(price, "一般前売", "Ordinary Reservation")
	price = strings.ReplaceAll(price, "予約", "Reservation")
	price = strings.ReplaceAll(price, "事前", "Reservation")
	price = strings.ReplaceAll(price, "ドリンク別", "Drinks sold separately")
	price = strings.ReplaceAll(price, "ドリンク", "Drink")
	price = strings.ReplaceAll(price, "Sチケット", "S-Ticket")
	price = strings.ReplaceAll(price, "高校生以下", "High School Students and Below")
	price = strings.ReplaceAll(price, "高校生", "High School Students")
	price = strings.ReplaceAll(price, "大学生・専門学生", "College Students")
	price = strings.ReplaceAll(price, "一般", "Ordinary Ticket")
	price = strings.ReplaceAll(price, "無料", "Free")
	price = strings.ReplaceAll(price, "入場", "Entry")
	price = strings.ReplaceAll(price, "イベント", "Event")
	price = strings.ReplaceAll(price, "チケット", " Ticket")
	price = strings.ReplaceAll(price, "学生", "Students")
	price = strings.ReplaceAll(price, "女性", "Women")
	price = strings.ReplaceAll(price, "男性", "Men")
	price = strings.ReplaceAll(price, "込み", "Included")
	price = strings.ReplaceAll(price, "無制限飲み放題", "Unlimited drinks")
	price = strings.ReplaceAll(price, "飲み放題", "All-you-can-drink")
	price = strings.ReplaceAll(price, "別途", "Separately")
	price = strings.ReplaceAll(price, "2D別", "2 Drink purchases required")
	price = strings.ReplaceAll(price, "1D別", "1 Drink purchase required")
	price = strings.ReplaceAll(price, "D別", "Drinks sold separately")
	price = strings.ReplaceAll(price, "別", "Separately")
	price = strings.ReplaceAll(price, "未定", "TBA")
	price = strings.ReplaceAll(price, "カメラ登録料", "Camera fee")
	price = strings.ReplaceAll(price, "前方エリア", "Front area")
	price = strings.ReplaceAll(price, "優先", "Priority entry")
	return
}
