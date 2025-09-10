package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	m "test_kode/internal/models"
)

/* Предполагаемые таблицы в SQL БД:
1) Пациенты
CREATE TABLE patients (
  	user_id				BIGINT PRIMARY KEY,
  	created_at			TIMESTAMPTZ NOT NULL DEFAULT now()
);

2) Лекарства
CREATE TABLE pills (
  	pill_id				BIGSERIAL PRIMARY KEY,
  	name				TEXT NOT NULL UNIQUE,
  	created_at			TIMESTAMPTZ NOT NULL DEFAULT now()
);

3) Рецепты (по одному лекарству на рецепт; просто и достаточно)
CREATE TABLE schedules (
  	schedule_id			BIGSERIAL PRIMARY KEY,
  	user_id				BIGINT NOT NULL REFERENCES patients(user_id) ON DELETE CASCADE,
  	pill_id				BIGINT NOT NULL REFERENCES drugs(drug_id),
  	period_minutes		SMALLINT NOT NULL,
  	start_date			DATE NOT NULL DEFAULT CURRENT_DATE, <--> старт курса
  	end_date			DATE,  <-------------------------------> конец курса, NULL = постоянный рецепт
  	created_at 			TIMESTAMPTZ NOT NULL DEFAULT now(),
);
*/

var ErrNotFound = errors.New("not found")

type Crud interface {
	CreateSchedule(ctx context.Context, p m.Schedule) (int64, error)
	ListScheduleIDs(ctx context.Context, userID int64) ([]int64, error)
	GetSchedule(ctx context.Context, userID, scheduleID int64) (m.Schedule, error)
	ListSchedulesByUser(ctx context.Context, userID int64) ([]m.Schedule, error)
}

type DB struct{
	sql *sql.DB
}

func New(sqlDB *sql.DB) *DB {
	return &DB{sql: sqlDB}
}

func (d *DB) CreateSchedule(ctx context.Context, p m.Schedule) (int64, error) {

	const q = `
		INSERT INTO schedules (user_id, pill_id, period_minutes, start_date, end_date, created_at)
		VALUES ($1, $2, $3, $4::date, $5::date, now())
		RETURNING schedule_id
	`

	var id int64
	if err := d.sql.QueryRowContext(ctx, q, p.UserID, p.PillID, p.PeriodMinutes, p.StartDate, p.EndDate).Scan(&id); err != nil {
		return 0, fmt.Errorf("ошибка создания рецепта: %w", err)
	}
	return id, nil
}

func (d *DB) ListScheduleIDs(ctx context.Context, userID int64) ([]int64, error) {

	const q = `
		SELECT schedule_id
		FROM schedules
		WHERE user_id = $1 AND (end_date IS NULL OR end_date >= CURRENT_DATE)
	`

	rows, err := d.sql.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list ids: %w", err)
	}
	defer rows.Close()

	var out []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}

func (d *DB) GetSchedule(ctx context.Context, userID, scheduleID int64) (m.Schedule, error) {
	
	const q = `
		SELECT schedule_id, user_id, pill_id, period_minutes, start_date, end_date, created_at
		FROM schedules 
		WHERE schedule_id = $1 AND user_id = $2 AND (end_date IS NULL OR end_date >= CURRENT_DATE)
	`
	
	var p m.Schedule
	var end sql.NullTime
	if err := d.sql.QueryRowContext(ctx, q, scheduleID, userID).
		Scan(&p.ScheduleID, &p.UserID, &p.PillID, &p.PeriodMinutes, &p.StartDate, &end, &p.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) { return m.Schedule{}, ErrNotFound }
		return m.Schedule{}, fmt.Errorf("get: %w", err)
	}
	return p, nil
}

func (d *DB) ListSchedulesByUser(ctx context.Context, userID int64) ([]m.Schedule, error) {
	const q = `
		SELECT schedule_id, user_id, pill_id, period_minutes, start_date, end_date, created_at
		FROM schedules
		WHERE user_id = $1 AND (end_date IS NULL OR end_date >= CURRENT_DATE)
	`
	rows, err := d.sql.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list by user: %w", err)
	}

	defer rows.Close()

	var out []m.Schedule

	for rows.Next() {
		var p m.Schedule

		if err := rows.Scan(
			&p.ScheduleID, &p.UserID, &p.PillID, &p.PeriodMinutes, &p.StartDate, &p.EndDate, &p.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
