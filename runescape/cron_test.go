package runescape

import (
	"testing"
	"time"
)

func Test_ShouldRunHourlySameHour(t *testing.T) {
	last := time.Date(2026, 7, 7, 10, 0, 30, 0, time.UTC)
	now := time.Date(2026, 7, 7, 10, 58, 0, 0, time.UTC)

	if shouldRunHourly(last, now) {
		t.Error("Expecting false within the same hour")
	}
}

func Test_ShouldRunHourlyNextHour(t *testing.T) {
	last := time.Date(2026, 7, 7, 10, 58, 0, 0, time.UTC)
	now := time.Date(2026, 7, 7, 11, 0, 5, 0, time.UTC)

	if !shouldRunHourly(last, now) {
		t.Error("Expecting true in a new hour")
	}
}

func Test_ShouldRunHourlyNextDaySameHourValue(t *testing.T) {
	last := time.Date(2026, 7, 7, 10, 30, 0, 0, time.UTC)
	now := time.Date(2026, 7, 8, 10, 30, 0, 0, time.UTC)

	if !shouldRunHourly(last, now) {
		t.Error("Expecting true a day later even with the same hour value")
	}
}
