package server

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
	"write/logger"
	"write/models"
	"write/service"

	"go.opentelemetry.io/otel"
)

const TIMEZONE_OFFSET_MS = 19800000
const TIMEZONE_OFFSET_S = 19800

var (
	serverTracer = otel.Tracer("write-service")
)

func WriteReading(w http.ResponseWriter, r *http.Request) {
	var reqBody []byte
	var reading models.Reading

	ctx := r.Context()
	ctx, span := serverTracer.Start(ctx, "server.WriteReading")
	defer span.End()

	logEntry := logger.
		WithTracingFields(ctx).
		WithContext(ctx)

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		logEntry.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ErrorToResponse(err))
		return
	}

	err = json.Unmarshal(reqBody, &reading)
	if err != nil {
		logEntry.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ErrorToResponse(err))
		return
	}

	reading.DateTime = time.Now().Unix()

	err = service.WriteReading(ctx, reading)
	if err != nil {
		logEntry.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ErrorToResponse(err))
		return
	}

	logEntry.Info("OK")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"msg": "OK"}`))
	return
}

// GetSetpoints -
func GetSetpoints(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{}`))
	return
}
