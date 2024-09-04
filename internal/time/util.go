package time

// TODO: Change to timeutil.

import (
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

const (
	rfc3339TimeFormat string = "2006-01-02T15:04:05Z"
)

type DateTime time.Time

func (d DateTime) MarshalJSON() ([]byte, error) {
	jsonStr := "\"" + time.Time(d).Format(rfc3339TimeFormat) + "\""
	return []byte(jsonStr), nil
}

func (d *DateTime) UnmarshalJSON(b []byte) error {
	if len(b) < 2 || b[0] != '"' || b[len(b)-1] != '"' {
		return errors.New("not a json string")
	}

	// Strip the double quotes from the JSON string.
	b = b[1 : len(b)-1]

	// Parse the result using date time format.
	t, err := time.Parse(rfc3339TimeFormat, string(b))
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}

	*d = DateTime(t)
	return nil
}

func (d DateTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(time.Time(d).Format(rfc3339TimeFormat))
}

func (d *DateTime) UnmarshalBSONValue(typ bsontype.Type, data []byte) error {

	var dateString string
	if err := bson.UnmarshalValue(typ, data, &dateString); err != nil {
		return err
	}

	// Parse the result using date time format.
	t, err := time.Parse(rfc3339TimeFormat, dateString)
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}

	*d = DateTime(t)
	return nil
}

func (t DateTime) After(u DateTime) bool {
	return time.Time(t).After(time.Time(u))
}

func (t DateTime) Add(secs int64) DateTime {
	return DateTime(time.Time(t).Add(time.Duration(secs) * time.Second))
}

func (t DateTime) AddYears(years int) DateTime {
	return DateTime(time.Time(t).AddDate(years, 0, 0))
}

func (t DateTime) Unix() int64 {
	return time.Time(t).Unix()
}

func Now() DateTime {
	return DateTime(time.Now().UTC())
}
