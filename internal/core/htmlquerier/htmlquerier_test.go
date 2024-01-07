package htmlquerier_test

import (
	"testing"

	"github.com/antchfx/htmlquery"
	"github.com/yayuyokitano/livefetcher/internal/core/htmlquerier"
	"golang.org/x/net/html"
)

func createQuerier(t *testing.T, s string) (q htmlquerier.Querier, n *html.Node) {
	t.Helper()
	doc, err := htmlquery.LoadDoc("./test.html")
	if err != nil {
		t.Error(err)
		return
	}
	n, err = htmlquery.Query(doc, "//body")
	if err != nil {
		t.Error(err)
		return
	}
	q = *htmlquerier.Q(s)
	return
}

func TestAfter(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='splitter']")
	arr, err := q.After(" - ").Execute(n)
	if err != nil {
		t.Error(err)
	}
	testSliceLength(t, 1, arr)
	testStringEquals(t, "two - three", arr[0])
}

func TestAfterEmpty(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='empty']")
	arr, err := q.After(" - ").Execute(n)
	if err != nil {
		t.Error(err)
	}
	testSliceLength(t, 1, arr)
	testStringEquals(t, "", arr[0])
}

func TestPrefix(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='splitter']")
	arr, err := q.Prefix("zero - ").Execute(n)
	if err != nil {
		t.Error(err)
	}
	testSliceLength(t, 1, arr)
	testStringEquals(t, "zero - one - two - three", arr[0])
}

func TestBefore(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='splitter']")
	arr, err := q.Before(" - ").Execute(n)
	if err != nil {
		t.Error(err)
	}
	testSliceLength(t, 1, arr)
	testStringEquals(t, "one", arr[0])
}

func TestBeforeEmpty(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='empty']")
	arr, err := q.After(" - ").Execute(n)
	if err != nil {
		t.Error(err)
	}
	testSliceLength(t, 1, arr)
	testStringEquals(t, "", arr[0])
}

func TestReplaceAll(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='multisplit']")
	arr, err := q.ReplaceAll("-", "").Execute(n)
	if err != nil {
		t.Error(err)
	}
	testSliceLength(t, 1, arr)
	testStringEquals(t, "onetwo/threefour", arr[0])
}

func TestSplit(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='splitter']")
	arr, err := q.Split(" - ").Execute(n)
	if err != nil {
		t.Error(err)
	}
	testSliceLength(t, 3, arr)
	testStringEquals(t, "one", arr[0])
	testStringEquals(t, "two", arr[1])
	testStringEquals(t, "three", arr[2])
}

func TestAfterSplit(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='multisplit']")
	arr, err := q.Split("/").After("-").Execute(n)
	if err != nil {
		t.Error(err)
	}
	testSliceLength(t, 2, arr)
	testStringEquals(t, "two", arr[0])
	testStringEquals(t, "four", arr[1])
}

func TestSplitIgnoreWithin(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='splitignorewithin']")
	arr, err := q.SplitIgnoreWithin(" / ", '（', '）').Execute(n)
	if err != nil {
		t.Error(err)
	}
	testSliceLength(t, 3, arr)
	testStringEquals(t, "one", arr[0])
	testStringEquals(t, "two", arr[1])
	testStringEquals(t, "three（four / five）", arr[2])
}

func TestSplitRegex(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='multisplit']")
	arr, err := q.SplitRegex("[/-]").Execute(n)
	if err != nil {
		t.Error(err)
	}
	testSliceLength(t, 4, arr)
	testStringEquals(t, "one", arr[0])
	testStringEquals(t, "two", arr[1])
	testStringEquals(t, "three", arr[2])
	testStringEquals(t, "four", arr[3])
}

func TestSplitIndex(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='splitter']")
	arr, err := q.SplitIndex(" - ", 1).Execute(n)
	if err != nil {
		t.Error(err)
	}
	testSliceLength(t, 1, arr)
	testStringEquals(t, "two", arr[0])
}

func TestSplitRegexIndex(t *testing.T) {
	q, n := createQuerier(t, "//p[@id='multisplit']")
	arr, err := q.SplitRegexIndex("[/-]", 2).Execute(n)
	if err != nil {
		t.Error(err)
	}
	testSliceLength(t, 1, arr)
	testStringEquals(t, "three", arr[0])
}

func testSliceLength[T any](t *testing.T, expectedLen int, arr []T) {
	t.Helper()
	if len(arr) != expectedLen {
		t.Errorf("Expected slice to have length %d, was length %d", expectedLen, len(arr))
	}
}

func testStringEquals(t *testing.T, expected, res string) {
	t.Helper()
	if expected != res {
		t.Errorf("Expected result to be %s, was %s", expected, res)
	}
}
