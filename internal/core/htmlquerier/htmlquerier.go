// Package htmlquerier provides a standardized declarative api for creating a querier
// that is able to fetch a string from a selector, apply filters to it, and send the result on for processing in the core.
package htmlquerier

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"golang.org/x/text/width"
)

// Querier contains details telling livefetcher where to fetch a piece of information,
// and how to process the text fetched to get it ready for use in livefetcher.
//
// It is recommended to not initialize struct directly.
// Instead, use Q or QAll.
type Querier struct {
	// Initialized specifies whether the querier has been initialized or not.
	Initialized bool
	// endSelector specifies a selector to stop before.
	// Any text in or after endSelector will not be included in results.
	endSelector string
	// arr denotes a slice of strings being returned. If only using basic filters and .Q() this will be a slice with only one string.
	arr []string
	// selector denotes the primary selector to fetch the initial string(s) from.
	selector string
	// selectAll determines fetch strategy. If selectAll is true, all matches are fetched and sent to arr; if selectAll is false, only first match is fetched and added to first index of array.
	selectAll bool
	// filters contains the filters to apply to modify result.
	filters []func([]string) []string
}

// QAll creates a pointer to a Querier struct.
// A Querier struct initialized using QAll will fetch all instances of the selector, get the string within, and assign them all to arr.
//
// Any basic filters specified will be applied individually on each match
func QAll(selector string) *Querier {
	return &Querier{
		selector:    selector,
		selectAll:   true,
		Initialized: true,
	}
}

// Q creates a pointer to a Querier struct.
// A Querier struct initialized using Q will only select the first match and get the string from that.
func Q(selector string) *Querier {
	return &Querier{
		selector:    selector,
		Initialized: true,
	}
}

// trim is a splitter function that trims whitespace from the beginning and end of the string.
func trim(s string) string {
	return strings.TrimSpace(s)
}

// Trim adds a filter to the querier that removes any leading and trailing whitespace.
func (q *Querier) Trim() *Querier {
	return q.AddFilter(func(s string) string {
		return trim(s)
	})
}

// TrimPrefix adds a filter to the querier that removes a specific prefix from the string.
func (q *Querier) TrimPrefix(prefix string) *Querier {
	return q.AddFilter(func(s string) string {
		return strings.TrimPrefix(s, prefix)
	})
}

// TrimSuffix adds a filter to the querier that removes a specific suffix from the string.
func (q *Querier) TrimSuffix(suffix string) *Querier {
	return q.AddFilter(func(s string) string {
		return strings.TrimSuffix(s, suffix)
	})
}

// CutWrapper adds a filter to the querier that removes a wrapping prefix and suffix only if both are present.
func (q *Querier) CutWrapper(prefix, suffix string) *Querier {
	return q.AddFilter(func(s string) string {
		if strings.HasPrefix(s, prefix) && strings.HasSuffix(s, suffix) && len(s) >= len(suffix)+len(prefix) {
			s = s[len(prefix) : len(s)-len(suffix)]
		}
		return s
	})
}

// BeforeSelector sets an endSelector, and will ensure that only text before the selector specified is selected.
func (q *Querier) BeforeSelector(selector string) *Querier {
	q.endSelector = selector
	return q
}

// Execute executes the query. This is only used internally in the core, please do not call this in connectors.
func (q *Querier) Execute(n *html.Node) (a []string, err error) {
	if n == nil {
		a = []string{""}
		err = fmt.Errorf("node is nil for selector %s", q.selector)
		return
	}

	if q.selectAll {
		var res []*html.Node
		res, err = htmlquery.QueryAll(n, q.selector)
		if err != nil || res == nil || len(res) == 0 {
			a = []string{""}
			return
		}
		strs := make([]string, 0)
		for _, artistNode := range res {
			strs = append(strs, htmlquery.InnerText(artistNode))
		}
		q.arr = strs
	} else {
		var res *html.Node
		res, err = htmlquery.Query(n, q.selector)
		if err != nil || res == nil {
			a = []string{""}
			return
		}
		q.arr = []string{htmlquery.InnerText(res)}

		if q.endSelector != "" {
			var end *html.Node
			end, err = htmlquery.Query(res, q.endSelector)
			if err == nil && end != nil {
				for i := range q.arr {
					q.arr[i] = strings.Split(q.arr[i], htmlquery.InnerText(end))[0]
				}
			}
			err = nil
		}
	}

	for _, filter := range q.filters {
		q.arr = filter(q.arr)
	}

	newArr := make([]string, 0)
	for _, str := range q.arr {
		newStr := trim(str)
		if newStr != "" {
			newArr = append(newArr, newStr)
		}
	}
	q.arr = newArr

	a = q.arr
	if len(a) == 0 {
		a = []string{""}
	}
	return
}

