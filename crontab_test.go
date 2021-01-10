package cronv

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseCrontab(t *testing.T) {
	line := "01 04 1 2 3 /usr/bin/somedirectory/somecommand1"
	r, _, err := parseCrontab(line, "UTC", "UTC")
	assert.NotNil(t, r)
	assert.Nil(t, err)
	assert.Equal(t, r.Line, line)
	assert.Equal(t, r.Schedule.Minute, "01")
	assert.Equal(t, r.Schedule.Hour, "04")
	assert.Equal(t, r.Schedule.DayOfMonth, "1")
	assert.Equal(t, r.Schedule.Month, "2")
	assert.Equal(t, r.Schedule.DayOfWeek, "3")
	assert.Equal(t, r.Job, "/usr/bin/somedirectory/somecommand1")
}

func TestParseCrontabConvert1(t *testing.T) {
	line := "01 04 1 1 1 /usr/bin/somedirectory/somecommand1"
	r, _, err := parseCrontab(line, "UTC", "JST")
	assert.NotNil(t, r)
	assert.Nil(t, err)
	assert.Equal(t, r.Line, line)
	assert.Equal(t, r.Schedule.Minute, "01")
	assert.Equal(t, r.Schedule.Hour, "13")
	assert.Equal(t, r.Schedule.DayOfMonth, "1")
	assert.Equal(t, r.Schedule.Month, "1")
	assert.Equal(t, r.Schedule.DayOfWeek, "1")
	assert.Equal(t, r.Job, "/usr/bin/somedirectory/somecommand1")
}

func TestParseCrontabConvert2(t *testing.T) {
	line := "01 04 * * * /usr/bin/somedirectory/somecommand2"
	r, _, err := parseCrontab(line, "UTC", "JST")
	assert.NotNil(t, r)
	assert.Nil(t, err)
	assert.Equal(t, r.Line, line)
	assert.Equal(t, r.Schedule.Minute, "01")
	assert.Equal(t, r.Schedule.Hour, "13")
	assert.Equal(t, r.Schedule.DayOfMonth, "1")
	assert.Equal(t, r.Schedule.Month, "1")
	assert.Equal(t, r.Schedule.DayOfWeek, "1")
	assert.Equal(t, r.Job, "/usr/bin/somedirectory/somecommand1")
}

//
//01,31 04,05 1-15 1,6 * /usr/bin/somedirectory/somecommand3
//*/5 * * * * /usr/bin/somedirectory/somecommand4
//*/18 * * * * /usr/bin/somedirectory/somecommand5
//0,10,20,30,40,50 * * * * /usr/bin/somedirectory/somecommand6
//00 01 * * * rusty /home/rusty/rusty-list-files.sh
//00 06 * * * env DISPLAY=:0 gui_appname
//45 04 * * * /usr/bin/updatedb
//*/1 * * * * /usr/bin/updatedb2
//*/3 * * * * /usr/bin/updatedb3
//* * * * * date >> /tmp/aaa
//@reboot root /usr/bin/rebooted arg1 arg2
//@reboot root /usr/bin/rebooted --arg1=val2
//@hourly /usr/bin/somedirectory/some-hourly-command arg1 arg2
//
//# duplicated jobs
//00 10 * * * user /usr/bin/dup > /dev/null
//10 17 * * * user /usr/bin/dup > /dev/null

func TestParseCrontabInvalidTask(t *testing.T) {
	line := "MAILTO=example.com"
	_, _, err := parseCrontab(line, "UTC", "UTC")
	assert.NotNil(t, err)
}

func TestIsRunningEveryMinutesFalseCase(t *testing.T) {
	c, _, _ := parseCrontab("3 * * * *", "UTC", "UTC")
	assert.False(t, c.isRunningEveryMinutes())

	c2, _, _ := parseCrontab("* * * * 1", "UTC", "UTC")
	assert.False(t, c2.isRunningEveryMinutes())
}

func TestIsRunningEveryMinutesTrueCase(t *testing.T) {
	c, _, _ := parseCrontab("* * * * *", "UTC", "UTC")
	assert.True(t, c.isRunningEveryMinutes())

	c2, _, _ := parseCrontab("*/1 * * * *", "UTC", "UTC")
	assert.True(t, c2.isRunningEveryMinutes())
}

func TestAlias(t *testing.T) {
	line := "@hourly /path/to/do/something arg1"
	c, _, err := parseCrontab(line, "UTC", "UTC")
	assert.Nil(t, err)
	assert.Equal(t, c.Schedule.Alias, "@hourly")
	assert.Equal(t, c.Job, "/path/to/do/something arg1")
}

func TestExtra(t *testing.T) {
	line := "@reboot root /path/to/do/something arg1 arg2 arg3"
	c, e, _ := parseCrontab(line, "UTC", "UTC")
	assert.Nil(t, c)
	assert.NotNil(t, e)
	assert.Equal(t, e.Line, line)
	assert.Equal(t, e.Label, "@reboot")
	assert.Equal(t, e.Job, "root /path/to/do/something arg1 arg2 arg3")
}

func TestConvertTZ_UTC_UTC(t *testing.T) {
	schedule := &Schedule{}
	schedule.Minute = "01"
	schedule.Hour = "04"
	schedule.DayOfMonth = "1"
	schedule.Month = "1"
	schedule.DayOfWeek = "1"
	s, err := convertTZ(schedule, "UTC", "UTC")
	assert.Nil(t, err)
	assert.Equal(t, s.Minute, "01")
	assert.Equal(t, s.Hour, "04")
	assert.Equal(t, s.DayOfMonth, "1")
	assert.Equal(t, s.Month, "1")
	assert.Equal(t, s.DayOfWeek, "1")
}

func TestConvertTZ_UTC_UTC_Empty_Some_Values(t *testing.T) {
	schedule := &Schedule{}
	schedule.Minute = "01"
	schedule.Hour = "04"
	schedule.DayOfMonth = "1"
	schedule.Month = "1"
	schedule.DayOfWeek = "1"
	schedule.Year = ""
	schedule.Alias = ""
	s, err := convertTZ(schedule, "UTC", "UTC")
	assert.Nil(t, err)
	assert.Equal(t, s.Minute, "01")
	assert.Equal(t, s.Hour, "04")
	assert.Equal(t, s.DayOfMonth, "1")
	assert.Equal(t, s.Month, "1")
	assert.Equal(t, s.DayOfWeek, "1")
}

func TestGetOffset(t *testing.T) {
	o, err := getOffset("JST")
	assert.Nil(t, err)
	assert.Equal(t, 32400, o)
}

func TestGetOffsetDifference(t *testing.T) {
	od, err := getOffsetDifference("UTC", "JST")
	assert.Nil(t, err)
	assert.Equal(t, 32400, od)
}

func TestConvertTZ_UTC_JST(t *testing.T) {
	schedule := &Schedule{}
	schedule.Minute = "01"
	schedule.Hour = "04"
	schedule.DayOfMonth = "1"
	schedule.Month = "1"
	schedule.DayOfWeek = "1"
	s, err := convertTZ(schedule, "UTC", "JST")
	if err != nil {
		t.Log(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, "01", s.Minute)
	assert.Equal(t, "13", s.Hour)
	assert.Equal(t, "1", s.DayOfMonth)
	assert.Equal(t, "1", s.Month)
	assert.Equal(t, "1", s.DayOfWeek)
}
