package server

import (
	"encoding/json"
	"net/http"
	"time"

	c "test_kode/internal/config"
	s "test_kode/internal/service"
	m "test_kode/internal/models"
	h "test_kode/internal/helpers"
)

type Server struct {
	mux *http.ServeMux
	service *s.Service
	config *c.Config
}

func New(service *s.Service, config *c.Config) *Server {
	s := &Server{mux: http.NewServeMux(), service: service, config: config}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux.HandleFunc("POST /schedule", s.handleCreateSchedule)
	s.mux.HandleFunc("GET /schedules", s.handleListSchedules)
	s.mux.HandleFunc("GET /schedule", s.handleGetScheduleForDay)
	s.mux.HandleFunc("GET /next_takings", s.handleNextTakings)
}

// реализуем эндпоинты
// для http лог ошибок оставим на английском

// POST /schedule
func (s *Server) handleCreateSchedule(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var in m.CreateScheduleIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	if in.UserID == 0 || in.PillID == 0 || in.PeriodMinutes <= 0 {
		http.Error(w, "user_id, pill_id, period_minutes required", http.StatusBadRequest)
		return
	}

	start, err := h.ParseDate(in.StartDate)
	if err != nil {
		http.Error(w, "bad start_date (YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	var endPtr *time.Time
	if in.EndDate != nil && *in.EndDate != "" {
		et, err := h.ParseDate(*in.EndDate)
		if err != nil {
			http.Error(w, "bad end_date (YYYY-MM-DD)", http.StatusBadRequest)
			return
		}	
		endPtr = &et
	}

	id, err := s.service.CreateSchedule(r.Context(), m.Schedule{
		UserID: in.UserID,
		PillID: in.PillID,
		PeriodMinutes: in.PeriodMinutes,
		StartDate: start,
		EndDate: endPtr,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.WriteJSON(w, http.StatusOK, m.CreateScheduleOut{ScheduleID: id})
}

// GET /schedules?user_id=
func (s *Server) handleListSchedules(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.ParseInt64(r.URL.Query().Get("user_id"))
	if err != nil || userID == 0 {
		http.Error(w, "bad user_id", http.StatusBadRequest)
		return
	}

	ids, err := s.service.ListScheduleIDs(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError); return
	}

	h.WriteJSON(w, http.StatusOK, m.ListSchedulesOut{ScheduleIDs: ids})
}

// GET /schedule?user_id=&schedule_id=
func (s *Server) handleGetScheduleForDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()
	userID, err := h.ParseInt64(q.Get("user_id"))
	if err != nil || userID == 0 {
		http.Error(w, "bad user_id", http.StatusBadRequest)
		return
	}

	scheduleID, err := h.ParseInt64(q.Get("schedule_id"))
	if err != nil || scheduleID == 0 {
		http.Error(w, "bad schedule_id", http.StatusBadRequest)
		return
	}

	res, err := s.service.DaySchedule(r.Context(), userID, scheduleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out := m.NextTaking{
		ScheduleID:    res.Schedule.ScheduleID,
		UserID:        res.Schedule.UserID,
		PillID:        res.Schedule.PillID,
		Time:          h.ToTimeSteps(res.Slots),
	}

	h.WriteJSON(w, http.StatusOK, out)
}

// GET /next_takings?user_id=
func (s *Server) handleNextTakings(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()
	userID, err := h.ParseInt64(q.Get("user_id"))
	if err != nil || userID == 0 {
		http.Error(w, "bad user_id", http.StatusBadRequest)
		return
	}

	out, err := s.service.NextTakings(r.Context(), s.config.SystemPeriod, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// отдаём как есть
	h.WriteJSON(w, http.StatusOK, out)
}
