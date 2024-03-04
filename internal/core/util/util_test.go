package util

import (
	"testing"
	"time"
)

func TestMax(t *testing.T) {
	if Max(1, 2) != 2 {
		t.Errorf("Max(1, 2) != 2, got %d", Max(1, 2))
	}
	if Max(2, 1) != 2 {
		t.Errorf("Max(2, 1) != 2, got %d", Max(2, 1))
	}
	if Max(1, 1) != 1 {
		t.Errorf("Max(1, 1) != 1, got %d", Max(1, 1))
	}
	if Max(-1, -2) != -1 {
		t.Errorf("Max(-1, -2) != -1, got %d", Max(-1, -2))
	}
}

func TestMin(t *testing.T) {
	if Min(1, 2) != 1 {
		t.Errorf("Min(1, 2) != 1, got %d", Min(1, 2))
	}
	if Min(2, 1) != 1 {
		t.Errorf("Min(2, 1) != 1, got %d", Min(2, 1))
	}
	if Min(1, 1) != 1 {
		t.Errorf("Min(1, 1) != 1, got %d", Min(1, 1))
	}
	if Min(-1, -2) != -2 {
		t.Errorf("Min(-1, -2) != -2, got %d", Min(-1, -2))
	}
}

func TestStripNonNumeric(t *testing.T) {
	if stripNonNumeric("a1b2c3") != "123" {
		t.Errorf("stripNonNumeric(\"a1b2c3\") != \"123\", got %s", stripNonNumeric("a1b2c3"))
	}
	if stripNonNumeric("a1b2c3d4e5f6g7h8i9j0") != "1234567890" {
		t.Errorf("stripNonNumeric(\"a1b2c3d4e5f6g7h8i9j0\") != \"1234567890\", got %s", stripNonNumeric("a1b2c3d4e5f6g7h8i9j0"))
	}
	if stripNonNumeric("1234567890") != "1234567890" {
		t.Errorf("stripNonNumeric(\"1234567890\") != \"1234567890\", got %s", stripNonNumeric("1234567890"))
	}
	if stripNonNumeric("abc") != "" {
		t.Errorf("stripNonNumeric(\"abc\") != \"\", got %s", stripNonNumeric("abc"))
	}
	if stripNonNumeric("") != "" {
		t.Errorf("stripNonNumeric(\"\") != \"\", got %s", stripNonNumeric(""))
	}
}

func TestGetTimeFromString(t *testing.T) {
	if hour, min, nextDay := GetTimeFromString("12:34"); hour != "12" || min != "34" || nextDay != false {
		t.Errorf("GetTimeFromString(\"12:34\") != \"12\", \"34\", false, got %s, %s, %t", hour, min, nextDay)
	}
	if hour, min, nextDay := GetTimeFromString("1:2"); hour != "01" || min != "02" || nextDay != false {
		t.Errorf("GetTimeFromString(\"1:2\") != \"01\", \"02\", false, got %s, %s, %t", hour, min, nextDay)
	}
	if hour, min, nextDay := GetTimeFromString("1:02"); hour != "01" || min != "02" || nextDay != false {
		t.Errorf("GetTimeFromString(\"1:02\") != \"01\", \"02\", false, got %s, %s, %t", hour, min, nextDay)
	}
	if hour, min, nextDay := GetTimeFromString("01:2"); hour != "01" || min != "02" || nextDay != false {
		t.Errorf("GetTimeFromString(\"01:2\") != \"01\", \"02\", false, got %s, %s, %t", hour, min, nextDay)
	}
	if hour, min, nextDay := GetTimeFromString("sdkl ndms  asdkjn dsfm dsm , sd 01:02 asd  adsf ads"); hour != "01" || min != "02" || nextDay != false {
		t.Errorf("GetTimeFromString(\"sdkl ndms  asdkjn dsfm dsm , sd 01:02 asd  adsf ads\") != \"01\", \"02\", false, got %s, %s, %t", hour, min, nextDay)
	}
	if hour, min, nextDay := GetTimeFromString("not a time"); hour != "03" || min != "24" || nextDay != false {
		t.Errorf("GetTimeFromString(\"not a time\") != \"03\", \"34\", false, got %s, %s, %t", hour, min, nextDay)
	}
	if hour, min, nextDay := GetTimeFromString("25:30"); hour != "01" || min != "30" || nextDay != true {
		t.Errorf("GetTimeFromString(\"25:30\") != \"01\", \"30\", true, got %s, %s, %t", hour, min, nextDay)
	}
}