// AddFilter adds a simple filter to the Querier struct.
// Simple filter will run once on each entry in slice, replacing each entry with the filtered version.
func (q *Querier) AddFilter(fn func(string) string) *Querier {
	q.AddComplexFilter(func(arr []string) []string {
		for i := range arr {
			arr[i] = fn(arr[i])
		}
		return arr
	})
	return q
}

// AddComplexFilter adds a filter that takes the full slice of strings, and returns a new slice.
// This should only be used if you need the full context of the array, or if you want to be able to entirely remove entries.
//
// Make sure not to return an empty slice, at minimum return slice containing a single entry with empty string.
func (q *Querier) AddComplexFilter(fn func([]string) []string) *Querier {
	q.filters = append(q.filters, fn)
	return q
}

// Split adds a splitter that splits on a given separator string sep.
func (q *Querier) Split(sep string) *Querier {
	return q.AddSplitter(func(s string) []string {
		return strings.Split(s, sep)
	})
}

// SplitIgnoreWithin adds a splitter that splits using a given separator, while ignoring that separator if its within a set of left and right brackets.
//
// For instance, often, the slash character "/" is used as a separator between artists on websites.
// However, slash may also appear often in parentheses on individual artists to denote things like features etc.
//
// In this case, we can use SplitIgnoreWithin to separate on "/", while ensuring splitting does not occur within the parentheses used by the site.
func (q *Querier) SplitIgnoreWithin(sep string, l, r rune) *Querier {
	return q.AddSplitter(func(s string) (arr []string) {
		re, err := regexp.Compile(fmt.Sprintf("[%s].*?[%s]|(%s)", string(l), string(r), sep))
		if err != nil {
			arr = []string{s}
			return
		}

		depth := 0
		prev := 0
		res := []string{""}
		for i, c := range s {
			if c == l {
				if depth == 0 {
					strs := re.Split(s[prev:i], -1)
					prev = i
					for i, str := range strs {
						if i == 0 {
							res[len(res)-1] += str
						} else {
							res = append(res, str)
						}
					}
				}
				depth++
			}
			if c == r {
				depth--
				if depth == 0 {
					res[len(res)-1] += s[prev : i+len(string(r))]
					prev = i + len(string(r))
				}
			}
		}
		if prev < len(s) {
			if depth == 0 {
				strs := re.Split(s[prev:], -1)
				for i, str := range strs {
					if i == 0 {
						res[len(res)-1] += str
					} else {
						res = append(res, str)
					}
				}
			} else {
				res[len(res)-1] += s[prev:]
			}
		}
		return res
	})
}

// SplitRegex adds a splitter that splits using a regular expression exp.
func (q *Querier) SplitRegex(exp string) *Querier {
	return q.AddSplitter(func(s string) []string {
		re, err := regexp.Compile(exp)
		if err != nil {
			return []string{s}
		}
		return re.Split(s, -1)
	})
}

// SplitIndex adds a splitter that splits using a string, but only returns the entry at index i, or empty string if index i doesnt exist.
func (q *Querier) SplitIndex(sep string, i int) *Querier {
	return q.AddSplitter(func(s string) []string {
		arr := strings.Split(s, sep)
		if i < len(arr) {
			return []string{arr[i]}
		}
		return []string{""}
	})
}

// SplitRegexIndex works like SplitIndex, except using regex.
func (q *Querier) SplitRegexIndex(exp string, i int) *Querier {
	return q.AddSplitter(func(s string) []string {
		re, err := regexp.Compile(exp)
		if err != nil {
			return []string{s}
		}
		arr := re.Split(s, -1)
		if i < len(arr) {
			return []string{arr[i]}
		}
		return []string{""}
	})
}

// AddSplitter adds a splitter filter for the Querier, which iterates over the slice, and may or may not turn the entry into multiple entries.
func (q *Querier) AddSplitter(fn func(string) []string) *Querier {
	q.AddComplexFilter(func(old []string) []string {
		newArr := make([]string, 0)
		for _, s := range old {
			newArr = append(newArr, fn(s)...)
		}
		return newArr
	})
	return q
}

// After adds a filter that removes any text before and including the first instance of given separator sep.
func (q *Querier) After(sep string) *Querier {
	return q.AddFilter(func(s string) string {
		arr := strings.SplitN(s, sep, 2)
		if len(arr) < 2 {
			return s
		}
		return arr[1]
	})
}

// Before adds a filter that removes any text after and including the first instance of given separator sep.
func (q *Querier) Before(sep string) *Querier {
	return q.AddFilter(func(s string) string {
		arr := strings.SplitN(s, sep, 2)
		if len(arr) == 0 {
			return s
		}
		return arr[0]
	})
}

// HalfWidth adds a filter that forces fullwidth alphanumeric characters to halfwidth characters.
// This is typically useful for sites that use fullwidth numbers for dates.
func (q *Querier) HalfWidth() *Querier {
	return q.AddFilter(func(s string) string {
		return width.Narrow.String(s)
	})
}

