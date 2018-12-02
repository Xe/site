package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/Unknwon/i18n"
	"github.com/Xe/ln"
)

func logTemplateTime(name string, from time.Time) {
	now := time.Now()
	ln.Log(context.Background(), ln.F{"action": "template_rendered", "dur": now.Sub(from).String(), "name": name})
}

func (s *Site) renderTemplatePage(templateFname string, data interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer logTemplateTime(templateFname, time.Now())
		s.tlock.RLock()
		defer s.tlock.RUnlock()

		const fallbackLang = `en-US`

		getTranslation := func(group, key string, vals ...interface{}) string {
			var lang string
			locale, err := GetPreferredLocale(r)
			if err != nil {
				ln.Error(r.Context(), err)
				lang = fallbackLang
				goto skip
			}

			if !i18n.IsExist(locale.Lang) {
				lang = fallbackLang
				goto skip
			}

			lang = locale.Lang
		skip:
			return i18n.Tr(lang, group+"."+key, vals...)
		}

		funcMap := template.FuncMap{
			"trans": getTranslation,
		}

		var t = template.New(templateFname).Funcs(funcMap)
		var err error

		if s.templates[templateFname] == nil {
			t, err = template.ParseFiles("templates/base.html", "templates/"+templateFname)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				ln.Error(context.Background(), err, ln.F{"action": "renderTemplatePage", "page": templateFname})
				fmt.Fprintf(w, "error: %v", err)
			}

			ln.Log(context.Background(), ln.F{"action": "loaded_new_template", "fname": templateFname})

			s.tlock.RUnlock()
			s.tlock.Lock()
			s.templates[templateFname] = t
			s.tlock.Unlock()
			s.tlock.RLock()
		} else {
			t = s.templates[templateFname]
		}

		err = t.Execute(w, data)
		if err != nil {
			panic(err)
		}
	})
}

func (s *Site) showPost(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/blog/" {
		http.Redirect(w, r, "/blog", http.StatusSeeOther)
		return
	}

	var p *Post
	for _, pst := range s.Posts {
		if pst.Link == r.RequestURI[1:] {
			p = pst
		}
	}

	if p == nil {
		w.WriteHeader(http.StatusNotFound)
		s.renderTemplatePage("error.html", "no such post found: "+r.RequestURI).ServeHTTP(w, r)
		return
	}

	s.renderTemplatePage("blogpost.html", p).ServeHTTP(w, r)
}
