package api

// TODO: Move this.
import (
	"encoding/json"
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

const DateTimeFormat = "2006-01-02T15:04:05Z"

type DateTime struct {
	time.Time
}

func (d DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format(DateTimeFormat))
}

func (d *DateTime) UnmarshalJSON(data []byte) error {
	var dateStr string
	err := json.Unmarshal(data, &dateStr)
	if err != nil {
		return err
	}
	parsed, err := time.Parse(DateTimeFormat, dateStr)
	if err != nil {
		return err
	}
	d.Time = parsed
	return nil
}

func (d DateTime) String() string {
	return d.Time.Format(DateTimeFormat)
}

func (d *DateTime) UnmarshalText(data []byte) error {
	parsed, err := time.Parse(DateTimeFormat, string(data))
	if err != nil {
		return err
	}
	d.Time = parsed
	return nil
}

func NewDateTime(t time.Time) DateTime {
	return DateTime{
		Time: t,
	}
}

func NewDate(t time.Time) openapi_types.Date {
	return openapi_types.Date{
		Time: t,
	}
}

func DateTimeNow() DateTime {
	return NewDateTime(time.Now().UTC())
}

func DateNow() openapi_types.Date {
	return openapi_types.Date{
		Time: time.Now().UTC(),
	}
}
