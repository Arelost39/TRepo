package helpers_test

import (
	"testing"
	"time"

	h "test_kode/internal/helpers"
)

func TestGenDaySlots_Basic(t *testing.T) {
	// период 60 минут
	slots := h.GenDaySlots(60)

	if len(slots) == 0 {
		t.Fatal("ожидали непустой список слотов")
	}

	first := slots[0]
	if first.Hour() != 8 || first.Minute() != 0 {
		t.Fatalf("ожидали  08:00, получили %02d:%02d", first.Hour(), first.Minute())
	}

	// проверим последний слот до 22:00
	last := slots[len(slots)-1]
	if !last.Before(time.Date(last.Year(), last.Month(), last.Day(), 22, 0, 0, 0, last.Location())) {
		t.Fatalf("должен быть до 22:00, получили %v", last)
	}
}