func TestGetDate(t *testing.T) {
	if month, day, err := GetDate([]rune("01/02"), '/'); month != "01" || day != "02" || err != nil {
		t.Errorf("GetDate(\"01/02\", \"/\") != \"01\", \"02\", nil, got %s, %s, %s", month, day, err)
	}
	if month, day, err := GetDate([]rune("1/2"), '/'); month != "01" || day != "02" || err != nil {
		t.Errorf("GetDate(\"1/2\", \"/\") != \"01\", \"02\", nil, got %s, %s, %s", month, day, err)
	}
	if month, day, err := GetDate([]rune("1/02"), '/'); month != "01" || day != "02" || err != nil {
		t.Errorf("GetDate(\"1/02\", \"/\") != \"01\", \"02\", nil, got %s, %s, %s", month, day, err)
	}
	if month, day, err := GetDate([]rune("01/2"), '/'); month != "01" || day != "02" || err != nil {
		t.Errorf("GetDate(\"01/2\", \"/\") != \"01\", \"02\", nil, got %s, %s, %s", month, day, err)
	}
	if month, day, err := GetDate([]rune("1月2日"), '月'); month != "01" || day != "02" || err != nil {
		t.Errorf("GetDate(\"1月2日\", \"月\") != \"01\", \"02\", nil, got %s, %s, %s", month, day, err)
	}
	if month, day, err := GetDate([]rune("12"), '/'); err == nil {
		t.Errorf("GetDate(\"12\", \"/\") != \"\", \"\", error, got %s, %s, %s", month, day, err)
	}
}

func TestGetYearMonth(t *testing.T) {
	if year, month, err := GetYearMonth([]rune("2023年12月"), '年'); year != "2023" || month != "12" || err != nil {
		t.Errorf("GetYearMonth(\"2023年12月\", '年') != \"2023\", \"12\", nil, got %s, %s, %s", year, month, err)
	}
	if year, month, err := GetYearMonth([]rune("2023年2月"), '年'); year != "2023" || month != "02" || err != nil {
		t.Errorf("GetYearMonth(\"2023年2月\", '年') != \"2023\", \"02\", nil, got %s, %s, %s", year, month, err)
	}
	if year, month, err := GetYearMonth([]rune("kl asdf2023年2月 aosdof"), '年'); year != "2023" || month != "02" || err != nil {
		t.Errorf("GetYearMonth(\"kl asdf2023年2月 aosdof\", '年') != \"2023\", \"02\", nil, got %s, %s, %s", year, month, err)
	}
	if year, month, err := GetYearMonth([]rune("invalid"), '年'); err == nil {
		t.Errorf("GetYearMonth(\"invalid\", '月') != error, got %s, %s, %s", year, month, err)
	}
}

