// internal/service/service.go
package service

import (
	"context"

	"time"

	"test_kode/internal/config"
	"test_kode/internal/db"
	h "test_kode/internal/helpers"
	m "test_kode/internal/models"
)

type Service struct {
	crud db.Crud
	cfg  *config.Config // systemPeriod <- системное "окно"
}

type Config struct {
    DatabaseURL string
    LookAhead   time.Duration
}

func New(crud db.Crud, cfg *config.Config) *Service {
	return &Service{crud: crud, cfg: cfg}
}

// поскольку тут по ТЗ простые операции - оставлю проброс до мотодов db
func (s *Service) CreateSchedule(ctx context.Context, p m.Schedule) (int64, error) {
	return s.crud.CreateSchedule(ctx, p)
}

func (s *Service) ListScheduleIDs(ctx context.Context, userID int64) ([]int64, error) {
	return s.crud.ListScheduleIDs(ctx, userID)
}

// расписание на день
func (s *Service) DaySchedule(ctx context.Context, userID, scheduleID int64) (m.ScheduleForDay, error) {
	p, err := s.crud.GetSchedule(ctx, userID, scheduleID) // <- протестировано в db
	if err != nil {
		return m.ScheduleForDay{}, err
	}

	// функция получения среза окрегленного расписания
	slots := h.GenDaySlots(p.PeriodMinutes) // <- протестировано в helpers
	return m.ScheduleForDay{Schedule: p, Slots: slots}, nil
}

// приемы в окне из конфига
func (s *Service) NextTakings(ctx context.Context, systemPeriod, userID int64) ([]m.NextTaking, error) {

	ps, err := s.crud.ListSchedulesByUser(ctx, userID)// <- протестировано в модуле db
	if err != nil {
		return nil, err
	}

	// часть логики вынесена для облегчения тестировани
	out := TimeSteps(ps, systemPeriod)

	return out, nil
}

func TimeSteps(ps []m.Schedule, systemPeriod int64) []m.NextTaking {
	//now := time.Date(2025, 9, 9, 8, 0, 0, 0, time.Local)  //для тестирования дата и время 2025.09.09 8:00
	now := time.Now() // <- будем считать, что время задано на проде									
	var out []m.NextTaking
	dayEnd	:= time.Date(now.Year(), now.Month(), now.Day(), 22, 0, 0, 0, now.Location())
	
	for _, p := range ps {

		// емли период приема таблеток меньше 15, то для выполнения условий округления 
		// примем его за минимальный шаг - 15
		var step time.Duration
		if p.PeriodMinutes < 15 {
			step = 15 * time.Minute
		} else {
			step = time.Duration(p.PeriodMinutes) * time.Minute
		}

		// составляем срез структур с шагом прима с учетом округления на ближайшее время
		// начало схемы - начало дня
		var arr []m.TimeSteps
		board := now.Add(time.Duration(systemPeriod) * time.Minute)
		for t := now; t.Before(dayEnd) && !t.After(board); t = t.Add(step) {
			rondedTime := h.CeilToQuarter(t)
			arr = append(arr, m.TimeSteps{Hour: rondedTime.Hour(), Minute: rondedTime.Minute()})
		}
		// возвращаем схему приема с расписанием на ближайшее время
		if len(arr) > 0 {
			out = append(out, m.NextTaking{
				ScheduleID: p.ScheduleID,
				UserID:     p.UserID,
				PillID:     p.PillID,
				Time:       arr,
			})
		}
	}
	return out
}