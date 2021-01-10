package cronv

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseCrontab(t *testing.T) {
	line := "01 04 1 2 3	/usr/bin/somedirectory/somecommand1"
	r, _, _ := parseCrontab(line, "UTC", "UTC")
	assert.NotNil(t, r)
	assert.Equal(t, r.Line, line)
	assert.Equal(t, r.Schedule.Minute, "01")
	assert.Equal(t, r.Schedule.Hour, "04")
	assert.Equal(t, r.Schedule.DayOfMonth, "1")
	assert.Equal(t, r.Schedule.Month, "2")
	assert.Equal(t, r.Schedule.DayOfWeek, "3")
	assert.Equal(t, r.Job, "/usr/bin/somedirectory/somecommand1")
}

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

func TestConvertTZ(t *testing.T) {
	schedule := &Schedule{}
	schedule.Minute = "01"
	schedule.Hour = "04"
	schedule.DayOfMonth = "1"
	schedule.Month = "1"
	schedule.DayOfWeek = "1"
	s, err := convertTZ(schedule, "UTC", "UTC")
	if err != nil {
		t.Log(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, s.Minute, "01")
	assert.Equal(t, s.Hour, "04")
	assert.Equal(t, s.DayOfMonth, "1")
	assert.Equal(t, s.Month, "1")
	assert.Equal(t, s.DayOfWeek, "1")
}