func TestGetYearMonthDay(t *testing.T) {
	if year, month, day, err := GetYearMonthDay([]rune("2023年12月31日"), '年', '月'); year != "2023" || month != "12" || day != "31" || err != nil {
		t.Errorf("GetYearMonthDay(\"2023年12月31日\", '年') != \"2023\", \"12\", \"31\", nil, got %s, %s, %s, %s", year, month, day, err)
	}
	if year, month, day, err := GetYearMonthDay([]rune("2023年2月3日"), '年', '月'); year != "2023" || month != "02" || day != "03" || err != nil {
		t.Errorf("GetYearMonthDay(\"2023年2月3日\", '年') != \"2023\", \"02\", \"03\", nil, got %s, %s, %s, %s", year, month, day, err)
	}
	if year, month, day, err := GetYearMonthDay([]rune("23年2月3日"), '年', '月'); year != "2023" || month != "02" || day != "03" || err != nil {
		t.Errorf("GetYearMonthDay(\"23年2月3日\", '年') != \"2023\", \"02\", \"03\", nil, got %s, %s, %s, %s", year, month, day, err)
	}
	if year, month, day, err := GetYearMonthDay([]rune("23.2.3 mon"), '.', '.'); year != "2023" || month != "02" || day != "03" || err != nil {
		t.Errorf("GetYearMonthDay(\"23.2.3 mon\", '.') != \"2023\", \"02\", \"03\", nil, got %s, %s, %s, %s", year, month, day, err)
	}
	if year, month, day, err := GetYearMonthDay([]rune("kl asdf2023年2月3日 aosdof"), '年', '月'); year != "2023" || month != "02" || day != "03" || err != nil {
		t.Errorf("GetYearMonthDay(\"kl asdf2023年2月3日 aosdof\", '年') != \"2023\", \"02\", \"03\", nil, got %s, %s, %s, %s", year, month, day, err)
	}
}

func TestSpacedPriceTimeFetcher(t *testing.T) {
	if price, open, start, err := SpacedPriceTimeFetcher("2023-06-01", "OPEN 17:30 START 18:00 ADV ¥2000 DOOR ¥2500 1D別"); price != "ADV ¥2000 DOOR ¥2500 1D別" || open.Unix() != 1685608200 || start.Unix() != 1685610000 || err != nil {
		t.Errorf("SpacedPriceTimeFetcher(\"2023-06-01\", \"OPEN 17:30 START 18:00 ADV ¥2000 DOOR ¥2500 1D別\") != \"ADV ¥2000 DOOR ¥2500 1D別\", 1685608200, 1685610000, nil, got %s, %d, %d, %s", price, open.Unix(), start.Unix(), err)
	}
	if price, open, start, err := SpacedPriceTimeFetcher("2023-06-01", "kjagdsn"); price != "" || open.Unix() != 1685557440 || start.Unix() != 1685557440 || err != nil {
		t.Errorf("SpacedPriceTimeFetcher(\"2023-06-01\", \"kjagdsn\") != \"\", 1685557440, 1685557440, nil, got %s, %d, %d, %s", price, open.Unix(), start.Unix(), err)
	}
	if price, open, start, err := SpacedPriceTimeFetcher("invalid date", "aksfls"); err == nil {
		t.Errorf("SpacedPriceTimeFetcher(\"invalid date\", \"aksfls\") != error, got %s, %d, %d, %s", price, open.Unix(), start.Unix(), err)
	}
}

func TestFindTime(t *testing.T) {
	str := `9 PARTY
	at LIVE HAUS SHIMOKITAZAWA
	11/9 (thu)
	open start 19:00
	当日券のみ¥2,300 +1drink
	
	哲学対話 : ホスト 永井玲衣
	
	LIVE :
	CHABE
	Nozomi Nobody
	ユッコ（MaCWORRY HILLBILLIES）
	スガナミユウ
	
	DJ :
	natsume
	ushi`

	testStringEquivalence("19:00", FindTime(str, "start"), t)
	testStringEquivalence("03:24", FindTime(str, "open"), t)
	testStringEquivalence("当日券のみ¥2,300", FindPrice(str), t)
}

func TestGetRelevantYear(t *testing.T) {
	testIntEquivalence(time.Now().Year(), GetRelevantYear(int(time.Now().Month())), t)
	testIntEquivalence(time.Now().Year()+1, GetRelevantYear(int(time.Now().Month())-1), t)
}

func testStringEquivalence(expected string, actual string, t *testing.T) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func testIntEquivalence(expected int, actual int, t *testing.T) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}
