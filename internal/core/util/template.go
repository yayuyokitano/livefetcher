package util

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"reflect"
	"time"

	i18nloader "github.com/yayuyokitano/livefetcher/internal/i18n"
)

func BuildTemplate(w io.Writer, r *http.Request, user AuthUser, funcMap template.FuncMap, paths ...string) (tmpl *template.Template, err error) {
	tmpl = template.New("layout")
	tmpl, err = tmpl.Funcs(template.FuncMap{
		"T": i18nloader.GetLocalizer(r).Localize,
		"ParseDate": func(t time.Time) string {
			return i18nloader.ParseDate(t, i18nloader.GetLanguages(w, r))
		},
		"Lang": func() string { return i18nloader.GetMainLanguage(w, r) },
		"GetUser": func() AuthUser {
			return user
		},
		"TemplateIfExists": func(name string, pipeline interface{}) (template.HTML, error) {
			t := tmpl.Lookup(name)
			if t == nil {
				return "", nil
			}

			buf := &bytes.Buffer{}
			err := t.Execute(buf, pipeline)
			if err != nil {
				return "", err
			}

			return template.HTML(buf.String()), nil
		},
		"HasField": func(name string, data interface{}) bool {
			v := reflect.ValueOf(data)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			if v.Kind() != reflect.Struct {
				return false
			}
			return v.FieldByName(name).IsValid()
		},
	}).Funcs(funcMap).ParseFiles(paths...)
	return
}