// ReplaceAll adds a filter that replaces all instances of a string old with string new.
func (q *Querier) ReplaceAll(old, new string) *Querier {
	return q.AddFilter(func(s string) string {
		return strings.ReplaceAll(s, old, new)
	})
}

// ReplaceAll adds a filter that replaces all instances of a regular expression exp with string new.
// ReplaceAll uses regexp.ReplaceAllString under the hood, so use $1, $2, etc for groups.
func (q *Querier) ReplaceAllRegex(exp, new string) *Querier {
	return q.AddFilter(func(s string) string {
		re, err := regexp.Compile(exp)
		if err != nil {
			return s
		}
		return re.ReplaceAllString(s, new)
	})
}

// Prefix adds a filter that adds a prefix p in front of string.
func (q *Querier) Prefix(p string) *Querier {
	return q.AddFilter(func(s string) string {
		return fmt.Sprintf(p + s)
	})
}

// DeleteFrom adds a complex filter that deletes every item starting at an item with specific value
func (q *Querier) DeleteFrom(s string) *Querier {
	return q.AddComplexFilter(func(old []string) []string {
		new := make([]string, 0)
		for _, cur := range old {
			if s == cur {
				break
			}
			new = append(new, cur)
		}
		return new
	})
}

// DeleteUntil adds a complex filter that deletes every item until and including an item with specific value
func (q *Querier) DeleteUntil(s string) *Querier {
	return q.AddComplexFilter(func(old []string) []string {
		new := make([]string, 0)
		shouldDelete := true
		for _, cur := range old {
			if !shouldDelete {
				new = append(new, cur)
				continue
			}
			if s == cur {
				shouldDelete = false
			}
		}
		return new
	})
}

func stringHasTitleIndicator(s string) bool {
	indicators := []string{
		"album",
		"presents",
		"live",
		"one man",
		"oneman",
		"tour",
		"ツアー",
		"ワンマン",
		"ツーマン",
	}
	lower := strings.ToLower(s)
	for _, indicator := range indicators {
		if strings.Contains(lower, indicator) {
			return true
		}
	}
	return false
}

func getArtistIndex(a []string, exp string, i int) int {
	if stringHasTitleIndicator(a[0]) {
		return 1
	}
	if stringHasTitleIndicator(a[1]) {
		return 0
	}

	re, err := regexp.Compile(exp)
	if err != nil {
		return i
	}
	if re.MatchString(a[0]) {
		return 0
	}
	if re.MatchString(a[1]) {
		return 1
	}
	if strings.Contains(a[0], strings.TrimSpace(a[1])) {
		return 1
	}
	if strings.Contains(a[1], strings.TrimSpace(a[0])) {
		return 0
	}
	return i
}

func getTitle(a []string, exp string, i int) string {
	if len(a) == 0 {
		return ""
	}
	if len(a) == 1 {
		return a[0]
	}
	return a[1-getArtistIndex(a, exp, i)]
}

func getArtist(a []string, exp string, i int) string {
	if len(a) == 0 {
		return ""
	}
	if len(a) == 1 {
		return a[0]
	}
	return a[getArtistIndex(a, exp, i)]
}

// FilterTitle is meant to be run on a querier that has fetched title and artist, without knowing which.
// It will then try to return only the title to the best of its ability.
//
// exp is expected separator regex for artists
//
// i is the most common index for title to have (fallback)
func (q *Querier) FilterTitle(exp string, i int) *Querier {
	return q.AddComplexFilter(func(old []string) []string {
		return []string{getTitle(old, exp, 1-i)}
	})
}

// FilterArtist is meant to be run on a querier that has fetched title and artist, without knowing which.
// It will then try to return only the artist to the best of its ability.
//
// exp is expected separator regex (FilterArtist will NOT split, you must do that separately after)
//
// i is the most common index for artist to have (fallback)
func (q *Querier) FilterArtist(exp string, i int) *Querier {
	return q.AddComplexFilter(func(old []string) []string {
		return []string{getArtist(old, exp, i)}
	})
}

// KeepIndex keeps only the element at specific index, or empty string if does not exist. Negative index will get index starting from last index.
func (q *Querier) KeepIndex(i int) *Querier {
	return q.AddComplexFilter(func(old []string) []string {
		if i < 0 {
			if len(old) >= -i {
				return []string{old[len(old)+i]}
			}
			return []string{""}
		}
		if len(old) > i {
			return []string{old[i]}
		}
		return []string{""}
	})
}

// Concat concatenates all the strings from the slice to one using a separator sep
func (q *Querier) Join(sep string) *Querier {
	return q.AddComplexFilter(func(old []string) []string {
		return []string{strings.Join(old, sep)}
	})
}
