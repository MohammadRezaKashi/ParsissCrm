package render

import (
	"ParsissCrm/internal/config"
	"ParsissCrm/internal/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	ptime "github.com/yaa110/go-persian-calendar"

	"github.com/jackc/pgtype"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig
var pathToTemplates = "./templates"
var language = "en"

var functions = template.FuncMap{
	"humanDate":              HumanDate,
	"humanPgtypeDate":        HumanPgtypeDate,
	"humanPgtypeDatePersian": HumanPgtypeDatePersian,
	"timestampToTime":        TimestampToTime,
	"i18n":                   I18n,
}

func I18n(key string) string {
	jsonFile, err := os.Open("./static/i18n/i18n.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]map[string]string
	json.Unmarshal([]byte(byteValue), &result)

	return result[language][key]
}

func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func HumanPgtypeDate(t pgtype.Date) string {
	return t.Time.Format("2006/01/02")
}

func HumanPgtypeDatePersian(t pgtype.Date) string {
	if t.Status == pgtype.Undefined {
		return ""
	}

	pt := ptime.New(t.Time)
	if pt.Year() < 0 || pt.Month() < 0 || pt.Day() < 0 {
		return ""
	}
	result := pt.Format("yyyy/MM/dd")
	return result
}

func TimestampToTime(t pgtype.Timestamp) string {
	time := t.Time.Format("15:04")
	return time
}

func NewRenderer(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Language = language
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}

	if app.Session.Exists(r.Context(), "access_level") {
		accessLevel := app.Session.Get(r.Context(), "access_level")
		td.AccessLevel = accessLevel.(int)
	}
	return td
}

func Template(rw http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	setDefaultLanguage(r, rw)

	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		return errors.New("can't get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(rw)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return err
	}

	return nil
}

func setDefaultLanguage(r *http.Request, rw http.ResponseWriter) {
	lang, err := r.Cookie("lang")
	if err != nil {
		cookie := &http.Cookie{
			Name:  "lang",
			Value: "en",
			Path:  "/",
		}
		http.SetCookie(rw, cookie)
		language = "en"
	} else {
		language = lang.Value
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCash := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
	if err != nil {
		return myCash, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		fmt.Println("Page is currently", page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCash, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil {
			return myCash, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if err != nil {
				return myCash, err
			}
		}

		myCash[name] = ts
	}

	return myCash, nil
}
