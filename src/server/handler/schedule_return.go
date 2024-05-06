package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/agavris/june-academy-go/src/algorithm/scheduler"
	"golang.org/x/exp/slog"
)

type ScheduleHandler struct {
	Scheduler *scheduler.Scheduler
}

func NewScheduleHandler() *ScheduleHandler {
	return &ScheduleHandler{
		Scheduler: scheduler.NewScheduler(),
	}
}

func (h *ScheduleHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	iterations, err := strconv.Atoi(request.URL.Query().Get("iterations"))
	if err != nil {
		http.Error(writer, "Invalid iterations parameter", http.StatusBadRequest)
		slog.Error("failed to parse iterations", err)
		return
	}

	resultSchedule := h.Scheduler.Run(iterations)

	response, err := json.Marshal(resultSchedule)
	if err != nil {
		http.Error(writer, "Failed to marshal schedule", http.StatusInternalServerError)
		slog.Error("failed to marshal schedule", err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, writeErr := writer.Write(response)
	if writeErr != nil {
		slog.Error("failed to write response", err)
	}
}
