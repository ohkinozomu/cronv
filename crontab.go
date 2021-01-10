package cronv

import (
	"fmt"
	"strings"
	"text/template"
)

type Schedule struct {
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	DayOfWeek  string
	Year       string
	Alias      string
}

func (self *Schedule) toCrontab() string {
	if self.Alias != "" {
		return self.Alias
	}
	dest := strings.Join([]string{self.Minute, self.Hour, self.DayOfMonth,
		self.Month, self.DayOfWeek, self.Year}, " ")
	return strings.Trim(dest, " ")
}

type Extra struct {
	Line  string
	Label string
	Job   string
}

type Crontab struct {
	Line     string
	Schedule *Schedule
	Job      string
}

func (c *Crontab) isRunningEveryMinutes() bool {
	for i, v := range strings.Split(c.Schedule.toCrontab(), " ") {
		if v != "*" && (i > 0 || v != "*/1") {
			return false
		}
	}
	return true
}

type InvalidTaskError struct {
	Line string
}

func (e *InvalidTaskError) Error() string {
	return fmt.Sprintf("Invalid task: '%s'", e.Line)
}

func parseCrontab(line string, crontabTZ string, outputTZ string) (*Crontab, *Extra, error) {
	// TODO use regrex to parse: https://gist.github.com/istvanp/310203
	parts := strings.Fields(line)

	schedule := &Schedule{}
	job := []string{}

	if strings.HasPrefix(parts[0], "@") {
		if len(parts) < 2 {
			return nil, nil, &InvalidTaskError{line}
		}

		// @reboot /something/to/do
		if parts[0] == "@reboot" {
			extra := &Extra{
				Line:  line,
				Label: parts[0],
				Job:   strings.Join(parts[1:], " "),
			}
			return nil, extra, nil
		}

		schedule.Alias = parts[0]
		job = parts[1:]
	} else {
		if len(parts) < 5 {
			return nil, nil, &InvalidTaskError{line}
		}

		// https://en.wikipedia.org/wiki/Cron#Predefined_scheduling_definitions
		c := 0
		for _, v := range parts {
			if len(v) == 0 {
				continue
			}
			switch c {
			case 0:
				schedule.Minute = v
			case 1:
				schedule.Hour = v
			case 2:
				schedule.DayOfMonth = v
			case 3:
				schedule.Month = v
			case 4:
				schedule.DayOfWeek = v
			default:
				job = append(job, v)
			}
			c++
		}
	}

	cs, err := convertTZ(schedule, crontabTZ, outputTZ)
	if err != nil {
		return nil, nil, err
	}

	crontab := &Crontab{
		Line:     line,
		Schedule: cs,
		Job:      template.JSEscapeString(strings.Join(job, " ")),
	}

	return crontab, nil, nil
}

func convertTZ(s *Schedule, crontabTZ string, outputTZ string) (*Schedule, error) {
	schedule := &Schedule{}
	schedule.Minute = s.Minute
	schedule.Hour = s.Hour
	schedule.DayOfMonth = s.DayOfMonth
	schedule.Month = s.Month
	schedule.DayOfWeek = s.DayOfWeek
	schedule.Alias = s.Alias
	return schedule, nil
}
