package cronv

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/tkuchiki/go-timezone"
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

func getOffset(key string) (int, error) {
	tz := timezone.New()
	tzAbbrInfos, err := tz.GetTzAbbreviationInfo(key)
	if err != nil {
		return 0, err
	}
	return tzAbbrInfos[0].Offset(), nil
}

func getOffsetDifference(TZ1 string, TZ2 string) (int, error) {
	o1, err := getOffset(TZ1)
	if err != nil {
		return 0, err
	}
	o2, err := getOffset(TZ2)
	if err != nil {
		return 0, err
	}
	diff := o2 - o1
	return diff, nil
}

func addZero(key string) string {
	i, _ := strconv.Atoi(key)
	if i >= 0 && i < 10 {
		return "0" + key
	}
	return key
}

func isScheduleRunningEveryMinutes(s *Schedule) bool {
	for i, v := range strings.Split(s.toCrontab(), " ") {
		if v != "*" && (i > 0 || v != "*/1") {
			return false
		}
	}
	return true
}

func convertTZ(s *Schedule, crontabTZ string, outputTZ string) (*Schedule, error) {
	if isScheduleRunningEveryMinutes(s) {
		return s, nil
	}
	offsetDifference, err := getOffsetDifference(crontabTZ, outputTZ)
	if err != nil {
		return nil, err
	}

	//TODO: error handling
	y, _ := strconv.Atoi(s.Year)
	mo, _ := strconv.Atoi(s.Month)
	dom, _ := strconv.Atoi(s.DayOfMonth)
	h, _ := strconv.Atoi(s.Hour)
	mi, _ := strconv.Atoi(s.Minute)

	t1 := time.Date(y, time.Month(mo), dom, h, mi, 0, 0, time.UTC)
	t2 := t1.Add(time.Duration(offsetDifference) * time.Second)

	schedule := &Schedule{}

	// this is little hacky...
	schedule.Year = s.Year

	schedule.Month = strconv.Itoa(int(t2.Month()))
	schedule.Minute = addZero(strconv.Itoa(t2.Minute()))
	schedule.Hour = addZero(strconv.Itoa(t2.Hour()))
	schedule.DayOfMonth = strconv.Itoa(t2.Day())
	schedule.DayOfWeek = s.DayOfWeek
	schedule.Alias = s.Alias

	return schedule, nil
}
