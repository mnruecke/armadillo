package model

import "time"

const modelTimeFormat = time.RFC3339Nano
const modalTimePrecision = time.Millisecond
const modalTimeZone = "UTC"

type modelTime struct {
	time.Time
}

func (mt modelTime) format() modelTime {
	location, _ := time.LoadLocation(modalTimeZone)
	return modelTime{mt.Time.In(location).Round(modalTimePrecision)}
}

func (mt modelTime) formatString() string {
	return mt.format().Time.Format(modelTimeFormat)
}

func (mt modelTime) String() string {
	return mt.formatString()
}

func (mt modelTime) MarshalText() ([]byte, error) {
	return []byte(mt.formatString()), nil
}

func (mt modelTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + mt.formatString() + `"`), nil
}
