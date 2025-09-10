package helpers

import (
	"net/http"
	"encoding/json"
	"time"
	"strconv"
	"errors"

	m "test_kode/internal/models"
)

// вспомогательные функции для server
func ParseInt64(s string) (int64, error) {
	if s == "" {
		return 0, errors.New("пустое значение")
	}
	return strconv.ParseInt(s, 10, 64)
}

func ParseDate(s string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", s, time.Local)
}

func WriteJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteErr(w http.ResponseWriter, code int, msg string) {
	WriteJSON(w, code, map[string]string{"error": msg})
}

// вспомогательные функции для модуля db

func GenDaySlots(periodMinutes int) []time.Time {

	if periodMinutes <= 0 {
		return nil
	}
	
	day:= time.Now()

	dayStart := time.Date(day.Year(), day.Month(), day.Day(), 8, 0, 0, 0, time.Local)
	dayEnd   := time.Date(day.Year(), day.Month(), day.Day(), 22, 0, 0, 0, time.Local)

	t := CeilToQuarter(dayStart)
	var out []time.Time
	step := time.Duration(periodMinutes) * time.Minute

	for t.Before(dayEnd) {
		out = append(out, t)
		t = CeilToQuarter(t.Add(step))
	}
	return out
}

func CeilToQuarter(t time.Time) time.Time {
	t = t.Truncate(time.Minute)
	min := t.Minute()
	add := (15 - (min % 15)) % 15
	if add == 0 {
		return t
	}
	return t.Add(time.Duration(add) * time.Minute)
}

func DateOnly(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func ToTimeSteps(ts []time.Time) []m.TimeSteps {
	out := make([]m.TimeSteps, len(ts))
	for i, t := range ts {
		out[i] = m.TimeSteps{Hour: t.Hour(), Minute: t.Minute()}
	}
	return out
}