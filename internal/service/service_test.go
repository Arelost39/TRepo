package service_test

import (
	"reflect"
	m "test_kode/internal/models"
	"test_kode/internal/service"
	"testing"
)

// тесты расчета расписаний для разных временных интервалов
func TestNextTaking_15(t *testing.T) {

	time := []m.TimeSteps{
		{Hour: 8, Minute: 0},
		{Hour: 8, Minute: 15},
		{Hour: 8, Minute: 30},
		{Hour: 8, Minute: 45},
		{Hour: 9, Minute: 00},
	}

	expected := m.NextTaking{
		ScheduleID: 1,
		UserID: 1,
		PillID: 1,
		Time: time,
	}

	sch1 := m.Schedule{
		ScheduleID: 1,
		UserID: 1,
		PillID: 1,
		PeriodMinutes: 15,
	}

	schedule := []m.Schedule{
		sch1,
	}

	result := service.TimeSteps(schedule, 60)
	if !reflect.DeepEqual(expected, result[0]) {
		t.Fatal(expected, result[0])
	}
}
func TestNextTaking_16(t *testing.T) {

	time := []m.TimeSteps{
		{Hour: 8, Minute: 0},
		{Hour: 8, Minute: 30},
		{Hour: 8, Minute: 45},
		{Hour: 9, Minute: 00},
	}

	expected := m.NextTaking{
		ScheduleID: 1,
		UserID: 1,
		PillID: 1,
		Time: time,
	}

	sch1 := m.Schedule{
		ScheduleID: 1,
		UserID: 1,
		PillID: 1,
		PeriodMinutes: 16,
	}

	schedule := []m.Schedule{
		sch1,
	}

	result := service.TimeSteps(schedule, 60)
	if !reflect.DeepEqual(expected, result[0]) {
		t.Fatal(expected, result[0])
	}
}
func TestNextTaking_25(t *testing.T) {

	time := []m.TimeSteps{
		{Hour: 8, Minute: 0},
		{Hour: 8, Minute: 30},
		{Hour: 9, Minute: 00},
		{Hour: 9, Minute: 15},
		{Hour: 9, Minute: 45},
	}

	expected := m.NextTaking{
		ScheduleID: 1,
		UserID: 1,
		PillID: 1,
		Time: time,
	}

	sch1 := m.Schedule{
		ScheduleID: 1,
		UserID: 1,
		PillID: 1,
		PeriodMinutes: 25,
	}

	schedule := []m.Schedule{
		sch1,
	}

	result := service.TimeSteps(schedule, 120)
	if !reflect.DeepEqual(expected, result[0]) {
		t.Fatal(expected, result[0])
	}
}