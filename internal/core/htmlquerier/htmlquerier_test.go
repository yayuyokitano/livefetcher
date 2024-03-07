package htmlquerier_test

import (
	"testing"

	"github.com/antchfx/htmlquery"
	"github.com/yayuyokitano/livefetcher/internal/core/htmlquerier"
	"golang.org/x/net/html"
)

func createBaseQuerier() (n *html.Node, err error) {
	doc, err := htmlquery.LoadDoc("./test.html")
	if err != nil {
		return
	}
	n, err = htmlquery.Query(doc, "//body")
	return
}

func createQuerier(t *testing.T, s string) (q htmlquerier.Querier, n *html.Node) {
	n, err := createBaseQuerier()
	if err != nil {
		t.Error(err)
		return
	}
	q = *htmlquerier.Q(s)
	return
}

func createQuerierAll(t *testing.T, s string) (q htmlquerier.Querier, n *html.Node) {
	n, err := createBaseQuerier()
	if err != nil {
		t.Error(err)
		return
	}
	q = *htmlquerier.QAll(s)
	return
}

func TestQAll(t *testing.T) {
	q, n := createQuerierAll(t, "//p[@id='complex']/text()")
	arr, err := q.Split("-").Trim().Execute(n)
	if err != nil {
		t.Error(err)
	}
	testStringSliceEquals(t, []string{"one", "two", "three", "four"}, arr)
}

func TestTrim(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='splitter']")
	arr, err := q.Split("-").Trim().Execute(n)
	if err != nil {
		t.Error(err)
	}
	testStringSliceEquals(t, []string{"one", "two", "three"}, arr)
}

func TestTrimSuffix(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='splitter']")
	arr, err := q.Split(" - ").TrimSuffix("ee").Execute(n)
	if err != nil {
		t.Error(err)
	}
	testStringSliceEquals(t, []string{"one", "two", "thr"}, arr)
}

func TestTrimPrefix(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='splitter']")
	arr, err := q.Split(" - ").TrimPrefix("thr").Execute(n)
	if err != nil {
		t.Error(err)
	}
	testStringSliceEquals(t, []string{"one", "two", "ee"}, arr)
}

func TestCutWrapper(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='wrapper']")
	arr, err := q.CutWrapper("「", "」").Execute(n)
	if err != nil {
		t.Error(err)
	}
	testStringSliceEquals(t, []string{"one - two"}, arr)

	q2, n := createQuerier(t, "//p[@id='wrapperfail']")
	arr2, err2 := q2.CutWrapper("「", "」").Execute(n)
	if err2 != nil {
		t.Error(err2)
	}
	testStringSliceEquals(t, []string{"one 「two」"}, arr2)
}

func TestBeforeSelector(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='complex']")
	arr, err := q.BeforeSelector("//span").Execute(n)
	if err != nil {
		t.Error(err)
	}
	testStringSliceEquals(t, []string{"onetwo"}, arr)
}

func TestAfter(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='splitter']")
	arr, err := q.After(" - ").Execute(n)
	if err != nil {
		t.Error(err)
	}
	testStringSliceEquals(t, []string{"two - three"}, arr)

	q2, n := createQuerier(t, "//p[@id='splitter']")
	arr2, err2 := q2.After("hehe").Execute(n)
	if err2 != nil {
		t.Error(err2)
	}
	testStringSliceEquals(t, []string{"one - two - three"}, arr2)
}

func TestPrefix(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='splitter']")
	arr, err := q.Prefix("zero - ").Execute(n)
	if err != nil {
		t.Error(err)
	}
	testStringSliceEquals(t, []string{"zero - one - two - three"}, arr)
}

func TestBefore(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='splitter']")
	arr, err := q.Before(" - ").Execute(n)
	if err != nil {
		t.Error(err)
	}
	testStringSliceEquals(t, []string{"one"}, arr)

	q2, n := createQuerier(t, "//p[@id='splitter']")
	arr2, err2 := q2.Before("hehe").Execute(n)
	if err2 != nil {
		t.Error(err2)
	}
	testStringSliceEquals(t, []string{"one - two - three"}, arr2)
}

func TestHalfWidth(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='fullwidth']")
	arr, err := q.HalfWidth().Execute(n)
	if err != nil {
		t.Error(err)
	}
	testStringSliceEquals(t, []string{"ONE23"}, arr)
}

func TestReplaceAll(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='multisplit']")
	arr, err := q.ReplaceAll("-", "").Execute(n)
	if err != nil {
		t.Error(err)
	}
	testStringSliceEquals(t, []string{"onetwo/threefour"}, arr)
}

