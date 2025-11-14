package acp

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
)

func decodeJSON(body io.ReadCloser, v any) error {
	defer func() { _ = body.Close() }()
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("request body required")
		}
		return err
	}
	if dec.More() {
		return errors.New("unexpected data after JSON body")
	}
	return nil
}

func writeServiceError(w http.ResponseWriter, err error) {
	var httpErr *Error
	if errors.As(err, &httpErr) {
		writeJSONError(w, httpErr)
		return
	}
	writeJSONError(w, NewProcessingError("internal server error"))
}

func writeJSONError(w http.ResponseWriter, payload *Error) {
	if payload == nil {
		payload = NewProcessingError("internal server error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("API-Version", APIVersion)
	if seconds := retryAfterSeconds(payload.RetryAfter()); seconds > 0 {
		w.Header().Set("Retry-After", strconv.FormatInt(seconds, 10))
	}
	w.WriteHeader(payload.status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("API-Version", APIVersion)
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(payload)
}

func retryAfterSeconds(d time.Duration) int64 {
	if d <= 0 {
		return 0
	}
	seconds := d / time.Second
	if d%time.Second != 0 {
		seconds++
	}
	if seconds <= 0 {
		return 1
	}
	return int64(seconds)
}
