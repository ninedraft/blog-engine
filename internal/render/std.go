package render

import (
	"strings"
	"text/template"
	"time"
)

var Funcs = template.FuncMap{
	"now": func(format ...string) string {
		var now = time.Now()
		if len(format) == 0 {
			return now.String()
		}

		var layout = format[0]
		switch strings.ToLower(layout) {
		case "rfc3339":
			return now.Format(time.RFC3339)
		case "kitchen":
			return now.Format(time.Kitchen)
		case "datetime":
			return now.Format(time.DateTime)
		case "dateonly":
			return now.Format(time.DateOnly)
		case "unixdate":
			return now.Format(time.UnixDate)
		case "timeonly":
			return now.Format(time.TimeOnly)
		default:
			return now.Format(layout)
		}
	},

	"season": func() string {
		month := time.Now().Month()
		switch month {
		case time.December, time.January, time.February:
			return "❄️"
		case time.March, time.April, time.May:
			return "🌱"
		case time.June, time.July, time.August:
			return "☀️"
		default:
			return "🍂"
		}
	},
}
