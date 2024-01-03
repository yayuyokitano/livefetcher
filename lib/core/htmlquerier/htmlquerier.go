package htmlquerier

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"golang.org/x/text/width"
)

type Querier struct {
	Initialized  bool
	endSelector  string
	str          string
	arr          []string
	selector     string
	selectAll    bool
	presplitmod  []func(string) string
	splitter     func(string) []string
	postsplitmod []func(string) string
}

func QAll(selector string) *Querier {
	return &Querier{
		selector:    selector,
		selectAll:   true,
		Initialized: true,
	}
}

func Q(selector string) *Querier {
	return &Querier{
		selector:    selector,
		Initialized: true,
	}
}

func dontSplit(s string) []string {
	return []string{s}
}

func trim(s string) string {
	return strings.TrimSpace(s)
}

func (q *Querier) Trim() *Querier {
	return q.AddFilter(func(s string) string {
		return trim(s)
	})
}

func (q *Querier) TrimPrefix(prefix string) *Querier {
	return q.AddFilter(func(s string) string {
		return strings.TrimPrefix(s, prefix)
	})
}

func (q *Querier) TrimSuffix(suffix string) *Querier {
	return q.AddFilter(func(s string) string {
		return strings.TrimSuffix(s, suffix)
	})
}

func (q *Querier) BeforeSelector(selector string) *Querier {
	(*q).endSelector = selector
	return q
}

func (q *Querier) Execute(n *html.Node) (a []string, err error) {
	if n == nil {
		a = []string{""}
		err = fmt.Errorf("node is nil for selector %s", (*q).selector)
		return
	}
	if (*q).splitter == nil {
		(*q).SetSplitter(dontSplit)
	}
	(*q).postsplitmod = append((*q).postsplitmod, trim)

	if (*q).selectAll {
		var res []*html.Node
		res, err = htmlquery.QueryAll(n, (*q).selector)
		if err != nil || res == nil || len(res) == 0 {
			a = []string{""}
			return
		}
		strs := make([]string, 0)
		for _, artistNode := range res {
			//fmt.Println(htmlquery.InnerText(artistNode))
			strs = append(strs, htmlquery.InnerText(artistNode))
		}
		(*q).arr = strs
	} else {
		var res *html.Node
		res, err = htmlquery.Query(n, (*q).selector)
		if err != nil || res == nil {
			a = []string{""}
			return
		}
		(*q).str = htmlquery.InnerText(res)

		if (*q).endSelector != "" {
			var end *html.Node
			end, err = htmlquery.Query(res, (*q).endSelector)
			if err == nil && end != nil {
				(*q).str = strings.Split((*q).str, htmlquery.InnerText(end))[0]
			}
			err = nil
		}

		for _, mod := range (*q).presplitmod {
			(*q).str = mod((*q).str)
		}
		(*q).arr = (*q).splitter((*q).str)
	}

	for _, mod := range (*q).postsplitmod {
		newArr := make([]string, 0)
		for _, str := range (*q).arr {
			newStr := mod(str)
			if newStr != "" {
				newArr = append(newArr, newStr)
			}
		}
		(*q).arr = newArr
	}

	a = (*q).arr
	if len(a) == 0 {
		a = []string{""}
	}
	return
}

func (q *Querier) AddFilter(fn func(string) string) *Querier {
	if (*q).splitter != nil || (*q).selectAll {
		(*q).postsplitmod = append((*q).postsplitmod, fn)
	} else {
		(*q).presplitmod = append((*q).presplitmod, fn)
	}
	return q
}

func (q *Querier) Split(sep string) *Querier {
	return q.SetSplitter(func(s string) []string {
		return strings.Split(s, sep)
	})
}

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

func (q *Querier) SplitRegex(exp string) *Querier {
	return q.SetSplitter(func(s string) []string {
		re, err := regexp.Compile(exp)
		if err != nil {
			return []string{s}
		}
		return re.Split(s, -1)
	})
}

func (q *Querier) SplitIndex(sep string, i int) *Querier {
	return q.SetSplitter(func(s string) []string {
		arr := strings.Split(s, sep)
		if i < len(arr) {
			return []string{arr[i]}
		}
		return []string{""}
	})
}

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

func (q *Querier) SetSplitter(fn func(string) []string) *Querier {
	if !(*q).selectAll {
		(*q).splitter = fn
	}
	return q
}

func (q *Querier) After(sep string) *Querier {
	return q.AddFilter(func(s string) string {
		arr := strings.SplitN(s, sep, 2)
		if len(arr) < 2 {
			return s
		}
		return arr[1]
	})
}

func (q *Querier) Before(sep string) *Querier {
	return q.AddFilter(func(s string) string {
		arr := strings.SplitN(s, sep, 2)
		if len(arr) == 0 {
			return s
		}
		return arr[0]
	})
}

func (q *Querier) HalfWidth() *Querier {
	return q.AddFilter(func(s string) string {
		return width.Narrow.String(s)
	})
}

func (q *Querier) ReplaceAll(old, new string) *Querier {
	return q.AddFilter(func(s string) string {
		return strings.ReplaceAll(s, old, new)
	})
}

func (q *Querier) ReplaceAllRegex(exp, new string) *Querier {
	return q.AddFilter(func(s string) string {
		re, err := regexp.Compile(exp)
		if err != nil {
			return s
		}
		return re.ReplaceAllString(s, new)
	})
}

func (q *Querier) Prefix(p string) *Querier {
	return q.AddFilter(func(s string) string {
		return fmt.Sprintf(p + s)
	})
}
