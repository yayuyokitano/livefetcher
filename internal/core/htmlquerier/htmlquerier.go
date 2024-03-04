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
	// str denotes the string return in its current state. This may be modified by filters as the query executes.
	str string
	// arr denotes an array of strings being returned. This will be created through a QAll query or a splitter function.
	// Typically used for fetching and returning multiple artists in a single live.
	arr []string
	// selector denotes the primary selector to fetch the initial string(s) from.
	selector string
	// selectAll denotes whether only string from first match should be assigned to str, or if string from all matches should be assigned to arr.
	selectAll bool
	// presplitmod denotes a filter that has been implemented before any splitter (if there is a splitter).
	// These filters will be executed in order on str, before splitter is executed.
	presplitmod []func(string) string
	// splitter is a splitter function, that takes an initial single input string (stored in str), and splits it into an array of strings (stored in arr)
	splitter func(string) []string
	// postsplitmod denotes a filter that has been implemented after any splitter (or if QAll is used, all filters are postsplitmod)
	// These filters will be applied after splitter, and will be applied separately to every entry in the string array arr.
	postsplitmod []func(string) string
	// sliceFilter is applied immediately after a splitter, and will filter the slice itself
	// This is only needed when you need context of other items in slice.
	// If you do not need context of other items in slice, just use postsplitmod.
	sliceFilter []func([]string) []string
}

// QAll creates a pointer to a Querier struct.
// A Querier struct initialized using QAll will fetch all instances of the selector, get the string within, and assign them all to arr.
//
// Any filters specified on QAll Querier will behave as if they are being executed after a splitter.
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

// dontSplit is a placeholder splitter function that does nothing, used when no splitter was used.
func dontSplit(s string) []string {
	return []string{s}
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
	if q.splitter == nil {
		q.SetSplitter(dontSplit)
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
			//fmt.Println(htmlquery.InnerText(artistNode))
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
		q.str = htmlquery.InnerText(res)

		if q.endSelector != "" {
			var end *html.Node
			end, err = htmlquery.Query(res, q.endSelector)
			if err == nil && end != nil {
				q.str = strings.Split(q.str, htmlquery.InnerText(end))[0]
			}
			err = nil
		}

		for _, mod := range q.presplitmod {
			q.str = mod(q.str)
		}
		q.arr = q.splitter(q.str)
	}

	for _, filter := range q.sliceFilter {
		q.arr = filter(q.arr)
	}

	for _, mod := range q.postsplitmod {
		newArr := make([]string, 0)
		for _, str := range q.arr {
			newStr := mod(str)
			if newStr != "" {
				newArr = append(newArr, newStr)
			}
		}
		q.arr = newArr
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

// AddFilter adds a filter to the Querier struct.
//
// If a splitter has been specified or Querier was made using QAll, the filter will be individually applied to every string in the string array created by splitter.
//
// If a splitter has not yet been specified, filter will be applied to the single string.
//
// The order that filters and splitter are added to the Querier matters.
func (q *Querier) AddFilter(fn func(string) string) *Querier {
	if q.splitter != nil || q.selectAll {
		q.postsplitmod = append(q.postsplitmod, fn)
	} else {
		q.presplitmod = append(q.presplitmod, fn)
	}
	return q
}

func (q *Querier) AddSliceFilter(fn func([]string) []string) *Querier {
	q.sliceFilter = append(q.sliceFilter, fn)
	return q
}

// Split sets a splitter that splits on a given separator string sep.
func (q *Querier) Split(sep string) *Querier {
	return q.SetSplitter(func(s string) []string {
		return strings.Split(s, sep)
	})
}

// SplitIgnoreWithin sets a splitter that splits using a given separator, while ignoring that separator if its within a set of left and right brackets.
//
// For instance, often, the slash character "/" is used as a separator between artists on websites.
// However, slash may also appear often in parentheses on individual artists to denote things like features etc.
//
// In this case, we can use SplitIgnoreWithin to separate on "/", while ensuring splitting does not occur within the parentheses used by the site.
func (q *Querier) SplitIgnoreWithin(sep string, l, r rune) *Querier {
	return q.SetSplitter(func(s string) (arr []string) {
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

// SplitRegex sets a splitter that splits using a regular expression exp.
func (q *Querier) SplitRegex(exp string) *Querier {
	return q.SetSplitter(func(s string) []string {
		re, err := regexp.Compile(exp)
		if err != nil {
			return []string{s}
		}
		return re.Split(s, -1)
	})
}

// SplitIndex sets a splitter that splits using a string, but only returns the entry at index i, or empty string if index i doesnt exist.
func (q *Querier) SplitIndex(sep string, i int) *Querier {
	return q.SetSplitter(func(s string) []string {
		arr := strings.Split(s, sep)
		if i < len(arr) {
			return []string{arr[i]}
		}
		return []string{""}
	})
}

// SplitRegexIndex works like SplitIndex, except using regex.
func (q *Querier) SplitRegexIndex(exp string, i int) *Querier {
	return q.SetSplitter(func(s string) []string {
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

// SetSplitter sets a splitter function for the Querier, splitting the string into a slice of strings using the function given.
//
// Any filter added after SetSpliter is called will be executed on each individual entries of this slice.
func (q *Querier) SetSplitter(fn func(string) []string) *Querier {
	if !q.selectAll {
		q.splitter = fn
	}
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

// DeleteFrom adds a slice filter that deletes every item starting at an item with specific value
func (q *Querier) DeleteFrom(s string) *Querier {
	return q.AddSliceFilter(func(old []string) []string {
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

// DeleteUntil adds a slice filter that deletes every item until and including an item with specific value
func (q *Querier) DeleteUntil(s string) *Querier {
	return q.AddSliceFilter(func(old []string) []string {
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

// KeepIndex keeps only the element at specific index, or empty string if does not exist
func (q *Querier) KeepIndex(i int) *Querier {
	return q.AddSliceFilter(func(old []string) []string {
		if len(old) > i {
			return []string{old[i]}
		}
		return []string{""}
	})
}
