package config

import (
	"database/sql/driver"
	"time"
)

type SinkData struct {
	SinkID          string          `json:"sink_id"`
	OwnerID         string          `json:"owner_id"`
	Url             string          `json:"remote_host"`
	User            string          `json:"username"`
	Password        string          `json:"password"`
	OpenTelemetry   string          `json:"opentelemetry"`
	State           PrometheusState `json:"state,omitempty"`
	Msg             string          `json:"msg,omitempty"`
	LastRemoteWrite time.Time       `json:"last_remote_write,omitempty"`
}

const (
	Unknown PrometheusState = iota
	Active
	Error
	Idle
)

type PrometheusState int

var promStateMap = [...]string{
	"unknown",
	"active",
	"error",
	"idle",
}

var promStateRevMap = map[string]PrometheusState{
	"unknown": Unknown,
	"active":  Active,
	"error":   Error,
	"idle":    Idle,
}

func (p PrometheusState) String() string {
	return promStateMap[p]
}

func (p *PrometheusState) Scan(value interface{}) error {
	*p = promStateRevMap[string(value.([]byte))]
	return nil
}

func (p PrometheusState) Value() (driver.Value, error) {
	return p.String(), nil
}
