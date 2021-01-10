package cronv

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

type Cronv struct {
	Crontab         *Crontab
	expr            cron.Schedule
	startTime       time.Time
	durationMinutes float64
	location        string
}

func NewCronv(line string, startTime time.Time, durationMinutes float64, crontabTZ string, outputTZ string) (*Cronv, *Extra, error) {
	crontab, extra, err := parseCrontab(line, crontabTZ, outputTZ)
	if err != nil {
		return nil, nil, err
	}

	// Maybe the line was extra (@reboot, ENV etc ...)
	if crontab == nil {
		return nil, extra, nil
	}

	expr, err := cron.ParseStandard(crontab.Schedule.toCrontab())
	if err != nil {
		return nil, nil, err
	}

	cronv := &Cronv{
		Crontab:         crontab,
		expr:            expr,
		startTime:       startTime,
		durationMinutes: durationMinutes,
	}

	return cronv, extra, nil
}

type Exec struct {
	Start time.Time
	End   time.Time
}

func (self *Cronv) iter() <-chan *Exec {
	ch := make(chan *Exec)
	eneTime := self.startTime.Add(time.Duration(self.durationMinutes) * time.Minute)
	next := self.expr.Next(self.startTime)
	go func() {
		for next.Equal(eneTime) || eneTime.After(next) {
			ch <- &Exec{
				Start: next,
				End:   next.Add(time.Duration(1) * time.Minute),
			}
			next = self.expr.Next(next)
		}
		close(ch)
	}()
	return ch
}

type CronvCtx struct {
	Opts            *Command
	TimeFrom        time.Time
	TimeTo          time.Time
	CronEntries     []*Cronv
	Extras          []*Extra
	durationMinutes float64
	CrontabTZ       string
	OutputTZ        string
}

func NewCtx(opts *Command) (*CronvCtx, error) {
	timeFrom, err := opts.toFromTime()
	if err != nil {
		return nil, err
	}

	durationMinutes, err := opts.toDurationMinutes()
	if err != nil {
		return nil, err
	}

	return &CronvCtx{
		Opts:            opts,
		TimeFrom:        timeFrom,
		TimeTo:          timeFrom.Add(time.Duration(durationMinutes) * time.Minute),
		durationMinutes: durationMinutes,
		CrontabTZ:       opts.CrontabTZ,
		OutputTZ:        opts.OutputTZ,
	}, nil
}

func (self *CronvCtx) AppendNewLine(line string) (bool, error) {
	trimed := strings.TrimSpace(line)
	if len(trimed) == 0 || string(trimed[0]) == "#" {
		return false, nil
	}

	cronv, extra, err := NewCronv(trimed, self.TimeFrom, self.durationMinutes, self.CrontabTZ, self.OutputTZ)
	if err != nil {
		switch err.(type) {
		case *InvalidTaskError:
			return false, nil // pass
		default:
			return false, fmt.Errorf("Failed to analyze cron '%s': %s", line, err)
		}
	}

	if cronv != nil {
		self.CronEntries = append(self.CronEntries, cronv)
	}
	if extra != nil {
		self.Extras = append(self.Extras, extra)
	}

	return true, nil
}

func (self *CronvCtx) Dump() error {
	if err := os.Mkdir("assets", 0644); !os.IsExist(err) {
		return err
	}
	jsf, err := os.Create("assets/graph.js")
	if err != nil {
		return err
	}
	js, err := buildJS()
	if err != nil {
		return err
	}
	js.Execute(jsf, self)

	output, err := os.Create("assets/" + self.Opts.OutputFilePath)
	if err != nil {
		return err
	}
	t, err := makeTemplate()
	if err != nil {
		return err
	}
	t.Execute(output, self)
	return nil
}