func TestReplaceAllRegex(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='multisplit']")
	arr, err := q.ReplaceAllRegex("[-/]", "").Execute(n)
	if err != nil {
		t.Error(err)
	}
	testStringSliceEquals(t, []string{"onetwothreefour"}, arr)
}

func TestSplit(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='splitter']")
	arr, err := q.Split(" - ").Execute(n)
	if err != nil {
		t.Error(err)
	}
	testStringSliceEquals(t, []string{"one", "two", "three"}, arr)
}

func TestAfterSplit(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='multisplit']")
	arr, err := q.Split("/").After("-").Execute(n)
	if err != nil {
		t.Error(err)
	}
	testStringSliceEquals(t, []string{"two", "four"}, arr)
}

func TestSplitIgnoreWithin(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='splitignorewithin']")
	arr, err := q.SplitIgnoreWithin(" / ", '（', '）').Execute(n)
	if err != nil {
		t.Error(err)
	}

	testStringSliceEquals(t, []string{"one", "two", "three（four / five （six / seven） eight / nine）", "ten"}, arr)
}

func TestSplitRegex(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='multisplit']")
	arr, err := q.SplitRegex("[/-]").Execute(n)
	if err != nil {
		t.Error(err)
	}
	testStringSliceEquals(t, []string{"one", "two", "three", "four"}, arr)
}

func TestSplitIndex(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='splitter']")
	arr, err := q.SplitIndex(" - ", 1).Execute(n)
	if err != nil {
		t.Error(err)
	}
	testStringSliceEquals(t, []string{"two"}, arr)

	q2, n := createQuerier(t, "//p[@id='splitter']")
	arr2, err2 := q2.SplitIndex(" - ", 3).Execute(n)
	if err2 != nil {
		t.Error(err2)
	}
	testStringSliceEquals(t, []string{""}, arr2)
}

func TestSplitRegexIndex(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='multisplit']")
	arr, err := q.SplitRegexIndex("[/-]", 2).Execute(n)
	if err != nil {
		t.Error(err)
	}
	testStringSliceEquals(t, []string{"three"}, arr)

	q2, n := createQuerier(t, "//p[@id='multisplit']")
	arr2, err2 := q2.SplitRegexIndex("[/-]", 4).Execute(n)
	if err2 != nil {
		t.Error(err2)
	}
	testStringSliceEquals(t, []string{""}, arr2)
}

func TestDeleteFrom(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='splitterlong']")
	arr, err := q.Split(" - ").DeleteFrom("three").Execute(n)
	if err != nil {
		t.Error(err)
	}
	testStringSliceEquals(t, []string{"one", "two"}, arr)

	q2, n := createQuerier(t, "//p[@id='splitterlong']")
	arr2, err2 := q2.Split(" - ").DeleteFrom("one").Execute(n)
	if err2 != nil {
		t.Error(err2)
	}
	testStringSliceEquals(t, []string{""}, arr2)
}

func TestDeleteUntil(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='splitterlong']")
	arr, err := q.Split(" - ").DeleteUntil("two").Execute(n)
	if err != nil {
		t.Error(err)
	}
	testStringSliceEquals(t, []string{"three", "four"}, arr)

	q2, n := createQuerier(t, "//p[@id='splitterlong']")
	arr2, err2 := q2.Split(" - ").DeleteUntil("four").Execute(n)
	if err2 != nil {
		t.Error(err2)
	}
	testStringSliceEquals(t, []string{""}, arr2)
}

func TestKeepIndex(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='splitter']")
	arr, err := q.Split(" - ").KeepIndex(1).Execute(n)
	if err != nil {
		t.Error(err)
	}
	testStringSliceEquals(t, []string{"two"}, arr)

	q2, n := createQuerier(t, "//p[@id='splitter']")
	arr2, err2 := q2.Split(" - ").KeepIndex(3).Execute(n)
	if err2 != nil {
		t.Error(err2)
	}
	testStringSliceEquals(t, []string{""}, arr2)
}

func TestJoin(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='splitter']")
	arr, err := q.Split(" - ").Join("|").Execute(n)
	if err != nil {
		t.Error(err)
	}
	testStringSliceEquals(t, []string{"one|two|three"}, arr)
}

func testStringSliceEquals(t *testing.T, expected, res []string) {
	t.Helper()
	if len(expected) != len(res) {
		t.Errorf("Expected length of %v (%d) to be equal to length of %v (%d)", res, len(res), expected, len(expected))
	}
	for i := 0; i < len(expected); i++ {
		if expected[i] != res[i] {
			t.Errorf("Expected %v to equal %v: on index %d expected %s, got %s", res, expected, i, expected[i], res[i])
		}
	}
}
