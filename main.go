package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// MeterMeasurementValue ...
type MeterMeasurementValue struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

// MeterMeasurement ...
type MeterMeasurement struct {
	MeasurementID string                  `json:"measurement_id"`
	Values        []MeterMeasurementValue `json:"values"`
}

// MeterReading ...
type MeterReading struct {
	MeterInternalID string             `json:"meter_internal_id"`
	Registers       []MeterMeasurement `json:"registers"`
}

// MetersReadings ...
type MetersReadings struct {
	Meters []MeterReading `json:"meters"`
}

func isHTTPMethodValid(
	r *http.Request,
	w http.ResponseWriter,
	method string) bool {

	if r.Method != method {
		errMsg := fmt.Sprintf("%s http method is not supported on %s resource", r.Method, r.URL.Path)
		writeErrorResponse(r, w, http.StatusMethodNotAllowed, errMsg)
		return false
	}

	return true
}

func writeJSONResponse(
	r *http.Request,
	w http.ResponseWriter,
	statusCode int,
	body interface{}) {

	log.Println(fmt.Sprintf(
		"%s %s - writing SUCCESS http %d response with json data: %+v",
		r.Method,
		r.URL.Path,
		statusCode,
		body))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)

}

func writeErrorResponse(
	r *http.Request,
	w http.ResponseWriter,
	statusCode int,
	errMsg string) {

	log.Println(fmt.Sprintf(
		"%s %s - writing ERROR http %d response with message: %s",
		r.Method,
		r.URL.Path,
		statusCode, errMsg))
	w.WriteHeader(statusCode)
	w.Write([]byte(errMsg))

}

func metersReadingsPost(w http.ResponseWriter, r *http.Request) {
	if !isHTTPMethodValid(r, w, http.MethodPost) {
		return
	}
	decoder := json.NewDecoder(r.Body)
	var metersReadings MetersReadings
	err := decoder.Decode(&metersReadings)
	if err != nil {
		errMsg := fmt.Sprintf("error parsing request body: %v", err)
		writeErrorResponse(r, w, http.StatusBadRequest, errMsg)
		return
	}
	writeJSONResponse(r, w, http.StatusOK, &metersReadings)
}

func main() {
	http.HandleFunc("/api/readings", metersReadingsPost)
	if err := http.ListenAndServe(":9999", nil); err != nil {
		panic(err)
	}
}
