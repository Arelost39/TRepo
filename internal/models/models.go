package models

import (
	"time"
)

type CreateScheduleIn struct {
	PillID			int64		`json:"pill_id"`
	PeriodMinutes	int			`json:"period_minutes"`
	StartDate		string		`json:"start_date"`
	EndDate 		*string		`json:"end_date,omitempty"` // мб пустой
	UserID			int64		`json:"user_id"`
}

type CreateScheduleOut struct {
	ScheduleID 		int64		`json:"schedule_id"`
}

type ListSchedulesOut struct {
	ScheduleIDs []int64			`json:"schedule_ids"`
}

type Schedule struct {
	ScheduleID		int64
	UserID			int64
	PillID			int64
	PeriodMinutes	int // интервал между приёмами, мин
	StartDate		time.Time // дата назначения рецепта
	EndDate 		*time.Time 	// NULL = постоянный рецепт
	CreatedAt		time.Time
}

// в тз не указан формат возвращаемых врменных меток
// будем использовать следующие варианты
// (разные просто для демонстрации разных подходов)

type NextTaking struct {
	ScheduleID		int64		`json:"schedule_id"`
	UserID			int64		`json:"user_id"`
	PillID 			int64		`json:"pill_id"`
	Time			[]TimeSteps	`json:"time"`
}

type TimeSteps struct {
	Hour   int					`json:"hour"`
	Minute int					`json:"minute"`
}

type ScheduleForDay struct {
	Schedule
	Slots  			[]time.Time // локальные метки в пределах дневного окна и кратные 15
}