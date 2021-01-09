package cronv

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"time"
)

func makeTemplate() (*template.Template, error) {
	f, err := os.Open("template/template.html")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	TEMPLATE := string(b)

	funcMap := template.FuncMap{
		"CronvIter": func(cronv *Cronv) <-chan *Exec {
			return cronv.iter()
		},
		"JSEscapeString": func(v string) string {
			return template.JSEscapeString(strings.TrimSpace(v))
		},
		"NewJsDate": func(v time.Time) string {
			return fmt.Sprintf("new Date(%d,%d,%d,%d,%d)", v.Year(), v.Month(), v.Day(), v.Hour(), v.Minute())
		},
		"DateFormat": func(v time.Time, format string) string {
			return v.Format(format)
		},
		"IsRunningEveryMinutes": func(c *Crontab) bool {
			return c.isRunningEveryMinutes()
		},
	}
	return template.Must(template.New("").Funcs(funcMap).Parse(TEMPLATE)), nil
}
