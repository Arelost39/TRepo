package db_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"
	"reflect"

	"github.com/DATA-DOG/go-sqlmock"

	"test_kode/internal/db"
	m "test_kode/internal/models"
)

// helper
func newMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()
	d, m, err := sqlmock.New()
	if err != nil { t.Fatalf("sqlmock: %v", err) }
	return d, m
}

func ReturnDate(y int, m time.Month, d int) time.Time {
	var loc *time.Location = time.UTC
	return time.Date(y, m, d, 0, 0, 0, 0, loc)
}

// проверяем создание рецепта
func TestCreateSchedule(t *testing.T) {
	sqlDB, mock := newMockDB(t)
	defer sqlDB.Close()

	DBinit := db.New(sqlDB)

	start := ReturnDate(2025, 9, 1)

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO schedules (user_id, pill_id, period_minutes, start_date, end_date, created_at)
		VALUES ($1, $2, $3, $4::date, $5::date, now())
		RETURNING schedule_id
	`)).WithArgs(int64(123), int64(5), 60, start, nil).
	WillReturnRows(sqlmock.NewRows([]string{"schedule_id"}).
	AddRow(int64(42)))

	id, err := DBinit.CreateSchedule(context.Background(), m.Schedule{
		UserID: 123, PillID: 5, PeriodMinutes: 60, StartDate: start, EndDate: nil,
	})
	if err != nil { t.Fatalf("Ошибка создания рецепта: %v", err) }
	if id != 42 { t.Fatalf("ожидаем id=42, получили %d", id) }

	if err := mock.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}

// проверяем вывод рецепта по ID
func TestListPrescriptionID(t *testing.T) {
	sqlDB, mock := newMockDB(t)
	defer sqlDB.Close()

	DBinit := db.New(sqlDB)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT schedule_id
		FROM schedules
		WHERE user_id = $1 AND (end_date IS NULL OR end_date >= CURRENT_DATE)
	`)).WithArgs(int64(123)).
	WillReturnRows(sqlmock.NewRows([]string{"schedule_id"}).
	AddRow(int64(3)).AddRow(int64(2)).AddRow(int64(1)))

	ids, err := DBinit.ListScheduleIDs(context.Background(), 123)
	if err != nil { t.Fatalf("ListScheduleIDs err: %v", err) }
	if len(ids) != 3 || ids[0] != 3 || ids[2] != 1 {
		t.Fatalf("получили %#v", ids)
	}
	if err := mock.ExpectationsWereMet(); err != nil { t.Fatal(err) }
}

// проверяем вывод ошибки
func TestGetPrescription_Error(t *testing.T) {
	sqlDB, mock := newMockDB(t)
	defer sqlDB.Close()

	DBinit := db.New(sqlDB)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT schedule_id, user_id, pill_id, period_minutes, start_date, end_date, created_at
		FROM schedules
		WHERE schedule_id = $1 AND user_id = $2 AND (end_date IS NULL OR end_date >= CURRENT_DATE)
	`)).
	WithArgs(int64(7), int64(123)).
	WillReturnError(sql.ErrNoRows)

	_, err := DBinit.GetSchedule(context.Background(), 123, 7)
	if err == nil || err != db.ErrNotFound {
		t.Fatalf("ждали ErrNotFound, получили %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}


// тут проверям и извлечение рецептов
func TestListSchedulesByUser_Window(t *testing.T) {
	sqlDB, mock := newMockDB(t)
	defer sqlDB.Close()

	DBinit := db.New(sqlDB)


	created := ReturnDate(2025, 9, 1)
	start := ReturnDate(2025, 9, 1)
	//ended := ReturnDate(2025, 9, 2)
	future := ReturnDate(2026, 9, 2)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT schedule_id, user_id, pill_id, period_minutes, start_date, end_date, created_at
		FROM schedules
		WHERE user_id = $1 AND (end_date IS NULL OR end_date >= CURRENT_DATE)
	`)).
	WithArgs(int64(1)).
	WillReturnRows(sqlmock.NewRows([]string{
		"schedule_id", "user_id", "pill_id", "period_minutes", "start_date", "end_date", "created_at",}).
	AddRow(int64(1), int64(1), int64(1), int64(60), start, nil, created).
	AddRow(int64(2), int64(1), int64(2), int64(90), start, future, created))

	list, err := DBinit.ListSchedulesByUser(context.Background(), 1)
	if err != nil {
		t.Fatalf("ошибка загрузки рецептов %v", err)
	}
	sch1 := m.Schedule{
		ScheduleID: 1,
		UserID: 1,
		PillID: 1,
		PeriodMinutes: 60,
		StartDate: start,
		EndDate: nil,
		CreatedAt: created,
	}

	sch2 := m.Schedule{
		ScheduleID: 2,
		UserID: 1,
		PillID: 2,
		PeriodMinutes: 90,
		StartDate: start,
		EndDate: &future,
		CreatedAt: created,
	}

	var expected = []m.Schedule{sch1, sch2}

	if !reflect.DeepEqual(expected, list) {
		t.Fatalf("Ошибка извлечения рецептов %v \n %v", expected, list)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
